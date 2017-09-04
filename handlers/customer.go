package handlers

import (
	"log"
	"strings"

	"../models"
	"../restapi/operations/customer"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	_ "gopkg.in/cq.v1"
)

func AddCustomer(params customer.AddCustomerParams) middleware.Responder {

	neo4jURL := `http://192.168.20.54:7474`

	cypher := `create(c:Customer { name: {0} }) RETURN ID(c)`

	if len(swag.StringValue(getCustomerByName(swag.StringValue(params.Body.Name)).Name)) > 0 {
		log.Println("customer already exists !")
		return customer.NewAddCustomerConflict().WithPayload(models.APIResponse{Message: "customer already exists"})
	}

	db, err := sqlx.Connect("neo4j-cypher", neo4jURL)
	if err != nil {
		log.Println("error connecting to neo4j:", err)
	}
	defer db.Close()

	stmt, err := db.Prepare(cypher)

	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(params.Body.Name)

	log.Printf("%#v", res)
	if err != nil {
		log.Printf("error creating customer name(%s): %v", params.Body.Name, err)
		return customer.NewAddCustomerInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	customerAdded := getCustomerByName(swag.StringValue(params.Body.Name))

	return customer.NewAddCustomerCreated().WithPayload(customerAdded.ID)
}

/*
func GetCustomerByID(params providertype.GetProviderTypeByIDParams) middleware.Responder {

	neo4jURL := `http://192.168.20.54:7474`

	cypher := `MATCH (p:ProviderType)
							WHERE ID(p) = {0}
							RETURN ID(p) as id,
											p.name as name,
											p.auth_url as auth_url,
											p.domain_name as domain_name,
											p.username as username,
											p.password as password`

	db, err := sqlx.Connect("neo4j-cypher", neo4jURL)
	if err != nil {
		log.Println("error connecting to neo4j:", err)
	}
	defer db.Close()
	db.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)

	providerType := models.ProviderType{}

	err = db.Get(&providerType, cypher, params.ProvidertypeID)

	if err != nil {
		log.Printf("error getting providertype id(%d): %v", params.ProvidertypeID, err)
		return providertype.NewGetProviderTypeByIDNotFound()
	}

	return providertype.NewGetProviderTypeByIDOK().WithPayload(&providerType)
}
*/
func GetCustomerByName(params customer.FindCustomerByNameParams) middleware.Responder {

	Customer := getCustomerByName(params.CustomerName)

	if len(swag.StringValue(Customer.Name)) <= 0 {
		return customer.NewFindCustomerByNameNotFound()
	}

	return customer.NewFindCustomerByNameOK().WithPayload(Customer)
}

func getCustomerByName(customerName string) *models.Customer {

	neo4jURL := `http://192.168.20.54:7474`

	cypher := `MATCH (c:Customer)
							WHERE c.name =~ {0}
							RETURN ID(c) as id,
											c.name as name`

	db, err := sqlx.Connect("neo4j-cypher", neo4jURL)
	db.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)
	if err != nil {
		log.Println("error connecting to neo4j:", err)
	}
	defer db.Close()

	customer := models.Customer{}

	err = db.Get(&customer, cypher, customerName)
	if err != nil &&
		!(err.Error() == "sql: Scan error on column index 1: unsupported Scan, storing driver.Value type <nil> into type *string" ||
			err.Error() == "sql: no rows in result set") {
		log.Printf("error getting customer by name(%s): %v", customerName, err)
		log.Printf("customer(%#v)", customer)
		return (&models.Customer{})
	}

	log.Printf("customer(%#v)", customer)
	return &customer
}
