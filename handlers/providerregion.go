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
	"configManager/restapi/operations/providerregion"
	"configManager/util"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/go-openapi/runtime/middleware"
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

	providerType, err := _getProviderTypeByID(ctx.rt, params.ProvidertypeID)
	if err != nil {
		ctxLogger.Errorf("getting providertype, %s", err)
		return providerregion.NewAddProviderRegionInternalServerError()
	}

	if providerType == nil {
		ctxLogger.Error("providertype does not exists")
		return providerregion.NewAddProviderRegionInternalServerError().WithPayload(&models.APIResponse{Message: "providertype does not exists"})
	}

	err = _addProviderRegion(ctx.rt, providerType.Name, *params.Body.Name)

	if err != nil {
		ctxLogger.Error("Adding Provider Region: %s", err)
		return providerregion.NewAddProviderRegionInternalServerError()
	}

	ctxLogger.Info("OK")
	return providerregion.NewAddProviderRegionCreated()
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

func _addProviderRegion(rt *configManager.Runtime, ProvidertypeName string, RegionName string) error {

	cypher := `MATCH (provider:ProviderType {name: {provider_name}})
							MERGE (provider)-[:HAS]->(region:ProviderRegion {
									 id: {id},
									name: {name}})
								RETURN region.id`

	/*
	 * Validate Provider Type
	 */
	providerType, err := _getProviderTypeByName(rt, ProvidertypeName)

	if err != nil {
		return fmt.Errorf("getting providertype, %s", err)
	}

	if providerType == nil {
		return fmt.Errorf("providertype does not exists")
	}

	/*
	 * Validate Provider Region
	 */
	providerRegion, err := _getProviderRegionByName(rt, ProvidertypeName, RegionName)

	if err != nil {
		return fmt.Errorf("getting provider region, %s", err)
	}

	if providerRegion != nil {
		return fmt.Errorf("provider region already exists")
	}

	db, err := rt.DB().OpenPool()
	if err != nil {
		return fmt.Errorf("error connecting to neo4j: %s", err)
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		return fmt.Errorf("An error occurred preparing statement: %s", err)
	}
	defer stmt.Close()

	ulid := configManager.GetULID()

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"provider_name": ProvidertypeName,
		"id":            ulid,
		"name":          RegionName})

	if err != nil {
		return fmt.Errorf("An error occurred querying Neo: %s", err)
	}

	_, _, err = rows.NextNeo()
	if err != nil {
		return fmt.Errorf("An error occurred getting next row: %s", err)
	}

	return nil
}

func _getProviderRegionByName(rt *configManager.Runtime, ProvidertypeName string, RegionName string) (*models.ProviderRegion, error) {

	var providerRegion *models.ProviderRegion

	query := `MATCH (provider:ProviderType {name: {provider_name}})-[:HAS]->
										(region:ProviderRegion {name: {region_name}})
							RETURN region {.*}`

	params := map[string]interface{}{
		"provider_name": ProvidertypeName,
		"region_name":   RegionName}

	output, err := rt.QueryDB(&query, &params)

	if err != nil {
		return providerRegion, err
	}

	if len(output) > 0 {
		providerRegion = new(models.ProviderRegion)
		util.FillStruct(providerRegion, output[0].(map[string]interface{}))
	}

	return providerRegion, nil
}
