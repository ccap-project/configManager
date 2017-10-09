package handlers

import (
	"log"

	"../models"
	"../restapi/operations/host"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	driver "github.com/johnnadratowski/golang-neo4j-bolt-driver"
)

func AddCellHost(params host.AddCellHostParams, principal *models.Customer) middleware.Responder {

	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->(cell:Cell)
							WHERE id(cell) = {cell_id}
							CREATE (cell)-[:HAS]->(host:Host {name: {host_name}} )
								RETURN id(host) as id`

	if getCellHostByName(principal.Name, params.CellID, params.Body.Name) != nil {
		log.Println("host already exists !")
		return host.NewAddCellHostConflict().WithPayload(models.APIResponse{Message: "host already exists"})
	}

	db, err := driver.NewDriver().OpenNeo("bolt://192.168.20.54:7687")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return host.NewAddCellHostInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Printf("An error occurred beginning transaction: %s", err)
		return host.NewAddCellHostInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}
	defer tx.Rollback()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return host.NewAddCellHostInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID,
		"host_name":     swag.StringValue(params.Body.Name)})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return host.NewAddCellHostInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		log.Printf("An error occurred getting next row: %s", err)
		return host.NewAddCellHostInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	stmt.Close()

	err = addCellHostOptions(principal.Name, params.CellID, params.Body.Name, params.Body.Options, db)
	if err != nil {
		log.Printf("An error occurred adding Host options: %s", err)
		return host.NewAddCellHostInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}
	tx.Commit()

	return host.NewAddCellHostCreated().WithPayload(output[0].(int64))
}

func FindCellHosts(params host.FindCellHostsParams, principal *models.Customer) middleware.Responder {

	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->
							(cell:Cell)-[:HAS]->(host:Host)
							WHERE id(cell) = {cell_id}
								RETURN id(host) as id,
												host.name as name`

	db, err := driver.NewDriver().OpenNeo("bolt://192.168.20.54:7687")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return host.NewFindCellHostsInternalServerError()
	}
	defer db.Close()

	data, _, _, err := db.QueryNeoAll(cypher, map[string]interface{}{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID})

	log.Printf("= data(%#v)", data)

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return host.NewFindCellHostsInternalServerError()

	}

	res := make([]*models.Host, len(data))

	for idx, row := range data {
		_name := row[1].(string)

		res[idx] = &models.Host{
			ID:   row[0].(int64),
			Name: &_name}
	}

	log.Printf("= Res(%#v)", res)

	return host.NewFindCellHostsOK().WithPayload(res)
}

func addCellHostOptions(customer *string, cellID int64, hostName *string, options []*models.Parameter, db driver.Conn) error {

	log.Printf("= Output(%#v)", options)

	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->
							(cell:Cell)-[:HAS]->(host:Host {name: {host_name}})
							WHERE id(cell) = {cell_id}
							CREATE (host)-[:OPT]->(option:Option {name: {opt_name}, value: {opt_val}} )
								RETURN id(option) as id`

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return err
	}

	// add parameters
	for _, option := range options {
		log.Printf("name(%s) val(%s)", option.Name, option.Value)
		//log.Printf("param(%#v)", param)
		_, err := stmt.ExecNeo(map[string]interface{}{
			"customer_name": swag.StringValue(customer),
			"cell_id":       cellID,
			"host_name":     swag.StringValue(hostName),
			"opt_name":      swag.StringValue(option.Name),
			"opt_val":       swag.StringValue(option.Value)})

		if err != nil {
			log.Printf("An error occurred querying Neo: %s", err)
			return err
		}
	}

	return nil
}

func getCellHostByName(customer *string, cellID int64, hostName *string) *models.Host {

	var host *models.Host
	host = nil

	cypher := `MATCH (customer:Customer {name: {customer_name} })-[:OWN]->
							(cell:Cell)-[:HAS]->(host:Host {name: {host_name}})
							WHERE id(cell) = {cell_id}
								RETURN ID(host) as id,
												host.name as name`

	db, err := driver.NewDriver().OpenNeo("bolt://192.168.20.54:7687")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return host
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return host
	}
	defer stmt.Close()

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"customer_name": swag.StringValue(customer),
		"cell_id":       cellID,
		"host_name":     swag.StringValue(hostName)})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return host
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		//log.Printf("An error occurred fetching row: %s", err)
		return host
	}
	_name := output[1].(string)

	host = &models.Host{
		ID:   output[0].(int64),
		Name: &_name}

	log.Printf("\n here => (%#v)\n", host)

	return host
}
