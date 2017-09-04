package handlers

/*
func EnableConstraints() error {

	constraints := []string{
		"CREATE INDEX ON :Cell(name)",
		"CREATE INDEX ON :Customer(name)",
		"CREATE INDEX ON :KeyPair(name)",
		"CREATE INDEX ON :Provider(name)",
		"CREATE INDEX ON :ProviderType(name)",
		"CREATE INDEX ON :Role(name)",
		"CREATE INDEX ON :RoleParameter(name)",
		"CREATE CONSTRAINT ON (customer:Customer) ASSERT customer.name IS UNIQUE",
	}

	db, err := driver.NewDriver().OpenNeo("bolt://192.168.20.54:7687")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return err
	}
	defer db.Close()

	for cmd range constraints
	err = db.ExecNeo(cypher, params)
	rows, err := stmt.QueryNeo(map[string]interface{}{"name": swag.StringValue(principal.Name),
		"cell_name": swag.StringValue(params.Body.Name)})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return err
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		log.Printf("An error occurred getting next row: %s", err)
		return cell.NewAddCellInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	log.Printf("= Output(%#v)", output)

	log.Printf("customer(%s) name(%s)", swag.StringValue(principal.Name), swag.StringValue(params.Body.Name))

	return cell.NewAddCellCreated().WithPayload(output[0].(int64))
}
*/
