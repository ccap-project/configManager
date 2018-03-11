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
	"configManager/restapi/operations/keypair"

	"github.com/Sirupsen/logrus"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
)

func NewAddCellKeypair(rt *configManager.Runtime) keypair.AddCellKeypairHandler {
	return &addCellKeypair{rt: rt}
}

type addCellKeypair struct {
	rt *configManager.Runtime
}

func (ctx *addCellKeypair) Handle(params keypair.AddCellKeypairParams, principal *models.Customer) middleware.Responder {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	cypher := `MATCH (c:Customer {name: {customer_name}})-[:OWN]->(cell:Cell {id: {cell_id}}),
							(c:Customer {name: {customer_name}})-[:HAS]->(keypair:Keypair {name: {keypair_name}})
							CREATE (cell)-[:DEPLOY_WITH]->(keypair)
							RETURN	keypair.id AS id`

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(principal.Name)})

	if getKeypairByName(ctx.rt, principal.Name, &params.KeypairName) == nil {
		ctxLogger.Error("keypair does not exists !")
		return keypair.NewAddCellKeypairConflict()
	}

	ctxLogger = ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(principal.Name),
		"keypair_name":  params.KeypairName})

	if getCellKeypair(ctx.rt, principal.Name, &params.CellID) != nil {
		ctxLogger.Error("This Cell already has a keypair")
		return keypair.NewAddCellKeypairNotFound()
	}

	ctxLogger = ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID,
		"keypair_name":  params.KeypairName})

	db, err := ctx.rt.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j:", err)
		return keypair.NewAddCellKeypairInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: %s", err)
		return keypair.NewAddCellKeypairInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"customer_name": swag.StringValue(principal.Name),
		"keypair_name":  params.KeypairName,
		"cell_id":       params.CellID})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: %s", err)
		return keypair.NewAddCellKeypairInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	output, _, err := rows.NextNeo()

	if err != nil {
		ctxLogger.Error("An error occurred getting next row: %s", err)
		return keypair.NewAddCellKeypairInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	ctxLogger.Info("OK")
	return keypair.NewAddCellKeypairCreated().WithPayload(output[0].(int64))
}

func NewAddKeypair(rt *configManager.Runtime) keypair.AddKeypairHandler {
	return &addKeypair{rt: rt}
}

type addKeypair struct {
	rt *configManager.Runtime
}

func (ctx *addKeypair) Handle(params keypair.AddKeypairParams, principal *models.Customer) middleware.Responder {

	cypher := `MATCH (c:Customer {name: {name} })
							CREATE (c)-[:HAS]->(k:Keypair { id: {id}, name: {kname}, public_key: {public_key} })
							RETURN	k.id AS id,
											k.name AS name,
											k.public_key AS public_key`

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(principal.Name)})

	if getKeypairByName(ctx.rt, principal.Name, params.Body.Name) != nil {
		ctxLogger.Error("keypair already exists !")
		return keypair.NewAddKeypairConflict().WithPayload(models.APIResponse{Message: "keypair name already exists"})
	}

	ctxLogger = ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(principal.Name),
		"keypair_name":  params.Body.Name})

	db, err := ctx.rt.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j:", err)
		return keypair.NewAddKeypairInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: %s", err)
		return keypair.NewAddKeypairInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	ulid := configManager.GetULID()

	ctxLogger = ctxLogger.WithFields(logrus.Fields{
		"keypair_id": ulid})

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"id":         ulid,
		"name":       swag.StringValue(principal.Name),
		"kname":      swag.StringValue(params.Body.Name),
		"public_key": swag.StringValue(params.Body.PublicKey)})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: %s", err)
		return keypair.NewAddKeypairInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		ctxLogger.Error("An error occurred getting next row: %s", err)
		return keypair.NewAddKeypairInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	ctxLogger.Info("OK")

	return keypair.NewAddKeypairCreated().WithPayload(models.ULID(output[0].(string)))
}

func NewGetKeypairByID(rt *configManager.Runtime) keypair.GetKeypairByIDHandler {
	return &getKeypairByID{rt: rt}
}

type getKeypairByID struct {
	rt *configManager.Runtime
}

