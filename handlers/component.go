package handlers

import (
	"log"

	"../models"
	"../restapi/operations/component"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	driver "github.com/johnnadratowski/golang-neo4j-bolt-driver"
)

func AddCellComponent(params component.AddComponentParams, principal *models.Customer) middleware.Responder {

	cypher := `MATCH (c:Customer {name: {name} })-[:OWN]->(cell:Cell)
							WHERE id(cell) = {cell_id}
							CREATE (cell)-[:PROVIDES]->(component:Component { name: {component_name} })
							RETURN	id(component) AS id,
											component.name AS name`

	if getComponentByName(principal.Name, params.CellID, params.Body.Name) != nil {
		log.Println("component already exists !")
		return component.NewAddComponentConflict().WithPayload(models.APIResponse{Message: "component already exists"})
	}

	db, err := driver.NewDriver().OpenNeo("bolt://192.168.20.54:7687")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return component.NewAddComponentInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return component.NewAddComponentInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{"name": swag.StringValue(principal.Name),
		"cell_id":        params.CellID,
		"component_name": swag.StringValue(params.Body.Name)})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return component.NewAddComponentInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		log.Printf("An error occurred getting next row: %s", err)
		return component.NewAddComponentInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	log.Printf("= Output(%#v)", output)

	log.Printf("customer(%s) name(%s) ", swag.StringValue(principal.Name), swag.StringValue(params.Body.Name))

	return component.NewAddComponentCreated().WithPayload(output[0].(int64))
}

/*
func GetComponentByID(params component.GetComponentByIDParams, principal *models.Customer) middleware.Responder {

	cypher := `MATCH (c:Customer {name: {name} })-[:HAS]->(k:Component)
								WHERE ID(k) = {kid}
								RETURN ID(c) as id,
												k.name as name,
												k.public_key as public_key`

	db, err := driver.NewDriver().OpenNeo("bolt://192.168.20.54:7687")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return component.NewGetComponentByIDInternalServerError()
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return component.NewGetComponentByIDInternalServerError()
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{"name": swag.StringValue(principal.Name),
		"kid": params.ComponentID})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return component.NewGetComponentByIDInternalServerError()
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		return component.NewGetComponentByIDNotFound()
	}
	_name := output[1].(string)
	_pubkey := output[2].(string)

	_component := &models.Component{ID: output[0].(int64),
		Name:      &_name,
		PublicKey: &_pubkey}

	stmt.Close()

	return component.NewGetComponentByIDOK().WithPayload(_component)
}
*/
func FindCellComponents(params component.FindCellComponentsParams, principal *models.Customer) middleware.Responder {
	cypher := `MATCH (c:Customer {name: {name} })-[:OWN]->(cell:Cell)-[:PROVIDES]->(component)
								WHERE id(cell) = {cell_id}
								RETURN ID(component) as id,
												component.name as name`

	db, err := driver.NewDriver().OpenNeo("bolt://192.168.20.54:7687")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return component.NewFindCellComponentsInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	data, _, _, err := db.QueryNeoAll(cypher, map[string]interface{}{"name": swag.StringValue(principal.Name),
		"cell_id": params.CellID})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return component.NewFindCellComponentsInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})

	} else if len(data) == 0 {
		return component.NewFindCellComponentsNotFound()
	}

	res := make([]*models.Component, len(data))

	for idx, row := range data {
		_name := row[1].(string)

		res[idx] = &models.Component{
			ID:   row[0].(int64),
			Name: &_name}
	}

	return component.NewFindCellComponentsOK().WithPayload(res)
}

func getComponentByName(customerName *string, CellID int64, componentName *string) *models.Component {

	var component *models.Component
	component = nil

	cypher := `MATCH (c:Customer {name: {name} })-[:OWN]->(cell:Cell)-[:PROVIDES]->(component:Component)
							WHERE id(cell) = {cell_id} AND component.name = {component_name}
								RETURN ID(component) as id,
												component.name as name`

	db, err := driver.NewDriver().OpenNeo("bolt://192.168.20.54:7687")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return component
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return component
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{"name": swag.StringValue(customerName),
		"cell_id":        CellID,
		"component_name": swag.StringValue(componentName)})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return component
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		return component
	}
	_name := output[1].(string)

	component = &models.Component{ID: output[0].(int64),
		Name: &_name}

	stmt.Close()

	return component
}
