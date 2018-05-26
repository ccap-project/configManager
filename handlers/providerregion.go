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
	"configManager/restapi/operations/providerregion"

	"github.com/Sirupsen/logrus"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
)

func NewAddProviderRegion(rt *configManager.Runtime) providerregion.AddProviderRegionHandler {
	return &addProviderRegion{rt: rt}
}

type addProviderRegion struct {
	rt *configManager.Runtime
}

func (ctx *addProviderRegion) Handle(params providerregion.AddProviderRegionParams) middleware.Responder {

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"provider_region": params.Body.Name})

	cypher := `MATCH (provider:ProviderType {id: {provider_id}})
							MERGE (provider)-[:HAS]->(region:ProviderRegion {
									 id: {id},
									name: {name}})
								RETURN region.id`

	if GetProviderRegionByName(ctx.rt, *params.Body.Name) != nil {
		ctxLogger.Error("providerregion already exists !")
		return providerregion.NewAddProviderRegionInternalServerError().WithPayload(&models.APIResponse{Message: "providerregion already exists"})
	}

	db, err := ctx.rt.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j:", err)
		return providerregion.NewAddProviderRegionInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: %s", err)
		return providerregion.NewAddProviderRegionInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer stmt.Close()

	ulid := configManager.GetULID()

	ctxLogger = ctx.rt.Logger().WithFields(logrus.Fields{
		"provider_id": params.ProvidertypeID})

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"provider_id": params.ProvidertypeID,
		"id":          ulid,
		"name":        swag.StringValue(params.Body.Name)})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: %s", err)
		return providerregion.NewAddProviderRegionInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	_, _, err = rows.NextNeo()
	if err != nil {
		log.Printf("An error occurred getting next row: %s", err)
		return providerregion.NewAddProviderRegionInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	ctxLogger.Info("OK")
	return providerregion.NewAddProviderRegionCreated().WithPayload(models.ULID(ulid))
}

func NewGetProviderRegionByID(rt *configManager.Runtime) providerregion.GetProviderRegionByIDHandler {
	return &getProviderRegionByID{rt: rt}
}

type getProviderRegionByID struct {
	rt *configManager.Runtime
}

func (ctx *getProviderRegionByID) Handle(params providerregion.GetProviderRegionByIDParams) middleware.Responder {

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"provider_region_id": params.ProviderRegionID})

	cypher := `MATCH (p:ProviderRegion {id: {id} })
							RETURN p.id as id,
											p.name as name`

	db, err := ctx.rt.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j:", err)
		return providerregion.NewGetProviderRegionByIDInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: %s", err)
		return providerregion.NewGetProviderRegionByIDInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer stmt.Close()

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"customer_name": params.ProviderRegionID})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: %s", err)
		return providerregion.NewGetProviderRegionByIDInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	if rows == nil {
		return providerregion.NewGetProviderRegionByIDNotFound()
	}

	row, _, err := rows.NextNeo()

	if err != nil {
		ctxLogger.Error("An error occurred getting next row: %s", err)
		return providerregion.NewListProviderRegionsInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	_name := row[1].(string)

	provider := &models.ProviderRegion{
		ID:   models.ULID(row[0].(string)),
		Name: &_name}

	return providerregion.NewGetProviderRegionByIDOK().WithPayload(provider)
}

func GetProviderRegionByName(rt *configManager.Runtime, providerregionName string) *models.ProviderRegion {

	var providerRegion *models.ProviderRegion

	ctxLogger := rt.Logger().WithFields(logrus.Fields{
		"provider_region": providerregionName})

	cypher := `MATCH (p:ProviderRegion)
							WHERE EXISTS(p.id) AND p.name = {name}
							RETURN p.id as id,
											p.name as name`

	db, err := rt.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j:", err)
		return providerRegion
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: %s", err)
		return providerRegion
	}
	defer stmt.Close()

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"name": providerregionName})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: %s", err)
		return providerRegion
	}

	if rows == nil {
		return providerRegion
	}

	row, _, err := rows.NextNeo()

	if err != nil {
		ctxLogger.Error("An error occurred getting next row: %s", err)
		return providerRegion
	}

	_name := row[1].(string)

	providerRegion = new(models.ProviderRegion)
	providerRegion.ID = models.ULID(row[0].(string))
	providerRegion.Name = &_name

	return providerRegion
}

func NewListProviderRegions(rt *configManager.Runtime) providerregion.ListProviderRegionsHandler {
	return &listProviderRegions{rt: rt}
}

type listProviderRegions struct {
	rt *configManager.Runtime
}

func (ctx *listProviderRegions) Handle(params providerregion.ListProviderRegionsParams) middleware.Responder {

	cypher := `MATCH (p:ProviderRegion)
							RETURN p.id as id,
											p.name as name`

	db, err := ctx.rt.DB().OpenPool()
	if err != nil {
		ctx.rt.Logger().Error("error connecting to neo4j:", err)
		return providerregion.NewListProviderRegionsInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	data, _, _, err := db.QueryNeoAll(cypher, nil)

	if err != nil {
		ctx.rt.Logger().Error("An error occurred querying Neo: %s", err)
		return providerregion.NewListProviderRegionsInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	res := make([]*models.ProviderRegion, len(data))

	for idx, row := range data {
		_name := row[1].(string)

		res[idx] = &models.ProviderRegion{
			ID:   models.ULID(row[0].(string)),
			Name: &_name}
	}

	return providerregion.NewListProviderRegionsOK().WithPayload(res)
}

/*
func InitProviderType(rt *configManager.Runtime) {

	rt.Logger().Info("Checking provider types...")

	if err := _addProviderType(rt, "Openstack", []string{"auth_url", "domain_name", "username", "password"}); err != nil {
		rt.Logger().Error("Error Initializing provider types, ", err)
	}
	if err := _addProviderType(rt, "AWS", []string{"access_key", "secret_key", "region"}); err != nil {
		rt.Logger().Error("Error Initializing provider types, ", err)
	}
}
*/

func _addProviderRegion(rt *configManager.Runtime, name string, fields []string) error {

	var allFields []string

	if GetProviderRegionByName(rt, name) != nil {
		rt.Logger().Warnf("Provider %s already exists", name)
		return nil
	}

	createTmpl := `Create (p:ProviderRegion { id: '%s', name: '%s', %s })`

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
