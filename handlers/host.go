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
	"configManager/restapi/operations/host"

	"github.com/Sirupsen/logrus"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
)

func NewAddCellHost(rt *configManager.Runtime) host.AddCellHostHandler {
	return &addCellHost{rt: rt}
}

type addCellHost struct {
	rt *configManager.Runtime
}

func (ctx *addCellHost) Handle(params host.AddCellHostParams, principal *models.Customer) middleware.Responder {

	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->(cell:Cell {id: {cell_id}})
							CREATE (cell)-[:HAS]->(host:Host {id: {host_id}, name: {host_name}} )
								RETURN host.id as id`

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer":  swag.StringValue(principal.Name),
		"cell_id":   params.CellID,
		"host_name": params.Body.Name})

	if getCellHostByName(ctx.rt.DB(), principal.Name, &params.CellID, params.Body.Name) != nil {
		ctxLogger.Warn("host already exists !")
		return host.NewAddCellHostConflict().WithPayload(models.APIResponse{Message: "host already exists"})
	}

	db, err := ctx.rt.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j:", err)
		return host.NewAddCellHostInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		ctxLogger.Error("An error occurred beginning transaction: %s", err)
		return host.NewAddCellHostInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}
	defer tx.Rollback()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: %s", err)
		return host.NewAddCellHostInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	ulid := configManager.GetULID()

	ctxLogger = ctxLogger.WithFields(logrus.Fields{
		"host_id": ulid})

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID,
		"host_id":       ulid,
		"host_name":     swag.StringValue(params.Body.Name)})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: ", err)
		return host.NewAddCellHostInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		ctxLogger.Error("An error occurred getting next row: ", err)
		return host.NewAddCellHostInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	stmt.Close()

	err = addCellHostOptions(db, ctxLogger, principal.Name, &params.CellID, params.Body.Name, params.Body.Options)
	if err != nil {
		ctxLogger.Error("An error occurred adding Host options: ", err)
		return host.NewAddCellHostInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}
	tx.Commit()

	return host.NewAddCellHostCreated().WithPayload(models.ULID(output[0].(string)))
}

/*
func FindCellHosts(params host.FindCellHostsParams, principal *models.Customer) middleware.Responder {

	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->
							(cell:Cell {id: {cell_id}})-[:HAS]->(host:Host)
								RETURN host.id as id,
												host.name as name`

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer": swag.StringValue(principal.Name),
		"cell_id":  params.CellID})

	db, err := neo4j.Connect("")
	if err != nil {
		ctxLogger.Error("error connecting to neo4j: ", err)
		return host.NewFindCellHostsInternalServerError()
	}
	defer db.Close()

	data, _, _, err := db.QueryNeoAll(cypher, map[string]interface{}{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: ", err)
		return host.NewFindCellHostsInternalServerError()
	}

	res := make([]*models.Host, len(data))

	for idx, row := range data {
		_name := row[1].(string)

		res[idx] = &models.Host{
			ID:   row[0].(int64),
			Name: &_name}
	}

	return host.NewFindCellHostsOK().WithPayload(res)
}*/

func addCellHostOptions(db neo4j.Conn, ctxLogger *logrus.Entry, customer *string, cellID *string, hostName *string, options []*models.Parameter) error {

	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->
							(cell:Cell {id: {cell_id}})-[:HAS]->(host:Host {name: {host_name}})
							CREATE (host)-[:OPT]->(option:Option {
									id: {opt_id},
									name: {opt_name},
									value: {opt_val}} )
								RETURN option.id as id`

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: ", err)
		return err
	}

	// add parameters
	for _, option := range options {
		ulid := configManager.GetULID()
		_, err := stmt.ExecNeo(map[string]interface{}{
			"customer_name": swag.StringValue(customer),
			"cell_id":       cellID,
			"host_name":     swag.StringValue(hostName),
			"opt_id":        ulid,
			"opt_name":      swag.StringValue(option.Name),
			"opt_val":       swag.StringValue(option.Value)})

		if err != nil {
			ctxLogger.Error("An error occurred querying Neo: %s", err)
			return err
		}
	}

	return nil
}

func getCellHostByName(conn neo4j.ConnPool, customer *string, cellID *string, hostName *string) *models.Host {

	var host *models.Host
	host = nil

	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->
							(cell:Cell {id: {cell_id}})-[:HAS]->(host:Host {name: {host_name}})
								RETURN host.id as id,
												host.name as name`

	db, err := conn.OpenPool()
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return host
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return host
	}
	defer stmt.Close()

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"customer_name": swag.StringValue(customer),
		"cell_id":       cellID,
		"host_name":     swag.StringValue(hostName)})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return host
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		//log.Printf("An error occurred fetching row: %s", err)
		return host
	}
	_name := output[1].(string)

	host = &models.Host{
		ID:   models.ULID(output[0].(string)),
		Name: &_name}

	return host
}
