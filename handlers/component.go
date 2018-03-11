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

	"github.com/Sirupsen/logrus"
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

	cypher := `MATCH (c:Customer {name: {name} })-[:OWN]->(cell:Cell {id: {cell_id}})
							CREATE (cell)-[:PROVIDES]->(component:Component { name: {component_name}, order: {component_order} })
							RETURN	id(component) AS id,
											component.name AS name`

	if _getComponentByName(ctx.rt.DB(), principal.Name, &params.CellID, params.Body.Name) != nil {
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

func NewAddCellComponentRelationship(rt *configManager.Runtime) component.AddComponentRelationshipHandler {
	return &addComponentRelationship{rt: rt}
}

type addComponentRelationship struct {
	rt *configManager.Runtime
}

func (ctx *addComponentRelationship) Handle(params component.AddComponentRelationshipParams, principal *models.Customer) middleware.Responder {

	var cypher string
	entityType := _getEntityType(ctx.rt.DB(), &params.CellID, params.EntityID)

	switch entityType {
	case "Loadbalancer":
		cypher = `
			MATCH (customer:Customer {name: {customer_name}})-[:OWN]->(cell:Cell {id: {cell_id}})-[:PROVIDES]->(component:Component)
			WHERE id(component) = {component_id}
			MATCH (cell)-[:HAS]->(lb:Loadbalancer)
			WHERE id(lb) = {entity_id}
			MERGE (component)-[:CONNECT_TO]->(lb)
			RETURN *`
	case "Component":
		cypher = `
			MATCH (customer:Customer {name: {customer_name}})-[:OWN]->(cell:Cell {id: {cell_id}})-[:PROVIDES]->(component:Component)
			WHERE id(component) = {component_id}
			MATCH (cell)-[:PROVIDES]->(component_t:Component)-[:LISTEN_ON]->(listener:Listener)
			WHERE id(component_t) = {entity_id}
			MERGE (component)-[:CONNECT_TO]->(listener)
			RETURN *`
	default:
		return component.NewAddComponentRelationshipNotFound().WithPayload(models.APIResponse{Message: "entity not found"})
	}

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID,
		"component_id":  params.ComponentID,
		"entity_id":     params.EntityID})

	db, err := ctx.rt.DB().OpenPool()

	if err != nil {
		ctxLogger.Warn("error connecting to neo4j: ", err)
		return component.NewAddComponentRelationshipInternalServerError().WithPayload(models.APIResponse{Message: "failure creating relationship"})
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Warn("An error occurred preparing statement: ", err)
		return component.NewAddComponentRelationshipInternalServerError().WithPayload(models.APIResponse{Message: "failure creating relationship"})
	}

	defer stmt.Close()

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID,
		"component_id":  params.ComponentID,
		"entity_id":     params.EntityID})

	ctxLogger.Info("rows", rows)

	if err != nil {
		ctxLogger.Warn("An error occurred querying Neo: ", err)
		return component.NewAddComponentRelationshipInternalServerError().WithPayload(models.APIResponse{Message: "failure creating relationship"})
	}

	_, _, err = rows.NextNeo()
	if err != nil {
		ctxLogger.Warn("An error occurred getting next row: ", err)
		return component.NewAddComponentRelationshipInternalServerError().WithPayload(models.APIResponse{Message: "failure creating relationship"})
	}

	return component.NewAddComponentRelationshipCreated().WithPayload(1)
}

func NewDeleteCellComponentRelationship(rt *configManager.Runtime) component.DeleteComponentRelationshipHandler {
	return &deleteComponentRelationship{rt: rt}
}

type deleteComponentRelationship struct {
	rt *configManager.Runtime
}

