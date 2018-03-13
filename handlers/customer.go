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
	"configManager/restapi/operations/customer"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
)

func NewAddCustomer(rt *configManager.Runtime) customer.AddCustomerHandler {
	return &addCustomer{rt: rt}
}

type addCustomer struct {
	rt *configManager.Runtime
}

func (ctx *addCustomer) Handle(params customer.AddCustomerParams) middleware.Responder {

	cypher := `create(c:Customer { id: {id}, name: {name} }) RETURN c.id`

	if _getCustomerByName(ctx.rt.DB(), params.Body.Name) != nil {
		log.Println("customer already exists !")
		return customer.NewAddCustomerConflict().WithPayload(models.APIResponse{Message: "customer name already exists"})
	}

	db, err := ctx.rt.DB().OpenPool()
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return customer.NewAddCustomerInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return customer.NewAddCustomerInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"id":   configManager.GetULID(),
		"name": swag.StringValue(params.Body.Name)})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return customer.NewAddCustomerInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	output, _, err := rows.NextNeo()

	if err != nil {
		log.Printf("An error occurred getting next row: %s", err)
		return customer.NewAddCustomerInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	return customer.NewAddCustomerCreated().WithPayload(models.ULID(output[0].(string)))
}

func NewDeleteCustomer(rt *configManager.Runtime) customer.DeleteCustomerHandler {
	return &deleteCustomer{rt: rt}
}

type deleteCustomer struct {
	rt *configManager.Runtime
}

func (ctx *deleteCustomer) Handle(params customer.DeleteCustomerParams) middleware.Responder {

	cypher := `MATCH (c:Customer { id: {id}}) DELETE c`

	db, err := ctx.rt.DB().OpenPool()
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return customer.NewDeleteCustomerInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return customer.NewDeleteCustomerInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	_, err = stmt.ExecNeo(map[string]interface{}{
		"id": params.CustomerID})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return customer.NewDeleteCustomerInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	return customer.NewDeleteCustomerOK()
}

func NewFindCustomerByName(rt *configManager.Runtime) customer.FindCustomerByNameHandler {
	return &getCustomerByName{rt: rt}
}

type getCustomerByName struct {
	rt *configManager.Runtime
}

func (ctx *getCustomerByName) Handle(params customer.FindCustomerByNameParams) middleware.Responder {

	Customer := _getCustomerByName(ctx.rt.DB(), &params.CustomerName)

	if len(swag.StringValue(Customer.Name)) <= 0 {
		return customer.NewFindCustomerByNameNotFound()
	}

	return customer.NewFindCustomerByNameOK().WithPayload(Customer)
}

func _getCustomerByName(conn neo4j.ConnPool, customerName *string) *models.Customer {

	var customer *models.Customer

	cypher := `MATCH (c:Customer)
							WHERE c.name =~ {customer_name}
								AND EXISTS(c.id)
 							RETURN c.id as id,
											c.name as name`

	db, err := conn.OpenPool()

	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return customer
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return customer
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"customer_name": swag.StringValue(customerName)})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return customer
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		return customer
	}
	_name := output[1].(string)

	customer = &models.Customer{
		ID:   models.ULID(output[0].(string)),
		Name: &_name}

	stmt.Close()

	return customer
}
