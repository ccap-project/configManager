package handlers

import (
	"log"

	"configManager/models"
	"configManager/restapi/operations/component"

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

func GetCellComponent(params component.GetCellComponentParams, principal *models.Customer) middleware.Responder {

	cellComponent, err := getCellComponent(principal.Name, params.CellID, params.ComponentID)

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return component.NewGetCellComponentInternalServerError()
	}

	if cellComponent == nil {
		return component.NewGetCellComponentOK()
	}

	return component.NewGetCellComponentOK().WithPayload(cellComponent)
}

func FindCellComponents(params component.FindCellComponentsParams, principal *models.Customer) middleware.Responder {

	cellComponents, err := findCellComponents(principal.Name, params.CellID)

	if err != nil {
		return component.NewFindCellComponentsInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	return component.NewFindCellComponentsOK().WithPayload(cellComponents)
}

func findCellComponents(customerName *string, CellID int64) ([]*models.Component, error) {
	cypher := `MATCH (c:Customer {name: {name} })-[:OWN]->(cell:Cell)-[:PROVIDES]->(component)
								WHERE id(cell) = {cell_id}
								RETURN ID(component) as id,
												component.name as name`

	db, err := driver.NewDriver().OpenNeo("bolt://192.168.20.54:7687")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return nil, err
	}
	defer db.Close()

	data, _, _, err := db.QueryNeoAll(cypher, map[string]interface{}{
		"name":    swag.StringValue(customerName),
		"cell_id": CellID})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return nil, err

	} else if len(data) == 0 {
		return nil, nil
	}

	res := make([]*models.Component, len(data))

	for idx, row := range data {
		res[idx], _ = getCellComponent(customerName, CellID, row[0].(int64))
		//_name := row[1].(string)
		//_roles, _ := _FindComponentRoles(params.CellID, row[0].(int64), principal)
		//_hostgroups, _ := _FindComponentHostgroups(principal.Name, params.CellID, row[0].(int64))

		//res[idx] = &models.Component{
		//	ID:         row[0].(int64),
		//	Name:       &_name,
		//	Roles:      _roles,
		//	Hostgroups: _hostgroups}
	}

	return res, nil
}

func getCellComponent(customerName *string, CellID int64, ComponentID int64) (*models.Component, error) {
	var component *models.Component
	component = nil

	cypher := `MATCH (c:Customer {name: {name} })-[:OWN]->(cell:Cell)-[:PROVIDES]->(component:Component)
							WHERE id(cell) = {cell_id} AND id(component) = {component_id}
								RETURN ID(component) as id,
												component.name as name`

	db, err := driver.NewDriver().OpenNeo("bolt://192.168.20.54:7687")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return component, err
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return component, err
	}

	defer stmt.Close()

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"name":         swag.StringValue(customerName),
		"cell_id":      CellID,
		"component_id": ComponentID})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return component, err
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		return component, err
	}

	_name := output[1].(string)
	_hostgroups, _ := _FindComponentHostgroups(customerName, CellID, ComponentID)

	component = &models.Component{
		ID:         output[0].(int64),
		Name:       &_name,
		Hostgroups: _hostgroups}

	return component, nil
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
