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

	"configManager"
	"configManager/models"
	"configManager/restapi/operations/regionaz"
	"configManager/util"

	"github.com/Sirupsen/logrus"
	"github.com/go-openapi/runtime/middleware"
)

func NewAddRegionAZ(rt *configManager.Runtime) regionaz.AddRegionAZHandler {
	return &addRegionAZ{rt: rt}
}

type addRegionAZ struct {
	rt *configManager.Runtime
}

func (ctx *addRegionAZ) Handle(params regionaz.AddRegionAZParams) middleware.Responder {

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"provider_region":    params.Body.Name,
		"provider_type_id":   params.ProvidertypeID,
		"provider_region_id": params.ProviderRegionID})

	regionAZ, err := _getRegionAZByName(ctx.rt, params.ProvidertypeID, params.ProviderRegionID, *params.Body.Name)
	if err != nil {
		ctxLogger.Errorf("getting region az, %s", err)
		return regionaz.NewAddRegionAZInternalServerError()
	}

	if regionAZ != nil {
		ctxLogger.Error("region az already exists !")
		return regionaz.NewAddRegionAZInternalServerError().WithPayload(&models.APIResponse{Message: "region az already exists"})
	}

	ulid, err := _addRegionAZ(ctx.rt, params.ProvidertypeID, params.ProviderRegionID, *params.Body.Name)

	if err != nil {
		ctxLogger.Errorf("adding region az, %s", err)
		return regionaz.NewAddRegionAZInternalServerError()
	}

	ctxLogger.Info("OK")
	return regionaz.NewAddRegionAZCreated().WithPayload(models.ULID(*ulid))
}

func NewGetRegionAZByID(rt *configManager.Runtime) regionaz.GetRegionAZByIDHandler {
	return &getRegionAZByID{rt: rt}
}

type getRegionAZByID struct {
	rt *configManager.Runtime
}

func (ctx *getRegionAZByID) Handle(params regionaz.GetRegionAZByIDParams) middleware.Responder {

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"provider_type_id":   params.ProvidertypeID,
		"provider_region_id": params.ProviderRegionID})

	cypher := `MATCH (provider:ProviderType {id: {provider_id}})
							-[:HAS]->(region:ProviderRegion{id: {region_id}})
							-[:HAS]->(az:RegionAZ{id: {az_id}})
							RETURN az.id as id,
											az.name as name`

	db, err := ctx.rt.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j:", err)
		return regionaz.NewGetRegionAZByIDInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: %s", err)
		return regionaz.NewGetRegionAZByIDInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer stmt.Close()

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"provider_id": params.ProvidertypeID,
		"region_id":   params.ProviderRegionID,
		"az_id":       params.RegionAzID})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: %s", err)
		return regionaz.NewGetRegionAZByIDInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	if rows == nil {
		return regionaz.NewGetRegionAZByIDNotFound()
	}

	row, _, err := rows.NextNeo()

	if err != nil {
		ctxLogger.Error("An error occurred getting next row: %s", err)
		return regionaz.NewListRegionAZsInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	_name := row[1].(string)

	provider := &models.RegionAZ{
		ID:   models.ULID(row[0].(string)),
		Name: &_name}

	return regionaz.NewGetRegionAZByIDOK().WithPayload(provider)
}

func NewListRegionAZs(rt *configManager.Runtime) regionaz.ListRegionAZsHandler {
	return &listRegionAZs{rt: rt}
}

type listRegionAZs struct {
	rt *configManager.Runtime
}

func (ctx *listRegionAZs) Handle(params regionaz.ListRegionAZsParams) middleware.Responder {

	azs, err := _listRegionAZs(ctx.rt, &params.ProvidertypeID, &params.ProviderRegionID)

	if err != nil {
		return regionaz.NewListRegionAZsInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	return regionaz.NewListRegionAZsOK().WithPayload(azs)
}

func _addRegionAZ(rt *configManager.Runtime, ProviderID string, RegionID string, AZName string) (*string, error) {

	/*
	 * Check if Region already exists
	 */
	regionAZ, err := _getRegionAZByName(rt, ProviderID, RegionID, AZName)
	if err != nil {
		return nil, fmt.Errorf("getting region az, %s", err)
	}

	if regionAZ != nil {
		return nil, fmt.Errorf("region az already exists !")
	}

	query := `MATCH (provider:ProviderType {id: {provider_id}})
							-[:HAS]->(region:ProviderRegion{id: {region_id}})
							MERGE (region)-[:HAS]->(az:RegionAZ {
									 id: {az_id},
									name: {az_name}})
								RETURN az.id`

	ulid := configManager.GetULID()

	params := map[string]interface{}{
		"provider_id": ProviderID,
		"region_id":   RegionID,
		"az_id":       ulid,
		"az_name":     AZName}

	_, err = rt.QueryDB(&query, &params)

	if err != nil {
		return nil, err
	}

	return &ulid, nil
}

func _getRegionAZByName(rt *configManager.Runtime, provider_id string, region_id string, az string) (*models.RegionAZ, error) {

	var regionAZ *models.RegionAZ

	query := `MATCH (provider:ProviderType {id: {provider_id}})
							-[:HAS]->(region:ProviderRegion{id: {region_id}})
							-[:HAS]->(az:RegionAZ{name: {az_name}})
							RETURN az {.*}`

	params := map[string]interface{}{
		"provider_id": provider_id,
		"region_id":   region_id,
		"az_name":     az}

	output, err := rt.QueryDB(&query, &params)

	if err != nil {
		return regionAZ, err
	}

	if len(output) > 0 {
		regionAZ = new(models.RegionAZ)
		util.FillStruct(regionAZ, output[0].(map[string]interface{}))
	}

	return regionAZ, nil
}

func _listRegionAZs(rt *configManager.Runtime, provider_id *string, region_id *string) ([]*models.RegionAZ, error) {

	var azs []*models.RegionAZ

	cypher := `MATCH (provider:ProviderType {id: {provider_id}})
							-[:HAS]->(region:ProviderRegion{id: {region_id}})
							-[:HAS]->(az:RegionAZ)
							RETURN az.id as id,
											az.name as name`

	ctxLogger := rt.Logger().WithFields(logrus.Fields{
		"provider_type_id": provider_id,
		"region_id":        region_id})

	db, err := rt.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j: ", err)
		return azs, err
	}
	defer db.Close()

	data, _, _, err := db.QueryNeoAll(cypher, map[string]interface{}{
		"provider_id": *provider_id,
		"region_id":   *region_id})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: ", err)
		return azs, err
	}

	for _, row := range data {
		_name := row[1].(string)

		az := &models.RegionAZ{
			ID:   models.ULID(row[0].(string)),
			Name: &_name}
		azs = append(azs, az)
	}

	return azs, nil
}
