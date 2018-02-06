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
	"strings"

	"configManager"
	"configManager/models"
	"configManager/neo4j"
	"configManager/restapi/operations/role"

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

	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->(cell:Cell)-[:PROVIDES]->(component:Component)
							WHERE id(cell) = {cell_id} AND id(component) = {component_id}
							CREATE (component)-[:USE]->(role:Role {name: {role_name}, url: {role_url}, version: {role_version}, order: {role_order}} )
								RETURN id(role) as id`

	//log.Printf("= getRoleByName(%s), (%#v)", params.Body.Name, _getComponentRoleByName(principal.Name, params.CellID, params.ComponentID, params.Body.Name))

	if _getComponentRoleByName(ctx.rt.DB(), principal.Name, params.CellID, params.ComponentID, params.Body.Name) != nil {
		log.Println("role already exists !")
		return role.NewAddComponentRoleConflict().WithPayload(models.APIResponse{Message: "role already exists"})
	}

	db, err := ctx.rt.DB().OpenPool()
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return role.NewAddComponentRoleInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Printf("An error occurred beginning transaction: %s", err)
		return role.NewAddComponentRoleInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}
	defer tx.Rollback()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return role.NewAddComponentRoleInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	if params.Body.Order == nil {
		params.Body.Order = swag.Int64(99)
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID,
		"component_id":  params.ComponentID,
		"role_name":     swag.StringValue(params.Body.Name),
		"role_url":      swag.StringValue(params.Body.URL),
		"role_version":  swag.StringValue(params.Body.Version),
		"role_order":    swag.Int64Value(params.Body.Order)})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return role.NewAddComponentRoleInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		log.Printf("An error occurred getting next row: %s", err)
		return role.NewAddComponentRoleInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	log.Printf("= Output(%#v)", output)

	log.Printf("name(%s)", swag.StringValue(params.Body.Name))

	stmt.Close()

	err = addComponentRoleParameters(ctx.rt.DB(), principal.Name, params.CellID, params.ComponentID, params.Body.Name, params.Body.Params)
	if err != nil {
		log.Printf("An error occurred adding Role parameters: %s", err)
		return role.NewAddComponentRoleInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}
	tx.Commit()

	return role.NewAddComponentRoleCreated().WithPayload(output[0].(int64))
}

func NewDeleteComponentRole(rt *configManager.Runtime) role.DeleteComponentRoleHandler {
	return &deleteComponentRole{rt: rt}
}

type deleteComponentRole struct {
	rt *configManager.Runtime
}

func (ctx *deleteComponentRole) Handle(params role.DeleteComponentRoleParams, principal *models.Customer) middleware.Responder {

	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->
							(cell:Cell)-[:PROVIDES]->(component:Component)-[:USE]->(role:Role {name: {role_name}})
							OPTIONAL MATCH (role)-[r:PARAM]->(p)
							WHERE id(cell) = {cell_id} AND id(component) = {component_id}
							DETACH DELETE role, r, p`

	if _getComponentRoleByName(ctx.rt.DB(), principal.Name, params.CellID, params.ComponentID, &params.RoleName) == nil {
		log.Println("role does not exists !")
		return role.NewDeleteComponentRoleNotFound()
	}

	db, err := ctx.rt.DB().OpenPool()
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return role.NewDeleteComponentRoleInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return role.NewDeleteComponentRoleInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	defer stmt.Close()

	_, err = stmt.ExecNeo(map[string]interface{}{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID,
		"component_id":  params.ComponentID,
		"role_name":     params.RoleName})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return role.NewAddComponentRoleInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
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

	res, err := _findComponentRoles(ctx.rt.DB(), params.ComponentID)

	if err != nil {
		return err
	}

	log.Printf("= Res(%#v)", res)

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
							(cell:Cell)-[:PROVIDES]->(component:Component)-[:USE]->(role:Role{name: {role_current_name}})-[:PARAM]->(param:Parameter)
							WHERE id(cell) = {cell_id} AND id(component) = {component_id}
							SET role.name={role_new_name}, role.url={role_url}, role.version={role_version}, role.order={role_order}
							DETACH DELETE param`

	//log.Printf("= getRoleByName(%s), (%v)", params.Body.Name, _getComponentRoleByName(principal.Name, params.CellID, params.ComponentID, params.Body.Name))

	if _getComponentRoleByName(ctx.rt.DB(), principal.Name, params.CellID, params.ComponentID, &params.RoleName) == nil {
		log.Println("role does not exists !")
		return role.NewUpdateComponentRoleNotFound()
	}

	if strings.Compare(params.RoleName, *params.Body.Name) != 0 &&
		_getComponentRoleByName(ctx.rt.DB(), principal.Name, params.CellID, params.ComponentID, params.Body.Name) != nil {
		log.Println("role target name already exists !")
		return role.NewUpdateComponentRoleConflict()
	}

	db, err := ctx.rt.DB().OpenPool()
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return role.NewAddComponentRoleInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Printf("An error occurred beginning transaction: %s", err)
		return role.NewAddComponentRoleInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}
	defer tx.Rollback()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return role.NewAddComponentRoleInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
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
		log.Printf("-> An error occurred querying Neo: %s", err)
		return role.NewAddComponentRoleInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	stmt.Close()

	if len(params.Body.Params) > 0 {
		err = addComponentRoleParameters(ctx.rt.DB(), principal.Name, params.CellID, params.ComponentID, params.Body.Name, params.Body.Params)
		if err != nil {
			log.Printf("An error occurred adding Role parameters: %s", err)
			return role.NewAddComponentRoleInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
		}
	}

	tx.Commit()

	return role.NewUpdateComponentRoleOK()
}

func addComponentRoleParameters(conn neo4j.ConnPool, customer *string, cellID int64, componentID int64, roleName *string, params []*models.Parameter) error {

	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->
							(cell:Cell)-[:PROVIDES]->(component:Component)-[:USE]->(role:Role{name: {role_name}})
							WHERE id(cell) = {cell_id} AND id(component) = {component_id}
							CREATE (role)-[:PARAM]->(param:Parameter {name: {param_name}, value: {param_val}} )
								RETURN id(param) as id`

	db, err := conn.OpenPool()
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return err
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return err
	}

	// add parameters
	for _, param := range params {
		log.Printf("name(%s) val(%s)", param.Name, param.Value)
		//log.Printf("param(%#v)", param)
		_, err := stmt.ExecNeo(map[string]interface{}{
			"customer_name": swag.StringValue(customer),
			"cell_id":       cellID,
			"component_id":  componentID,
			"role_name":     swag.StringValue(roleName),
			"param_name":    swag.StringValue(param.Name),
			"param_val":     swag.StringValue(param.Value)})

		if err != nil {
			log.Printf("An error occurred querying Neo: %s", err)
			return err
		}
	}

	return nil
}

