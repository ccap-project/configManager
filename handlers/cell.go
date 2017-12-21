package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"configManager/models"
	"configManager/restapi/operations/cell"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	driver "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/structures/graph"
)

func AddCell(params cell.AddCellParams, principal *models.Customer) middleware.Responder {

	cypher := `MATCH (c:Customer {name: {name} })
							CREATE (c)-[:OWN]->(cell:Cell { name: {cell_name} })
							RETURN	id(cell) AS id,
											cell.name AS name`

	if getCellByName(principal.Name, params.Body.Name) != nil {
		log.Println("cell already exists !")
		return cell.NewAddCellConflict().WithPayload(models.APIResponse{Message: "cell already exists"})
	}

	db, err := driver.NewDriver().OpenNeo("bolt://192.168.20.54:7687")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return cell.NewAddCellInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return cell.NewAddCellInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{"name": swag.StringValue(principal.Name),
		"cell_name": swag.StringValue(params.Body.Name)})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return cell.NewAddCellInternalServerError().WithPayload(models.APIResponse{Message: err.Error()})
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

func DeployCell(params cell.DeployCellByIDParams, principal *models.Customer) middleware.Responder {

	Cell := getCellByID(principal.Name, params.CellID)

	if Cell == nil {
		log.Println("cell does not exists !")
		return cell.NewDeployCellByIDNotFound()
	}

	EntireCell := getCellRecursive(principal.Name, params.CellID)

	log.Printf("DeployCell(%#v)", EntireCell)

	if EntireCell == nil {
		log.Print("cell is empty")
		return cell.NewDeployCellByIDNoContent()
	}

	jsonOut, err := json.Marshal(EntireCell)
	if err != nil {
		log.Println("decoding cell, ", err)
		return cell.NewDeployCellByIDInternalServerError()
	}

	jsonString := strings.NewReader(string(jsonOut))
	log.Println(jsonString)

	requestRes, err := http.Post("http://127.0.0.1:8080/v1/deploy", "application/json", jsonString)

	if err != nil {
		log.Println("deploying cell, ", err)
		return cell.NewDeployCellByIDInternalServerError()
	}
	defer requestRes.Body.Close()

	response := cell.NewDeployCellByIDOK()

	buf := new(bytes.Buffer)
	buf.ReadFrom(requestRes.Body)

	response.Payload.Message = buf.String()

	if err != nil {
		log.Println("reading deploy cell response, ", err)
		return cell.NewDeployCellByIDInternalServerError()
	}
	return response
}

func DeployCellApp(params cell.DeployCellAppByIDParams, principal *models.Customer) middleware.Responder {

	Cell := getCellByID(principal.Name, params.CellID)

	if Cell == nil {
		log.Println("cell does not exists !")
		return cell.NewDeployCellAppByIDNotFound()
	}

	EntireCell := getCellRecursive(principal.Name, params.CellID)

	if EntireCell == nil {
		log.Print("cell is empty")
		return cell.NewDeployCellAppByIDNoContent()
	}

	jsonOut, err := json.Marshal(EntireCell)
	if err != nil {
		log.Println("decoding cell, ", err)
		return cell.NewDeployCellAppByIDInternalServerError()
	}

	jsonString := strings.NewReader(string(jsonOut))
	log.Println(jsonString)

	requestRes, err := http.Post("http://127.0.0.1:8080/v1/application/deploy", "application/json", jsonString)

	if err != nil {
		log.Println("deploying cell, ", err)
		return cell.NewDeployCellAppByIDInternalServerError()
	}
	defer requestRes.Body.Close()

	response := cell.NewDeployCellAppByIDOK()

	buf := new(bytes.Buffer)
	buf.ReadFrom(requestRes.Body)

	response.Payload.Message = buf.String()

	if err != nil {
		log.Println("reading deploy cell response, ", err)
		return cell.NewDeployCellAppByIDInternalServerError()
	}
	return response
}

