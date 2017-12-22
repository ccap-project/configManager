// Code generated by go-swagger; DO NOT EDIT.

package hostgroup

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"

	"configManager/models"
)

// DeleteComponentHostgroupHandlerFunc turns a function with the right signature into a delete component hostgroup handler
type DeleteComponentHostgroupHandlerFunc func(DeleteComponentHostgroupParams, *models.Customer) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteComponentHostgroupHandlerFunc) Handle(params DeleteComponentHostgroupParams, principal *models.Customer) middleware.Responder {
	return fn(params, principal)
}

// DeleteComponentHostgroupHandler interface for that can handle valid delete component hostgroup params
type DeleteComponentHostgroupHandler interface {
	Handle(DeleteComponentHostgroupParams, *models.Customer) middleware.Responder
}

// NewDeleteComponentHostgroup creates a new http.Handler for the delete component hostgroup operation
func NewDeleteComponentHostgroup(ctx *middleware.Context, handler DeleteComponentHostgroupHandler) *DeleteComponentHostgroup {
	return &DeleteComponentHostgroup{Context: ctx, Handler: handler}
}

/*DeleteComponentHostgroup swagger:route DELETE /cell/{cell_id}/component/{component_id}/hostgroup/{hostgroup_id} hostgroup deleteComponentHostgroup

Deletes a hostgroup from component

*/
type DeleteComponentHostgroup struct {
	Context *middleware.Context
	Handler DeleteComponentHostgroupHandler
}

func (o *DeleteComponentHostgroup) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewDeleteComponentHostgroupParams()

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