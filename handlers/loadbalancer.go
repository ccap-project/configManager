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
	"configManager/restapi/operations/loadbalancer"

	"github.com/Sirupsen/logrus"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
)

func NewAddCellLoadbalancer(rt *configManager.Runtime) loadbalancer.AddLoadbalancerHandler {
	return &addCellLoadbalancer{rt: rt}
}

type addCellLoadbalancer struct {
	rt *configManager.Runtime
}

func (ctx *addCellLoadbalancer) Handle(params loadbalancer.AddLoadbalancerParams, principal *models.Customer) middleware.Responder {

	cypher := `MATCH (c:Customer {name: {name} })-[:OWN]->(cell:Cell {id: {cell_id}})
							CREATE (cell)-[:HAS]->(loadbalancer:Loadbalancer {
								id: {loadbalancer_id},
								name: {loadbalancer_name},
								port: {loadbalancer_port},
							 	protocol: {loadbalancer_protocol},
								algorithm: {loadbalancer_algorithm}})
							RETURN	loadbalancer.id AS id`

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name":     swag.StringValue(principal.Name),
		"cell_id":           params.CellID,
		"loadbalancer_name": swag.StringValue(params.Body.Name)})

	// XXX: Consistency check should have more than only name...
	if _getLoadbalancerByName(ctx.rt.DB(), principal.Name, &params.CellID, params.Body.Name) != nil {
		ctxLogger.Warn("loadbalancer already exists !")
		return loadbalancer.NewAddLoadbalancerConflict().WithPayload(&models.APIResponse{Message: "loadbalancer already exists"})
	}

	db, err := ctx.rt.DB().OpenPool()

	if err != nil {
		ctxLogger.Error("error connecting to neo4j: ", err)
		return loadbalancer.NewAddLoadbalancerInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: ", err)
		return loadbalancer.NewAddLoadbalancerInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	ulid := configManager.GetULID()

	ctxLogger = ctxLogger.WithFields(logrus.Fields{
		"loadbalancer_id": ulid})

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"name":                   swag.StringValue(principal.Name),
		"cell_id":                params.CellID,
		"loadbalancer_id":        ulid,
		"loadbalancer_name":      swag.StringValue(params.Body.Name),
		"loadbalancer_port":      swag.Int64Value(params.Body.Port),
		"loadbalancer_protocol":  swag.StringValue(params.Body.Protocol),
		"loadbalancer_algorithm": swag.StringValue(params.Body.Algorithm)})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: ", err)
		return loadbalancer.NewAddLoadbalancerInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		ctxLogger.Error("An error occurred getting next row: ", err)
		return loadbalancer.NewAddLoadbalancerInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	ctxLogger.Info("OK")

	return loadbalancer.NewAddLoadbalancerCreated().WithPayload(models.ULID(output[0].(string)))
}

func NewAddLoadbalancerRelationship(rt *configManager.Runtime) loadbalancer.AddLoadbalancerRelationshipHandler {
	return &addLoadbalancerRelationship{rt: rt}
}

type addLoadbalancerRelationship struct {
	rt *configManager.Runtime
}

func (ctx *addLoadbalancerRelationship) Handle(params loadbalancer.AddLoadbalancerRelationshipParams, principal *models.Customer) middleware.Responder {

	if _getComponentListenerByID(ctx.rt, principal.Name, &params.CellID, &params.ListenerID) == nil {
		return loadbalancer.NewAddLoadbalancerRelationshipInternalServerError().WithPayload(&models.APIResponse{Message: "listener not found"})
	}

	cellLoadbalancer, err := _getCellLoadbalancer(ctx.rt, principal.Name, &params.CellID, &params.LoadbalancerID)

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return loadbalancer.NewAddLoadbalancerRelationshipInternalServerError()
	}

	if cellLoadbalancer == nil {
		return loadbalancer.NewAddLoadbalancerRelationshipInternalServerError().WithPayload(&models.APIResponse{Message: "loadbalancer not found"})
	}

	cypher := `
		MATCH (customer:Customer {name: {customer_name}})-[:OWN]->
					(cell:Cell {id: {cell_id}})-[:HAS]->
					(lb:Loadbalancer {id: {loadbalancer_id}})
		MATCH (cell {id: {cell_id}})-[:PROVIDES]->
			(component:Component)-[:LISTEN_ON]->
			(listener:Listener {id: {listener_id}})
		MERGE (lb)-[:CONNECT_TO]->(listener)
		RETURN *`

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name":   swag.StringValue(principal.Name),
		"cell_id":         params.CellID,
		"loadbalancer_id": params.LoadbalancerID,
		"listener_id":     params.ListenerID})

	db, err := ctx.rt.DB().OpenPool()

	if err != nil {
		ctxLogger.Warn("error connecting to neo4j: ", err)
		return loadbalancer.NewAddLoadbalancerRelationshipInternalServerError().WithPayload(&models.APIResponse{Message: "failure creating relationship"})
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Warn("An error occurred preparing statement: ", err)
		return loadbalancer.NewAddLoadbalancerRelationshipInternalServerError().WithPayload(&models.APIResponse{Message: "failure creating relationship"})
	}

	defer stmt.Close()

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"customer_name":   swag.StringValue(principal.Name),
		"cell_id":         params.CellID,
		"loadbalancer_id": params.LoadbalancerID,
		"listener_id":     params.ListenerID})

	ctxLogger.Info("rows", rows)

	if err != nil {
		ctxLogger.Warn("An error occurred querying Neo: ", err)
		return loadbalancer.NewAddLoadbalancerRelationshipInternalServerError().WithPayload(&models.APIResponse{Message: "failure creating relationship"})
	}

	_, _, err = rows.NextNeo()
	if err != nil {
		ctxLogger.Warn("An error occurred getting next row: ", err)
		return loadbalancer.NewAddLoadbalancerRelationshipInternalServerError().WithPayload(&models.APIResponse{Message: "failure creating relationship"})
	}

	return loadbalancer.NewAddCellLoadbalancerRelationshipOK()
}

