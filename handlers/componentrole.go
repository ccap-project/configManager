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
	"strings"

	"configManager"
	"configManager/models"
	"configManager/neo4j"
	"configManager/restapi/operations/role"

	"github.com/Sirupsen/logrus"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
)

func NewAddComponentRole(rt *configManager.Runtime) role.AddComponentRoleHandler {
	return &addComponentRole{rt: rt}
}

type addComponentRole struct {
	rt *configManager.Runtime
}

func (ctx *addComponentRole) Handle(params role.AddComponentRoleParams, principal *models.Customer) middleware.Responder {

	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->
							(cell:Cell {id: {cell_id}})-[:PROVIDES]->
							(component:Component {id: {component_id}})
						CREATE (component)-[:USE]->(role:Role {
							id: {role_id},
							name: {role_name},
							url: {role_url},
							version: {role_version},
							order: {role_order}} )
						RETURN role.id as id`

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID,
		"component_id":  params.ComponentID})

	if _getComponentRoleByName(ctx.rt, principal.Name, &params.CellID, &params.ComponentID, params.Body.Name) != nil {
		ctxLogger.Warn("role already exists !")
		return role.NewAddComponentRoleConflict().WithPayload(&models.APIResponse{Message: "role already exists"})
	}

	db, err := ctx.rt.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j: ", err)
		return role.NewAddComponentRoleInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		ctxLogger.Error("An error occurred beginning transaction: ", err)
		return role.NewAddComponentRoleInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer tx.Rollback()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: ", err)
		return role.NewAddComponentRoleInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer stmt.Close()

	if params.Body.Order == nil {
		params.Body.Order = swag.Int64(99)
	}
	ulid := configManager.GetULID()

	ctxLogger = ctx.rt.Logger().WithFields(logrus.Fields{
		"role_id": ulid})

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID,
		"component_id":  params.ComponentID,
		"role_id":       ulid,
		"role_name":     swag.StringValue(params.Body.Name),
		"role_url":      swag.StringValue(params.Body.URL),
		"role_version":  swag.StringValue(params.Body.Version),
		"role_order":    swag.Int64Value(params.Body.Order)})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: ", err)
		return role.NewAddComponentRoleInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		ctxLogger.Error("An error occurred getting next row: ", err)
		return role.NewAddComponentRoleInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	stmt.Close()

	err = addComponentRoleParameters(ctx.rt.Logger(), db, principal.Name, &params.CellID, &params.ComponentID, params.Body.Name, params.Body.Params)
	if err != nil {
		ctxLogger.Error("An error occurred adding Role parameters: ", err)
		return role.NewAddComponentRoleInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	tx.Commit()

	ctxLogger.Info("OK")

	return role.NewAddComponentRoleCreated().WithPayload(models.ULID(output[0].(string)))
}

func NewDeleteComponentRole(rt *configManager.Runtime) role.DeleteComponentRoleHandler {
	return &deleteComponentRole{rt: rt}
}

type deleteComponentRole struct {
	rt *configManager.Runtime
}

func (ctx *deleteComponentRole) Handle(params role.DeleteComponentRoleParams, principal *models.Customer) middleware.Responder {

	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->
							(cell:Cell {id: {cell_id}})-[:PROVIDES]->
							(component:Component {id: {component_id}})-[:USE]->(role:Role {name: {role_name}})
							OPTIONAL MATCH (role)-[r:PARAM]->(p)
							DETACH DELETE role, r, p`

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID,
		"role_name":     params.RoleName,
		"component_id":  params.ComponentID})

	if _getComponentRoleByName(ctx.rt, principal.Name, &params.CellID, &params.ComponentID, &params.RoleName) == nil {
		ctxLogger.Warn("role does not exists !")
		return role.NewDeleteComponentRoleNotFound()
	}

	db, err := ctx.rt.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j: ", err)
		return role.NewDeleteComponentRoleInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: ", err)
		return role.NewDeleteComponentRoleInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	defer stmt.Close()

	_, err = stmt.ExecNeo(map[string]interface{}{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID,
		"component_id":  params.ComponentID,
		"role_name":     params.RoleName})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: ", err)
		return role.NewAddComponentRoleInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	return role.NewDeleteComponentRoleOK()
}

func NewFindComponentRoles(rt *configManager.Runtime) role.FindComponentRolesHandler {
	return &findComponentRoles{rt: rt}
}

type findComponentRoles struct {
	rt *configManager.Runtime
}