func GetCellByID(params cell.GetCellByIDParams, principal *models.Customer) middleware.Responder {

	cypher := `MATCH (c:Customer {name: {name} })-[:HAS]->(k:Cell)
								WHERE ID(k) = {kid}
								RETURN ID(c) as id,
												k.name as name,
												k.public_key as public_key`

	db, err := driver.NewDriver().OpenNeo("bolt://192.168.20.54:7687")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return cell.NewGetCellByIDInternalServerError()
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return cell.NewGetCellByIDInternalServerError()
	}
	defer stmt.Close()

	rows, err := stmt.QueryNeo(map[string]interface{}{"name": swag.StringValue(principal.Name),
		"kid": params.CellID})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return cell.NewGetCellByIDInternalServerError()
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		return cell.NewGetCellByIDNotFound()
	}
	_name := output[1].(string)
	_cell := &models.Cell{ID: output[0].(int64),
		Name: &_name}

	return cell.NewGetCellByIDOK().WithPayload(_cell)
}

func GetCellFullByID(params cell.GetCellFullByIDParams, principal *models.Customer) middleware.Responder {

	Cell := getCellByID(principal.Name, params.CellID)

	if Cell == nil {
		log.Println("cell does not exists !")
		return cell.NewDeployCellByIDNotFound()
	}

	FullCell := getCellFull(principal.Name, params.CellID)

	if FullCell == nil {
		log.Print("cell is empty")
		return cell.NewDeployCellByIDNotFound()
	}

	return cell.NewGetCellFullByIDOK().WithPayload(FullCell)
}

func FindCellByCustomer(params cell.FindCellByCustomerParams, principal *models.Customer) middleware.Responder {
	cypher := `MATCH (c:Customer {name: {name} })-[:OWN]->(cell:Cell)
								RETURN ID(cell) as id,
												cell.name as name`

	db, err := driver.NewDriver().OpenNeo("bolt://192.168.20.54:7687")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return cell.NewFindCellByCustomerInternalServerError()
	}
	defer db.Close()

	data, _, _, err := db.QueryNeoAll(cypher, map[string]interface{}{"name": swag.StringValue(principal.Name)})

	log.Printf("= data(%#v)", data)

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return cell.NewFindCellByCustomerInternalServerError()

	} else if len(data) == 0 {
		return cell.NewFindCellByCustomerNotFound()
	}

	res := make([]*models.Cell, len(data))

	for idx, row := range data {
		_name := row[1].(string)

		res[idx] = &models.Cell{
			ID:   row[0].(int64),
			Name: &_name}
	}

	log.Printf("= Res(%#v)", res)

	return cell.NewFindCellByCustomerOK().WithPayload(res)
}

func getCellByName(customerName *string, cellName *string) *models.Cell {

	var cell *models.Cell
	cell = nil

	cypher := `MATCH (c:Customer {name: {name} })-[:OWN]->(cell:Cell)
								WHERE cell.name = {cell_name}
								RETURN ID(cell) as id,
												cell.name as name`

	db, err := driver.NewDriver().OpenNeo("bolt://192.168.20.54:7687")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return cell
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return cell
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{"name": swag.StringValue(customerName),
		"cell_name": swag.StringValue(cellName)})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return cell
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		return cell
	}
	_name := output[1].(string)

	cell = &models.Cell{ID: output[0].(int64),
		Name: &_name}

	stmt.Close()

	return cell
}

func getCellByID(customerName *string, cellID int64) *models.Cell {

	var cell *models.Cell
	cell = nil

	cypher := `MATCH (c:Customer {name: {name} })-[:OWN]->(cell:Cell)
								WHERE id(cell) = {cell_id}
								RETURN ID(cell) as id,
												cell.name as name`

	db, err := driver.NewDriver().OpenNeo("bolt://192.168.20.54:7687")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return cell
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		log.Printf("An error occurred preparing statement: %s", err)
		return cell
	}
	defer stmt.Close()

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"name":    swag.StringValue(customerName),
		"cell_id": cellID})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return cell
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		return cell
	}
	_name := output[1].(string)

	cell = &models.Cell{
		ID:   output[0].(int64),
		Name: &_name}

	return cell
}

/*
 * Return cell structure in ui format
 */
