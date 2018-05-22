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
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"configManager"
	"configManager/models"
	"configManager/restapi/operations/cell"

	"github.com/Sirupsen/logrus"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/structures/graph"
)

func NewAddCell(rt *configManager.Runtime) cell.AddCellHandler {
	return &addCell{rt: rt}
}

type addCell struct {
	rt *configManager.Runtime
}

func (ctx *addCell) Handle(params cell.AddCellParams, principal *models.Customer) middleware.Responder {

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(principal.Name)})

	cypher := `MATCH (c:Customer {name: {name}})
							CREATE (c)-[:OWN]->(cell:Cell { name: {cell_name}, id: {cell_id} })
							RETURN	cell.id AS id`

	if getCellByName(ctx.rt, principal.Name, params.Body.Name) != nil {
		ctxLogger.Error("cell already exists !")
		return cell.NewAddCellConflict().WithPayload(&models.APIResponse{Message: "cell already exists"})
	}

	db, err := ctx.rt.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j: ", err)
		return cell.NewAddCellInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: ", err)
		return cell.NewAddCellInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}
	defer stmt.Close()

	ulid := configManager.GetULID()

	ctxLogger = ctxLogger.WithFields(logrus.Fields{
		"cell_id": ulid})

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"name":      swag.StringValue(principal.Name),
		"cell_name": swag.StringValue(params.Body.Name),
		"cell_id":   ulid})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: ", err)
		return cell.NewAddCellInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		ctxLogger.Error("An error occurred getting next row: ", err)
		return cell.NewAddCellInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	ctxLogger.Info("OK")

	return cell.NewAddCellCreated().WithPayload(models.ULID(output[0].(string)))
}

func NewDeployCellByID(rt *configManager.Runtime) cell.DeployCellByIDHandler {
	return &deployCellByID{rt: rt}
}

type deployCellByID struct {
	rt *configManager.Runtime
}

func (ctx *deployCellByID) Handle(params cell.DeployCellByIDParams, principal *models.Customer) middleware.Responder {

	Cell := _getCellByID(ctx.rt, principal.Name, &params.CellID)

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID})

	if Cell == nil {
		ctxLogger.Error("cell does not exists !")
		return cell.NewDeployCellByIDNotFound()
	}

	EntireCell, cellErr := getCellRecursive(ctx.rt, principal.Name, Cell)

	log.Printf("DeployCell(%#v)", EntireCell)

	if cellErr != nil {
		return cell.NewDeployCellAppByIDInternalServerError().WithPayload(&models.APIResponse{Message: *cellErr})
	}

	if EntireCell == nil {
		ctxLogger.Warn("cell is empty")
		return cell.NewDeployCellByIDNoContent()
	}

	jsonOut, err := json.Marshal(EntireCell)
	if err != nil {
		ctxLogger.Error("decoding cell, ", err)
		return cell.NewDeployCellByIDInternalServerError()
	}

	jsonString := strings.NewReader(string(jsonOut))
	log.Println(jsonString)

	requestRes, err := http.Post("http://127.0.0.1:8080/v1/deploy", "application/json", jsonString)

	if err != nil {
		ctxLogger.Error("deploying cell, ", err)
		return cell.NewDeployCellByIDInternalServerError()
	}
	defer requestRes.Body.Close()

	response := cell.NewDeployCellByIDOK()

	buf := new(bytes.Buffer)
	buf.ReadFrom(requestRes.Body)

	response.Payload.Message = buf.String()

	if err != nil {
		ctxLogger.Error("reading deploy cell response, ", err)
		return cell.NewDeployCellByIDInternalServerError()
	}
	return response
}

func NewDeployCellAppByID(rt *configManager.Runtime) cell.DeployCellAppByIDHandler {
	return &deployCellAppByID{rt: rt}
}

type deployCellAppByID struct {
	rt *configManager.Runtime
}