func (ctx *findComponentRoles) Handle(params role.FindComponentRolesParams, principal *models.Customer) middleware.Responder {

	res, err := _findComponentRoles(ctx.rt, &params.ComponentID)

	if err != nil {
		return err
	}

	return role.NewFindComponentRolesOK().WithPayload(res)
}

func NewUpdateComponentRole(rt *configManager.Runtime) role.UpdateComponentRoleHandler {
	return &updateComponentRole{rt: rt}
}

type updateComponentRole struct {
	rt *configManager.Runtime
}

func (ctx *updateComponentRole) Handle(params role.UpdateComponentRoleParams, principal *models.Customer) middleware.Responder {

	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->
							(cell:Cell {id: {cell_id}})-[:PROVIDES]->
							(component:Component {id: {component_id}})-[:USE]->
							(role:Role{name: {role_current_name}})-[:PARAM]->(param:Parameter)
						SET role.name={role_new_name}, role.url={role_url}, role.version={role_version}, role.order={role_order}
						DETACH DELETE param`

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID,
		"role_name":     params.RoleName,
		"component_id":  params.ComponentID})

	if _getComponentRoleByName(ctx.rt, principal.Name, &params.CellID, &params.ComponentID, &params.RoleName) == nil {
		ctxLogger.Warn("role does not exists !")
		return role.NewUpdateComponentRoleNotFound()
	}

	if strings.Compare(params.RoleName, *params.Body.Name) != 0 &&
		_getComponentRoleByName(ctx.rt, principal.Name, &params.CellID, &params.ComponentID, params.Body.Name) != nil {
		ctxLogger.Warn("role target name already exists !")
		return role.NewUpdateComponentRoleConflict()
	}

	db, err := ctx.rt.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j: ", err)
		return role.NewAddComponentRoleInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		ctxLogger.Error("An error occurred beginning transaction: ", err)
		return role.NewAddComponentRoleInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer tx.Rollback()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: ", err)
		return role.NewAddComponentRoleInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer stmt.Close()

	_, err = stmt.ExecNeo(map[string]interface{}{
		"customer_name":     swag.StringValue(principal.Name),
		"cell_id":           params.CellID,
		"component_id":      params.ComponentID,
		"role_current_name": params.RoleName,
		"role_new_name":     swag.StringValue(params.Body.Name),
		"role_url":          swag.StringValue(params.Body.URL),
		"role_version":      swag.StringValue(params.Body.Version),
		"role_order":        swag.Int64Value(params.Body.Order)})

	if err != nil {
		ctxLogger.Error("Error occurred querying Neo: ", err)
		return role.NewAddComponentRoleInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	stmt.Close()

	if len(params.Body.Params) > 0 {
		err = addComponentRoleParameters(ctx.rt.Logger(), db, principal.Name, &params.CellID, &params.ComponentID, params.Body.Name, params.Body.Params)
		if err != nil {
			ctxLogger.Error("An error occurred adding Role parameters: %s", err)
			return role.NewAddComponentRoleInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
		}
	}

	tx.Commit()

	return role.NewUpdateComponentRoleOK()
}

func addComponentRoleParameters(
	logger logrus.FieldLogger,
	db neo4j.Conn,
	customer *string,
	cellID *string,
	componentID *string,
	roleName *string,
	params []*models.Parameter) error {

	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->
							(cell:Cell {id: {cell_id}})-[:PROVIDES]->
							(component:Component {id: {component_id}})-[:USE]->(role:Role{name: {role_name}})
						CREATE (role)-[:PARAM]->(param:Parameter {id: {param_id}, name: {param_name}, value: {param_val}} )
							RETURN param.id as id`

	ctxLogger := logger.WithFields(logrus.Fields{
		"customer_name": swag.StringValue(customer),
		"cell_id":       swag.StringValue(cellID),
		"component_id":  swag.StringValue(componentID),
		"roleName":      swag.StringValue(roleName)})

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: ", err)
		return err
	}
	defer stmt.Close()

	// add parameters
	for _, param := range params {
		ulid := configManager.GetULID()
		_, err := stmt.ExecNeo(map[string]interface{}{
			"customer_name": swag.StringValue(customer),
			"cell_id":       swag.StringValue(cellID),
			"component_id":  swag.StringValue(componentID),
			"role_name":     swag.StringValue(roleName),
			"param_id":      ulid,
			"param_name":    swag.StringValue(param.Name),
			"param_val":     swag.StringValue(param.Value)})

		if err != nil {
			ctxLogger.Error("An error occurred querying Neo: ", err)
			return err
		} else {
			ctxLogger.WithFields(logrus.Fields{
				"param_name": swag.StringValue(param.Name),
				"param_id":   ulid}).Info("OK")
		}
	}

	return nil
}

