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

	"configManager/models"
	"configManager/neo4j"
	"configManager/restapi/operations/keypair"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
)

func AddCellKeypair(params keypair.AddCellKeypairParams, principal *models.Customer) middleware.Responder {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	cypher := `MATCH (c:Customer {name: {customer_name}})-[:OWN]->(cell:Cell),
							(c:Customer {name: {customer_name}})-[:HAS]->(keypair:Keypair {name: {keypair_name}})
							WHERE id(cell) = {cell_id}
							CREATE (cell)-[:DEPLOY_WITH]->(keypair)
							RETURN	id(keypair) AS id`

	if getKeypairByName(principal.Name, &params.KeypairName) == nil {
		log.Println("keypair does not exists !")
		return keypair.NewAddCellKeypairConflict()
	}

	if getCellKeypair(principal.Name, params.CellID) != nil {
		log.Println("This Cell already has a keypair")
		return keypair.NewAddCellKeypairNotFound()
	}

	db, err := neo4j.Connect("")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return keypair.NewAddCellKeypairInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return keypair.NewAddCellKeypairInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"customer_name": swag.StringValue(principal.Name),
		"keypair_name":  params.KeypairName,
		"cell_id":       params.CellID})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return keypair.NewAddCellKeypairInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	output, _, err := rows.NextNeo()
	log.Printf("Rows(%#v)", rows)

	if err != nil {
		log.Printf("An error occurred getting next row: %s", err)
		return keypair.NewAddCellKeypairInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	log.Printf("= Output(%#v)", output)

	return keypair.NewAddCellKeypairCreated().WithPayload(output[0].(int64))
}

func AddKeypair(params keypair.AddKeypairParams, principal *models.Customer) middleware.Responder {

	cypher := `MATCH (c:Customer {name: {name} })
							CREATE (c)-[:HAS]->(k:Keypair { name: {kname}, public_key: {public_key} })
							RETURN	id(k) AS id,
											k.name AS name,
											k.public_key AS public_key`

	if getKeypairByName(principal.Name, params.Body.Name) != nil {
		log.Println("keypair already exists !")
		return keypair.NewAddKeypairConflict().WithPayload(models.APIResponse{Message: "keypair already exists"})
	}

	db, err := neo4j.Connect("")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return keypair.NewAddKeypairInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return keypair.NewAddKeypairInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"name":       swag.StringValue(principal.Name),
		"kname":      swag.StringValue(params.Body.Name),
		"public_key": swag.StringValue(params.Body.PublicKey)})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return keypair.NewAddKeypairInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		log.Printf("An error occurred getting next row: %s", err)
		return keypair.NewAddKeypairInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	log.Printf("= Output(%#v)", output)

	log.Printf("customer(%s) name(%s) key(%s)", swag.StringValue(principal.Name), swag.StringValue(params.Body.Name), swag.StringValue(params.Body.PublicKey))

	return keypair.NewAddKeypairCreated().WithPayload(output[0].(int64))
}

func GetKeypairByID(params keypair.GetKeypairByIDParams, principal *models.Customer) middleware.Responder {

	cypher := `MATCH (c:Customer {name: {name} })-[:HAS]->(k:Keypair)
								WHERE ID(k) = {kid}
								RETURN ID(c) as id,
												k.name as name,
												k.public_key as public_key`

	db, err := neo4j.Connect("")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return keypair.NewGetKeypairByIDInternalServerError()
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return keypair.NewGetKeypairByIDInternalServerError()
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{"name": swag.StringValue(principal.Name),
		"kid": params.KeypairID})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return keypair.NewGetKeypairByIDInternalServerError()
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		return keypair.NewGetKeypairByIDNotFound()
	}
	_name := output[1].(string)
	_pubkey := output[2].(string)

	_keypair := &models.Keypair{ID: output[0].(int64),
		Name:      &_name,
		PublicKey: &_pubkey}

	stmt.Close()

	return keypair.NewGetKeypairByIDOK().WithPayload(_keypair)
}

func FindKeypairByCustomer(params keypair.FindKeypairByCustomerParams, principal *models.Customer) middleware.Responder {
	cypher := `MATCH (c:Customer {name: {name} })-[:HAS]->(k:Keypair)
								RETURN ID(c) as id,
												k.name as name,
												k.public_key as public_key`

	db, err := neo4j.Connect("")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return keypair.NewFindKeypairByCustomerInternalServerError()
	}
	defer db.Close()

	data, _, _, err := db.QueryNeoAll(cypher, map[string]interface{}{
		"name": swag.StringValue(principal.Name)})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return keypair.NewFindKeypairByCustomerInternalServerError()

	} else if len(data) == 0 {
		return keypair.NewFindKeypairByCustomerNotFound()
	}

	res := make([]*models.Keypair, len(data))

	for idx, row := range data {
		_name := row[1].(string)
		_pubkey := row[2].(string)

		res[idx] = &models.Keypair{
			ID:        row[0].(int64),
			Name:      &_name,
			PublicKey: &_pubkey}
	}

	return keypair.NewFindKeypairByCustomerOK().WithPayload(res)
}

func getCellKeypair(customerName *string, CellID int64) *models.Keypair {

	var keypair *models.Keypair
	keypair = nil

	log.Printf("CustomerName(%v) CellID(%v)", customerName, CellID)

	cypher := `MATCH (c:Customer {name: {customer_name} })-[:OWN]->(cell:Cell)-[:DEPLOY_WITH]->(keypair)
							WHERE id(cell) = {cell_id}
								RETURN ID(keypair) as id,
									keypair.name as name`

	db, err := neo4j.Connect("")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return keypair
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return keypair
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"customer_name": swag.StringValue(customerName),
		"cell_id":       CellID})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return keypair
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		return keypair
	}
	_name := output[1].(string)

	keypair = &models.Keypair{
		ID:   output[0].(int64),
		Name: &_name}

	stmt.Close()

	log.Printf("Keypair => %#v", keypair)

	return keypair
}

func getKeypairByName(customerName *string, keypairName *string) *models.Keypair {

	var keypair *models.Keypair
	keypair = nil

	cypher := `MATCH (c:Customer {name: {name} })-[:HAS]->(k:Keypair)
								WHERE k.name = {kname}
								RETURN ID(c) as id,
												c.name as name,
												k.public_key as public_key`

	db, err := neo4j.Connect("")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return keypair
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return keypair
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{"name": swag.StringValue(customerName),
		"kname": swag.StringValue(keypairName)})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return keypair
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		return keypair
	}
	_name := output[1].(string)
	_pubkey := output[2].(string)

	keypair = &models.Keypair{ID: output[0].(int64),
		Name:      &_name,
		PublicKey: &_pubkey}

	stmt.Close()

	return keypair
}