func (ctx *deployCellAppByID) Handle(params cell.DeployCellAppByIDParams, principal *models.Customer) middleware.Responder {

	Cell := _getCellByID(ctx.rt, principal.Name, &params.CellID)

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID})

	if Cell == nil {
		ctxLogger.Error("cell does not exists !")
		return cell.NewDeployCellAppByIDNotFound()
	}

	EntireCell, cellErr := getCellRecursive(ctx.rt, principal.Name, Cell)

	log.Printf("DeployCell(%#v)", EntireCell)

	if cellErr != nil {
		return cell.NewDeployCellAppByIDInternalServerError().WithPayload(&models.APIResponse{Message: *cellErr})
	}

	if EntireCell == nil {
		ctxLogger.Error("cell is empty")
		return cell.NewDeployCellAppByIDNoContent()
	}

	jsonOut, err := json.Marshal(EntireCell)
	if err != nil {
		ctxLogger.Error("decoding cell, ", err)
		return cell.NewDeployCellAppByIDInternalServerError()
	}

	jsonString := strings.NewReader(string(jsonOut))
	log.Println(jsonString)

	requestRes, err := http.Post("http://127.0.0.1:8080/v1/application/deploy", "application/json", jsonString)

	if err != nil {
		ctxLogger.Error("deploying cell, ", err)
		return cell.NewDeployCellAppByIDInternalServerError()
	}
	defer requestRes.Body.Close()

	response := cell.NewDeployCellAppByIDOK()

	buf := new(bytes.Buffer)
	buf.ReadFrom(requestRes.Body)

	response.Payload.Message = buf.String()

	if err != nil {
		ctxLogger.Error("reading deploy cell response, ", err)
		return cell.NewDeployCellAppByIDInternalServerError()
	}
	return response
}

func NewGetCellByID(rt *configManager.Runtime) cell.GetCellByIDHandler {
	return &getCellByID{rt: rt}
}

type getCellByID struct {
	rt *configManager.Runtime
}

func (ctx *getCellByID) Handle(params cell.GetCellByIDParams, principal *models.Customer) middleware.Responder {

	cypher := `MATCH (c:Customer {name: {customer_name} })-[:OWN]->(cell:Cell {id: {cell_id}})
								RETURN cell.id as id,
												cell.name as name,
												cell.public_key as public_key`

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID})

	db, err := ctx.rt.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j:", err)
		return cell.NewAddCellInternalServerError().WithPayload(&models.APIResponse{Message: err.Error()})
	}

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Error("An error occurred preparing statement: ", err)
		return cell.NewGetCellByIDInternalServerError()
	}
	defer stmt.Close()

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID})

	if err != nil {
		ctxLogger.Error("An error occurred querying Neo: ", err)
		return cell.NewGetCellByIDInternalServerError()
	}

	output, _, err := rows.NextNeo()

	if err != nil {
		return cell.NewGetCellByIDNotFound()
	}
	_name := output[1].(string)
	_cell := &models.Cell{
		ID:   models.ULID(output[0].(string)),
		Name: &_name}

	return cell.NewGetCellByIDOK().WithPayload(_cell)
}

func NewGetCellFullByID(rt *configManager.Runtime) cell.GetCellFullByIDHandler {
	return &getCellFullByID{rt: rt}
}

type getCellFullByID struct {
	rt *configManager.Runtime
}

func (ctx *getCellFullByID) Handle(params cell.GetCellFullByIDParams, principal *models.Customer) middleware.Responder {

	Cell := _getCellByID(ctx.rt, principal.Name, &params.CellID)

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(principal.Name),
		"cell_id":       params.CellID})

	if Cell == nil {
		ctxLogger.Error("cell does not exists !")
		return cell.NewDeployCellByIDNotFound()
	}

	FullCell := getCellFull(ctx.rt, principal.Name, &params.CellID)

	if FullCell == nil {
		ctxLogger.Error("cell is empty")
		return cell.NewDeployCellByIDNotFound()
	}

	return cell.NewGetCellFullByIDOK().WithPayload(FullCell)
}

func NewFindCellByCustomer(rt *configManager.Runtime) cell.FindCellByCustomerHandler {
	return &findCellByCustomer{rt: rt}
}

type findCellByCustomer struct {
	rt *configManager.Runtime
}

