package handlers

import (
	"log"

	"configManager/models"
	"configManager/neo4j"
	"configManager/restapi/operations/hostgroup"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
)

func AddComponentHostgroup(params hostgroup.AddComponentHostgroupParams, principal *models.Customer) middleware.Responder {

	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->(cell:Cell)-[:PROVIDES]->(component:Component)
							WHERE id(cell) = {cell_id} AND id(component) = {component_id}
							CREATE (component)-[:DEPLOYED_ON]->(hostgroup:Hostgroup {
								name: {hostgroup_name},
								image: {hostgroup_image},
								flavor: {hostgroup_flavor},
								username: {hostgroup_username},
								bootstrap_command: {hostgroup_bootstrap_command},
								count: {hostgroup_count},
								network: {hostgroup_network},
								order: {hostgroup_order} } )
								RETURN id(hostgroup) as id`

	db, err := neo4j.Connect("")

	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return hostgroup.NewAddComponentHostgroupInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	if err != nil {
		log.Printf("An error occurred beginning transaction: %s", err)
		return hostgroup.NewAddComponentHostgroupInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return hostgroup.NewAddComponentHostgroupInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	if params.Body.Order == nil {
		params.Body.Order = swag.Int64(99)
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"customer_name":               swag.StringValue(principal.Name),
		"cell_id":                     params.CellID,
		"component_id":                params.ComponentID,
		"hostgroup_name":              swag.StringValue(params.Body.Name),
		"hostgroup_image":             swag.StringValue(params.Body.Image),
		"hostgroup_flavor":            swag.StringValue(params.Body.Flavor),
		"hostgroup_username":          swag.StringValue(params.Body.Username),
		"hostgroup_bootstrap_command": swag.StringValue(&params.Body.BootstrapCommand),
		"hostgroup_count":             swag.Int64Value(params.Body.Count),
		"hostgroup_network":           swag.StringValue(params.Body.Network),
		"hostgroup_order":             params.Body.Order})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return hostgroup.NewAddComponentHostgroupInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		log.Printf("An error occurred getting next row: %s", err)
		return hostgroup.NewAddComponentHostgroupInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	log.Printf("= Output(%#v)", output)

	log.Printf("name(%s)", swag.StringValue(params.Body.Name))

	stmt.Close()

	return hostgroup.NewAddComponentHostgroupCreated().WithPayload(output[0].(int64))
}

func DeleteComponentHostgroup(params hostgroup.DeleteComponentHostgroupParams, principal *models.Customer) middleware.Responder {

	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->
							(cell:Cell)-[:PROVIDES]->(component:Component)-[:DEPLOYED_ON]->(hostgroup:Hostgroup)
							WHERE id(cell) = {cell_id}
								AND id(component) = {component_id}
								AND id(hostgroup) = {hostgroup_id}
							DETACH DELETE hostgroup`

	if getComponentHostgroupByID(principal.Name, params.CellID, params.ComponentID, params.HostgroupID) == nil {
		log.Println("hostgroup does not exists !")
		return hostgroup.NewDeleteComponentHostgroupNotFound()
	}

	db, err := neo4j.Connect("")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return hostgroup.NewDeleteComponentHostgroupInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return hostgroup.NewDeleteComponentHostgroupInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	defer stmt.Close()

	_, err = stmt.ExecNeo(map[string]interface{}{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID,
		"component_id":  params.ComponentID,
		"hostgroup_id":  params.HostgroupID})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return hostgroup.NewAddComponentHostgroupInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	//if res.Count.(int) <= 0 {
	//	return hostgroup.NewDeleteComponentHostgroupNotFound()
	//}

	return hostgroup.NewDeleteComponentHostgroupOK()
}

func GetComponentHostgroupByID(params hostgroup.GetComponentHostgroupByIDParams, principal *models.Customer) middleware.Responder {

	Hostgroup := getComponentHostgroupByID(principal.Name, params.CellID, params.ComponentID, params.HostgroupID)
	if Hostgroup == nil {
		return hostgroup.NewGetComponentHostgroupByIDNotFound()
	}

	return hostgroup.NewGetComponentHostgroupByIDOK().WithPayload(Hostgroup)
}

func FindComponentHostgroups(params hostgroup.FindComponentHostgroupsParams, principal *models.Customer) middleware.Responder {

	data, err := _FindComponentHostgroups(principal.Name, params.CellID, params.ComponentID)

	log.Printf("= data(%#v)", data)

	if err != nil {
		return err
	}

	return hostgroup.NewFindComponentHostgroupsOK().WithPayload(data)
}

