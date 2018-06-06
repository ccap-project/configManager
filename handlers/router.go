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
	"fmt"

	"configManager"
	"configManager/models"
	"configManager/restapi/operations/router"
	"configManager/util"

	"github.com/Sirupsen/logrus"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
)

func NewAddCellRouter(rt *configManager.Runtime) router.AddRouterHandler {
	return &addCellRouter{rt: rt}
}

type addCellRouter struct {
	rt *configManager.Runtime
}

func (ctx *addCellRouter) Handle(params router.AddRouterParams, principal *models.Customer) middleware.Responder {

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID,
		"router_name":   swag.StringValue(params.Body.Name)})

	// XXX: Consistency check should check cidr also
	_router, err := _getRouterByName(ctx.rt, principal.Name, &params.CellID, params.Body.Name)
	if _router != nil {
		ctxLogger.Warn("router already exists !")
		return router.NewAddRouterConflict().WithPayload(&models.APIResponse{Message: "router already exists"})

	} else if err != nil {
		ctxLogger.Error(">Failure getting router: ", err)
		return router.NewAddRouterInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	ulid := configManager.GetULID()
	ctxLogger = ctxLogger.WithFields(logrus.Fields{"router_id": ulid})

	cypher := `MATCH (c:Customer {name: {customer_name} })-[:OWN]->(cell:Cell {id: {cell_id}})
							CREATE (cell)-[:HAS]->(router:Router {
								id: {router_id},
								%s})
							RETURN router.id AS id`

	_Query := fmt.Sprintf(cypher, util.BuildQuery(&params.Body, "router", "merge", []string{"ID"}))
	_Params := util.BuildParams(params.Body, "router",
		map[string]interface{}{
			"customer_name": swag.StringValue(principal.Name),
			"cell_id":       params.CellID,
			"router_id":     ulid},
		[]string{"ID"})

	output, err := ctx.rt.QueryDB(&_Query, &_Params)

	if err != nil {
		ctxLogger.Error(err)
		return router.NewAddRouterInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	if len(output) <= 0 {
		ctxLogger.Error("router not added")
		return router.NewAddRouterInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	ctxLogger.Info("OK")

	return router.NewAddRouterCreated().WithPayload(models.ULID(output[0].(string)))
}

func NewDeleteCellRouter(rt *configManager.Runtime) router.DeleteCellRouterHandler {
	return &deleteCellRouter{rt: rt}
}

type deleteCellRouter struct {
	rt *configManager.Runtime
}

func (ctx *deleteCellRouter) Handle(params router.DeleteCellRouterParams, principal *models.Customer) middleware.Responder {

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID,
		"router_id":     params.RouterID})

	_router, err := _getCellRouter(ctx.rt, principal.Name, &params.CellID, &params.RouterID)
	if err != nil {
		ctxLogger.Error("Failure getting router: ", err)
		return router.NewDeleteCellRouterInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})

	} else if _router == nil {
		ctxLogger.Error("router does not exists !")
		return router.NewDeleteCellRouterNotFound()
	}

	query := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->
							(cell:Cell{id: {cell_id}})-[:HAS]->
							(router:Router {id: {router_id}})
						DETACH DELETE router`

	_params := map[string]interface{}{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID,
		"router_id":     params.RouterID}

	_, err = ctx.rt.ExecDB(&query, &_params)

	if err != nil {
		ctxLogger.Error("Deleting router: ", err)
		return router.NewDeleteCellRouterInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	return router.NewDeleteCellRouterOK()
}

func NewGetCellRouter(rt *configManager.Runtime) router.GetCellRouterHandler {
	return &getCellRouter{rt: rt}
}

type getCellRouter struct {
	rt *configManager.Runtime
}

func (ctx *getCellRouter) Handle(params router.GetCellRouterParams, principal *models.Customer) middleware.Responder {

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID,
		"router_id":     params.RouterID})

	cellRouter, err := _getCellRouter(ctx.rt, principal.Name, &params.CellID, &params.RouterID)

	if err != nil {
		ctxLogger.Error("getting router: %s", err)
		return router.NewGetCellRouterInternalServerError()
	}

	if cellRouter == nil {
		ctxLogger.Warn("router does not exists !")
		return router.NewGetCellRouterNotFound()
	}

	return router.NewGetCellRouterOK().WithPayload(cellRouter)
}

func NewFindCellRouters(rt *configManager.Runtime) router.FindCellRoutersHandler {
	return &findCellRouters{rt: rt}
}

type findCellRouters struct {
	rt *configManager.Runtime
}

func (ctx *findCellRouters) Handle(params router.FindCellRoutersParams, principal *models.Customer) middleware.Responder {

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID})

	cellRouters, err := _findCellRouters(ctx.rt, principal.Name, &params.CellID)

	if err != nil {
		ctxLogger.Error("finding routers: ", err)
		return router.NewFindCellRoutersInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	if cellRouters == nil {
		ctxLogger.Warn("there are no routers")
	}

	return router.NewFindCellRoutersOK().WithPayload(cellRouters)
}

func _findCellRouters(rt *configManager.Runtime, customerName *string, CellID *string) ([]*models.Router, error) {

	var routers []*models.Router

	query := `MATCH (c:Customer {name: {name} })-[:OWN]->(cell:Cell {id: {cell_id}})-[:HAS]->(router:Router)
								RETURN router.id as id`

	params := map[string]interface{}{
		"name":    swag.StringValue(customerName),
		"cell_id": swag.StringValue(CellID)}

	data, err := rt.QueryAllDB(&query, &params)

	if err != nil {
		return nil, err

	} else if len(data) == 0 {
		return nil, nil
	}

	for _, row := range data {
		router_id := row[0].(string)
		router, _ := _getCellRouter(rt, customerName, CellID, &router_id)

		routers = append(routers, router)
	}

	return routers, nil
}

func _getCellRouter(rt *configManager.Runtime, customerName *string, CellID *string, RouterID *string) (*models.Router, error) {
	var router *models.Router

	query := `MATCH (c:Customer {name: {name} })-[:OWN]->
										(cell:Cell {id: {cell_id}})-[:HAS]->
										(router:Router {id: {router_id}})
								RETURN router {.*}`

	params := map[string]interface{}{
		"name":      swag.StringValue(customerName),
		"cell_id":   swag.StringValue(CellID),
		"router_id": swag.StringValue(RouterID)}

	output, err := rt.QueryDB(&query, &params)

	if err != nil {
		return router, err
	}

	router = new(models.Router)
	util.FillStruct(router, output[0].(map[string]interface{}))

	return router, nil
}

func _getRouterByID(rt *configManager.Runtime, RouterID *string) (*models.Router, error) {

	var router *models.Router

	query := `MATCH (router:Router {id: {router_id}})
								RETURN router {.*}`

	params := map[string]interface{}{
		"router_id": swag.StringValue(RouterID)}

	output, err := rt.QueryDB(&query, &params)

	if err != nil {
		return router, err
	}

	if len(output) > 0 {
		router = new(models.Router)
		util.FillStruct(router, output[0].(map[string]interface{}))
	}

	return router, nil
}

func _getRouterByName(rt *configManager.Runtime, customerName *string, CellID *string, routerName *string) (*models.Router, error) {

	var router *models.Router

	query := `MATCH (c:Customer {name: {name} })-[:OWN]->
										(cell:Cell {id: {cell_id}})-[:HAS]->
										(router:Router {name: {router_name}})
								RETURN router {.*}`

	params := map[string]interface{}{
		"name":        swag.StringValue(customerName),
		"cell_id":     swag.StringValue(CellID),
		"router_name": swag.StringValue(routerName)}

	output, err := rt.QueryDB(&query, &params)

	if err != nil {
		return router, err
	}

	if len(output) > 0 {
		router = new(models.Router)
		util.FillStruct(router, output[0].(map[string]interface{}))
	}

	return router, nil
}
