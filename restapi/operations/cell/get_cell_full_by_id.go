package cell

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"

	"../../../models"
)

// GetCellFullByIDHandlerFunc turns a function with the right signature into a get cell full by Id handler
type GetCellFullByIDHandlerFunc func(GetCellFullByIDParams, *models.Customer) middleware.Responder

// Handle executing the request and returning a response
func (fn GetCellFullByIDHandlerFunc) Handle(params GetCellFullByIDParams, principal *models.Customer) middleware.Responder {
	return fn(params, principal)
}

// GetCellFullByIDHandler interface for that can handle valid get cell full by Id params
type GetCellFullByIDHandler interface {
	Handle(GetCellFullByIDParams, *models.Customer) middleware.Responder
}

// NewGetCellFullByID creates a new http.Handler for the get cell full by Id operation
func NewGetCellFullByID(ctx *middleware.Context, handler GetCellFullByIDHandler) *GetCellFullByID {
	return &GetCellFullByID{Context: ctx, Handler: handler}
}

/*GetCellFullByID swagger:route GET /cell/{cell_id}/full cell getCellFullById

Get full cell by ID

Returns full cell definition

*/
type GetCellFullByID struct {
	Context *middleware.Context
	Handler GetCellFullByIDHandler
}

func (o *GetCellFullByID) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetCellFullByIDParams()

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