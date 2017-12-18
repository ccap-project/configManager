package restapi

import (
	"crypto/tls"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	//"github.com/go-swagger/go-swagger/examples/authentication/models"
	graceful "github.com/tylerb/graceful"

	"../handlers"
	"../models"
	"./operations"
	"./operations/cell"
	"./operations/component"
	"./operations/customer"
	"./operations/host"
	"./operations/hostgroup"
	"./operations/keypair"
	"./operations/provider"
	"./operations/providertype"
	"./operations/role"
)

// This file is safe to edit. Once it exists it will not be overwritten

//go:generate swagger generate server --target .. --name  --spec ../swagger.yml

func configureFlags(api *operations.ConfigManagerAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.ConfigManagerAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// s.api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.UrlformConsumer = runtime.DiscardConsumer

	api.JSONProducer = runtime.JSONProducer()

	// Applies when the "x-api-token" header is set
	api.APIKeyHeaderAuth = func(token string) (*models.Customer, error) {

		Customer := new(models.Customer)
		Customer.Name = new(string)
		*Customer.Name = "customer1"
		Customer.ID = 84

		return Customer, nil
		return nil, errors.NotImplemented("api key auth (APIKeyHeader) x-api-token from header param [x-api-token] has not yet been implemented")
	}

	// Provider Type
	api.ProvidertypeAddProviderTypeHandler = providertype.AddProviderTypeHandlerFunc(handlers.AddProviderType)
	api.ProvidertypeDeleteProviderTypeHandler = providertype.DeleteProviderTypeHandlerFunc(func(params providertype.DeleteProviderTypeParams) middleware.Responder {
		return middleware.NotImplemented("operation providertype.DeleteProviderType has not yet been implemented")
	})
	api.ProvidertypeGetProviderTypeByIDHandler = providertype.GetProviderTypeByIDHandlerFunc(handlers.GetProviderTypeByID)
	api.ProvidertypeListProviderTypesHandler = providertype.ListProviderTypesHandlerFunc(handlers.ListProviderTypes)

	// Customer
	api.CustomerAddCustomerHandler = customer.AddCustomerHandlerFunc(handlers.AddCustomer)
	api.CustomerDeleteCustomerHandler = customer.DeleteCustomerHandlerFunc(func(params customer.DeleteCustomerParams) middleware.Responder {
		return middleware.NotImplemented("operation customer.DeleteCustomer has not yet been implemented")
	})
	api.CustomerFindCustomerByNameHandler = customer.FindCustomerByNameHandlerFunc(handlers.GetCustomerByName)
	api.CustomerGetCustomerByIDHandler = customer.GetCustomerByIDHandlerFunc(func(params customer.GetCustomerByIDParams) middleware.Responder {
		return middleware.NotImplemented("operation customer.GetCustomerByID has not yet been implemented")
	})

	// Key Pair
	api.KeypairAddKeypairHandler = keypair.AddKeypairHandlerFunc(handlers.AddKeypair)
	api.KeypairGetKeypairByIDHandler = keypair.GetKeypairByIDHandlerFunc(handlers.GetKeypairByID)
	api.KeypairFindKeypairByCustomerHandler = keypair.FindKeypairByCustomerHandlerFunc(handlers.FindKeypairByCustomer)
	api.KeypairAddCellKeypairHandler = keypair.AddCellKeypairHandlerFunc(handlers.AddCellKeypair)

	// Cell
	api.CellAddCellHandler = cell.AddCellHandlerFunc(handlers.AddCell)
	api.CellFindCellByCustomerHandler = cell.FindCellByCustomerHandlerFunc(handlers.FindCellByCustomer)
	api.CellGetCellByIDHandler = cell.GetCellByIDHandlerFunc(handlers.GetCellByID)
	api.CellGetCellFullByIDHandler = cell.GetCellFullByIDHandlerFunc(handlers.GetCellFullByID)

	// Host
	api.HostAddCellHostHandler = host.AddCellHostHandlerFunc(handlers.AddCellHost)

	// Deploy
	api.CellDeployCellByIDHandler = cell.DeployCellByIDHandlerFunc(handlers.DeployCell)
	api.CellDeployCellAppByIDHandler = cell.DeployCellAppByIDHandlerFunc(handlers.DeployCellApp)

	// Component
	api.ComponentAddComponentHandler = component.AddComponentHandlerFunc(handlers.AddCellComponent)
	api.ComponentGetCellComponentHandler = component.GetCellComponentHandlerFunc(handlers.GetCellComponent)
	api.ComponentFindCellComponentsHandler = component.FindCellComponentsHandlerFunc(handlers.FindCellComponents)

	// Hostgroup
	api.HostgroupAddComponentHostgroupHandler = hostgroup.AddComponentHostgroupHandlerFunc(handlers.AddComponentHostgroup)
	api.HostgroupDeleteComponentHostgroupHandler = hostgroup.DeleteComponentHostgroupHandlerFunc(handlers.DeleteComponentHostgroup)
	api.HostgroupFindComponentHostgroupsHandler = hostgroup.FindComponentHostgroupsHandlerFunc(handlers.FindComponentHostgroups)
	api.HostgroupGetComponentHostgroupByIDHandler = hostgroup.GetComponentHostgroupByIDHandlerFunc(handlers.GetComponentHostgroupByID)
	api.HostgroupUpdateComponentHostgroupHandler = hostgroup.UpdateComponentHostgroupHandlerFunc(handlers.UpdateComponentHostgroup)

	// Roles
	api.RoleAddComponentRoleHandler = role.AddComponentRoleHandlerFunc(handlers.AddComponentRole)
	api.RoleDeleteComponentRoleHandler = role.DeleteComponentRoleHandlerFunc(handlers.DeleteComponentRole)
	api.RoleFindComponentRolesHandler = role.FindComponentRolesHandlerFunc(handlers.FindComponentRoles)
	api.RoleUpdateComponentRoleHandler = role.UpdateComponentRoleHandlerFunc(handlers.UpdateComponentRole)

	//api.RoleAddRoleHandler = role.AddRoleHandlerFunc(handlers.AddRole)
	//api.RoleGetRoleByIDHandler = role.GetRoleByIDHandlerFunc(handlers.GetRoleByID)
	//api.RoleFindRolesHandler = role.FindRolesHandlerFunc(handlers.FindRoles)

	// Provider
	api.ProviderAddProviderHandler = provider.AddProviderHandlerFunc(handlers.AddCellProvider)
	api.ProviderGetProviderHandler = provider.GetProviderHandlerFunc(handlers.GetCellProvider)
	api.ProviderUpdateProviderHandler = provider.UpdateProviderHandlerFunc(handlers.UpdateCellProvider)

	api.CustomerUpdateCustomerHandler = customer.UpdateCustomerHandlerFunc(func(params customer.UpdateCustomerParams) middleware.Responder {
		return middleware.NotImplemented("operation customer.UpdateCustomer has not yet been implemented")
	})
	api.CustomerUpdateCustomerWithFormHandler = customer.UpdateCustomerWithFormHandlerFunc(func(params customer.UpdateCustomerWithFormParams) middleware.Responder {
		return middleware.NotImplemented("operation customer.UpdateCustomerWithForm has not yet been implemented")
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *graceful.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
