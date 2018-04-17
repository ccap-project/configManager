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
	"configManager/restapi/operations/listener"

	"github.com/Sirupsen/logrus"
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

	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->
							(cell:Cell{id: {cell_id}})-[:PROVIDES]->
							(component:Component {id: {component_id}})
						CREATE (component)-[:LISTEN_ON]->(listener:Listener {
							id: {listener_id},
							name: {listener_name},
							port: {listener_port},
							protocol: {listener_protocol}} )
							RETURN listener.id as id`

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID,
		"component_id":  params.ComponentID})

	db, err := ctx.rt.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j: ", err)
		return listener.NewAddComponentListenerInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	if err != nil {
		ctxLogger.Error("An error occurred beginning transaction: ", err)
		return listener.NewAddComponentListenerInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: ", err)
		return listener.NewAddComponentListenerInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer stmt.Close()
	ulid := configManager.GetULID()

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"customer_name":     swag.StringValue(principal.Name),
		"cell_id":           params.CellID,
		"component_id":      params.ComponentID,
		"listener_id":       ulid,
		"listener_name":     swag.StringValue(params.Body.Name),
		"listener_port":     swag.Int64Value(params.Body.Port),
		"listener_protocol": swag.StringValue(params.Body.Protocol)})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: ", err)
		return listener.NewAddComponentListenerInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		ctxLogger.Error("An error occurred getting next row: ", err)
		return listener.NewAddComponentListenerInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	return listener.NewAddComponentListenerCreated().WithPayload(models.ULID(output[0].(string)))
}

func NewDeleteComponentListener(rt *configManager.Runtime) listener.DeleteComponentListenerHandler {
	return &deleteComponentListener{rt: rt}
}

type deleteComponentListener struct {
	rt *configManager.Runtime
}

func (ctx *deleteComponentListener) Handle(params listener.DeleteComponentListenerParams, principal *models.Customer) middleware.Responder {

	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->
							(cell:Cell{id: {cell_id}})-[:PROVIDES]->
							(component:Component {id: {component_id}})-[:LISTEN_ON]->
							(listener:Listener {id: {listernet_id}})
						DETACH DELETE listener`

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID,
		"component_id":  params.ComponentID,
		"listener_id":   params.ListenerID})

	if _getComponentListenerByID(ctx.rt, principal.Name, &params.CellID, &params.ListenerID) == nil {
		ctxLogger.Error("listener does not exists !")
		return listener.NewDeleteComponentListenerNotFound()
	}

	db, err := ctx.rt.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j: ", err)
		return listener.NewDeleteComponentListenerInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: ", err)
		return listener.NewDeleteComponentListenerInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer stmt.Close()

	_, err = stmt.ExecNeo(map[string]interface{}{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID,
		"component_id":  params.ComponentID,
		"listener_id":   params.ListenerID})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: ", err)
		return listener.NewAddComponentListenerInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
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

	Listener := _getComponentListenerByID(ctx.rt, principal.Name, &params.CellID, &params.ListenerID)
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

	data, err := _findComponentListeners(ctx.rt, principal.Name, &params.CellID, &params.ComponentID)

	log.Printf("= data(%#v)", data)

	if err != nil {
		return err
	}

	return listener.NewFindComponentListenersOK().WithPayload(data)
}

func NewUpdateComponentListener(rt *configManager.Runtime) listener.UpdateComponentListenerHandler {
	return &updateComponentListener{rt: rt}
}

type updateComponentListener struct {
	rt *configManager.Runtime
}