func _findComponentRoles(rt *configManager.Runtime, ComponentID *string) ([]*models.Role, middleware.Responder) {

	cypher := `MATCH (component:Component {id: {component_id}})-[:USE]->(role:Role)
								RETURN role.id as id,
												role.name as name,
												role.url as url,
												role.version as version,
												role.order as order`

	ctxLogger := rt.Logger().WithFields(logrus.Fields{
		"component_id": swag.StringValue(ComponentID)})

	db, err := rt.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j: ", err)
		return nil, role.NewFindComponentRolesInternalServerError()
	}
	defer db.Close()

	data, _, _, err := db.QueryNeoAll(cypher, map[string]interface{}{
		"component_id": swag.StringValue(ComponentID)})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: ", err)
		return nil, role.NewFindComponentRolesInternalServerError()

	}

	res := make([]*models.Role, len(data))

	for idx, row := range data {
		_name := row[1].(string)
		_url := ""
		_version := ""

		var _order int64
		_order = 99

		if row[2] != nil {
			_url = row[2].(string)
		}

		if row[3] != nil {
			_version = row[3].(string)
		}

		if row[4] != nil {
			_order = row[4].(int64)
		}
		role_id := row[0].(string)
		_params, _ := findComponentRoleParameters(rt, &role_id)

		res[idx] = &models.Role{
			ID:      models.ULID(role_id),
			Name:    &_name,
			Version: &_version,
			URL:     &_url,
			Order:   &_order,
			Params:  _params}
	}

	return res, nil
}

func findComponentRoleParameters(rt *configManager.Runtime, roleID *string) ([]*models.Parameter, middleware.Responder) {

	cypher := `MATCH (role:Role {id: {role_id}})-[:PARAM]->(param:Parameter)
							RETURN param.id as id,
											param.name as name,
											param.value as value`

	ctxLogger := rt.Logger().WithFields(logrus.Fields{
		"role_id": swag.StringValue(roleID)})

	db, err := rt.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j: ", err)
		return nil, role.NewFindComponentRolesInternalServerError()
	}
	defer db.Close()

	data, _, _, err := db.QueryNeoAll(cypher, map[string]interface{}{
		"role_id": swag.StringValue(roleID)})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: ", err)
		return nil, role.NewFindComponentRolesInternalServerError()

	}

	res := make([]*models.Parameter, len(data))

	for idx, row := range data {
		_name := row[1].(string)
		_value := ""

		if row[2] != nil {
			_value = row[2].(string)
		}

		res[idx] = &models.Parameter{
			ID:    models.ULID(row[0].(string)),
			Name:  &_name,
			Value: &_value}
	}
	return res, nil
}

func _getComponentRoleByName(rt *configManager.Runtime, customer *string, cellID *string, componentID *string, roleName *string) *models.Role {

	var role *models.Role
	role = nil

	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->
							(cell:Cell {id: {cell_id}})-[:PROVIDES]->
							(component:Component {id: {component_id}})-[:USE]->
							(role:Role {name: {role_name}})
						RETURN role.id as id,
										role.name as name,
										role.url as url,
										role.version as version,
										role.order as order`

	ctxLogger := rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(customer),
		"cell_id":       swag.StringValue(cellID),
		"role_name":     swag.StringValue(roleName),
		"component_id":  swag.StringValue(componentID)})

	db, err := rt.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j: ", err)
		return role
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: ", err)
		return role
	}
	defer stmt.Close()

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"customer_name": swag.StringValue(customer),
		"cell_id":       swag.StringValue(cellID),
		"component_id":  swag.StringValue(componentID),
		"role_name":     swag.StringValue(roleName)})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: ", err)
		return role
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		//log.Printf("An error occurred fetching row: %s", err)
		return role
	}

	var _order int64

	_name := output[1].(string)
	_url := ""
	_version := output[3].(string)

	if output[2] != nil {
		_url = output[2].(string)
	}

	if output[4] != nil {
		_order = output[4].(int64)
	}

	role = &models.Role{
		ID:      models.ULID(output[0].(string)),
		Name:    &_name,
		URL:     &_url,
		Version: &_version,
		Order:   &_order}

	return role
}
