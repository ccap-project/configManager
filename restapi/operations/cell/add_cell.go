// Code generated by go-swagger; DO NOT EDIT.

package cell

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"

	"configManager/models"
)

// AddCellHandlerFunc turns a function with the right signature into a add cell handler
type AddCellHandlerFunc func(AddCellParams, *models.Customer) middleware.Responder

// Handle executing the request and returning a response
func (fn AddCellHandlerFunc) Handle(params AddCellParams, principal *models.Customer) middleware.Responder {
	return fn(params, principal)
}

// AddCellHandler interface for that can handle valid add cell params
type AddCellHandler interface {
	Handle(AddCellParams, *models.Customer) middleware.Responder
}

// NewAddCell creates a new http.Handler for the add cell operation
func NewAddCell(ctx *middleware.Context, handler AddCellHandler) *AddCell {
	return &AddCell{Context: ctx, Handler: handler}
}

/*AddCell swagger:route POST /cell cell addCell

Add a new cell

*/
type AddCell struct {
	Context *middleware.Context
	Handler AddCellHandler
}

func (o *AddCell) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewAddCellParams()

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