func (ctx *findCellByCustomer) Handle(params cell.FindCellByCustomerParams, principal *models.Customer) middleware.Responder {

	var res []*models.Cell

	cypher := `MATCH (c:Customer {name: {name} })-[:OWN]->(cell:Cell)
								WHERE EXISTS(cell.id)
								RETURN cell.id as id,
												cell.name as name`

	ctxLogger := ctx.rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(principal.Name)})

	db, err := ctx.rt.DB().OpenPool()

	if err != nil {
		ctxLogger.Error("error connecting to neo4j: ", err)
		return cell.NewFindCellByCustomerInternalServerError()
	}
	defer db.Close()

	data, _, _, err := db.QueryNeoAll(cypher, map[string]interface{}{
		"name": swag.StringValue(principal.Name)})

	if err != nil {
		ctxLogger.Errorf("An error occurred querying Neo: ", err)
		return cell.NewFindCellByCustomerInternalServerError()

	} else if len(data) == 0 {
		return cell.NewFindCellByCustomerNotFound()
	}

	for _, row := range data {
		_name := row[1].(string)

		c := &models.Cell{
			ID:   models.ULID(row[0].(string)),
			Name: &_name}
		res = append(res, c)
	}

	return cell.NewFindCellByCustomerOK().WithPayload(res)
}

func getCellByName(rt *configManager.Runtime, customerName *string, cellName *string) *models.Cell {

	var cell *models.Cell
	cell = nil

	cypher := `MATCH (c:Customer {name: {name} })-[:OWN]->(cell:Cell)
								WHERE EXISTS(cell.id) AND cell.name = {cell_name}
								RETURN cell.id as id,
												cell.name as name`

	ctxLogger := rt.Logger().WithFields(logrus.Fields{
		"customer_name": swag.StringValue(customerName),
		"cell_name":     swag.StringValue(cellName)})

	db, err := rt.DB().OpenPool()
	if err != nil {
		ctxLogger.Error("error connecting to neo4j: ", err)
		return cell
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Errorf("An error occurred preparing statement: ", err)
		return cell
	}
	defer stmt.Close()

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"name":      swag.StringValue(customerName),
		"cell_name": swag.StringValue(cellName)})

	if err != nil {
		ctxLogger.Errorf("An error occurred querying Neo: ", err)
		return cell
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		return cell
	}
	_name := output[1].(string)

	cell = &models.Cell{ID: models.ULID(output[0].(string)),
		Name: &_name}

	return cell
}

func _getCellByID(rt *configManager.Runtime, customerName *string, cellID *string) *models.Cell {

	var cell *models.Cell
	cell = nil

	ctxLogger := rt.Logger().WithFields(logrus.Fields{
		"customer": swag.StringValue(customerName),
		"cell_id":  swag.StringValue(cellID)})

	cypher := `MATCH (c:Customer {name: {name} })-[:OWN]->(cell:Cell {id: {cell_id}})
								RETURN cell.id as id,
												cell.name as name`

	db, err := rt.DB().OpenPool()

	if err != nil {
		ctxLogger.Error("error connecting to neo4j: ", err)
		return cell
	}
	defer db.Close()

	stmt, err := db.PrepareNeo(cypher)
	if err != nil {
		ctxLogger.Errorf("An error occurred preparing statement: ", err)
		return cell
	}
	defer stmt.Close()

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"name":    swag.StringValue(customerName),
		"cell_id": swag.StringValue(cellID)})

	if err != nil {
		ctxLogger.Errorf("An error occurred querying Neo: ", err)
		return cell
	}

	output, _, err := rows.NextNeo()
	if err != nil {
		return cell
	}
	_name := output[1].(string)

	cell = &models.Cell{
		ID:   models.ULID(output[0].(string)),
		Name: &_name}

	return cell
}

/*
 * Return cell structure in ui format
 */