func NewDeleteLoadbalancerRelationship(rt *configManager.Runtime) loadbalancer.DeleteLoadbalancerRelationshipHandler {
	return &deleteLoadbalancerRelationship{rt: rt}
}

type deleteLoadbalancerRelationship struct {
	rt *configManager.Runtime
}

func (ctx *deleteLoadbalancerRelationship) Handle(params loadbalancer.DeleteLoadbalancerRelationshipParams, principal *models.Customer) middleware.Responder {

	if _getComponentListenerByID(ctx.rt, principal.Name, &params.CellID, &params.ListenerID) == nil {
		return loadbalancer.NewDeleteLoadbalancerRelationshipInternalServerError().WithPayload(&models.APIResponse{Message: "listener not found"})
	}

	cellLoadbalancer, err := _getCellLoadbalancer(ctx.rt, principal.Name, &params.CellID, &params.LoadbalancerID)

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return loadbalancer.NewDeleteLoadbalancerRelationshipInternalServerError()
	}

	if cellLoadbalancer == nil {
		return loadbalancer.NewDeleteLoadbalancerRelationshipInternalServerError().WithPayload(&models.APIResponse{Message: "loadbalancer not found"})
	}

	cypher := `
			MATCH (customer:Customer {name: {customer_name}})-[:OWN]->
				(cell:Cell {id: {cell_id}})-[:HAS]->
				(loadbalancer:Loadbalancer {id: loadbalancer_id})-[r:CONNECT_TO]->
				(listener:Listener {id: {listener_id}})
			delete r`

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name":   swag.StringValue(principal.Name),
		"cell_id":         params.CellID,
		"loadbalancer_id": params.LoadbalancerID,
		"listener_id":     params.ListenerID})

	db, err := ctx.rt.DB().OpenPool()

	if err != nil {
		ctxLogger.Warn("error connecting to neo4j: ", err)
		return loadbalancer.NewDeleteLoadbalancerRelationshipInternalServerError().WithPayload(&models.APIResponse{Message: "failure deleting relationship"})
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Warn("An error occurred preparing statement: ", err)
		return loadbalancer.NewDeleteLoadbalancerRelationshipInternalServerError().WithPayload(&models.APIResponse{Message: "failure deleting relationship"})
	}

	defer stmt.Close()

	_, err = stmt.ExecNeo(map[string]interface{}{
		"customer_name":   swag.StringValue(principal.Name),
		"cell_id":         params.CellID,
		"loadbalancer_id": params.LoadbalancerID,
		"listener_id":     params.ListenerID})

	if err != nil {
		ctxLogger.Warn("An error occurred querying Neo: ", err)
		return loadbalancer.NewDeleteLoadbalancerRelationshipInternalServerError().WithPayload(&models.APIResponse{Message: "failure deleting relationship"})
	}

	return loadbalancer.NewDeleteLoadbalancerRelationshipOK()
}

func NewGetCellLoadbalancer(rt *configManager.Runtime) loadbalancer.GetCellLoadbalancerHandler {
	return &getCellLoadbalancer{rt: rt}
}

type getCellLoadbalancer struct {
	rt *configManager.Runtime
}

func (ctx *getCellLoadbalancer) Handle(params loadbalancer.GetCellLoadbalancerParams, principal *models.Customer) middleware.Responder {

	cellLoadbalancer, err := _getCellLoadbalancer(ctx.rt, principal.Name, &params.CellID, &params.LoadbalancerID)

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return loadbalancer.NewGetCellLoadbalancerInternalServerError()
	}

	if cellLoadbalancer == nil {
		return loadbalancer.NewGetCellLoadbalancerOK()
	}

	return loadbalancer.NewGetCellLoadbalancerOK().WithPayload(cellLoadbalancer)
}

