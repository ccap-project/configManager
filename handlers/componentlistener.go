/*
 *
 * Copyright (c) 2018 Alexandre Biancalana <ale@biancalanas.net>.
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
	"configManager/restapi/operations/listener"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
)

func NewAddComponentListener(rt *configManager.Runtime) listener.AddComponentListenerHandler {
	return &addComponentListener{rt: rt}
}

type addComponentListener struct {
	rt *configManager.Runtime
}

func (ctx *addComponentListener) Handle(params listener.AddComponentListenerParams, principal *models.Customer) middleware.Responder {

	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->(cell:Cell)-[:PROVIDES]->(component:Component)
							WHERE id(cell) = {cell_id} AND id(component) = {component_id}
							CREATE (component)-[:LISTEN_ON]->(listener:Listener {
								name: {listener_name},
								port: {listener_port},
								protocol: {listener_protocol}} )
								RETURN id(listener) as id`

	db, err := ctx.rt.DB().OpenPool()
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return listener.NewAddComponentListenerInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	if err != nil {
		log.Printf("An error occurred beginning transaction: %s", err)
		return listener.NewAddComponentListenerInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return listener.NewAddComponentListenerInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"customer_name":     swag.StringValue(principal.Name),
		"cell_id":           params.CellID,
		"component_id":      params.ComponentID,
		"listener_name":     swag.StringValue(params.Body.Name),
		"listener_port":     swag.Int64Value(params.Body.Port),
		"listener_protocol": swag.StringValue(params.Body.Protocol)})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return listener.NewAddComponentListenerInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		log.Printf("An error occurred getting next row: %s", err)
		return listener.NewAddComponentListenerInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	log.Printf("= Output(%#v)", output)

	log.Printf("name(%s)", swag.StringValue(params.Body.Name))

	stmt.Close()

	return listener.NewAddComponentListenerCreated().WithPayload(output[0].(int64))
}

func NewDeleteComponentListener(rt *configManager.Runtime) listener.DeleteComponentListenerHandler {
	return &deleteComponentListener{rt: rt}
}

type deleteComponentListener struct {
	rt *configManager.Runtime
}

func (ctx *deleteComponentListener) Handle(params listener.DeleteComponentListenerParams, principal *models.Customer) middleware.Responder {

	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->
							(cell:Cell)-[:PROVIDES]->(component:Component)-[:LISTEN_ON]->(listener:Listener)
							WHERE id(cell) = {cell_id}
								AND id(component) = {component_id}
								AND id(listener) = {listener_id}
							DETACH DELETE listener`

	if _getComponentListenerByID(ctx.rt.DB(), principal.Name, params.CellID, params.ListenerID) == nil {
		log.Println("listener does not exists !")
		return listener.NewDeleteComponentListenerNotFound()
	}

	db, err := ctx.rt.DB().OpenPool()
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return listener.NewDeleteComponentListenerInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return listener.NewDeleteComponentListenerInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	defer stmt.Close()

	_, err = stmt.ExecNeo(map[string]interface{}{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID,
		"component_id":  params.ComponentID,
		"listener_id":   params.ListenerID})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return listener.NewAddComponentListenerInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	//if res.Count.(int) <= 0 {
	//	return listener.NewDeleteComponentListenerNotFound()
	//}

	return listener.NewDeleteComponentListenerOK()
}

func NewGetComponentListenerByID(rt *configManager.Runtime) listener.GetComponentListenerByIDHandler {
	return &getComponentListenerByID{rt: rt}
}

type getComponentListenerByID struct {
	rt *configManager.Runtime
}

func (ctx *getComponentListenerByID) Handle(params listener.GetComponentListenerByIDParams, principal *models.Customer) middleware.Responder {

	Listener := _getComponentListenerByID(ctx.rt.DB(), principal.Name, params.CellID, params.ListenerID)
	if Listener == nil {
		return listener.NewGetComponentListenerByIDNotFound()
	}

	return listener.NewGetComponentListenerByIDOK().WithPayload(Listener)
}

func NewFindComponentListeners(rt *configManager.Runtime) listener.FindComponentListenersHandler {
	return &findComponentListeners{rt: rt}
}

type findComponentListeners struct {
	rt *configManager.Runtime
}

func (ctx *findComponentListeners) Handle(params listener.FindComponentListenersParams, principal *models.Customer) middleware.Responder {

	data, err := _FindComponentListeners(ctx.rt.DB(), principal.Name, params.CellID, params.ComponentID)

	log.Printf("= data(%#v)", data)

	if err != nil {
		return err
	}

	return listener.NewFindComponentListenersOK().WithPayload(data)
}

func _FindComponentListeners(conn neo4j.ConnPool, customerName *string, CellID int64, ComponentID int64) ([]*models.Listener, middleware.Responder) {

	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->
							(cell:Cell)-[:PROVIDES]->(component:Component)-[:LISTEN_ON]->(listener:Listener)
							WHERE id(cell) = {cell_id} AND id(component) = {component_id}
								RETURN id(listener) as id,
												listener.name as name,
												listener.port as port,
												listener.protocol as protocol`

	db, err := conn.OpenPool()
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return nil, listener.NewFindComponentListenersInternalServerError()
	}
	defer db.Close()

	log.Println(" customerName =>>>>>", *customerName)

	data, _, _, err := db.QueryNeoAll(cypher, map[string]interface{}{
		"customer_name": *customerName,
		"cell_id":       CellID,
		"component_id":  ComponentID})

	log.Printf("= data(%#v)", data)

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return nil, listener.NewFindComponentListenersInternalServerError()
	}

	res := make([]*models.Listener, len(data))

	for idx, row := range data {

		var _port int64

		_name := row[1].(string)
		_port = row[2].(int64)
		_protocol := row[3].(string)

		res[idx] = &models.Listener{
			ID:       row[0].(int64),
			Name:     &_name,
			Port:     &_port,
			Protocol: &_protocol}
	}
	return res, nil
}

func NewUpdateComponentListener(rt *configManager.Runtime) listener.UpdateComponentListenerHandler {
	return &updateComponentListener{rt: rt}
}

type updateComponentListener struct {
	rt *configManager.Runtime
}

func (ctx *updateComponentListener) Handle(params listener.UpdateComponentListenerParams, principal *models.Customer) middleware.Responder {

	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->
								(cell:Cell)-[:PROVIDES]->(component:Component)-[:LISTEN_ON]->(listener:Listener)
							WHERE id(cell) = {cell_id} AND id(component) = {component_id} AND id(listener) = {listener_id}
							SET listener.name={listener_name},
									listener.port={listener_port},
									listener.protocol={listener_protocol}`

	db, err := ctx.rt.DB().OpenPool()
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return listener.NewAddComponentListenerInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return listener.NewAddComponentListenerInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}
	defer stmt.Close()

	/*log.Printf("customer_name(%s) cell_id(%s) component_id(%s) listener_id(%s) listener_name(%s) listener_image(%s) listener_flavor(%s)",
	swag.StringValue(principal.Name),
	params.CellID,
	params.ComponentID,
	params.ListenerID,
	swag.StringValue(params.Body.Name),
	swag.StringValue(params.Body.Image),
	swag.StringValue(params.Body.Flavor))*/

	_, err = stmt.ExecNeo(map[string]interface{}{
		"customer_name":     swag.StringValue(principal.Name),
		"cell_id":           params.CellID,
		"component_id":      params.ComponentID,
		"listener_id":       params.ListenerID,
		"listener_port":     swag.Int64Value(params.Body.Port),
		"listener_protocol": swag.StringValue(params.Body.Protocol)})

	if err != nil {
		log.Printf("-> An error occurred querying Neo: %s", err)
		return listener.NewAddComponentListenerInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	return listener.NewUpdateComponentListenerOK()
}

func _getComponentListenerByID(conn neo4j.ConnPool, customer *string, cellID int64, listenerID int64) *models.Listener {

	var listener *models.Listener
	listener = nil

	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->
							(cell:Cell)-[:PROVIDES]->(component:Component)-[:LISTEN_ON]->(listener:Listener)
							WHERE id(cell) = {cell_id}
								AND id(listener) = {listener_id}
							RETURN id(listener) as id,
											listener.name as name,
											listener.port as image,
											listener.protocol as flavor`

	db, err := conn.OpenPool()
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return listener
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return listener
	}
	defer stmt.Close()

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"customer_name": swag.StringValue(customer),
		"cell_id":       cellID,
		"listener_id":   listenerID})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return listener
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		return listener
	}

	_name := output[1].(string)
	_port := output[2].(int64)
	_protocol := output[3].(string)

	listener = &models.Listener{
		ID:       output[0].(int64),
		Port:     &_port,
		Name:     &_name,
		Protocol: &_protocol}

	log.Printf("here => (%#v)", listener)

	return listener
}
