package handlers

import (
	"log"

	"../models"
	"../restapi/operations/provider"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	driver "github.com/johnnadratowski/golang-neo4j-bolt-driver"
)

func AddCellProvider(params provider.AddProviderParams, principal *models.Customer) middleware.Responder {

	cypher := `MATCH (c:Customer {name: {name} })-[:OWN]->(cell:Cell),
										(providertype:ProviderType)
							WHERE id(cell) = {cell_id} AND id(providertype) = {providertype_id}
							CREATE (cell)-[:USE]->(provider:Provider {
								name: {provider_name},
							 	domain_name: {domain_name},
								tenantname: {tenant_name},
								auth_url: {auth_url},
								username: {username},
								password: {password}})-[:PROVIDER_IS]->(providertype)
							RETURN	id(provider) AS id,
											provider.name AS name`

	Provider := getProvider(principal.Name, params.CellID)
	log.Printf("Here =>>>> %#v\n", Provider)

	if Provider != nil {
		log.Println("provider already exists !")
		return provider.NewAddProviderConflict().WithPayload(models.APIResponse{Message: "provider already exists"})
	}

	db, err := driver.NewDriver().OpenNeo("bolt://192.168.20.54:7687")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return provider.NewAddProviderInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return provider.NewAddProviderInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"name":            swag.StringValue(principal.Name),
		"cell_id":         params.CellID,
		"provider_name":   swag.StringValue(params.Body.Name),
		"domain_name":     swag.StringValue(params.Body.DomainName),
		"tenant_name":     swag.StringValue(params.Body.TenantName),
		"auth_url":        swag.StringValue(params.Body.AuthURL),
		"username":        swag.StringValue(params.Body.Username),
		"password":        swag.StringValue(params.Body.Password),
		"providertype_id": swag.Int64Value(params.Body.ProvidertypeID)})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return provider.NewAddProviderInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		log.Printf("An error occurred getting next row: %s", err)
		return provider.NewAddProviderInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	log.Printf("= Output(%#v)", output)

	log.Printf("customer(%s) name(%s) ", swag.StringValue(principal.Name), swag.StringValue(params.Body.Name))

	return provider.NewAddProviderCreated().WithPayload(output[0].(int64))
}

func GetCellProvider(params provider.GetProviderParams, principal *models.Customer) middleware.Responder {

	Provider := getProvider(principal.Name, params.CellID)

	if Provider == nil {
		log.Println("provider does not exists !")
		return provider.NewGetProviderNotFound()
	}

	return provider.NewGetProviderOK().WithPayload(Provider)
}

func getProvider(customerName *string, CellID int64) *models.Provider {

	var provider *models.Provider
	provider = nil

	cypher := `MATCH (c:Customer {name: {name} })-[:OWN]->(cell:Cell)-[:USE]->(provider:Provider)
							WHERE id(cell) = {cell_id}
								RETURN ID(provider) as id,
												provider.name as name,
												provider.domain_name as domain_name,
												provider.tenantname as tenantname,
												provider.auth_url as auth_url,
												provider.providertype_id as providertype_id,
												provider.username as username,
												provider.password as password`

	db, err := driver.NewDriver().OpenNeo("bolt://192.168.20.54:7687")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return provider
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return provider
	}
	defer stmt.Close()

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"name":    swag.StringValue(customerName),
		"cell_id": CellID})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return provider
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		return provider
	}

	provider = new(models.Provider)
	provider.Name = new(string)
	provider.DomainName = new(string)
	provider.TenantName = new(string)
	provider.AuthURL = new(string)
	provider.ProvidertypeID = new(int64)
	provider.Username = new(string)
	provider.Password = new(string)

	provider.ID = output[0].(int64)
	*provider.Name = output[1].(string)
	*provider.DomainName = output[2].(string)
	*provider.TenantName = output[3].(string)
	*provider.AuthURL = output[4].(string)
	*provider.ProvidertypeID = output[5].(int64)
	*provider.Username = output[6].(string)
	*provider.Password = output[7].(string)

	return provider
}

func UpdateCellProvider(params provider.UpdateProviderParams, principal *models.Customer) middleware.Responder {

	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->
							(cell:Cell)-[rel:USE]->(provider:Provider)-[rel2:PROVIDER_IS]->(provider_type:ProviderType)
							WHERE id(cell) = {cell_id}
						MATCH (newProviderType:ProviderType)
							WHERE id(newProviderType) = {providertype_id}
							SET provider.name={name},
									provider.domain_name={domain_name},
									provider.tenantname={tenant_name},
									provider.auth_url={auth_url},
									provider.username={username},
									provider.password={password},
									provider.providertype_id={providertype_id}
							DELETE rel, rel2
							CREATE (cell)-[:USE]->(provider)-[:PROVIDER_IS]->(newProviderType)
							return provider`

	Provider := getProvider(principal.Name, params.CellID)
	log.Printf("UpdateCellProvider =>>>> %#v\n", Provider)

	if Provider == nil {
		log.Println("provider does not exists !")
		return provider.NewUpdateProviderNotFound().WithPayload(models.APIResponse{Message: "provider does not exists"})
	}

	db, err := driver.NewDriver().OpenNeo("bolt://192.168.20.54:7687")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return provider.NewUpdateProviderInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return provider.NewUpdateProviderInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"customer_name":   swag.StringValue(principal.Name),
		"cell_id":         params.CellID,
		"name":            swag.StringValue(params.Body.Name),
		"domain_name":     swag.StringValue(params.Body.DomainName),
		"tenant_name":     swag.StringValue(params.Body.TenantName),
		"auth_url":        swag.StringValue(params.Body.AuthURL),
		"username":        swag.StringValue(params.Body.Username),
		"password":        swag.StringValue(params.Body.Password),
		"providertype_id": swag.Int64Value(params.Body.ProvidertypeID)})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return provider.NewUpdateProviderInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		log.Printf("An error occurred getting next row: %s", err)
		return provider.NewUpdateProviderInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	log.Printf("= Output(%#v)", output)

	//log.Printf("customer(%s) name(%s) ", swag.StringValue(principal.Name), swag.StringValue(params.Body.Name))

	return provider.NewUpdateProviderOK()
}
