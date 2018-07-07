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
	"strings"

	"configManager"
	"configManager/models"
	"configManager/restapi/operations/providertype"
	"configManager/util"

	"github.com/Sirupsen/logrus"
	"github.com/go-openapi/runtime/middleware"
)

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

	if err := _addProviderType(rt, "GCP", []string{"tenant_name", "secret_key", "region"}); err != nil {
		rt.Logger().Error("Error Initializing provider types, ", err)
	}
	/*
		if err := _addProviderRegion(rt, "GCP", "us-west1"); err != nil {
			rt.Logger().Error("Error Initializing provider types, ", err)
		}
		if err := _addProviderRegion(rt, "GCP", "us-central1"); err != nil {
			rt.Logger().Error("Error Initializing provider types, ", err)
		}
	*/
}

func _addProviderType(rt *configManager.Runtime, name string, fields []string) error {

	var allFields []string

	providerType, err := _getProviderTypeByName(rt, name)
	if providerType != nil {
		return fmt.Errorf("providertype %s already exists", name)
	}

	if err != nil {
		return fmt.Errorf("getting providertype, %s", err)
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

func _getProviderTypeByName(rt *configManager.Runtime, providertypeName string) (*models.ProviderType, error) {

	var providertype *models.ProviderType

	query := `MATCH (p:ProviderType)
							WHERE EXISTS(p.id) AND p.name = {name}
							RETURN p {.*}`

	params := map[string]interface{}{
		"name": providertypeName}

	output, err := rt.QueryDB(&query, &params)

	if err != nil {
		return providertype, err
	}

	if len(output) > 0 {
		providertype = new(models.ProviderType)
		util.FillStruct(providertype, output[0].(map[string]interface{}))
	}

	return providertype, nil
}

func _getProviderTypeByID(rt *configManager.Runtime, providertypeID string) (*models.ProviderType, error) {

	var providertype *models.ProviderType

	query := `MATCH (p:ProviderType)
							WHERE EXISTS(p.id) AND p.id = {id}
							RETURN p {.*}`

	params := map[string]interface{}{
		"id": providertypeID}

	output, err := rt.QueryDB(&query, &params)

	if err != nil {
		return providertype, err
	}

	if len(output) > 0 {
		providertype = new(models.ProviderType)
		util.FillStruct(providertype, output[0].(map[string]interface{}))
	}

	return providertype, nil
}