func NewFindCellLoadbalancers(rt *configManager.Runtime) loadbalancer.FindCellLoadbalancersHandler {
	return &findCellLoadbalancers{rt: rt}
}

type findCellLoadbalancers struct {
	rt *configManager.Runtime
}

func (ctx *findCellLoadbalancers) Handle(params loadbalancer.FindCellLoadbalancersParams, principal *models.Customer) middleware.Responder {

	cellLoadbalancers, err := _findCellLoadbalancers(ctx.rt, principal.Name, &params.CellID)

	if err != nil {
		return loadbalancer.NewFindCellLoadbalancersInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	return loadbalancer.NewFindCellLoadbalancersOK().WithPayload(cellLoadbalancers)
}

func _findCellLoadbalancers(rt *configManager.Runtime, customerName *string, CellID *string) ([]*models.Loadbalancer, error) {

	var res []*models.Loadbalancer

	cypher := `MATCH (c:Customer {name: {name} })-[:OWN]->(cell:Cell {id: {cell_id}})-[:HAS]->(loadbalancer)
								RETURN loadbalancer.id as id,
												loadbalancer.name as name`

	db, err := rt.DB().OpenPool()

	ctxLogger := rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(customerName),
		"cell_id":       swag.StringValue(CellID)})

	if err != nil {
		ctxLogger.Error("error connecting to neo4j:", err)
		return nil, err
	}
	defer db.Close()

	data, _, _, err := db.QueryNeoAll(cypher, map[string]interface{}{
		"name":    swag.StringValue(customerName),
		"cell_id": swag.StringValue(CellID)})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: %s", err)
		return nil, err

	} else if len(data) == 0 {
		return nil, nil
	}

	for _, row := range data {
		lb_id := row[0].(string)
		lb, _ := _getCellLoadbalancer(rt, customerName, CellID, &lb_id)

		res = append(res, lb)
	}

	return res, nil
}

func _getCellLoadbalancer(rt *configManager.Runtime, customerName *string, CellID *string, LoadbalancerID *string) (*models.Loadbalancer, error) {
	var loadbalancer *models.Loadbalancer
	loadbalancer = nil

	cypher := `MATCH (c:Customer {name: {name} })-[:OWN]->
										(cell:Cell {id: {cell_id}})-[:HAS]->
										(loadbalancer:Loadbalancer {id: {loadbalancer_id}})
								RETURN loadbalancer.id as id,
												loadbalancer.name as name,
												loadbalancer.port as port,
												loadbalancer.protocol as protocol,
												loadbalancer.algorithm as algorithm`

	db, err := rt.DB().OpenPool()
	ctxLogger := rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(customerName),
		"cell_id":       swag.StringValue(CellID)})

	if err != nil {
		ctxLogger.Error("error connecting to neo4j:", err)
		return loadbalancer, err
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: %s", err)
		return loadbalancer, err
	}

	defer stmt.Close()

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"name":            swag.StringValue(customerName),
		"cell_id":         swag.StringValue(CellID),
		"loadbalancer_id": swag.StringValue(LoadbalancerID)})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: %s", err)
		return loadbalancer, err
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		return loadbalancer, err
	}

	_name := output[1].(string)
	_port := output[2].(int64)
	_protocol := output[3].(string)
	_algorithm := output[4].(string)

	loadbalancer = &models.Loadbalancer{
		ID:        models.ULID(output[0].(string)),
		Name:      &_name,
		Port:      &_port,
		Protocol:  &_protocol,
		Algorithm: &_algorithm}

	return loadbalancer, nil
}

func _getLoadbalancerByName(conn neo4j.ConnPool, customerName *string, CellID *string, loadbalancerName *string) *models.Loadbalancer {

	var loadbalancer *models.Loadbalancer
	loadbalancer = nil

	cypher := `MATCH (c:Customer {name: {name} })-[:OWN]->
										(cell:Cell {id: {cell_id}})-[:PROVIDES]->
										(loadbalancer:Loadbalancer)
							WHERE loadbalancer.name = {loadbalancer_name}
								RETURN loadbalancer. as id,
												loadbalancer.name as name`

	db, err := conn.OpenPool()

	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return loadbalancer
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return loadbalancer
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{"name": swag.StringValue(customerName),
		"cell_id":           CellID,
		"loadbalancer_name": swag.StringValue(loadbalancerName)})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return loadbalancer
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		return loadbalancer
	}
	_name := output[1].(string)

	loadbalancer = &models.Loadbalancer{ID: models.ULID(output[0].(string)),
		Name: &_name}

	stmt.Close()

	return loadbalancer
}
