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
	"configManager/restapi/operations/hostgroup"

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
							name: {hostgroup_name},
							image: {hostgroup_image},
							flavor: {hostgroup_flavor},
							username: {hostgroup_username},
							bootstrap_command: {hostgroup_bootstrap_command},
							count: {hostgroup_count},
							order: {hostgroup_order} } )
							RETURN hostgroup.id as id`

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID})

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
		params.Body.Order = new(int64)
		*params.Body.Order = 99
	}

	ulid := configManager.GetULID()

	ctxLogger = ctxLogger.WithFields(logrus.Fields{
		"hostgroup_id":   ulid,
		"hostgroup_name": swag.StringValue(params.Body.Name)})

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"customer_name":               swag.StringValue(principal.Name),
		"cell_id":                     params.CellID,
		"component_id":                params.ComponentID,
		"hostgroup_id":                ulid,
		"hostgroup_name":              swag.StringValue(params.Body.Name),
		"hostgroup_image":             swag.StringValue(params.Body.Image),
		"hostgroup_flavor":            swag.StringValue(params.Body.Flavor),
		"hostgroup_username":          swag.StringValue(params.Body.Username),
		"hostgroup_bootstrap_command": swag.StringValue(&params.Body.BootstrapCommand),
		"hostgroup_count":             swag.Int64Value(params.Body.Count),
		"hostgroup_order":             swag.Int64Value(params.Body.Order)})

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
		_connectToNetwork(db, tx, ulid, string(net.ID))
	}

	tx.Commit()

	return hostgroup.NewAddComponentHostgroupCreated().WithPayload(models.ULID(output[0].(string)))
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

func _findComponentHostgroups(rt *configManager.Runtime, customerName *string, CellID *string, ComponentID *string) ([]*models.Hostgroup, middleware.Responder) {

	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->
							(cell:Cell{id: {cell_id}})-[:PROVIDES]->
							(component:Component {id: {component_id}})-[:DEPLOYED_ON]->
							(hostgroup:Hostgroup)
						RETURN hostgroup.id as id,
										hostgroup.name as name,
										hostgroup.image as image,
										hostgroup.flavor as flavor,
										hostgroup.username as username,
										hostgroup.bootstrap_command as bootstrap_command,
										hostgroup.count as count,
										hostgroup.order as order`

	ctxLogger := rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(customerName),
		"cell_id":       swag.StringValue(CellID),
		"component_id":  swag.StringValue(ComponentID)})

	db, err := rt.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j: ", err)
		return nil, hostgroup.NewFindComponentHostgroupsInternalServerError()
	}
	defer db.Close()

	data, _, _, err := db.QueryNeoAll(cypher, map[string]interface{}{
		"customer_name": swag.StringValue(customerName),
		"cell_id":       swag.StringValue(CellID),
		"component_id":  swag.StringValue(ComponentID)})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: ", err)
		return nil, hostgroup.NewFindComponentHostgroupsInternalServerError()
	}

	res := make([]*models.Hostgroup, len(data))

	for idx, row := range data {

		var _order int64
		var _nets []string

		_id := row[0].(string)
		_name := row[1].(string)
		_image := row[2].(string)
		_flavor := row[3].(string)
		_username := row[4].(string)
		_bootstrap_command := ""
		_count := row[6].(int64)

		if row[7] == nil {
			_order = 99
		} else {
			_order = row[7].(int64)
		}

		if row[5] != nil {
			_bootstrap_command = row[5].(string)
		}

		_networks, _ := _findHostgroupNetworks(rt, customerName, &_id)

		for _, n := range _networks {
			_nets = append(_nets, *n.Name)
		}

		res[idx] = &models.Hostgroup{
			ID:               models.ULID(_id),
			Count:            &_count,
			Name:             &_name,
			Image:            &_image,
			Flavor:           &_flavor,
			Username:         &_username,
			BootstrapCommand: _bootstrap_command,
			Network:          _nets,
			Order:            &_order}
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
		"hostgroup_count":             swag.Int64Value(params.Body.Count),
		"hostgroup_order":             swag.Int64Value(params.Body.Order)})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: ", err)
		return hostgroup.NewAddComponentHostgroupInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	stmt.Close()

	for _, net := range cellNetworks {
		ctxLogger.Infof("Connecting (%s) to (%s)", params.HostgroupID, string(net.ID))
		if err := _connectToNetwork(db, tx, params.HostgroupID, string(net.ID)); err != nil {
			tx.Rollback()
			return hostgroup.NewAddComponentHostgroupInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
		}
	}

	tx.Commit()

	return hostgroup.NewUpdateComponentHostgroupOK()
}

func _connectToNetwork(conn neo4j.Conn, tx neo4j.Tx, hostgroupID string, networkID string) error {

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
	hostgroup = nil

	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->
							(cell:Cell{id: {cell_id}})-[:PROVIDES]->
							(component:Component {id: {component_id}})-[:DEPLOYED_ON]->
							(hostgroup:Hostgroup {id: {hostgroup_id}})
							RETURN hostgroup.id as id,
											hostgroup.name as name,
											hostgroup.image as image,
											hostgroup.flavor as flavor,
											hostgroup.username as username,
											hostgroup.bootstrap_command as bootstrap_command,
											hostgroup.count as count,
											hostgroup.order as order`

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

	_name := output[1].(string)
	_image := output[2].(string)
	_flavor := output[3].(string)
	_username := output[4].(string)
	_bootstrap_command := ""
	_count := output[6].(int64)
	//_network := output[7].(string)

	var _order int64

	if output[7].(int64) <= 0 {
		_order = 99
	} else {
		_order = output[7].(int64)
	}

	if output[5] != nil {
		_bootstrap_command = output[5].(string)
	}

	hostgroup = &models.Hostgroup{
		ID:               models.ULID(output[0].(string)),
		Count:            &_count,
		Name:             &_name,
		Image:            &_image,
		Flavor:           &_flavor,
		Username:         &_username,
		BootstrapCommand: _bootstrap_command,
		//Network:          &_network,
		Order: &_order}

	log.Printf("here => (%#v)", hostgroup)

	return hostgroup
}
