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
	"io"

	"configManager"
	"configManager/models"
	"configManager/neo4j"
	"configManager/restapi/operations/hostgroup"
	"configManager/util"

	"github.com/Sirupsen/logrus"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
)

func NewAddComponentHostgroup(rt *configManager.Runtime) hostgroup.AddComponentHostgroupHandler {
	return &addComponentHostgroup{rt: rt}
}

type addComponentHostgroup struct {
	rt *configManager.Runtime
}

func (ctx *addComponentHostgroup) Handle(params hostgroup.AddComponentHostgroupParams, principal *models.Customer) middleware.Responder {

	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->
							(cell:Cell {id: {cell_id}})-[:PROVIDES]->
							(component:Component {id: {component_id}})
						MERGE (component)-[:DEPLOYED_ON]->(hostgroup:Hostgroup {
							id: {hostgroup_id},
							%s } )
							RETURN hostgroup.id as id`

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID})

	if _getComponentHostgroupByName(ctx.rt, principal.Name, &params.CellID, &params.ComponentID, params.Body.Name) != nil {
		ctxLogger.Error("hostgroup already exists !")
		return hostgroup.NewAddComponentHostgroupConflict().WithPayload(&models.APIResponse{Message: "hostgroup already exists"})
	}

	// Check if required networks exists
	var cellNetworks []*models.Network
	for _, n := range params.Body.Network {
		_net := _getNetworkByName(ctx.rt, principal.Name, &params.CellID, &n)
		if _net == nil {
			ctxLogger.Errorf("network not found (%s)", n)
			return hostgroup.NewAddComponentHostgroupInternalServerError().WithPayload(&models.APIResponse{Message: "network not found"})
		}
		cellNetworks = append(cellNetworks, _net)
	}

	db, err := ctx.rt.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j: ", err)
		return hostgroup.NewAddComponentHostgroupInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		ctxLogger.Error("An error occurred beginning transaction: ", err)
		return hostgroup.NewAddComponentHostgroupInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer tx.Rollback()

	_params := util.BuildQuery(&params.Body, "", "merge", []string{"ID", "Network"})

	stmt, err := db.PrepareNeo(fmt.Sprintf(cypher, _params))
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: ", err)
		return hostgroup.NewAddComponentHostgroupInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer stmt.Close()

	if params.Body.Order == nil {
		params.Body.Order = new(int64)
		*params.Body.Order = 99
	}

	ulid := configManager.GetULID()

	ctxLogger = ctxLogger.WithFields(logrus.Fields{
		"hostgroup_id":   ulid,
		"hostgroup_name": swag.StringValue(params.Body.Name)})

	rows, err := stmt.QueryNeo(util.BuildParams(params.Body, "",
		map[string]interface{}{
			"customer_name": swag.StringValue(principal.Name),
			"cell_id":       params.CellID,
			"component_id":  params.ComponentID,
			"hostgroup_id":  ulid}, []string{"ID", "Network"}))

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: ", err)
		return hostgroup.NewAddComponentHostgroupInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	output, _, err := rows.NextNeo()

	if err != nil {
		ctxLogger.Error("An error occurred getting next row: ", err)
		return hostgroup.NewAddComponentHostgroupInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	stmt.Close()

	for _, net := range cellNetworks {
		_connectToNetwork(db, ulid, string(net.ID))
	}

	tx.Commit()

	return hostgroup.NewAddComponentHostgroupCreated().WithPayload(models.ULID(output[0].(string)))
}

func NewConnectHostgroupToNetwork(rt *configManager.Runtime) hostgroup.ConnectHostgroupToNetworkHandler {
	return &connectHostgroupToNetwork{rt: rt}
}

type connectHostgroupToNetwork struct {
	rt *configManager.Runtime
}

func (ctx *connectHostgroupToNetwork) Handle(params hostgroup.ConnectHostgroupToNetworkParams, principal *models.Customer) middleware.Responder {

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID,
		"hostgroup_id":  params.HostgroupID})

	if _getComponentHostgroupByID(ctx.rt, principal.Name, &params.CellID, &params.ComponentID, &params.HostgroupID) == nil {
		ctxLogger.Error("hostgroup does not exists !")
		return hostgroup.NewConnectHostgroupToNetworkNotFound()
	}
	network, err := _getCellNetwork(ctx.rt, principal.Name, &params.CellID, &params.NetworkID)

	if err != nil && err != io.EOF {
		ctxLogger.Error("Getting network, ", err)
		return hostgroup.NewConnectHostgroupToNetworkInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	if network == nil {
		ctxLogger.Error("network does not exists !")
		return hostgroup.NewConnectHostgroupToNetworkNotFound()
	}

	db, err := ctx.rt.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j: ", err)
		return hostgroup.NewAddComponentHostgroupInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	if err := _connectToNetwork(db, params.HostgroupID, params.NetworkID); err != nil {
		return hostgroup.NewConnectHostgroupToNetworkInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	return hostgroup.NewConnectHostgroupToNetworkOK()
}

func NewDeleteComponentHostgroup(rt *configManager.Runtime) hostgroup.DeleteComponentHostgroupHandler {
	return &deleteComponentHostgroup{rt: rt}
}

type deleteComponentHostgroup struct {
	rt *configManager.Runtime
}

func (ctx *deleteComponentHostgroup) Handle(params hostgroup.DeleteComponentHostgroupParams, principal *models.Customer) middleware.Responder {

	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->
							(cell:Cell {id: {cell_id}})-[:PROVIDES]->
							(component:Component {id: {component_id}})-[:DEPLOYED_ON]->
							(hostgroup:Hostgroup {id: {hostgroup_id}})
						DETACH DELETE hostgroup`

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID,
		"hostgroup_id":  params.HostgroupID})

	if _getComponentHostgroupByID(ctx.rt, principal.Name, &params.CellID, &params.ComponentID, &params.HostgroupID) == nil {
		ctxLogger.Error("hostgroup does not exists !")
		return hostgroup.NewDeleteComponentHostgroupNotFound()
	}

	db, err := ctx.rt.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j: ", err)
		return hostgroup.NewDeleteComponentHostgroupInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: ", err)
		return hostgroup.NewDeleteComponentHostgroupInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	defer stmt.Close()

	_, err = stmt.ExecNeo(map[string]interface{}{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID,
		"component_id":  params.ComponentID,
		"hostgroup_id":  params.HostgroupID})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: ", err)
		return hostgroup.NewAddComponentHostgroupInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	return hostgroup.NewDeleteComponentHostgroupOK()
}

