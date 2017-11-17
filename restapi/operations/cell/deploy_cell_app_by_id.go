package cell

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"

	"../../../models"
)

// DeployCellAppByIDHandlerFunc turns a function with the right signature into a deploy cell app by Id handler
type DeployCellAppByIDHandlerFunc func(DeployCellAppByIDParams, *models.Customer) middleware.Responder

// Handle executing the request and returning a response
func (fn DeployCellAppByIDHandlerFunc) Handle(params DeployCellAppByIDParams, principal *models.Customer) middleware.Responder {
	return fn(params, principal)
}

// DeployCellAppByIDHandler interface for that can handle valid deploy cell app by Id params
type DeployCellAppByIDHandler interface {
	Handle(DeployCellAppByIDParams, *models.Customer) middleware.Responder
}

// NewDeployCellAppByID creates a new http.Handler for the deploy cell app by Id operation
func NewDeployCellAppByID(ctx *middleware.Context, handler DeployCellAppByIDHandler) *DeployCellAppByID {
	return &DeployCellAppByID{Context: ctx, Handler: handler}
}

/*DeployCellAppByID swagger:route GET /cell/{cell_id}/deploy/app cell deployCellAppById

Deploy cell apps by ID

Deploy a single cell apps

*/
type DeployCellAppByID struct {
	Context *middleware.Context
	Handler DeployCellAppByIDHandler
}

func (o *DeployCellAppByID) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewDeployCellAppByIDParams()

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