// Code generated by go-swagger; DO NOT EDIT.

package role

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"

	"configManager/models"
)

// AddComponentRoleHandlerFunc turns a function with the right signature into a add component role handler
type AddComponentRoleHandlerFunc func(AddComponentRoleParams, *models.Customer) middleware.Responder

// Handle executing the request and returning a response
func (fn AddComponentRoleHandlerFunc) Handle(params AddComponentRoleParams, principal *models.Customer) middleware.Responder {
	return fn(params, principal)
}

// AddComponentRoleHandler interface for that can handle valid add component role params
type AddComponentRoleHandler interface {
	Handle(AddComponentRoleParams, *models.Customer) middleware.Responder
}

// NewAddComponentRole creates a new http.Handler for the add component role operation
func NewAddComponentRole(ctx *middleware.Context, handler AddComponentRoleHandler) *AddComponentRole {
	return &AddComponentRole{Context: ctx, Handler: handler}
}

/*AddComponentRole swagger:route POST /cell/{cell_id}/component/{component_id}/role role addComponentRole

Add a new role to a component

*/
type AddComponentRole struct {
	Context *middleware.Context
	Handler AddComponentRoleHandler
}

func (o *AddComponentRole) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewAddComponentRoleParams()

	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		r = aCtx
	}
	var principal *models.Customer
	if uprinc != nil {
		principal = uprinc.(*models.Customer) // this is really a models.Customer, I promise
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