func NewDisconnectHostgroupFromNetwork(rt *configManager.Runtime) hostgroup.DisconnectHostgroupFromNetworkHandler {
	return &disconnectHostgroupFromNetwork{rt: rt}
}

type disconnectHostgroupFromNetwork struct {
	rt *configManager.Runtime
}

func (ctx *disconnectHostgroupFromNetwork) Handle(params hostgroup.DisconnectHostgroupFromNetworkParams, principal *models.Customer) middleware.Responder {

	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->
							(cell:Cell {id: {cell_id}})-[:PROVIDES]->
							(component:Component {id: {component_id}})-[:DEPLOYED_ON]->
							(hostgroup:Hostgroup {id: {hostgroup_id}})-[n_c:CONNECTED_ON]->
							(network:Network {id: {network_id}})
						DELETE n_c`

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID,
		"hostgroup_id":  params.HostgroupID})

	if _getComponentHostgroupByID(ctx.rt, principal.Name, &params.CellID, &params.ComponentID, &params.HostgroupID) == nil {
		ctxLogger.Error("hostgroup does not exists !")
		return hostgroup.NewDisconnectHostgroupFromNetworkNotFound()
	}

	db, err := ctx.rt.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j: ", err)
		return hostgroup.NewDisconnectHostgroupFromNetworkInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: ", err)
		return hostgroup.NewDisconnectHostgroupFromNetworkInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	defer stmt.Close()

	_, err = stmt.ExecNeo(map[string]interface{}{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID,
		"component_id":  params.ComponentID,
		"hostgroup_id":  params.HostgroupID,
		"network_id":    params.NetworkID})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: ", err)
		return hostgroup.NewDisconnectHostgroupFromNetworkInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	return hostgroup.NewDisconnectHostgroupFromNetworkOK()
}

func NewGetComponentHostgroupByID(rt *configManager.Runtime) hostgroup.GetComponentHostgroupByIDHandler {
	return &getComponentHostgroupByID{rt: rt}
}

type getComponentHostgroupByID struct {
	rt *configManager.Runtime
}

func (ctx *getComponentHostgroupByID) Handle(params hostgroup.GetComponentHostgroupByIDParams, principal *models.Customer) middleware.Responder {

	Hostgroup := _getComponentHostgroupByID(ctx.rt, principal.Name, &params.CellID, &params.ComponentID, &params.HostgroupID)
	if Hostgroup == nil {
		return hostgroup.NewGetComponentHostgroupByIDNotFound()
	}

	return hostgroup.NewGetComponentHostgroupByIDOK().WithPayload(Hostgroup)
}

func NewFindComponentHostgroups(rt *configManager.Runtime) hostgroup.FindComponentHostgroupsHandler {
	return &findComponentHostgroups{rt: rt}
}

type findComponentHostgroups struct {
	rt *configManager.Runtime
}

func (ctx *findComponentHostgroups) Handle(params hostgroup.FindComponentHostgroupsParams, principal *models.Customer) middleware.Responder {

	data, err := _findComponentHostgroups(ctx.rt, principal.Name, &params.CellID, &params.ComponentID)

	if err != nil {
		return err
	}

	return hostgroup.NewFindComponentHostgroupsOK().WithPayload(data)
}

func _findComponentHostgroups(rt *configManager.Runtime, customerName *string, cellID *string, componentID *string) ([]*models.Hostgroup, middleware.Responder) {

	var res []*models.Hostgroup
	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->
							(cell:Cell{id: {cell_id}})-[:PROVIDES]->
							(component:Component {id: {component_id}})-[:DEPLOYED_ON]->
							(hostgroup:Hostgroup)
						RETURN hostgroup.id`

	ctxLogger := rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(customerName),
		"cell_id":       swag.StringValue(cellID),
		"component_id":  swag.StringValue(componentID)})

	db, err := rt.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j: ", err)
		return nil, hostgroup.NewFindComponentHostgroupsInternalServerError()
	}
	defer db.Close()

	data, _, _, err := db.QueryNeoAll(cypher, map[string]interface{}{
		"customer_name": swag.StringValue(customerName),
		"cell_id":       swag.StringValue(cellID),
		"component_id":  swag.StringValue(componentID)})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: ", err)
		return nil, hostgroup.NewFindComponentHostgroupsInternalServerError()
	}

	for _, row := range data {

		hostgroupID := row[0].(string)

		h := _getComponentHostgroupByID(rt, customerName, cellID, componentID, &hostgroupID)

		res = append(res, h)
	}
	return res, nil
}

