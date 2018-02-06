/*
 *
 * Copyright (c) 2017, 2018 Alexandre Biancalana <ale@biancalanas.net>.
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *     * Redistributions of source code must retain the above copyright
 *       notice, this list of conditions and the following disclaimer.
 *     * Redistributions in binary form must reproduce the above copyright
 *       notice, this list of conditions and the following disclaimer in the
 *       documentation and/or other materials provided with the distribution.
 *     * Neither the name of the <organization> nor the
 *       names of its contributors may be used to endorse or promote products
 *       derived from this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL <COPYRIGHT HOLDER> BE LIABLE FOR ANY
 * DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
 * ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package handlers

import (
	"log"

	"configManager"
	"configManager/models"
	"configManager/neo4j"
	"configManager/restapi/operations/component"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
)

func NewAddCellComponent(rt *configManager.Runtime) component.AddComponentHandler {
	return &addCellComponent{rt: rt}
}

type addCellComponent struct {
	rt *configManager.Runtime
}

func (ctx *addCellComponent) Handle(params component.AddComponentParams, principal *models.Customer) middleware.Responder {

	cypher := `MATCH (c:Customer {name: {name} })-[:OWN]->(cell:Cell)
							WHERE id(cell) = {cell_id}
							CREATE (cell)-[:PROVIDES]->(component:Component { name: {component_name}, order: {component_order} })
							RETURN	id(component) AS id,
											component.name AS name`

	if _getComponentByName(ctx.rt.DB(), principal.Name, params.CellID, params.Body.Name) != nil {
		log.Println("component already exists !")
		return component.NewAddComponentConflict().WithPayload(models.APIResponse{Message: "component already exists"})
	}

	db, err := ctx.rt.DB().OpenPool()

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

	if params.Body.Order == nil {
		params.Body.Order = swag.Int64(99)
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{"name": swag.StringValue(principal.Name),
		"cell_id":         params.CellID,
		"component_name":  swag.StringValue(params.Body.Name),
		"component_order": swag.Int64Value(params.Body.Order)})

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

func NewGetCellComponent(rt *configManager.Runtime) component.GetCellComponentHandler {
	return &getCellComponent{rt: rt}
}

type getCellComponent struct {
	rt *configManager.Runtime
}

func (ctx *getCellComponent) Handle(params component.GetCellComponentParams, principal *models.Customer) middleware.Responder {

	cellComponent, err := _getCellComponent(ctx.rt.DB(), principal.Name, params.CellID, params.ComponentID)

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return component.NewGetCellComponentInternalServerError()
	}

	if cellComponent == nil {
		return component.NewGetCellComponentOK()
	}

	return component.NewGetCellComponentOK().WithPayload(cellComponent)
}

func NewFindCellComponents(rt *configManager.Runtime) component.FindCellComponentsHandler {
	return &findCellComponents{rt: rt}
}

type findCellComponents struct {
	rt *configManager.Runtime
}

func (ctx *findCellComponents) Handle(params component.FindCellComponentsParams, principal *models.Customer) middleware.Responder {

	cellComponents, err := _findCellComponents(ctx.rt.DB(), principal.Name, params.CellID)

	if err != nil {
		return component.NewFindCellComponentsInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	return component.NewFindCellComponentsOK().WithPayload(cellComponents)
}

func _findCellComponents(conn neo4j.ConnPool, customerName *string, CellID int64) ([]*models.Component, error) {
	cypher := `MATCH (c:Customer {name: {name} })-[:OWN]->(cell:Cell)-[:PROVIDES]->(component)
								WHERE id(cell) = {cell_id}
								RETURN ID(component) as id,
												component.name as name`

	db, err := conn.OpenPool()

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
		res[idx], _ = _getCellComponent(conn, customerName, CellID, row[0].(int64))
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

func _getCellComponent(conn neo4j.ConnPool, customerName *string, CellID int64, ComponentID int64) (*models.Component, error) {
	var component *models.Component
	component = nil

	cypher := `MATCH (c:Customer {name: {name} })-[:OWN]->(cell:Cell)-[:PROVIDES]->(component:Component)
							WHERE id(cell) = {cell_id} AND id(component) = {component_id}
								RETURN ID(component) as id,
												component.name as name`

	db, err := conn.OpenPool()

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
	_hostgroups, _ := _FindComponentHostgroups(conn, customerName, CellID, ComponentID)
	_roles, _ := _findComponentRoles(conn, ComponentID)

	component = &models.Component{
		ID:         output[0].(int64),
		Name:       &_name,
		Hostgroups: _hostgroups,
		Roles:      _roles}

	return component, nil
}

func _getComponentByName(conn neo4j.ConnPool, customerName *string, CellID int64, componentName *string) *models.Component {

	var component *models.Component
	component = nil

	cypher := `MATCH (c:Customer {name: {name} })-[:OWN]->(cell:Cell)-[:PROVIDES]->(component:Component)
							WHERE id(cell) = {cell_id} AND component.name = {component_name}
								RETURN ID(component) as id,
												component.name as name`

	db, err := conn.OpenPool()

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
