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
	"fmt"
	"log"
	"strings"

	"configManager"
	"configManager/models"
	"configManager/neo4j"
	"configManager/restapi/operations/providertype"

	"github.com/go-openapi/runtime/middleware"
)

func NewAddProviderType(rt *configManager.Runtime) providertype.AddProviderTypeHandler {
	return &addProviderType{rt: rt}
}

type addProviderType struct {
	rt *configManager.Runtime
}

func (ctx *addProviderType) Handle(params providertype.AddProviderTypeParams) middleware.Responder {

	cypher := `create(p:ProviderType { name: {name},
																			auth_url: {auth_url},
																			domain_name: {domain_name},
																			username: {username},
																			password: {password} }) RETURN ID(p)`

	if len(GetProviderTypeByName(ctx.rt.DB(), params.Body.Name).Name) > 0 {
		log.Println("providertype already exists !")
		return providertype.NewAddProviderTypeInternalServerError().WithPayload(models.APIResponse{Message: "providertype already exists"})
	}

	db, err := ctx.rt.DB().OpenPool()
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return providertype.NewAddProviderTypeInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Printf("An error occurred beginning transaction: %s", err)
		return providertype.NewAddProviderTypeInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}
	defer tx.Rollback()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return providertype.NewAddProviderTypeInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	defer stmt.Close()

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"name":        params.Body.Name,
		"auth_url":    params.Body.AuthURL,
		"domain_name": params.Body.DomainName,
		"username":    params.Body.Username,
		"password":    params.Body.Password})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return providertype.NewAddProviderTypeInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	_, _, err = rows.NextNeo()
	if err != nil {
		log.Printf("An error occurred getting next row: %s", err)
		return providertype.NewAddProviderTypeInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	tx.Commit()

	return providertype.NewAddProviderTypeCreated().WithPayload("OK")
}

func NewGetProviderTypeByID(rt *configManager.Runtime) providertype.GetProviderTypeByIDHandler {
	return &getProviderTypeByID{rt: rt}
}

type getProviderTypeByID struct {
	rt *configManager.Runtime
}

func (ctx *getProviderTypeByID) Handle(params providertype.GetProviderTypeByIDParams) middleware.Responder {

	cypher := `MATCH (p:ProviderType)
							WHERE ID(p) = {id}
							RETURN ID(p) as id,
											p.name as name,
											p.auth_url as auth_url,
											p.domain_name as domain_name,
											p.username as username,
											p.password as password`

	db, err := ctx.rt.DB().OpenPool()
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return providertype.NewGetProviderTypeByIDInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return providertype.NewGetProviderTypeByIDInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	defer stmt.Close()

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"customer_name": params.ProvidertypeID})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return providertype.NewGetProviderTypeByIDInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	if rows == nil {
		return providertype.NewGetProviderTypeByIDNotFound()
	}

	row, _, err := rows.NextNeo()

	if err != nil {
		log.Printf("An error occurred getting next row: %s", err)
		return providertype.NewListProviderTypesInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	provider := &models.ProviderType{
		ID:         row[0].(int64),
		Name:       row[1].(string),
		AuthURL:    row[1].(string),
		DomainName: row[2].(string),
		Username:   row[3].(string),
		Password:   row[4].(string)}

	return providertype.NewGetProviderTypeByIDOK().WithPayload(provider)
}

func GetProviderTypeByName(conn neo4j.ConnPool, providertypeName string) models.ProviderType {

	var providerType models.ProviderType

	cypher := `MATCH (p:ProviderType)
							WHERE p.name = {name}
							RETURN ID(p) as id,
											p.name as name,
											p.auth_url as auth_url,
											p.domain_name as domain_name,
											p.username as username,
											p.password as password`

	db, err := conn.OpenPool()
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return providerType
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return providerType
	}

	defer stmt.Close()

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"name": providertypeName})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return providerType
	}

	if rows == nil {
		return providerType
	}

	row, _, err := rows.NextNeo()

	if err != nil {
		log.Printf("An error occurred getting next row: %s", err)
		log.Printf("%#v", row)
		return providerType
	}

	providerType.ID = row[0].(int64)
	providerType.Name = row[1].(string)
	providerType.AuthURL = row[2].(string)
	providerType.DomainName = row[3].(string)
	providerType.Username = row[4].(string)
	providerType.Password = row[5].(string)

	return providerType
}

func NewListProviderTypes(rt *configManager.Runtime) providertype.ListProviderTypesHandler {
	return &listProviderTypes{rt: rt}
}

type listProviderTypes struct {
	rt *configManager.Runtime
}

func (ctx *listProviderTypes) Handle(params providertype.ListProviderTypesParams) middleware.Responder {

	cypher := `MATCH (p:ProviderType)
							RETURN ID(p) as id,
											p.name as name,
											p.auth_url as auth_url,
											p.domain_name as domain_name,
											p.username as username,
											p.password as password`

	db, err := ctx.rt.DB().OpenPool()
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return providertype.NewListProviderTypesInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	data, _, _, err := db.QueryNeoAll(cypher, nil)

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return providertype.NewListProviderTypesInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	res := make([]*models.ProviderType, len(data))

	for idx, row := range data {
		res[idx] = &models.ProviderType{
			ID:         row[0].(int64),
			Name:       row[1].(string),
			AuthURL:    row[2].(string),
			DomainName: row[3].(string),
			Username:   row[4].(string),
			Password:   row[5].(string)}
	}

	return providertype.NewListProviderTypesOK().WithPayload(res)
}

func InitProviderType(rt *configManager.Runtime) {

	log.Printf("Checking provider types...")

	if err := _addProviderType(rt.DB(), "Openstack", []string{"auth_url", "domain_name", "username", "password"}); err != nil {
		log.Println("Error Initializing provider types, ", err)
	}
}

func _addProviderType(conn neo4j.ConnPool, name string, fields []string) error {

	var allFields []string

	if len(GetProviderTypeByName(conn, name).Name) > 0 {
		log.Printf("Provider %s already exists", name)
		return nil
	}

	createTmpl := `Create (p:ProviderType { name: '%s', %s })`

	lastField := len(fields)

	if lastField <= 0 {
		return fmt.Errorf("No fields specified !")
	} else {
		lastField -= 1
	}

	for i := 0; i < lastField; i++ {
		allFields = append(allFields, fmt.Sprintf("%s: '%s', ", fields[i], fields[i]))
	}

	allFields = append(allFields, fmt.Sprintf("%s: '%s'", fields[lastField], fields[lastField]))

	create := fmt.Sprintf(createTmpl, name, strings.Join(allFields, ""))

	db, err := conn.OpenPool()
	if err != nil {
		return fmt.Errorf("error connecting to neo4j:", err)
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(create)
	if err != nil {
		return fmt.Errorf("An error occurred preparing statement: %s", err)
	}

	defer stmt.Close()

	_, err = stmt.QueryNeo(nil)

	if err != nil {
		return fmt.Errorf("An error occurred querying Neo: %s", err)
	}

	log.Printf("Provider %s has been created", name)

	return nil
}
