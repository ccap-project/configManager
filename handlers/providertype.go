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
	"configManager/restapi/operations/providertype"

	"github.com/Sirupsen/logrus"
	"github.com/go-openapi/runtime/middleware"
)

func NewAddProviderType(rt *configManager.Runtime) providertype.AddProviderTypeHandler {
	return &addProviderType{rt: rt}
}

type addProviderType struct {
	rt *configManager.Runtime
}

func (ctx *addProviderType) Handle(params providertype.AddProviderTypeParams) middleware.Responder {

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"provider_type": params.Body.Name})

	cypher := `create(p:ProviderType { id: {id},
																			name: {name},
																			auth_url: {auth_url},
																			domain_name: {domain_name},
																			username: {username},
																			password: {password} }) RETURN {id}`

	if len(GetProviderTypeByName(ctx.rt, params.Body.Name).Name) > 0 {
		ctxLogger.Error("providertype already exists !")
		return providertype.NewAddProviderTypeInternalServerError().WithPayload(&models.APIResponse{Message: "providertype already exists"})
	}

	db, err := ctx.rt.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j:", err)
		return providertype.NewAddProviderTypeInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: %s", err)
		return providertype.NewAddProviderTypeInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer stmt.Close()

	ulid := configManager.GetULID()

	ctxLogger = ctx.rt.Logger().WithFields(logrus.Fields{
		"provider_type":    params.Body.Name,
		"provider_type_id": ulid})

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"id":          ulid,
		"name":        params.Body.Name,
		"auth_url":    params.Body.AuthURL,
		"domain_name": params.Body.DomainName,
		"username":    params.Body.Username,
		"password":    params.Body.Password})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: %s", err)
		return providertype.NewAddProviderTypeInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	_, _, err = rows.NextNeo()
	if err != nil {
		log.Printf("An error occurred getting next row: %s", err)
		return providertype.NewAddProviderTypeInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	ctxLogger.Info("OK")
	return providertype.NewAddProviderTypeCreated().WithPayload(models.ULID(ulid))
}

func NewGetProviderTypeByID(rt *configManager.Runtime) providertype.GetProviderTypeByIDHandler {
	return &getProviderTypeByID{rt: rt}
}

type getProviderTypeByID struct {
	rt *configManager.Runtime
}

func (ctx *getProviderTypeByID) Handle(params providertype.GetProviderTypeByIDParams) middleware.Responder {

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"provider_type_id": params.ProvidertypeID})

	cypher := `MATCH (p:ProviderType {id: {id} })
							RETURN p.id as id,
											p.name as name,
											p.auth_url as auth_url,
											p.domain_name as domain_name,
											p.username as username,
											p.password as password`

	db, err := ctx.rt.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j:", err)
		return providertype.NewGetProviderTypeByIDInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: %s", err)
		return providertype.NewGetProviderTypeByIDInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer stmt.Close()

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"customer_name": params.ProvidertypeID})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: %s", err)
		return providertype.NewGetProviderTypeByIDInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	if rows == nil {
		return providertype.NewGetProviderTypeByIDNotFound()
	}

	row, _, err := rows.NextNeo()

	if err != nil {
		ctxLogger.Error("An error occurred getting next row: %s", err)
		return providertype.NewListProviderTypesInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	provider := &models.ProviderType{
		ID:         models.ULID(row[0].(string)),
		Name:       row[1].(string),
		AuthURL:    row[1].(string),
		DomainName: row[2].(string),
		Username:   row[3].(string),
		Password:   row[4].(string)}

	return providertype.NewGetProviderTypeByIDOK().WithPayload(provider)
}

func GetProviderTypeByName(rt *configManager.Runtime, providertypeName string) *models.ProviderType {

	var providerType *models.ProviderType

	ctxLogger := rt.Logger().WithFields(logrus.Fields{
		"provider_type": providertypeName})

	cypher := `MATCH (p:ProviderType)
							WHERE EXISTS(p.id) AND p.name = {name}
							RETURN p.id as id,
											p.name as name,
											p.access_key as access_key,
											p.auth_url as auth_url,
											p.domain_name as domain_name,
											p.username as username,
											p.password as password,
											p.region as region,
											p.secret_key as secret_key`

	db, err := rt.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j:", err)
		return providerType
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: %s", err)
		return providerType
	}
	defer stmt.Close()

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"name": providertypeName})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: %s", err)
		return providerType
	}

	if rows == nil {
		return providerType
	}

	row, _, err := rows.NextNeo()

	if err != nil {
		ctxLogger.Error("An error occurred getting next row: %s", err)
		return providerType
	}

	providerType = new(models.ProviderType)
	providerType.ID = models.ULID(row[0].(string))
	providerType.Name = row[1].(string)

	if row[2] != nil {
		providerType.AccessKey = row[2].(string)
	}

	if row[3] != nil {
		providerType.AuthURL = row[3].(string)
	}

	if row[4] != nil {
		providerType.DomainName = row[4].(string)
	}

	if row[5] != nil {
		providerType.Username = row[5].(string)
	}

	if row[6] != nil {
		providerType.Password = row[6].(string)
	}
	if row[7] != nil {
		providerType.Region = row[7].(string)
	}

	if row[8] != nil {
		providerType.SecretKey = row[8].(string)
	}

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
							RETURN p.id as id,
											p.name as name,
											p.auth_url as auth_url,
											p.domain_name as domain_name,
											p.username as username,
											p.password as password`

	db, err := ctx.rt.DB().OpenPool()
	if err != nil {
		ctx.rt.Logger().Error("error connecting to neo4j:", err)
		return providertype.NewListProviderTypesInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	data, _, _, err := db.QueryNeoAll(cypher, nil)

	if err != nil {
		ctx.rt.Logger().Error("An error occurred querying Neo: %s", err)
		return providertype.NewListProviderTypesInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	res := make([]*models.ProviderType, len(data))

	for idx, row := range data {
		res[idx] = &models.ProviderType{
			ID:         models.ULID(row[0].(string)),
			Name:       row[1].(string),
			AuthURL:    row[2].(string),
			DomainName: row[3].(string),
			Username:   row[4].(string),
			Password:   row[5].(string)}
	}

	return providertype.NewListProviderTypesOK().WithPayload(res)
}

func InitProviderType(rt *configManager.Runtime) {

	rt.Logger().Info("Checking provider types...")

	if err := _addProviderType(rt, "Openstack", []string{"auth_url", "domain_name", "username", "password"}); err != nil {
		rt.Logger().Error("Error Initializing provider types, ", err)
	}
	if err := _addProviderType(rt, "AWS", []string{"access_key", "secret_key", "region"}); err != nil {
		rt.Logger().Error("Error Initializing provider types, ", err)
	}
}

func _addProviderType(rt *configManager.Runtime, name string, fields []string) error {

	var allFields []string

	if GetProviderTypeByName(rt, name) != nil {
		rt.Logger().Warnf("Provider %s already exists", name)
		return nil
	}

	createTmpl := `Create (p:ProviderType { id: '%s', name: '%s', %s })`

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

	create := fmt.Sprintf(createTmpl, configManager.GetULID(), name, strings.Join(allFields, ""))

	db, err := rt.DB().OpenPool()
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

	rt.Logger().Infof("Provider %s has been created", name)

	return nil
}