func getCellFull(customerName *string, cellID int64) *models.FullCell {
	cypher := `MATCH (customer:Customer{ name:{customer_name}})-[:OWN]->(cell:Cell)
							WHERE id(cell) = {cell_id}
							MATCH (cell)-[:DEPLOY_WITH]->(keypair:Keypair),
										(cell)-[:USE]->(provider:Provider),
										(provider)-[:PROVIDER_IS]->(provider_type:ProviderType),
										(cell)-[:PROVIDES]->(component:Component)-[:USE]->(role:Role)
							OPTIONAL MATCH (role)-->(parameter:Parameter)
							OPTIONAL MATCH (component)-->(hostgroup:Hostgroup)
							OPTIONAL MATCH (cell)-->(host)-->(option:Option)
							RETURN *`

	db, err := driver.NewDriver().OpenNeo("bolt://192.168.20.54:7687")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return nil
	}
	defer db.Close()

	data, _, _, err := db.QueryNeoAll(cypher, map[string]interface{}{
		"customer_name": swag.StringValue(customerName),
		"cell_id":       cellID})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return nil
	}

	res := new(models.FullCell)

	log.Printf("res(%v)", res)

	res.CustomerName = *customerName
	//res.Keypair = new(models.Keypair)
	res.Provider = new(models.Provider)

	for _, row := range data {

		var cellID int64
		if len(res.Name) == 0 {
			cellNode := getNodeByLabel(row, "Cell")

			if len(cellNode) > 0 {
				res.Name = cellNode["name"].(string)
			}
		}

		res.Keypair = getCellKeypair(customerName, cellID)
		/*
			if res.Keypair.Name == nil {
				keypairNode := getNodeByLabel(row, "Keypair")

				if len(keypairNode) > 0 {
					res.Keypair.Name = new(string)
					res.Keypair.PublicKey = new(string)

					*res.Keypair.Name = keypairNode["name"].(string)
					*res.Keypair.PublicKey = keypairNode["public_key"].(string)
				}
			}
		*/

		if res.Provider.Name == nil {
			providerNode := getNodeByLabel(row, "Provider")
			providerTypeNode := getNodeByLabel(row, "ProviderType")

			if len(providerNode) > 0 {
				res.Provider.AuthURL = copyString(providerNode["auth_url"])
				res.Provider.DomainName = copyString(providerNode["domain_name"])
				res.Provider.Name = copyString(providerNode["name"])
				res.Provider.Password = copyString(providerNode["password"])
				res.Provider.TenantName = copyString(providerNode["tenantname"])
				res.Provider.Username = copyString(providerNode["username"])
				res.Provider.Type = copyString(providerTypeNode["name"])
			}
		}

		// Component
		C, err := findCellComponents(customerName, cellID)
		res.Components = C
		log.Printf(">>>>>>>>>>>>>>>>> %#v<<<<<<  %v<<<<<<<<<<<<", C, err)
		//componentNode := getNodeByLabel(row, "Component")
		/*
			if len(componentNode) > 0 {
				var component *models.Component

				component = _getComponentByName(res.Components, componentNode["name"].(string))

				if component == nil {
					component = new(models.Component)

					component.Name = copyString(componentNode["name"])

					log.Printf("-------->>>>> %#v", componentNode)

					res.Components = append(res.Components, component)
				}
				// Hostgroup
				hostgroupNode := getNodeByLabel(row, "Hostgroup")

				if len(hostgroupNode) > 0 {

					component.Hostgroups, _ = _FindComponentHostgroups(&res.CustomerName, cellID, componentNode["id"].(int64))
					/*
						var hg *models.Hostgroup

						hg = getHostgroupByName(component.Hostgroups, hostgroupNode["name"].(string))

						if hg == nil {
							hg = new(models.Hostgroup)

							hg.Flavor = copyString(hostgroupNode["flavor"])
							hg.Image = copyString(hostgroupNode["image"])
							hg.Name = copyString(hostgroupNode["name"])
							hg.Network = copyString(hostgroupNode["network"])
							hg.Username = copyString(hostgroupNode["username"])
							hg.BootstrapCommand = *copyString(hostgroupNode["bootstrap_command"])

							hg.Count = new(int64)

							*hg.Count = hostgroupNode["count"].(int64)

							//component.Hostgroups = append(component.Hostgroups, hg)
							component.Hostgroups = handlers.
						}
				}

				// Roles
				roleNode := getNodeByLabel(row, "Role")

				if len(roleNode) > 0 {

					var role *models.Role

					role = getRoleByName(component.Roles, roleNode["name"].(string))

					if role == nil {
						role = new(models.Role)

						role.Name = copyString(roleNode["name"])
						role.URL = copyString(roleNode["url"])
						role.Version = copyString(roleNode["version"])

						component.Roles = append(component.Roles, role)
					}

					parameterNode := getNodeByLabel(row, "Parameter")

					if parameterNode != nil {
						var parameter *models.Parameter
						parameter = getParameterByName(role.Params, parameterNode["name"].(string))

						if parameter == nil {
							parameter = new(models.Parameter)

							parameter.Name = copyString(parameterNode["name"])
							parameter.Value = copyString(parameterNode["value"])

							role.Params = append(role.Params, parameter)
						}
					}
				}
			}
		*/
	}

	return (res)
}