func NewUpdateComponentHostgroup(rt *configManager.Runtime) hostgroup.UpdateComponentHostgroupHandler {
	return &updateComponentHostgroup{rt: rt}
}

type updateComponentHostgroup struct {
	rt *configManager.Runtime
}

func (ctx *updateComponentHostgroup) Handle(params hostgroup.UpdateComponentHostgroupParams, principal *models.Customer) middleware.Responder {

	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->
								(cell:Cell {id: {cell_id}})-[:PROVIDES]->
								(component:Component {id: {component_id}})-[:DEPLOYED_ON]->
								(hostgroup:Hostgroup {id: {hostgroup_id}})
							SET hostgroup.name={hostgroup_name},
									hostgroup.image={hostgroup_image},
									hostgroup.flavor={hostgroup_flavor},
									hostgroup.username={hostgroup_username},
									hostgroup.bootstrap_command={hostgroup_bootstrap_command},
									hostgroup.count={hostgroup_count},
									hostgroup.order={hostgroup_order}`

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID,
		"hostgroup_id":  params.HostgroupID,
		"component_id":  params.ComponentID})

	// Check if required networks exists
	var cellNetworks []*models.Network
	for _, n := range params.Body.Network {
		_net := _getNetworkByName(ctx.rt, principal.Name, &params.CellID, &n)
		if _net == nil {
			ctxLogger.Errorf("network not found (%s)", n)
			return hostgroup.NewAddComponentHostgroupInternalServerError().WithPayload(&models.APIResponse{Message: "network not found"})
		}
		cellNetworks = append(cellNetworks, _net)
	}

	db, err := ctx.rt.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j: ", err)
		return hostgroup.NewAddComponentHostgroupInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		ctxLogger.Error("An error occurred beginning transaction: ", err)
		return hostgroup.NewAddComponentHostgroupInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer tx.Rollback()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: ", err)
		return hostgroup.NewAddComponentHostgroupInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer stmt.Close()

	if params.Body.Order == nil {
		params.Body.Order = swag.Int64(99)
	}

	_, err = stmt.ExecNeo(map[string]interface{}{
		"customer_name":               swag.StringValue(principal.Name),
		"cell_id":                     params.CellID,
		"component_id":                params.ComponentID,
		"hostgroup_id":                params.HostgroupID,
		"hostgroup_name":              swag.StringValue(params.Body.Name),
		"hostgroup_image":             swag.StringValue(params.Body.Image),
		"hostgroup_flavor":            swag.StringValue(params.Body.Flavor),
		"hostgroup_username":          swag.StringValue(params.Body.Username),
		"hostgroup_bootstrap_command": swag.StringValue(&params.Body.BootstrapCommand),
		"hostgroup_desired_size":      swag.Int64Value(params.Body.DesiredSize),
		"hostgroup_order":             swag.Int64Value(params.Body.Order)})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: ", err)
		return hostgroup.NewAddComponentHostgroupInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	stmt.Close()

	for _, net := range cellNetworks {
		ctxLogger.Infof("Connecting (%s) to (%s)", params.HostgroupID, string(net.ID))
		if err := _connectToNetwork(db, params.HostgroupID, string(net.ID)); err != nil {
			tx.Rollback()
			return hostgroup.NewAddComponentHostgroupInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
		}
	}

	tx.Commit()

	return hostgroup.NewUpdateComponentHostgroupOK()
}

func _connectToNetwork(conn neo4j.Conn, hostgroupID string, networkID string) error {

	cypher := `MATCH (hg:Hostgroup {id: {hostgroup_id} })
							MATCH (network:Network {id: {network_id}})
							MERGE (hg)-[:CONNECTED_ON]->(network)
							RETURN hg`

	stmt, err := conn.PrepareNeo(cypher)
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"network_id":   networkID,
		"hostgroup_id": hostgroupID})

	if err != nil {
		return err
	}

	_, _, err = rows.NextNeo()
	if err != nil {
		return err
	}

	return nil
}

func _getComponentHostgroupByID(rt *configManager.Runtime, customer *string, cellID *string, componentID *string, hostgroupID *string) *models.Hostgroup {

	var hostgroup *models.Hostgroup
	var _nets []string

	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->
							(cell:Cell{id: {cell_id}})-[:PROVIDES]->
							(component:Component {id: {component_id}})-[:DEPLOYED_ON]->
							(hostgroup:Hostgroup {id: {hostgroup_id}})
							RETURN hostgroup {.*}`

	ctxLogger := rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(customer),
		"cell_id":       swag.StringValue(cellID),
		"component_id":  swag.StringValue(componentID),
		"hostgroup_id":  swag.StringValue(hostgroupID)})

	db, err := rt.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j: ", err)
		return hostgroup
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: ", err)
		return hostgroup
	}
	defer stmt.Close()

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"customer_name": swag.StringValue(customer),
		"cell_id":       swag.StringValue(cellID),
		"component_id":  swag.StringValue(componentID),
		"hostgroup_id":  swag.StringValue(hostgroupID)})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: ", err)
		return hostgroup
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		return hostgroup
	}

	hostgroup = &models.Hostgroup{}

	util.FillStruct(hostgroup, output[0].(map[string]interface{}))

	_networks, _ := _findHostgroupNetworks(rt, customer, hostgroupID)

	for _, n := range _networks {
		_nets = append(_nets, *n.Name)
	}
	hostgroup.Network = _nets

	return hostgroup
}

