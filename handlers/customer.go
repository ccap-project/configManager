package handlers

import (
	"log"

	"configManager/models"
	"configManager/neo4j"
	"configManager/restapi/operations/customer"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
)

func AddCustomer(params customer.AddCustomerParams) middleware.Responder {

	cypher := `create(c:Customer { name: {name} }) RETURN ID(c)`

	//if len(swag.StringValue(getCustomerByName(swag.StringValue(params.Body.Name)).Name)) > 0 {

	if getCustomerByName(params.Body.Name) != nil {
		log.Println("customer already exists !")
		return customer.NewAddCustomerConflict().WithPayload(models.APIResponse{Message: "customer already exists"})
	}

	db, err := neo4j.Connect("")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return customer.NewAddCustomerInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return customer.NewAddCustomerInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"name": swag.StringValue(params.Body.Name)})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return customer.NewAddCustomerInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		log.Printf("An error occurred getting next row: %s", err)
		return customer.NewAddCustomerInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	log.Printf("= Output(%#v)", output)

	return customer.NewAddCustomerCreated().WithPayload(output[0].(int64))
}

func GetCustomerByName(params customer.FindCustomerByNameParams) middleware.Responder {

	Customer := getCustomerByName(&params.CustomerName)

	if len(swag.StringValue(Customer.Name)) <= 0 {
		return customer.NewFindCustomerByNameNotFound()
	}

	return customer.NewFindCustomerByNameOK().WithPayload(Customer)
}

func getCustomerByName(customerName *string) *models.Customer {

	var customer *models.Customer

	cypher := `MATCH (c:Customer)
							WHERE c.name =~ {customer_name}
							RETURN ID(c) as id,
											c.name as name`

	db, err := neo4j.Connect("")

	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return customer
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return customer
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"customer_name": swag.StringValue(customerName)})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return customer
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		return customer
	}
	_name := output[1].(string)

	customer = &models.Customer{
		ID:   output[0].(int64),
		Name: &_name}

	stmt.Close()

	return customer
}
