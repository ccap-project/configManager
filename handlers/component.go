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
							CREATE (cell)-[:PROVIDES]->(component:Component {
								id: {component_id},
								name: {component_name},
								order: {component_order} })
							RETURN	component.id AS id,
											component.name AS name`

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID})

	if _getComponentByName(ctx.rt, principal.Name, &params.CellID, params.Body.Name) != nil {
		ctxLogger.Error("component already exists !")
		return component.NewAddComponentConflict().WithPayload(&models.APIResponse{Message: "component already exists"})
	}

	db, err := ctx.rt.DB().OpenPool()

	if err != nil {
		ctxLogger.Error("error connecting to neo4j:", err)
		return component.NewAddComponentInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Errorf("An error occurred preparing statement: %s", err)
		return component.NewAddComponentInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	defer stmt.Close()

	if params.Body.Order == nil {
		params.Body.Order = swag.Int64(99)
	}

	ulid := configManager.GetULID()

	ctxLogger = ctx.rt.Logger().WithFields(logrus.Fields{
		"component_id":   ulid,
		"component_name": params.Body.Name})

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"name":            swag.StringValue(principal.Name),
		"cell_id":         params.CellID,
		"component_id":    ulid,
		"component_name":  swag.StringValue(params.Body.Name),
		"component_order": swag.Int64Value(params.Body.Order)})

	if err != nil {
		ctxLogger.Errorf("An error occurred querying Neo: %s", err)
		return component.NewAddComponentInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		ctxLogger.Errorf("An error occurred getting next row: %s", err)
		return component.NewAddComponentInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	ctxLogger.Info("OK")

	return component.NewAddComponentCreated().WithPayload(models.ULID(output[0].(string)))
}

func NewAddCellComponentRelationship(rt *configManager.Runtime) component.AddComponentRelationshipHandler {
	return &addComponentRelationship{rt: rt}
}

type addComponentRelationship struct {
	rt *configManager.Runtime
}

func (ctx *addComponentRelationship) Handle(params component.AddComponentRelationshipParams, principal *models.Customer) middleware.Responder {

	var cypher string
	entityType := _getEntityType(ctx.rt, &params.CellID, &params.EntityID)

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID,
		"component_id":  params.ComponentID,
		"entity_id":     params.EntityID})

	switch entityType {
	case "Loadbalancer":
		cypher = `
			MATCH (customer:Customer {name: {customer_name}})-[:OWN]->
				(cell:Cell {id: {cell_id}})-[:PROVIDES]->(component:Component {id: {component_id}})
			MATCH (cell)-[:HAS]->(lb:Loadbalancer {id: {entity_id}})
			MERGE (component)-[:CONNECT_TO]->(lb)
			RETURN *`
	case "Listener":
		cypher = `
			MATCH (customer:Customer {name: {customer_name}})-[:OWN]->
				(cell:Cell {id: {cell_id}})-[:PROVIDES]->(component:Component {id: {component_id}})
			MATCH (cell)-[:PROVIDES]->(component_t:Component)-[:LISTEN_ON]->
				(listener:Listener {id: {entity_id}})
			MERGE (component)-[:CONNECT_TO]->(listener)
			RETURN *`
	default:
		ctxLogger.Infof("entityType(%s)", entityType)
		return component.NewAddComponentRelationshipNotFound().WithPayload(&models.APIResponse{Message: "entity not found"})
	}

	db, err := ctx.rt.DB().OpenPool()

	if err != nil {
		ctxLogger.Error("error connecting to neo4j: ", err)
		return component.NewAddComponentRelationshipInternalServerError().WithPayload(&models.APIResponse{Message: "failure creating relationship"})
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: ", err)
		return component.NewAddComponentRelationshipInternalServerError().WithPayload(&models.APIResponse{Message: "failure creating relationship"})
	}

	defer stmt.Close()

	_, err = stmt.QueryNeo(map[string]interface{}{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID,
		"component_id":  params.ComponentID,
		"entity_id":     params.EntityID})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: ", err)
		return component.NewAddComponentRelationshipInternalServerError().WithPayload(&models.APIResponse{Message: "failure creating relationship"})
	}

	//	_, _, err = rows.NextNeo()
	//	if err != nil {
	//ctxLogger.Error("An error occurred getting next row: ", err)
	//return component.NewAddComponentRelationshipInternalServerError().WithPayload(&models.APIResponse{Message: "failure creating relationship"})
	//	}

	return component.NewAddComponentRelationshipCreated()
}

func NewDeleteCellComponentRelationship(rt *configManager.Runtime) component.DeleteComponentRelationshipHandler {
	return &deleteComponentRelationship{rt: rt}
}

type deleteComponentRelationship struct {
	rt *configManager.Runtime
}

func (ctx *deleteComponentRelationship) Handle(params component.DeleteComponentRelationshipParams, principal *models.Customer) middleware.Responder {

	var cypher string
	entityType := _getEntityType(ctx.rt, &params.CellID, &params.EntityID)

	switch entityType {
	case "Loadbalancer":
		cypher = `
			MATCH (customer:Customer {name: {customer_name}})-[:OWN]->
			 (cell:Cell {id: {cell_id}})-[:PROVIDES]->
			 (component:Component {id: {component_id}})-[r:CONNECT_TO]->
			 (entity:Loadbalancer {id: {entity_id}})
			delete r`

	case "Component":
		cypher = `
			MATCH (customer:Customer {name: {customer_name}})-[:OWN]->
			 (cell:Cell {id: {cell_id}})-[:PROVIDES]->
			 (component:Component {id: {component_id}})-[r:CONNECT_TO]->
			 (Listener)<-[:LISTEN_ON]-(entity:Component {id: {entity_id}})
			delete r`

	default:
		return component.NewDeleteComponentRelationshipNotFound().WithPayload(&models.APIResponse{Message: "entity not found"})
	}

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID,
		"component_id":  params.ComponentID,
		"entity_id":     params.EntityID})

	db, err := ctx.rt.DB().OpenPool()

	if err != nil {
		ctxLogger.Error("error connecting to neo4j: ", err)
		return component.NewDeleteComponentRelationshipInternalServerError().WithPayload(&models.APIResponse{Message: "failure deleting relationship"})
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: ", err)
		return component.NewDeleteComponentRelationshipInternalServerError().WithPayload(&models.APIResponse{Message: "failure deleting relationship"})
	}

	defer stmt.Close()

	_, err = stmt.ExecNeo(map[string]interface{}{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID,
		"component_id":  params.ComponentID,
		"entity_id":     params.EntityID})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: ", err)
		return component.NewDeleteComponentRelationshipInternalServerError().WithPayload(&models.APIResponse{Message: "failure deleting relationship"})
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

	cellComponent, err := _getCellComponent(ctx.rt, principal.Name, &params.CellID, &params.ComponentID)

	if err != nil {
		//log.Printf("An error occurred querying Neo: %s", err)
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

	cellComponents, err := _findCellComponents(ctx.rt, principal.Name, &params.CellID)

	if err != nil {
		return component.NewFindCellComponentsInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	return component.NewFindCellComponentsOK().WithPayload(cellComponents)
}

func _findCellComponents(rt *configManager.Runtime, customerName *string, CellID *string) ([]*models.Component, error) {

	var res []*models.Component

	cypher := `MATCH (c:Customer {name: {name} })-[:OWN]->
							(cell:Cell {id: {cell_id}})-[:PROVIDES]->(component)
								RETURN component.id as id, component.name as name`

	db, err := rt.DB().OpenPool()

	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return nil, err
	}
	defer db.Close()

	data, _, _, err := db.QueryNeoAll(cypher, map[string]interface{}{
		"name":    swag.StringValue(customerName),
		"cell_id": swag.StringValue(CellID)})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return nil, err

	} else if len(data) == 0 {
		return nil, nil
	}

	for _, row := range data {
		id := row[0].(string)
		c, _ := _getCellComponent(rt, customerName, CellID, &id)
		res = append(res, c)
	}

	return res, nil
}

func _listCellComponents(rt *configManager.Runtime, customerName *string, CellID *string) (*[]string, error) {

	var res []string

	cypher := `MATCH (c:Customer {name: {name} })-[:OWN]->
							(cell:Cell {id: {cell_id}})-[:PROVIDES]->(component)
								RETURN component.name as name`

	db, err := rt.DB().OpenPool()

	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return nil, err
	}
	defer db.Close()

	data, _, _, err := db.QueryNeoAll(cypher, map[string]interface{}{
		"name":    swag.StringValue(customerName),
		"cell_id": swag.StringValue(CellID)})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return nil, err

	} else if len(data) == 0 {
		return nil, nil
	}

	for _, row := range data {
		name := row[0].(string)
		res = append(res, name)
	}

	return &res, nil
}

func _findCellComponentRelationships(rt *configManager.Runtime, customerName *string, CellID *string, ComponentID *string) ([]models.ULID, error) {

	var res []models.ULID

	cypher := `MATCH (c:Customer {name: {name} })-[:OWN]->
							(cell:Cell {id: {cell_id}})-[:PROVIDES]->
							(component:Component {id: {component_id}})-[:CONNECT_TO]->(listener:Listener)
								RETURN listener.id as listener_id`

	db, err := rt.DB().OpenPool()

	ctxLogger := rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(customerName),
		"cell_id":       swag.StringValue(CellID),
		"component_id":  swag.StringValue(ComponentID)})

	if err != nil {
		ctxLogger.Error("error connecting to neo4j:", err)
		return nil, err
	}
	defer db.Close()

	data, _, _, err := db.QueryNeoAll(cypher, map[string]interface{}{
		"name":         swag.StringValue(customerName),
		"cell_id":      swag.StringValue(CellID),
		"component_id": swag.StringValue(ComponentID)})

	if err != nil {
		ctxLogger.Errorf("An error occurred querying Neo: %s", err)
		return nil, err

	} else if len(data) == 0 {
		return nil, nil
	}

	for _, row := range data {
		if len(row[0].(string)) > 0 {
			res = append(res, models.ULID(row[0].(string)))
		}
	}

	return res, nil
}

func _getCellComponent(rt *configManager.Runtime, customerName *string, CellID *string, ComponentID *string) (*models.Component, error) {
	var component *models.Component
	component = nil

	ctxLogger := rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(customerName),
		"cell_id":       swag.StringValue(CellID),
		"component_id":  swag.StringValue(ComponentID)})

	cypher := `MATCH (c:Customer {name: {name} })-[:OWN]->
							(cell:Cell {id: {cell_id}})-[:PROVIDES]->
							(component:Component {id: {component_id}})
								RETURN component.id as id,
												component.name as name`

	db, err := rt.DB().OpenPool()

	if err != nil {
		ctxLogger.Error("error connecting to neo4j: ", err)
		return component, err
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: ", err)
		return component, err
	}

	defer stmt.Close()

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"name":         swag.StringValue(customerName),
		"cell_id":      swag.StringValue(CellID),
		"component_id": swag.StringValue(ComponentID)})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: ", err)
		return component, err
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		return component, err
	}

	_name := output[1].(string)
	_hostgroups, _ := _FindComponentHostgroups(rt, customerName, CellID, ComponentID)
	_roles, _ := _findComponentRoles(rt, ComponentID)
	_listeners, _ := _findComponentListeners(rt, customerName, CellID, ComponentID)
	//_relationships, _ := _findCellComponentRelationships(rt, customerName, CellID, ComponentID)

	component = &models.Component{
		ID:         models.ULID(output[0].(string)),
		Name:       &_name,
		Hostgroups: _hostgroups,
		Roles:      _roles,
		Listeners:  _listeners}

	component.Relationships, _ = _findCellComponentRelationships(rt, customerName, CellID, ComponentID)

	return component, nil
}

func _getComponentByName(rt *configManager.Runtime, customerName *string, CellID *string, componentName *string) *models.Component {

	var component *models.Component
	component = nil

	cypher := `MATCH (c:Customer {name: {name} })-[:OWN]->
							(cell:Cell {id: {cell_id}})-[:PROVIDES]->(component:Component)
							WHERE component.name = {component_name}
								RETURN component.id as id,
												component.name as name`

	ctxLogger := rt.Logger().WithFields(logrus.Fields{
		"customer_name":  swag.StringValue(customerName),
		"cell_id":        swag.StringValue(CellID),
		"component_name": swag.StringValue(componentName)})

	db, err := rt.DB().OpenPool()

	if err != nil {
		ctxLogger.Error("error connecting to neo4j: ", err)
		return component
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: ", err)
		return component
	}
	defer stmt.Close()

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"name":           swag.StringValue(customerName),
		"cell_id":        swag.StringValue(CellID),
		"component_name": swag.StringValue(componentName)})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: ", err)
		return component
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		return component
	}
	_name := output[1].(string)

	component = &models.Component{
		ID:   models.ULID(output[0].(string)),
		Name: &_name}

	return component
}

func _getEntityType(rt *configManager.Runtime, CellID *string, EntityID *string) string {

	cypher := `MATCH (cell:Cell{id: {cell_id}})-[*]->(entity {id: {entity_id}})
						RETURN labels(entity)`

	ctxLogger := rt.Logger().WithFields(logrus.Fields{
		"cell_id":   swag.StringValue(CellID),
		"entity_id": swag.StringValue(EntityID)})

	db, err := rt.DB().OpenPool()

	if err != nil {
		ctxLogger.Error("error connecting to neo4j: ", err)
		return ""
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: ", err)
		return ""
	}
	defer stmt.Close()

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"cell_id":   swag.StringValue(CellID),
		"entity_id": swag.StringValue(EntityID)})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: ", err)
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

	return ""
}