/*
 * Return cell structure in deploy format
 */
func getCellRecursive(customerName *string, cellID int64) *models.EntireCell {
	cypher := `MATCH (customer:Customer{ name:{customer_name}})-[:OWN]->(cell:Cell)
							WHERE id(cell) = {cell_id}
							MATCH (cell)-[:DEPLOY_WITH]->(keypair:Keypair),
										(cell)-[:USE]->(provider:Provider),
										(provider)-[:PROVIDER_IS]->(provider_type:ProviderType),
										(cell)-[:PROVIDES]->(component:Component)-[:USE]->(role:Role)
							OPTIONAL MATCH (role)-->(parameter:Parameter)
							OPTIONAL MATCH (component)-->(hostgroup:Hostgroup)
							OPTIONAL MATCH (cell)-->(host)-->(option:Option)
							RETURN *
							ORDER BY component.name, role.order`
	/*
		cypher := `MATCH (customer:Customer{ name:{customer_name}})-[:OWN]->(cell:Cell)
								WHERE id(cell) = {cell_id}
								MATCH (cell)-[:DEPLOY_WITH]->(keypair:Keypair)
								MATCH (cell)-[:HAS]->(host:Host)
								MATCH (cell)-[:USE]->(provider:Provider)
								MATCH (provider)-[:PROVIDER_IS]->(provider_type:ProviderType)
								MATCH (cell)-[:PROVIDES]->(component:Component)-[:USE]->(role:Role)
								OPTIONAL MATCH (role)-->(parameter:Parameter)
								OPTIONAL MATCH (component)-->(hostgroup:Hostgroup)
								OPTIONAL MATCH (host)-->(option:Option)
								RETURN *
								ORDER BY component.name, role.order`
	*/
	db, err := driver.NewDriver().OpenNeo("bolt://192.168.20.54:7687")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return nil
	}
	defer db.Close()

	data, _, _, err := db.QueryNeoAll(cypher, map[string]interface{}{
		"customer_name": swag.StringValue(customerName),
		"cell_id":       cellID})

	if err != nil {
		log.Printf("An error occurred querying Neo: %s", err)
		return nil
	}

	res := new(models.EntireCell)

	res.CustomerName = *customerName
	res.Keypair = new(models.Keypair)
	res.Provider = new(models.Provider)

	for _, row := range data {
		if len(res.Name) == 0 {
			cellNode := getNodeByLabel(row, "Cell")

			if len(cellNode) > 0 {
				res.Name = cellNode["name"].(string)
			}
		}

		if len(res.CustomerName) == 0 {
			customerNode := getNodeByLabel(row, "Customer")

			if len(customerNode) > 0 {
				res.CustomerName = customerNode["name"].(string)
			}
		}

		componentNode := getNodeByLabel(row, "Component")

		if res.Keypair.Name == nil {
			keypairNode := getNodeByLabel(row, "Keypair")

			if len(keypairNode) > 0 {
				res.Keypair.Name = new(string)
				res.Keypair.PublicKey = new(string)

				*res.Keypair.Name = keypairNode["name"].(string)
				*res.Keypair.PublicKey = keypairNode["public_key"].(string)
			}
		}

		if res.Provider.Name == nil {
			providerNode := getNodeByLabel(row, "Provider")
			providerTypeNode := getNodeByLabel(row, "ProviderType")

			if len(providerNode) > 0 {
				res.Provider.AuthURL = copyString(providerNode["auth_url"])
				res.Provider.DomainName = copyString(providerNode["domain_name"])
				res.Provider.Name = copyString(providerNode["name"])
				res.Provider.Password = copyString(providerNode["password"])
				res.Provider.TenantName = copyString(providerNode["tenantname"])
				res.Provider.Username = copyString(providerNode["username"])

				res.Provider.Type = copyString(providerTypeNode["name"])
			}
		}

		// Hosts
		hostNode := getNodeByLabel(row, "Host")

		if len(hostNode) > 0 {

			var h *models.Host

			h = getHostByName(res.Hosts, hostNode["name"].(string))

			if h == nil {
				h = new(models.Host)

				h.Name = copyString(hostNode["name"])

				res.Hosts = append(res.Hosts, h)
			}

			optionNode := getNodeByLabel(row, "Option")

			if optionNode != nil {
				var option *models.Parameter
				option = getParameterByName(h.Options, optionNode["name"].(string))

				if option == nil {
					option = new(models.Parameter)

					option.Name = copyString(optionNode["name"])
					option.Value = copyString(optionNode["value"])

					h.Options = append(h.Options, option)
				}
			}

		}

		// Hostgroup
		hostgroupNode := getNodeByLabel(row, "Hostgroup")

		if len(hostgroupNode) > 0 {

			var hg *models.Hostgroup

			hg = getHostgroupByName(res.Hostgroups, hostgroupNode["name"].(string))

			if hg == nil {
				hg = new(models.Hostgroup)

				hg.Flavor = copyString(hostgroupNode["flavor"])
				hg.Image = copyString(hostgroupNode["image"])
				hg.Name = copyString(hostgroupNode["name"])
				hg.Network = copyString(hostgroupNode["network"])
				hg.Username = copyString(hostgroupNode["username"])
				hg.BootstrapCommand = *copyString(hostgroupNode["bootstrap_command"])
				hg.Component = *copyString(componentNode["name"])
				hg.Count = new(int64)

				*hg.Count = hostgroupNode["count"].(int64)

				res.Hostgroups = append(res.Hostgroups, hg)

			}

			// Roles
			roleNode := getNodeByLabel(row, "Role")

			if len(roleNode) > 0 {

				var role *models.Role

				role = getRoleByName(hg.Roles, roleNode["name"].(string))

				if role == nil {
					role = new(models.Role)

					role.Name = copyString(roleNode["name"])
					role.URL = copyString(roleNode["url"])
					role.Version = copyString(roleNode["version"])

					hg.Roles = append(hg.Roles, role)
				}

				parameterNode := getNodeByLabel(row, "Parameter")

				if parameterNode != nil {
					var parameter *models.Parameter
					parameter = getParameterByName(role.Params, parameterNode["name"].(string))

					if parameter == nil {
						parameter = new(models.Parameter)

						parameter.Name = copyString(parameterNode["name"])
						parameter.Value = copyString(parameterNode["value"])

						role.Params = append(role.Params, parameter)
					}
				}
			}
		}
	}

	return (res)
}