func (ctx *updateComponentListener) Handle(params listener.UpdateComponentListenerParams, principal *models.Customer) middleware.Responder {

	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->
							(cell:Cell {id: {cell_id}})-[:PROVIDES]->
							(component:Component {id: {component_id}})-[:LISTEN_ON]->
							(listener:Listener {id: {listener_id}})
						SET listener.name={listener_name},
								listener.port={listener_port},
								listener.protocol={listener_protocol}`

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID,
		"listener_id":   params.ListenerID,
		"component_id":  params.ComponentID})

	db, err := ctx.rt.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j: ", err)
		return listener.NewAddComponentListenerInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: ", err)
		return listener.NewAddComponentListenerInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
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
		ctxLogger.Error("-> An error occurred querying Neo: ", err)
		return listener.NewAddComponentListenerInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	return listener.NewUpdateComponentListenerOK()
}

func _findComponentListeners(rt *configManager.Runtime, customerName *string, CellID *string, ComponentID *string) ([]*models.Listener, middleware.Responder) {

	var res []*models.Listener

	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->
							(cell:Cell {id: {cell_id}})-[:PROVIDES]->
							(component:Component {id: {component_id}})-[:LISTEN_ON]->(listener:Listener)
								RETURN listener.id as id,
												listener.name as name,
												listener.port as port,
												listener.protocol as protocol`

	ctxLogger := rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(customerName),
		"cell_id":       swag.StringValue(CellID),
		"component_id":  swag.StringValue(ComponentID)})

	db, err := rt.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j: ", err)
		return nil, listener.NewFindComponentListenersInternalServerError()
	}
	defer db.Close()

	data, _, _, err := db.QueryNeoAll(cypher, map[string]interface{}{
		"customer_name": swag.StringValue(customerName),
		"cell_id":       swag.StringValue(CellID),
		"component_id":  swag.StringValue(ComponentID)})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: ", err)
		return nil, listener.NewFindComponentListenersInternalServerError()
	}

	for _, row := range data {

		var _port int64

		_name := row[1].(string)
		_port = row[2].(int64)
		_protocol := row[3].(string)

		l := &models.Listener{
			ID:       models.ULID(row[0].(string)),
			Name:     &_name,
			Port:     &_port,
			Protocol: &_protocol}

		res = append(res, l)
	}

	return res, nil
}

func _getComponentListenerByID(rt *configManager.Runtime, customer *string, cellID *string, listenerID *string) *models.Listener {

	var listener *models.Listener
	listener = nil

	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->
							(cell:Cell {id: {cell_id}})-[:PROVIDES]->
							(component:Component)-[:LISTEN_ON]->
							(listener:Listener {id: {listener_id}})
						RETURN listener.id as id,
										listener.name as name,
										listener.port as port,
										listener.protocol as protocol`

	ctxLogger := rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(customer),
		"cell_id":       swag.StringValue(cellID),
		"listener_id":   swag.StringValue(listenerID)})

	db, err := rt.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j: ", err)
		return listener
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: ", err)
		return listener
	}
	defer stmt.Close()

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"customer_name": swag.StringValue(customer),
		"cell_id":       swag.StringValue(cellID),
		"listener_id":   swag.StringValue(listenerID)})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: ", err)
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
		ID:       models.ULID(output[0].(string)),
		Port:     &_port,
		Name:     &_name,
		Protocol: &_protocol}

	return listener
}

func getComponentListenerConnections(rt *configManager.Runtime, customer *string, cellID *string, listenerID *models.ULID) *[]string {

	var res []string

	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->
							(cell:Cell {id: {cell_id}})-[:PROVIDES]->
							(component:Component)-[:LISTEN_ON]->
							(listener:Listener {id: {listener_id}})<-[:CONNECT_TO]-
							(src_conn)
						RETURN src_conn.name as name`

	ctxLogger := rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(customer),
		"cell_id":       swag.StringValue(cellID),
		"listener_id":   string(*listenerID)})

	db, err := rt.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j: ", err)
		return &res
	}
	defer db.Close()

	data, _, _, err := db.QueryNeoAll(cypher, map[string]interface{}{
		"customer_name": swag.StringValue(customer),
		"cell_id":       swag.StringValue(cellID),
		"listener_id":   string(*listenerID)})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: ", err)
		return &res
	}

	for _, row := range data {

		_name := row[0].(string)

		res = append(res, _name)
	}
	return &res

}
