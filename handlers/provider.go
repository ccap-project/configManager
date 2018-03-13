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
	"configManager/restapi/operations/provider"

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
								auth_url: {auth_url},
								username: {username},
								password: {password}})-[:PROVIDER_IS]->(providertype)
							RETURN	provider.id AS id,
											provider.name AS name`

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(principal.Name),
		"provider_name": swag.StringValue(params.Body.Name),
		"cell_id":       params.CellID})

	Provider := getProvider(ctx.rt.DB(), principal.Name, &params.CellID)

	if Provider != nil {
		ctxLogger.Warn("provider already exists !")
		return provider.NewAddProviderConflict().WithPayload(models.APIResponse{Message: "provider already exists"})
	}

	db, err := ctx.rt.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j: ", err)
		return provider.NewAddProviderInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: ", err)
		return provider.NewAddProviderInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	ulid := configManager.GetULID()

	ctxLogger = ctxLogger.WithFields(logrus.Fields{
		"cell_id": ulid})

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"name":          swag.StringValue(principal.Name),
		"cell_id":       params.CellID,
		"provider_id":   ulid,
		"provider_name": swag.StringValue(params.Body.Name),
		"domain_name":   swag.StringValue(params.Body.DomainName),
		"tenant_name":   swag.StringValue(params.Body.TenantName),
		"auth_url":      swag.StringValue(params.Body.AuthURL),
		"username":      swag.StringValue(params.Body.Username),
		"password":      swag.StringValue(params.Body.Password),
		"providertype":  swag.StringValue(params.Body.Type)})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: ", err)
		return provider.NewAddProviderInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		ctxLogger.Error("An error occurred getting next row: ", err)
		return provider.NewAddProviderInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
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

	Provider := getProvider(ctx.rt.DB(), principal.Name, &params.CellID)

	if Provider == nil {
		ctxLogger.Warn("provider does not exists !")
		return provider.NewGetProviderNotFound()
	}

	return provider.NewGetProviderOK().WithPayload(Provider)
}

func getProvider(conn neo4j.ConnPool, customerName *string, CellID *string) *models.Provider {

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
												provider_type.name as provider_type_name`

	db, err := conn.OpenPool()
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return provider
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return provider
	}
	defer stmt.Close()

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"name":    swag.StringValue(customerName),
		"cell_id": CellID})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return provider
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		return provider
	}

	provider = new(models.Provider)
	provider.Name = new(string)
	provider.DomainName = new(string)
	provider.TenantName = new(string)
	provider.AuthURL = new(string)
	provider.Type = new(string)
	provider.Username = new(string)
	provider.Password = new(string)

	provider.ID = models.ULID(output[0].(string))
	*provider.Name = output[1].(string)
	*provider.DomainName = output[2].(string)
	*provider.TenantName = output[3].(string)
	*provider.AuthURL = output[4].(string)
	*provider.Username = output[6].(string)
	*provider.Password = output[7].(string)
	*provider.Type = output[8].(string)

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
							WHERE newProviderType.name = {providertype}
							SET provider.name={name},
									provider.domain_name={domain_name},
									provider.tenantname={tenant_name},
									provider.auth_url={auth_url},
									provider.username={username},
									provider.password={password},
									provider.providertype={providertype}
							DELETE rel, rel2
							CREATE (cell)-[:USE]->(provider)-[:PROVIDER_IS]->(newProviderType)
							return provider`

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(principal.Name),
		"provider_name": swag.StringValue(params.Body.Name),
		"cell_id":       params.CellID})

	Provider := getProvider(ctx.rt.DB(), principal.Name, &params.CellID)

	if Provider == nil {
		ctxLogger.Warn("provider does not exists !")
		return provider.NewUpdateProviderNotFound().WithPayload(models.APIResponse{Message: "provider does not exists"})
	}

	db, err := ctx.rt.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j: ", err)
		return provider.NewUpdateProviderInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: ", err)
		return provider.NewUpdateProviderInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID,
		"name":          swag.StringValue(params.Body.Name),
		"domain_name":   swag.StringValue(params.Body.DomainName),
		"tenant_name":   swag.StringValue(params.Body.TenantName),
		"auth_url":      swag.StringValue(params.Body.AuthURL),
		"username":      swag.StringValue(params.Body.Username),
		"password":      swag.StringValue(params.Body.Password),
		"type":          swag.StringValue(params.Body.Type)})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: ", err)
		return provider.NewUpdateProviderInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	_, _, err = rows.NextNeo()
	if err != nil {
		ctxLogger.Error("An error occurred getting next row: ", err)
		return provider.NewUpdateProviderInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	return provider.NewUpdateProviderOK()
}