func _getComponentHostgroupByName(rt *configManager.Runtime, customer *string, cellID *string, componentID *string, hostgroupName *string) *models.Hostgroup {

	var hostgroup *models.Hostgroup
	hostgroup = nil

	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->
							(cell:Cell{id: {cell_id}})-[:PROVIDES]->
							(component:Component {id: {component_id}})-[:DEPLOYED_ON]->
							(hostgroup:Hostgroup {name: {hostgroup_name}})
							RETURN hostgroup {.*}`

	ctxLogger := rt.Logger().WithFields(logrus.Fields{
		"customer_name":  swag.StringValue(customer),
		"cell_id":        swag.StringValue(cellID),
		"component_id":   swag.StringValue(componentID),
		"hostgroup_name": swag.StringValue(hostgroupName)})

	db, err := rt.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j: ", err)
		return hostgroup
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: ", err)
		return hostgroup
	}
	defer stmt.Close()

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"customer_name":  swag.StringValue(customer),
		"cell_id":        swag.StringValue(cellID),
		"component_id":   swag.StringValue(componentID),
		"hostgroup_name": swag.StringValue(hostgroupName)})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: ", err)
		return hostgroup
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		return hostgroup
	}

	hostgroup = &models.Hostgroup{}

	util.FillStruct(hostgroup, output[0].(map[string]interface{}))

	return hostgroup
}