func (ctx *getKeypairByID) Handle(params keypair.GetKeypairByIDParams, principal *models.Customer) middleware.Responder {

	cypher := `MATCH (c:Customer {name: {name} })-[:HAS]->(k:Keypair {id: {kid}})
								RETURN k.id as id,
												k.name as name,
												k.public_key as public_key`

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(principal.Name),
		"keypair_id":    params.KeypairID})

	db, err := ctx.rt.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j:", err)
		return keypair.NewGetKeypairByIDInternalServerError()
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: %s", err)
		return keypair.NewGetKeypairByIDInternalServerError()
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"name": swag.StringValue(principal.Name),
		"kid":  params.KeypairID})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: %s", err)
		return keypair.NewGetKeypairByIDInternalServerError()
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		ctxLogger.Error("not found")
		return keypair.NewGetKeypairByIDNotFound()
	}
	_name := output[1].(string)
	_pubkey := output[2].(string)

	_keypair := &models.Keypair{
		ID:        models.ULID(output[0].(string)),
		Name:      &_name,
		PublicKey: &_pubkey}

	stmt.Close()

	return keypair.NewGetKeypairByIDOK().WithPayload(_keypair)
}

func NewFindKeypairByCustomer(rt *configManager.Runtime) keypair.FindKeypairByCustomerHandler {
	return &findKeypairByCustomer{rt: rt}
}

type findKeypairByCustomer struct {
	rt *configManager.Runtime
}

func (ctx *findKeypairByCustomer) Handle(params keypair.FindKeypairByCustomerParams, principal *models.Customer) middleware.Responder {
	cypher := `MATCH (c:Customer {name: {name} })-[:HAS]->(k:Keypair)
							WHERE EXISTS(k.id)
							RETURN k.id as id,
											k.name as name,
											k.public_key as public_key`

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(principal.Name)})

	db, err := ctx.rt.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j:", err)
		return keypair.NewFindKeypairByCustomerInternalServerError()
	}
	defer db.Close()

	data, _, _, err := db.QueryNeoAll(cypher, map[string]interface{}{
		"name": swag.StringValue(principal.Name)})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: %s", err)
		return keypair.NewFindKeypairByCustomerInternalServerError()

	} else if len(data) == 0 {
		return keypair.NewFindKeypairByCustomerNotFound()
	}

	res := make([]*models.Keypair, len(data))

	for idx, row := range data {
		_name := row[1].(string)
		_pubkey := row[2].(string)

		res[idx] = &models.Keypair{
			ID:        models.ULID(row[0].(string)),
			Name:      &_name,
			PublicKey: &_pubkey}
	}

	return keypair.NewFindKeypairByCustomerOK().WithPayload(res)
}

func getCellKeypair(ctx *configManager.Runtime, customerName *string, CellID *string) *models.Keypair {

	var keypair *models.Keypair
	keypair = nil

	ctxLogger := ctx.Logger().WithFields(logrus.Fields{
		"customer_name": customerName,
		"cell_id":       CellID})

	cypher := `MATCH (c:Customer {name: {customer_name} })-[:OWN]->(cell:Cell{id: {cell_id}})-[:DEPLOY_WITH]->(keypair)
								RETURN keypair.id as id,
									keypair.name as name,
									keypair.public_key as public_key`

	db, err := ctx.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j:", err)
		return keypair
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: %s", err)
		return keypair
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"customer_name": swag.StringValue(customerName),
		"cell_id":       CellID})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: %s", err)
		return keypair
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		return keypair
	}
	_name := output[1].(string)
	_public_key := output[2].(string)

	keypair = &models.Keypair{
		ID:        models.ULID(output[0].(string)),
		Name:      &_name,
		PublicKey: &_public_key}

	stmt.Close()

	return keypair
}

func getKeypairByName(ctx *configManager.Runtime, customerName *string, keypairName *string) *models.Keypair {

	var keypair *models.Keypair
	keypair = nil

	ctxLogger := ctx.Logger().WithFields(logrus.Fields{
		"customer_name": *customerName,
		"keypair_name":  *keypairName})

	cypher := `MATCH (c:Customer {name: {name} })-[:HAS]->(k:Keypair {name: {kname}})
							WHERE EXISTS(k.id)
							RETURN k.id as id,
											k.name as name,
											k.public_key as public_key`

	db, err := ctx.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j:", err)
		return keypair
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: %s", err)
		return keypair
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"name":  swag.StringValue(customerName),
		"kname": swag.StringValue(keypairName)})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: %s", err)
		return keypair
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		ctxLogger.Error(err)
		return keypair
	}
	_name := output[1].(string)
	_pubkey := output[2].(string)

	keypair = &models.Keypair{
		ID:        models.ULID(output[0].(string)),
		Name:      &_name,
		PublicKey: &_pubkey}

	stmt.Close()

	return keypair
}
