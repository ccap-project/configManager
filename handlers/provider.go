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

	cypher := `MATCH (c:Customer {name: {name} })-[:OWN]->(cell:Cell {id: {cell_id}}),
										(providertype:ProviderType {name: {providertype}})
							CREATE (cell)-[:USE]->(provider:Provider {
								id: {provider_id},
								name: {provider_name},
							 	domain_name: {domain_name},
								tenantname: {tenant_name},
								access_key: {access_key},
								auth_url: {auth_url},
								username: {username},
								password: {password},
								region: {region},
								secret_key: {secret_key}})-[:PROVIDER_IS]->(providertype)
							RETURN	provider.id AS id,
											provider.name AS name`

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(principal.Name),
		"provider_name": params.Body.Name,
		"cell_id":       params.CellID})

	if GetProviderTypeByName(ctx.rt, params.Body.Type) == nil {
		ctxLogger.Error("provider type does not exists !")
		return provider.NewAddProviderBadRequest().WithPayload(&models.APIResponse{Message: "provider type does not exists"})
	}

	Provider := getProvider(ctx.rt, principal.Name, &params.CellID)

	if Provider != nil {
		ctxLogger.Warn("provider already exists !")
		return provider.NewAddProviderConflict().WithPayload(&models.APIResponse{Message: "provider already exists"})
	}

	db, err := ctx.rt.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j: ", err)
		return provider.NewAddProviderInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: ", err)
		return provider.NewAddProviderInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer stmt.Close()

	ulid := configManager.GetULID()

	ctxLogger = ctxLogger.WithFields(logrus.Fields{
		"cell_id": ulid})

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"name":          swag.StringValue(principal.Name),
		"cell_id":       params.CellID,
		"provider_id":   ulid,
		"provider_name": params.Body.Name,
		"domain_name":   params.Body.DomainName,
		"tenant_name":   params.Body.TenantName,
		"auth_url":      params.Body.AuthURL,
		"access_key":    params.Body.AccessKey,
		"username":      params.Body.Username,
		"password":      params.Body.Password,
		"providertype":  params.Body.Type,
		"region":        params.Body.Region,
		"secret_key":    params.Body.SecretKey})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: ", err)
		return provider.NewAddProviderInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	ctxLogger.Infoln(rows)

	output, _, err := rows.NextNeo()
	ctxLogger.Infoln(output)

	if err != nil {
		ctxLogger.Error("> An error occurred getting next row: ", err)
		return provider.NewAddProviderInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	return provider.NewAddProviderCreated().WithPayload(models.ULID(output[0].(string)))
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

	Provider := getProvider(ctx.rt, principal.Name, &params.CellID)

	if Provider == nil {
		ctxLogger.Warn("provider does not exists !")
		return provider.NewGetProviderNotFound()
	}

	return provider.NewGetProviderOK().WithPayload(Provider)
}

func getProvider(rt *configManager.Runtime, customerName *string, CellID *string) *models.Provider {

	var provider *models.Provider
	provider = nil

	cypher := `MATCH (c:Customer {name: {name} })-[:OWN]->(cell:Cell {id: {cell_id}})-[:USE]->(provider:Provider)
							MATCH (provider)-[:PROVIDER_IS]->(provider_type:ProviderType)
								RETURN provider.id as id,
												provider.name as name,
												provider.domain_name as domain_name,
												provider.tenantname as tenantname,
												provider.auth_url as auth_url,
												provider.providertype_id as providertype_id,
												provider.username as username,
												provider.password as password,
												provider.access_key as access_key,
												provider.region as region,
												provider.secret_key as secret_key,
												provider_type.name as provider_type_name`

	ctxLogger := rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(customerName),
		"cell_id":       swag.StringValue(CellID)})

	db, err := rt.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j: ", err)
		return provider
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: ", err)
		return provider
	}
	defer stmt.Close()

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"name":    swag.StringValue(customerName),
		"cell_id": swag.StringValue(CellID)})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: ", err)
		return provider
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		return provider
	}

	ctxLogger.Infoln(output)

	provider = new(models.Provider)

	provider.ID = models.ULID(output[0].(string))
	provider.Name = output[1].(string)
	provider.DomainName = output[2].(string)
	provider.TenantName = output[3].(string)
	provider.AuthURL = output[4].(string)
	provider.Username = output[6].(string)
	provider.Password = output[7].(string)
	provider.AccessKey = output[8].(string)
	provider.Region = output[9].(string)
	provider.SecretKey = output[10].(string)
	provider.Type = output[11].(string)

	return provider
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

	_params := util.BuildQuery(&params.Body, "provider", "ID")

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(principal.Name),
		"provider_name": params.Body.Name,
		"cell_id":       params.CellID})

	Provider := getProvider(ctx.rt, principal.Name, &params.CellID)

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
	ctxLogger.Infoln(util.BuildParams(params.Body,
		map[string]interface{}{
			"customer_name": swag.StringValue(principal.Name),
			"cell_id":       params.CellID}, "ID"))

	rows, err := stmt.QueryNeo(util.BuildParams(params.Body,
		map[string]interface{}{
			"customer_name": swag.StringValue(principal.Name),
			"cell_id":       params.CellID}, "ID"))

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
