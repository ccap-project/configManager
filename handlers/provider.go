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
	"configManager"
	"configManager/models"
	"configManager/restapi/operations/provider"
	"configManager/util"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
)

func NewAddProvider(rt *configManager.Runtime) provider.AddProviderHandler {
	return &addCellProvider{rt: rt}
}

type addCellProvider struct {
	rt *configManager.Runtime
}

func (ctx *addCellProvider) Handle(params provider.AddProviderParams, principal *models.Customer) middleware.Responder {

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(principal.Name),
		"provider_name": params.Body.Name,
		"cell_id":       params.CellID})

	/*
	 * Validate provider type
	 */
	providerType, err := _getProviderTypeByName(ctx.rt, params.Body.Type)
	if err != nil {
		ctxLogger.Errorf("getting provider type, %s", err)
		return provider.NewAddProviderInternalServerError()
	}

	if providerType == nil {
		ctxLogger.Error("provider type does not exists !")
		return provider.NewAddProviderBadRequest().WithPayload(&models.APIResponse{Message: "provider type does not exists"})
	}

	/*
	 * Validate provider
	 */
	Provider, err := _getProvider(ctx.rt, principal.Name, &params.CellID)
	if err != nil {
		ctxLogger.Errorf("getting provider, %s", err)
		return provider.NewAddProviderInternalServerError()
	}

	if Provider != nil {
		ctxLogger.Warn("provider already exists !")
		return provider.NewAddProviderConflict().WithPayload(&models.APIResponse{Message: "provider already exists"})
	}

	ulid := configManager.GetULID()

	ctxLogger = ctxLogger.WithFields(logrus.Fields{
		"cell_id": ulid})

	cypher := `MATCH (c:Customer {name: {customer_name} })-[:OWN]->(cell:Cell {id: {cell_id}}),
										(providertype:ProviderType {name: {providertype}})
							CREATE (cell)-[:USE]->(provider:Provider {
								id: {id},
								%s})-[:PROVIDER_IS]->(providertype)
							RETURN	provider.id AS id,
											provider.name AS name`

	_Query := fmt.Sprintf(cypher, util.BuildQuery(&params.Body, "", "merge", []string{"ID"}))

	_Params := util.BuildParams(params.Body, "",
		map[string]interface{}{
			"customer_name": swag.StringValue(principal.Name),
			"cell_id":       params.CellID,
			"providertype":  params.Body.Type,
			"id":            ulid},
		[]string{"ID"})

	output, err := ctx.rt.QueryDB(&_Query, &_Params)

	if err != nil {
		ctxLogger.Error("adding provider, ", err)
		return provider.NewAddProviderInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	ctxLogger.Infoln(output)

	if len(output) < 1 {
		ctxLogger.Error("network not added")
		return provider.NewAddProviderInternalServerError()
	}

	ctxLogger.Info("OK")

	return provider.NewAddProviderCreated().WithPayload(models.ULID(models.ULID(output[0].(string))))
}

func NewGetProvider(rt *configManager.Runtime) provider.GetProviderHandler {
	return &getCellProvider{rt: rt}
}

type getCellProvider struct {
	rt *configManager.Runtime
}

func (ctx *getCellProvider) Handle(params provider.GetProviderParams, principal *models.Customer) middleware.Responder {

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID})

	/*
	 * Validate provider
	 */

	Provider, err := _getProvider(ctx.rt, principal.Name, &params.CellID)
	if err != nil {
		ctxLogger.Errorf("getting provider, %s", err)
		return provider.NewGetProviderInternalServerError()
	}

	if Provider == nil {
		ctxLogger.Warn("provider does not exists !")
		return provider.NewGetProviderNotFound()
	}

	return provider.NewGetProviderOK().WithPayload(Provider)
}

func NewUpdateProvider(rt *configManager.Runtime) provider.UpdateProviderHandler {
	return &updateCellProvider{rt: rt}
}

type updateCellProvider struct {
	rt *configManager.Runtime
}

func (ctx *updateCellProvider) Handle(params provider.UpdateProviderParams, principal *models.Customer) middleware.Responder {

	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->
							(cell:Cell {id: {cell_id}})-[rel:USE]->(provider:Provider)-[rel2:PROVIDER_IS]->(provider_type:ProviderType)
						MATCH (newProviderType:ProviderType)
							WHERE newProviderType.name = {type}
							SET %s
							DELETE rel, rel2
							CREATE (cell)-[:USE]->(provider)-[:PROVIDER_IS]->(newProviderType)
							return provider`

	_params := util.BuildQuery(&params.Body, "provider", "update", []string{"ID"})

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(principal.Name),
		"provider_name": params.Body.Name,
		"cell_id":       params.CellID})

	Provider, err := _getProvider(ctx.rt, principal.Name, &params.CellID)
	if err != nil {
		ctxLogger.Errorf("getting provider, %s", err)
		return provider.NewUpdateProviderInternalServerError()
	}

	if Provider == nil {
		ctxLogger.Warn("provider does not exists !")
		return provider.NewUpdateProviderNotFound().WithPayload(&models.APIResponse{Message: "provider does not exists"})
	}

	db, err := ctx.rt.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j: ", err)
		return provider.NewUpdateProviderInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(fmt.Sprintf(cypher, _params))
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: ", err)
		return provider.NewUpdateProviderInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer stmt.Close()
	ctxLogger.Infoln(util.BuildParams(params.Body, "provider",
		map[string]interface{}{
			"customer_name": swag.StringValue(principal.Name),
			"cell_id":       params.CellID}, []string{"ID"}))

	rows, err := stmt.QueryNeo(util.BuildParams(params.Body, "provider",
		map[string]interface{}{
			"customer_name": swag.StringValue(principal.Name),
			"cell_id":       params.CellID}, []string{"ID"}))

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: ", err)
		return provider.NewUpdateProviderInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	_, _, err = rows.NextNeo()
	if err != nil {
		ctxLogger.Error("An error occurred getting next row: ", err)
		return provider.NewUpdateProviderInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	return provider.NewUpdateProviderOK()
}

func _getProvider(rt *configManager.Runtime, customerName *string, CellID *string) (*models.Provider, error) {

	var provider *models.Provider

	query := `MATCH (c:Customer {name: {name} })-[:OWN]->(cell:Cell {id: {cell_id}})-[:USE]->(provider:Provider)
							MATCH (provider)-[:PROVIDER_IS]->(provider_type:ProviderType)
								RETURN provider {.*}`

	params := map[string]interface{}{
		"name":    swag.StringValue(customerName),
		"cell_id": swag.StringValue(CellID)}

	output, err := rt.QueryDB(&query, &params)

	if err != nil {
		return provider, err
	}

	if len(output) > 0 {
		provider = new(models.Provider)
		util.FillStruct(provider, output[0].(map[string]interface{}))
	}

	return provider, nil
}