func _FindComponentHostgroups(customerName *string, CellID int64, ComponentID int64) ([]*models.Hostgroup, middleware.Responder) {

	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->
							(cell:Cell)-[:PROVIDES]->(component:Component)-[:DEPLOYED_ON]->(hostgroup:Hostgroup)
							WHERE id(cell) = {cell_id} AND id(component) = {component_id}
								RETURN id(hostgroup) as id,
												hostgroup.name as name,
												hostgroup.image as image,
												hostgroup.flavor as flavor,
												hostgroup.username as username,
												hostgroup.bootstrap_command as bootstrap_command,
												hostgroup.count as count,
												hostgroup.network as network,
												hostgroup.order as order`

	db, err := neo4j.Connect("")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return nil, hostgroup.NewFindComponentHostgroupsInternalServerError()
	}
	defer db.Close()

	log.Println(" customerName =>>>>>", *customerName)

	data, _, _, err := db.QueryNeoAll(cypher, map[string]interface{}{
		"customer_name": *customerName,
		"cell_id":       CellID,
		"component_id":  ComponentID})

	log.Printf("= data(%#v)", data)

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return nil, hostgroup.NewFindComponentHostgroupsInternalServerError()
	}

	res := make([]*models.Hostgroup, len(data))

	for idx, row := range data {

		var _order int64

		_name := row[1].(string)
		_image := row[2].(string)
		_flavor := row[3].(string)
		_username := row[4].(string)
		_bootstrap_command := ""
		_count := row[6].(int64)
		_network := row[7].(string)

		if row[8] == nil {
			_order = 99
		} else {
			_order = row[8].(int64)
		}

		if row[5] != nil {
			_bootstrap_command = row[5].(string)
		}

		res[idx] = &models.Hostgroup{
			ID:               row[0].(int64),
			Count:            &_count,
			Name:             &_name,
			Image:            &_image,
			Flavor:           &_flavor,
			Username:         &_username,
			BootstrapCommand: _bootstrap_command,
			Network:          &_network,
			Order:            &_order}
	}
	return res, nil
}

func UpdateComponentHostgroup(params hostgroup.UpdateComponentHostgroupParams, principal *models.Customer) middleware.Responder {

	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->
								(cell:Cell)-[:PROVIDES]->(component:Component)-[:DEPLOYED_ON]->(hostgroup:Hostgroup)
							WHERE id(cell) = {cell_id} AND id(component) = {component_id} AND id(hostgroup) = {hostgroup_id}
							SET hostgroup.name={hostgroup_name},
									hostgroup.image={hostgroup_image},
									hostgroup.flavor={hostgroup_flavor},
									hostgroup.username={hostgroup_username},
									hostgroup.bootstrap_command={hostgroup_bootstrap_command},
									hostgroup.count={hostgroup_count},
									hostgroup.network={hostgroup_network},
									hostgroup.order={hostgroup_order}`

	db, err := neo4j.Connect("")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return hostgroup.NewAddComponentHostgroupInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return hostgroup.NewAddComponentHostgroupInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}
	defer stmt.Close()

	log.Printf("customer_name(%s) cell_id(%s) component_id(%s) hostgroup_id(%s) hostgroup_name(%s) hostgroup_image(%s) hostgroup_flavor(%s)",
		swag.StringValue(principal.Name),
		params.CellID,
		params.ComponentID,
		params.HostgroupID,
		swag.StringValue(params.Body.Name),
		swag.StringValue(params.Body.Image),
		swag.StringValue(params.Body.Flavor))

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
		"hostgroup_network":           swag.StringValue(params.Body.Network),
		"hostgroup_order":             swag.Int64Value(params.Body.Order)})

	if err != nil {
		log.Printf("-> An error occurred querying Neo: %s", err)
		return hostgroup.NewAddComponentHostgroupInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	return hostgroup.NewUpdateComponentHostgroupOK()
}

func getComponentHostgroupByID(customer *string, cellID int64, componentID int64, hostgroupID int64) *models.Hostgroup {

	var hostgroup *models.Hostgroup
	hostgroup = nil

	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->
							(cell:Cell)-[:PROVIDES]->(component:Component)-[:DEPLOYED_ON]->(hostgroup:Hostgroup)
							WHERE id(cell) = {cell_id}
								AND id(component) = {component_id}
								AND id(hostgroup) = {hostgroup_id}
							RETURN id(hostgroup) as id,
											hostgroup.name as name,
											hostgroup.image as image,
											hostgroup.flavor as flavor,
											hostgroup.username as username,
											hostgroup.bootstrap_command as bootstrap_command,
											hostgroup.count as count,
											hostgroup.network as network,
											hostgroup.order as order`

	db, err := neo4j.Connect("")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return hostgroup
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return hostgroup
	}
	defer stmt.Close()

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"customer_name": swag.StringValue(customer),
		"cell_id":       cellID,
		"component_id":  componentID,
		"hostgroup_id":  hostgroupID})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
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
	_network := output[7].(string)

	var _order int64

	if output[8].(int64) <= 0 {
		_order = 99
	} else {
		_order = output[8].(int64)
	}

	if output[5] != nil {
		_bootstrap_command = output[5].(string)
	}

	hostgroup = &models.Hostgroup{
		ID:               output[0].(int64),
		Count:            &_count,
		Name:             &_name,
		Image:            &_image,
		Flavor:           &_flavor,
		Username:         &_username,
		BootstrapCommand: _bootstrap_command,
		Network:          &_network,
		Order:            &_order}

	log.Printf("here => (%#v)", hostgroup)

	return hostgroup
}
