package host

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"

	"../../../models"
)

// AddCellHostHandlerFunc turns a function with the right signature into a add cell host handler
type AddCellHostHandlerFunc func(AddCellHostParams, *models.Customer) middleware.Responder

// Handle executing the request and returning a response
func (fn AddCellHostHandlerFunc) Handle(params AddCellHostParams, principal *models.Customer) middleware.Responder {
	return fn(params, principal)
}

// AddCellHostHandler interface for that can handle valid add cell host params
type AddCellHostHandler interface {
	Handle(AddCellHostParams, *models.Customer) middleware.Responder
}

// NewAddCellHost creates a new http.Handler for the add cell host operation
func NewAddCellHost(ctx *middleware.Context, handler AddCellHostHandler) *AddCellHost {
	return &AddCellHost{Context: ctx, Handler: handler}
}

/*AddCellHost swagger:route POST /cell/{cell_id}/host host addCellHost

Add a new host

*/
type AddCellHost struct {
	Context *middleware.Context
	Handler AddCellHostHandler
}

func (o *AddCellHost) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewAddCellHostParams()

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