func (ctx *deleteComponentRelationship) Handle(params component.DeleteComponentRelationshipParams, principal *models.Customer) middleware.Responder {

	var cypher string
	entityType := _getEntityType(ctx.rt.DB(), &params.CellID, params.EntityID)

	switch entityType {
	case "Loadbalancer":
		cypher = `
			MATCH (customer:Customer {name: {customer_name}})-[:OWN]->
			 (cell:Cell {id: {cell_id}})-[:PROVIDES]->(component:Component)-[r:CONNECT_TO]->(entity:Loadbalancer)
			WHERE id(component) = {component_id} AND id(entity) = {entity_id}
			delete r`

	case "Component":
		cypher = `
			MATCH (customer:Customer {name: {customer_name}})-[:OWN]->
			 (cell:Cell {id: {cell_id}})-[:PROVIDES]->(component:Component)-[r:CONNECT_TO]->(Listener)<-[:LISTEN_ON]-(entity:Component)
			WHERE id(component) = {component_id} AND id(entity) = {entity_id}
			delete r`

	default:
		return component.NewDeleteComponentRelationshipNotFound().WithPayload(models.APIResponse{Message: "entity not found"})
	}

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID,
		"component_id":  params.ComponentID,
		"entity_id":     params.EntityID})

	db, err := ctx.rt.DB().OpenPool()

	if err != nil {
		ctxLogger.Warn("error connecting to neo4j: ", err)
		return component.NewDeleteComponentRelationshipInternalServerError().WithPayload(models.APIResponse{Message: "failure deleting relationship"})
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Warn("An error occurred preparing statement: ", err)
		return component.NewDeleteComponentRelationshipInternalServerError().WithPayload(models.APIResponse{Message: "failure deleting relationship"})
	}

	defer stmt.Close()

	_, err = stmt.ExecNeo(map[string]interface{}{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID,
		"component_id":  params.ComponentID,
		"entity_id":     params.EntityID})

	if err != nil {
		ctxLogger.Warn("An error occurred querying Neo: ", err)
		return component.NewDeleteComponentRelationshipInternalServerError().WithPayload(models.APIResponse{Message: "failure deleting relationship"})
	}

	return component.NewDeleteComponentRelationshipOK()
}

func NewGetCellComponent(rt *configManager.Runtime) component.GetCellComponentHandler {
	return &getCellComponent{rt: rt}
}

type getCellComponent struct {
	rt *configManager.Runtime
}

func (ctx *getCellComponent) Handle(params component.GetCellComponentParams, principal *models.Customer) middleware.Responder {

	cellComponent, err := _getCellComponent(ctx.rt.DB(), principal.Name, &params.CellID, params.ComponentID)

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

	cellComponents, err := _findCellComponents(ctx.rt.DB(), principal.Name, &params.CellID)

	if err != nil {
		return component.NewFindCellComponentsInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	return component.NewFindCellComponentsOK().WithPayload(cellComponents)
}

func _findCellComponents(conn neo4j.ConnPool, customerName *string, CellID *string) ([]*models.Component, error) {
	cypher := `MATCH (c:Customer {name: {name} })-[:OWN]->(cell:Cell {id: {cell_id}})-[:PROVIDES]->(component)
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

func _getCellComponent(conn neo4j.ConnPool, customerName *string, CellID *string, ComponentID int64) (*models.Component, error) {
	var component *models.Component
	component = nil

	cypher := `MATCH (c:Customer {name: {name} })-[:OWN]->(cell:Cell {id: {cell_id}})-[:PROVIDES]->(component:Component)
							WHERE id(component) = {component_id}
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

func _getComponentByName(conn neo4j.ConnPool, customerName *string, CellID *string, componentName *string) *models.Component {

	var component *models.Component
	component = nil

	cypher := `MATCH (c:Customer {name: {name} })-[:OWN]->(cell:Cell {id: {cell_id}})-[:PROVIDES]->(component:Component)
							WHERE component.name = {component_name}
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

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"name":           swag.StringValue(customerName),
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

func _getEntityType(conn neo4j.ConnPool, CellID *string, EntityID int64) string {

	cypher := `MATCH (cell:Cell{id: {cell_id}})-->(entity)
							WHERE id(entity) = {entity_id}
								RETURN labels(entity)`

	db, err := conn.OpenPool()

	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return ""
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return ""
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"cell_id":   CellID,
		"entity_id": EntityID})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return ""
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		return ""
	}
	_labels := output[0]

	switch x := _labels.(type) {
	case []interface{}:
		return (x[0].(string))
	}
	stmt.Close()

	return ""
}
