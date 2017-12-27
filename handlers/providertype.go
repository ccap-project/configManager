package handlers

import (
	"fmt"
	"log"
	"strings"

	"configManager/models"
	"configManager/neo4j"
	"configManager/restapi/operations/providertype"

	"github.com/go-openapi/runtime/middleware"
)

func AddProviderType(params providertype.AddProviderTypeParams) middleware.Responder {

	cypher := `create(p:ProviderType { name: {name},
																			auth_url: {auth_url},
																			domain_name: {domain_name},
																			username: {username},
																			password: {password} }) RETURN ID(p)`

	if len(GetProviderTypeByName(params.Body.Name).Name) > 0 {
		log.Println("providertype already exists !")
		return providertype.NewAddProviderTypeInternalServerError().WithPayload(models.APIResponse{Message: "providertype already exists"})
	}

	db, err := neo4j.Connect("")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return providertype.NewAddProviderTypeInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Printf("An error occurred beginning transaction: %s", err)
		return providertype.NewAddProviderTypeInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}
	defer tx.Rollback()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return providertype.NewAddProviderTypeInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	defer stmt.Close()

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"name":        params.Body.Name,
		"auth_url":    params.Body.AuthURL,
		"domain_name": params.Body.DomainName,
		"username":    params.Body.Username,
		"password":    params.Body.Password})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return providertype.NewAddProviderTypeInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	_, _, err = rows.NextNeo()
	if err != nil {
		log.Printf("An error occurred getting next row: %s", err)
		return providertype.NewAddProviderTypeInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	tx.Commit()

	return providertype.NewAddProviderTypeCreated().WithPayload("OK")
}

func GetProviderTypeByID(params providertype.GetProviderTypeByIDParams) middleware.Responder {

	cypher := `MATCH (p:ProviderType)
							WHERE ID(p) = {id}
							RETURN ID(p) as id,
											p.name as name,
											p.auth_url as auth_url,
											p.domain_name as domain_name,
											p.username as username,
											p.password as password`

	db, err := neo4j.Connect("")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return providertype.NewGetProviderTypeByIDInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return providertype.NewGetProviderTypeByIDInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	defer stmt.Close()

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"customer_name": params.ProvidertypeID})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return providertype.NewGetProviderTypeByIDInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	if rows == nil {
		return providertype.NewGetProviderTypeByIDNotFound()
	}

	row, _, err := rows.NextNeo()

	if err != nil {
		log.Printf("An error occurred getting next row: %s", err)
		return providertype.NewListProviderTypesInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	provider := &models.ProviderType{
		ID:         row[0].(int64),
		Name:       row[1].(string),
		AuthURL:    row[1].(string),
		DomainName: row[2].(string),
		Username:   row[3].(string),
		Password:   row[4].(string)}

	return providertype.NewGetProviderTypeByIDOK().WithPayload(provider)
}

func GetProviderTypeByName(providertypeName string) models.ProviderType {

	var providerType models.ProviderType

	cypher := `MATCH (p:ProviderType)
							WHERE p.name = {name}
							RETURN ID(p) as id,
											p.name as name,
											p.auth_url as auth_url,
											p.domain_name as domain_name,
											p.username as username,
											p.password as password`

	db, err := neo4j.Connect("")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return providerType
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return providerType
	}

	defer stmt.Close()

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"name": providertypeName})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return providerType
	}

	if rows == nil {
		return providerType
	}

	row, _, err := rows.NextNeo()

	if err != nil {
		log.Printf("An error occurred getting next row: %s", err)
		return providerType
	}

	providerType.ID = row[0].(int64)
	providerType.Name = row[1].(string)
	providerType.AuthURL = row[2].(string)
	providerType.DomainName = row[3].(string)
	providerType.Username = row[4].(string)
	providerType.Password = row[5].(string)

	return providerType
}

func ListProviderTypes(params providertype.ListProviderTypesParams) middleware.Responder {

	cypher := `MATCH (p:ProviderType)
							RETURN ID(p) as id,
											p.name as name,
											p.auth_url as auth_url,
											p.domain_name as domain_name,
											p.username as username,
											p.password as password`

	db, err := neo4j.Connect("")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return providertype.NewListProviderTypesInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	data, _, _, err := db.QueryNeoAll(cypher, nil)

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return providertype.NewListProviderTypesInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	res := make([]*models.ProviderType, len(data))

	for idx, row := range data {
		res[idx] = &models.ProviderType{
			ID:         row[0].(int64),
			Name:       row[1].(string),
			AuthURL:    row[2].(string),
			DomainName: row[3].(string),
			Username:   row[4].(string),
			Password:   row[5].(string)}
	}

	return providertype.NewListProviderTypesOK().WithPayload(res)
}

func InitProviderType() {

	log.Printf("Checking provider types...")

	if err := addProviderType("Openstack", []string{"auth_url", "domain_name", "username", "password"}); err != nil {
		log.Println("Error Initializing provider types, ", err)
	}
}

func addProviderType(name string, fields []string) error {

	var allFields []string

	if len(GetProviderTypeByName(name).Name) > 0 {
		log.Printf("Provider %s already exists", name)
		return nil
	}

	createTmpl := `Create (p:ProviderType { %s: '%s', %s })`

	lastField := len(fields)

	if lastField <= 0 {
		return fmt.Errorf("No fields specified !")
	} else {
		lastField -= 1
	}

	for i := 0; i < lastField; i++ {
		log.Println(fields[i])
		allFields = append(allFields, fmt.Sprintf("%s: '%s', ", fields[i], fields[i]))
	}

	allFields = append(allFields, fmt.Sprintf("%s: '%s'", fields[lastField], fields[lastField]))

	create := fmt.Sprintf(createTmpl, name, name, strings.Join(allFields, ""))

	db, err := neo4j.Connect("")
	if err != nil {
		return fmt.Errorf("error connecting to neo4j:", err)
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(create)
	if err != nil {
		return fmt.Errorf("An error occurred preparing statement: %s", err)
	}

	defer stmt.Close()

	_, err = stmt.QueryNeo(nil)

	if err != nil {
		return fmt.Errorf("An error occurred querying Neo: %s", err)
	}

	log.Printf("Provider %s has been created", name)

	return nil
}