func getNodeByLabel(row []interface{}, nodeName string) map[string]interface{} {
	for _, node := range row {
		if node == nil {
			continue
		}

		for _, label := range node.(graph.Node).Labels {
			if strings.Compare(nodeName, label) == 0 {
				return node.(graph.Node).Properties
			}
		}
	}

	var res map[string]interface{}

	return res
}

func _getComponentByName(components []*models.Component, componentName string) *models.Component {
	for _, component := range components {
		if strings.Compare(componentName, *component.Name) == 0 {
			return component
		}
	}

	return nil
}

func getHostByName(hosts []*models.Host, hostName string) *models.Host {
	for _, host := range hosts {
		if strings.Compare(hostName, *host.Name) == 0 {
			return host
		}
	}

	return nil
}

func getHostgroupByName(hostgroups []*models.Hostgroup, hostgroupName string) *models.Hostgroup {
	for _, hostgroup := range hostgroups {
		if strings.Compare(hostgroupName, *hostgroup.Name) == 0 {
			return hostgroup
		}
	}

	return nil
}

func getParameterByName(params []*models.Parameter, paramName string) *models.Parameter {
	for _, param := range params {
		if strings.Compare(paramName, *param.Name) == 0 {
			return param
		}
	}

	return nil
}

func getRoleByName(roles []*models.Role, roleName string) *models.Role {

	if roles != nil {
		for _, role := range roles {
			if strings.Compare(roleName, *role.Name) == 0 {
				return role
			}
		}
	}

	return nil
}

func copyString(key interface{}) *string {

	res := new(string)

	if key != nil {
		*res = key.(string)
	}

	return res
}