func getCellFull(rt *configManager.Runtime, customerName *string, cellID *string) *models.FullCell {
	cypher := `MATCH (customer:Customer{ name:{customer_name}})-[:OWN]->(cell:Cell {id: {cell_id}})
							MATCH (cell)-[:DEPLOY_WITH]->(keypair:Keypair),
										(cell)-[:USE]->(provider:Provider),
										(provider)-[:PROVIDER_IS]->(provider_type:ProviderType),
										(cell)-[:PROVIDES]->(component:Component)-[:USE]->(role:Role)
							OPTIONAL MATCH (role)-->(parameter:Parameter)
							OPTIONAL MATCH (component)-->(hostgroup:Hostgroup)
							OPTIONAL MATCH (component)-->(listener:Listener)
							OPTIONAL MATCH (cell)-->(host)-->(option:Option)
							OPTIONAL MATCH (cell)-->(loadbalancer)
							RETURN *`

	ctxLogger := rt.Logger().WithFields(logrus.Fields{
		"customer": swag.StringValue(customerName),
		"cell_id":  swag.StringValue(cellID)})

	conn := rt.DB()
	db, err := conn.OpenPool()

	if err != nil {
		ctxLogger.Error("error connecting to neo4j:", err)
		return nil
	}
	defer db.Close()

	data, _, _, err := db.QueryNeoAll(cypher, map[string]interface{}{
		"customer_name": swag.StringValue(customerName),
		"cell_id":       swag.StringValue(cellID)})

	if err != nil {
		ctxLogger.Errorf("An error occurred querying Neo: %s", err)
		return nil
	}

	res := new(models.FullCell)

	res.CustomerName = *customerName
	res.Keypair = getCellKeypair(rt, customerName, cellID)
	res.Provider = getProvider(rt, customerName, cellID)

	for _, row := range data {

		if len(res.Name) == 0 {
			cellNode := getNodeByLabel(row, "Cell")

			if len(cellNode) > 0 {
				res.Name = cellNode["name"].(string)
			}
		}
	}

	// Component
	res.Components, _ = _findCellComponents(rt, customerName, cellID)
	res.Loadbalancers, _ = _findCellLoadbalancers(rt, customerName, cellID)

	return (res)
}

/*
 * Return cell structure in deploy format
 */
func getCellRecursive(rt *configManager.Runtime, customerName *string, cell *models.Cell) (*models.EntireCell, *string) {

	cellID := string(cell.ID)
	ctxLogger := rt.Logger().WithFields(logrus.Fields{
		"customer": swag.StringValue(customerName),
		"cell_id":  cellID})

	res := new(models.EntireCell)

	res.CustomerName = *customerName
	res.Name = *cell.Name

	if res.Keypair = getCellKeypair(rt, customerName, &cellID); res.Keypair == nil {
		err_msg := "Cell needs an associated keypair"
		ctxLogger.Error(err_msg)
		return nil, &err_msg
	}

	if res.Provider = getProvider(rt, customerName, &cellID); res.Provider == nil {
		err_msg := "Cell needs an associated provider"
		ctxLogger.Error(err_msg)
		return nil, &err_msg
	}

	/*
	 * Loadbalancers
	 */
	loadbalancers, _ := _findCellLoadbalancers(rt, customerName, &cellID)
	for _, lb := range loadbalancers {
		lbID := string(lb.ID)

		// get lb members
		_, _, _, member := _getLoadbalancerMembers(rt, customerName, &cellID, &lbID)

		lb.Network = *_getLoadbalancerNetwork(rt, customerName, &cellID, &lbID)

		if member != nil {

			lb.Members = *member
			res.Loadbalancers = append(res.Loadbalancers, lb)
		}
	}

	/*
	 * Networks
	 */
	res.Networks, _ = _findCellNetworks(rt, customerName, &cellID)

	/*
	 * SecurityGroups
	 */
	components, _ := _findCellComponents(rt, customerName, &cellID)

	for _, component := range components {

		securityGroup := &models.Securitygroup{Name: *component.Name}

		for _, hg := range component.Hostgroups {

			hg.Securitygroups = append(hg.Securitygroups, *component.Name)

			// give ordering precedence to component value
			if (hg.Order == nil && component.Order != nil) ||
				(component.Order != nil && *component.Order <= *hg.Order) {
				hg.Order = new(int64)
				hg.Order = component.Order

			} else if hg.Order != nil {
				hg.Order = new(int64)
				hg.Order = hg.Order
			}
			hg.Roles = models.HostgroupRoles(component.Roles)
			res.Hostgroups = append(res.Hostgroups, hg)

		}

		// build SecurityRules
		for _, listener := range component.Listeners {

			connections := getComponentListenerConnections(rt, customerName, &cellID, &listener.ID)

			for _, conn := range *connections {
				var securityRule models.Securityrule

				securityRule.SourceSecuritygroup = conn
				securityRule.DestinationSecuritygroup = securityGroup.Name
				securityRule.Proto = *listener.Protocol
				securityRule.DestinationPort = strconv.Itoa(int(*listener.Port))

				securityGroup.Rules = append(securityGroup.Rules, &securityRule)
			}
		}
		res.Securitygroups = append(res.Securitygroups, securityGroup)
	}
	return res, nil
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

func copyString(key interface{}) *string {

	res := new(string)

	if key != nil {
		*res = key.(string)
	}

	return res
}