func _findComponentRoles(conn neo4j.ConnPool, ComponentID int64) ([]*models.Role, middleware.Responder) {

	cypher := `MATCH (component:Component)-[:USE]->(role:Role)
							WHERE id(component) = {component_id}
								RETURN id(role) as id,
												role.name as name,
												role.url as url,
												role.version as version,
												role.order as order`

	db, err := conn.OpenPool()
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return nil, role.NewFindComponentRolesInternalServerError()
	}
	defer db.Close()

	data, _, _, err := db.QueryNeoAll(cypher, map[string]interface{}{
		"component_id": ComponentID})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
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

		_params, _ := findComponentRoleParameters(conn, row[0].(int64))

		res[idx] = &models.Role{
			ID:      row[0].(int64),
			Name:    &_name,
			Version: &_version,
			URL:     &_url,
			Order:   &_order,
			Params:  _params}
	}

	return res, nil
}

func findComponentRoleParameters(conn neo4j.ConnPool, roleID int64) ([]*models.Parameter, middleware.Responder) {

	cypher := `MATCH (role:Role)-[:PARAM]->(param:Parameter)
							WHERE id(role) = {role_id}
							RETURN id(param) as id,
											param.name as name,
											param.value as value`

	db, err := conn.OpenPool()
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return nil, role.NewFindComponentRolesInternalServerError()
	}
	defer db.Close()

	data, _, _, err := db.QueryNeoAll(cypher, map[string]interface{}{
		"role_id": roleID})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
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
			ID:    row[0].(int64),
			Name:  &_name,
			Value: &_value}
	}
	return res, nil
}

func _getComponentRoleByName(conn neo4j.ConnPool, customer *string, cellID int64, componentID int64, roleName *string) *models.Role {

	var role *models.Role
	role = nil

	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->
							(cell:Cell)-[:PROVIDES]->(component:Component)-[:USE]->(role:Role {name: {role_name}})
							WHERE id(cell) = {cell_id} AND id(component) = {component_id}
								RETURN ID(role) as id,
												role.name as name,
												role.url as url,
												role.version as version,
												role.order as order`

	db, err := conn.OpenPool()
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return role
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return role
	}
	defer stmt.Close()

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"customer_name": swag.StringValue(customer),
		"cell_id":       cellID,
		"component_id":  componentID,
		"role_name":     swag.StringValue(roleName)})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
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

	role = &models.Role{ID: output[0].(int64),
		Name:    &_name,
		URL:     &_url,
		Version: &_version,
		Order:   &_order}

	return role
}
