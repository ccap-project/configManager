/*
 *
 * Copyright (c) 2018 Alexandre Biancalana <ale@biancalanas.net>.
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
	"configManager/restapi/operations/network"

	"github.com/Sirupsen/logrus"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
)

func NewAddCellNetwork(rt *configManager.Runtime) network.AddNetworkHandler {
	return &addCellNetwork{rt: rt}
}

type addCellNetwork struct {
	rt *configManager.Runtime
}

func (ctx *addCellNetwork) Handle(params network.AddNetworkParams, principal *models.Customer) middleware.Responder {

	cypher := `MATCH (c:Customer {name: {name} })-[:OWN]->(cell:Cell {id: {cell_id}})
							CREATE (cell)-[:HAS]->(network:Network {
								id: {network_id},
								name: {network_name},
								cidr: {network_cidr}})
							RETURN	network.id AS id`

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID,
		"network_name":  swag.StringValue(params.Body.Name)})

	// XXX: Consistency check should have more than only name...
	if _getNetworkByName(ctx.rt, principal.Name, &params.CellID, params.Body.Name) != nil {
		ctxLogger.Warn("network already exists !")
		return network.NewAddNetworkConflict().WithPayload(&models.APIResponse{Message: "network already exists"})
	}

	// Check if required az exists
	var networkAZs []*models.RegionAZ

	cellAZs, err := listCellAZs(ctx.rt, &params.CellID)

	if err != nil {
		ctxLogger.Error("error listing cell azs ", err)
		return network.NewAddNetworkInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	for _, r_az := range params.Body.RegionAz {
		exists := 0
		for _, az := range cellAZs {
			if r_az == *az.Name {
				exists++
				networkAZs = append(networkAZs, az)
				break
			}
		}

		if exists == 0 {
			ctxLogger.Warnf("region az (%s) does not exists !", r_az)
			return network.NewAddNetworkConflict().WithPayload(&models.APIResponse{Message: "region az does not exists"})
		}
	}

	db, err := ctx.rt.DB().OpenPool()

	if err != nil {
		ctxLogger.Error("error connecting to neo4j: ", err)
		return network.NewAddNetworkInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: ", err)
		return network.NewAddNetworkInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	ulid := configManager.GetULID()

	ctxLogger = ctxLogger.WithFields(logrus.Fields{
		"network_id": ulid})

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"name":         swag.StringValue(principal.Name),
		"cell_id":      params.CellID,
		"network_id":   ulid,
		"network_name": swag.StringValue(params.Body.Name),
		"network_cidr": swag.StringValue(params.Body.Cidr)})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: ", err)
		return network.NewAddNetworkInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		ctxLogger.Error("An error occurred getting next row: ", err)
		return network.NewAddNetworkInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	for _, az := range networkAZs {
		_connectToAZ(ctx.rt, ulid, string(az.ID))
	}

	ctxLogger.Info("OK")

	return network.NewAddNetworkCreated().WithPayload(models.ULID(output[0].(string)))
}

func NewDeleteCellNetwork(rt *configManager.Runtime) network.DeleteCellNetworkHandler {
	return &deleteCellNetwork{rt: rt}
}

type deleteCellNetwork struct {
	rt *configManager.Runtime
}

func (ctx *deleteCellNetwork) Handle(params network.DeleteCellNetworkParams, principal *models.Customer) middleware.Responder {

	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->
							(cell:Cell{id: {cell_id}})-[:HAS]->
							(network:Network {id: {network_id}})
						DETACH DELETE network`

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID,
		"network_id":    params.NetworkID})

	cell, _ := _getCellNetwork(ctx.rt, principal.Name, &params.CellID, &params.NetworkID)
	if cell == nil {
		ctxLogger.Error("network does not exists !")
		return network.NewDeleteCellNetworkNotFound()
	}

	db, err := ctx.rt.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j: ", err)
		return network.NewDeleteCellNetworkInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: ", err)
		return network.NewDeleteCellNetworkInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer stmt.Close()

	_, err = stmt.ExecNeo(map[string]interface{}{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID,
		"network_id":    params.NetworkID})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: ", err)
		return network.NewDeleteCellNetworkInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	return network.NewDeleteCellNetworkOK()
}

func NewGetCellNetwork(rt *configManager.Runtime) network.GetCellNetworkHandler {
	return &getCellNetwork{rt: rt}
}

type getCellNetwork struct {
	rt *configManager.Runtime
}

func (ctx *getCellNetwork) Handle(params network.GetCellNetworkParams, principal *models.Customer) middleware.Responder {

	cellNetwork, err := _getCellNetwork(ctx.rt, principal.Name, &params.CellID, &params.NetworkID)

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return network.NewGetCellNetworkInternalServerError()
	}

	if cellNetwork == nil {
		return network.NewGetCellNetworkOK()
	}

	return network.NewGetCellNetworkOK().WithPayload(cellNetwork)
}

func NewFindCellNetworks(rt *configManager.Runtime) network.FindCellNetworksHandler {
	return &findCellNetworks{rt: rt}
}

type findCellNetworks struct {
	rt *configManager.Runtime
}

func (ctx *findCellNetworks) Handle(params network.FindCellNetworksParams, principal *models.Customer) middleware.Responder {

	cellNetworks, err := _findCellNetworks(ctx.rt, principal.Name, &params.CellID)

	if err != nil {
		return network.NewFindCellNetworksInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	return network.NewFindCellNetworksOK().WithPayload(cellNetworks)
}

func _findCellNetworks(rt *configManager.Runtime, customerName *string, CellID *string) ([]*models.Network, error) {

	var res []*models.Network

	cypher := `MATCH (c:Customer {name: {name} })-[:OWN]->(cell:Cell {id: {cell_id}})-[:HAS]->(network:Network)
								RETURN network.id as id,
												network.name as name`

	db, err := rt.DB().OpenPool()

	ctxLogger := rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(customerName),
		"cell_id":       swag.StringValue(CellID)})

	if err != nil {
		ctxLogger.Error("error connecting to neo4j:", err)
		return nil, err
	}
	defer db.Close()

	data, _, _, err := db.QueryNeoAll(cypher, map[string]interface{}{
		"name":    swag.StringValue(customerName),
		"cell_id": swag.StringValue(CellID)})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: %s", err)
		return nil, err

	} else if len(data) == 0 {
		return nil, nil
	}

	for _, row := range data {
		net_id := row[0].(string)
		net, _ := _getCellNetwork(rt, customerName, CellID, &net_id)

		res = append(res, net)
	}

	return res, nil
}

func _findHostgroupNetworks(rt *configManager.Runtime, customerName *string, hostgroupID *string) ([]*models.Network, error) {

	var res []*models.Network

	cypher := `MATCH (c:Customer {name: {name} })-[:OWN]->
										(Cell)-[:PROVIDES]->
										(Component)-[:DEPLOYED_ON]->
										(h:Hostgroup {id: {hostgroup_id}})-[:CONNECTED_ON]->(network:Network)
								RETURN network.id as id,
												network.name as name`

	db, err := rt.DB().OpenPool()

	ctxLogger := rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(customerName),
		"hostgroup_id":  swag.StringValue(hostgroupID)})

	if err != nil {
		ctxLogger.Error("error connecting to neo4j:", err)
		return nil, err
	}
	defer db.Close()

	data, _, _, err := db.QueryNeoAll(cypher, map[string]interface{}{
		"name":         swag.StringValue(customerName),
		"hostgroup_id": swag.StringValue(hostgroupID)})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: %s", err)
		return nil, err

	} else if len(data) == 0 {
		return nil, nil
	}

	for _, row := range data {
		net_id := row[0].(string)
		net, _ := _getNetworkByID(rt, &net_id)

		res = append(res, net)
	}

	return res, nil
}

func _getCellNetwork(rt *configManager.Runtime, customerName *string, CellID *string, NetworkID *string) (*models.Network, error) {
	var network *models.Network
	network = nil

	cypher := `MATCH (c:Customer {name: {name} })-[:OWN]->
										(cell:Cell {id: {cell_id}})-[:HAS]->
										(network:Network {id: {network_id}})
								RETURN network.id as id,
												network.name as name,
												network.cidr as cidr`

	db, err := rt.DB().OpenPool()
	ctxLogger := rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(customerName),
		"cell_id":       swag.StringValue(CellID)})

	if err != nil {
		ctxLogger.Error("error connecting to neo4j:", err)
		return network, err
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: %s", err)
		return network, err
	}

	defer stmt.Close()

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"name":       swag.StringValue(customerName),
		"cell_id":    swag.StringValue(CellID),
		"network_id": swag.StringValue(NetworkID)})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: %s", err)
		return network, err
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		return network, err
	}

	_name := output[1].(string)
	_cidr := output[2].(string)

	network = &models.Network{
		ID:   models.ULID(output[0].(string)),
		Name: &_name,
		Cidr: &_cidr}

	return network, nil
}

func _getNetworkByID(rt *configManager.Runtime, NetworkID *string) (*models.Network, error) {
	var network *models.Network
	network = nil

	cypher := `MATCH (network:Network {id: {network_id}})
								RETURN network.id as id,
												network.name as name,
												network.cidr as cidr`

	db, err := rt.DB().OpenPool()
	ctxLogger := rt.Logger().WithFields(logrus.Fields{
		"network_id": swag.StringValue(NetworkID)})

	if err != nil {
		ctxLogger.Error("error connecting to neo4j:", err)
		return network, err
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: %s", err)
		return network, err
	}

	defer stmt.Close()

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"network_id": swag.StringValue(NetworkID)})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: %s", err)
		return network, err
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		return network, err
	}

	_name := output[1].(string)
	_cidr := output[2].(string)

	network = &models.Network{
		ID:   models.ULID(output[0].(string)),
		Name: &_name,
		Cidr: &_cidr}

	return network, nil
}

func _getNetworkByName(rt *configManager.Runtime, customerName *string, CellID *string, networkName *string) *models.Network {

	var network *models.Network

	cypher := `MATCH (c:Customer {name: {name} })-[:OWN]->
										(cell:Cell {id: {cell_id}})-[:HAS]->
										(network:Network {name: {network_name}})
								RETURN network.id as id,
												network.name as name,
												network.cidr as cidr`

	db, err := rt.DB().OpenPool()
	ctxLogger := rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(customerName),
		"cell_id":       swag.StringValue(CellID),
		"network_name":  swag.StringValue(networkName)})

	if err != nil {
		ctxLogger.Error("error connecting to neo4j: ", err)
		return network
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: ", err)
		return network
	}
	defer stmt.Close()

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"name":         swag.StringValue(customerName),
		"cell_id":      swag.StringValue(CellID),
		"network_name": swag.StringValue(networkName)})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: %s", err)
		return network
	}

	output, _, err := rows.NextNeo()

	if err != nil {
		ctxLogger.Errorf("An error occurred querying Neo: %s", err)
		return network
	}

	_name := output[1].(string)
	_cidr := output[2].(string)

	network = &models.Network{
		ID:   models.ULID(output[0].(string)),
		Name: &_name,
		Cidr: &_cidr}

	return network
}

func _connectToAZ(rt *configManager.Runtime, networkID string, regionAZID string) error {

	cypher := `MATCH (network:Network {id: {network_id} })
							MATCH (az:RegionAZ {id: {region_az_id}})
							MERGE (network)-[:DEPLOYED_ON]->(az)`

	db, err := rt.DB().OpenPool()

	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"network_id":   networkID,
		"region_az_id": regionAZID})

	if err != nil {
		return err
	}

	_, _, err = rows.NextNeo()
	if err != nil {
		return err
	}

	return nil
}
