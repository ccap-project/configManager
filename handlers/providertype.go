package handlers

import (
	"log"
	"strings"

	"../models"
	"../restapi/operations/providertype"

	"github.com/go-openapi/runtime/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	_ "gopkg.in/cq.v1"
)

func AddProviderType(params providertype.AddProviderTypeParams) middleware.Responder {

	//log.Println(params.Body.AuthURL)

	neo4jURL := `http://192.168.20.54:7474`

	cypher := `create(p:ProviderType { name: {0},
																			auth_url: {1},
																			domain_name: {2},
																			username: {3},
																			password: {4} }) RETURN ID(p)`

	if len(GetProviderTypeByName(params.Body.Name).Name) > 0 {
		log.Println("providertype already exists !")
		return providertype.NewAddProviderTypeOK().WithPayload(models.APIResponse{Message: "providertype already exists"})
	}

	db, err := sqlx.Connect("neo4j-cypher", neo4jURL)
	if err != nil {
		log.Println("error connecting to neo4j:", err)
	}
	defer db.Close()

	stmt, err := db.Prepare(cypher)

	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(params.Body.Name,
		params.Body.AuthURL,
		params.Body.DomainName,
		params.Body.Username,
		params.Body.Password)

	log.Printf("%#v", res)
	if err != nil {
		log.Printf("error creating providertype name(%s): %v", params.Body.Name, err)
		return providertype.NewAddProviderTypeInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	return providertype.NewAddProviderTypeCreated().WithPayload("OK")
}

func GetProviderTypeByID(params providertype.GetProviderTypeByIDParams) middleware.Responder {

	neo4jURL := `http://192.168.20.54:7474`

	cypher := `MATCH (p:ProviderType)
							WHERE p.name = {0}
							RETURN ID(p) as id,
											p.name as name,
											p.auth_url as auth_url,
											p.domain_name as domain_name,
											p.username as username,
											p.password as password`

	db, err := sqlx.Connect("neo4j-cypher", neo4jURL)
	if err != nil {
		log.Println("error connecting to neo4j:", err)
	}
	defer db.Close()
	db.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)

	providerType := models.ProviderType{}

	err = db.Get(&providerType, cypher, params.ProvidertypeID)

	if err != nil {
		log.Printf("error getting providertype id(%d): %v", params.ProvidertypeID, err)
		return providertype.NewGetProviderTypeByIDNotFound()
	}

	return providertype.NewGetProviderTypeByIDOK().WithPayload(&providerType)
}

func GetProviderTypeByName(providertypeName string) models.ProviderType {

	neo4jURL := `http://192.168.20.54:7474`

	cypher := `MATCH (p:ProviderType)
							WHERE p.name =~ {0}
							RETURN p.name as name,
											p.auth_url as auth_url,
											p.domain_name as domain_name,
											p.username as username,
											p.password as password`

	db, err := sqlx.Connect("neo4j-cypher", neo4jURL)
	db.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)
	if err != nil {
		log.Println("error connecting to neo4j:", err)
	}
	defer db.Close()

	providerType := models.ProviderType{}

	Err := db.Get(&providerType, cypher, providertypeName)
	if Err != nil &&
		!(Err.Error() == "sql: Scan error on column index 1: unsupported Scan, storing driver.Value type <nil> into type *string" ||
			Err.Error() == "sql: no rows in result set") {
		log.Printf("error getting providertype by name(%s): %v", providertypeName, Err)
		log.Printf("providertype(%#v)", providerType)
		return (models.ProviderType{})
	}

	log.Printf("providertype(%#v)", providerType)
	return providerType
}

func ListProviderTypes(params providertype.ListProviderTypesParams) middleware.Responder {

	neo4jURL := `http://192.168.20.54:7474`

	cypher := `MATCH (p:ProviderType)
							RETURN ID(p) as id,
											p.name as name,
											p.auth_url as auth_url,
											p.domain_name as domain_name,
											p.username as username,
											p.password as password`

	db, err := sqlx.Connect("neo4j-cypher", neo4jURL)
	db.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)
	if err != nil {
		log.Println("error connecting to neo4j:", err)
	}
	defer db.Close()

	providerTypesList := []*models.ProviderType{}

	err = db.Select(&providerTypesList, cypher)

	if err != nil {
		log.Printf("error listing providertypes: %v", err)
		return providertype.NewListProviderTypesInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	log.Printf("providertypes(%#v)", providerTypesList)

	return providertype.NewListProviderTypesOK().WithPayload(providerTypesList)
}
