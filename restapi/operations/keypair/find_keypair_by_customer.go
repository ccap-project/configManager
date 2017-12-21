// Code generated by go-swagger; DO NOT EDIT.

package keypair

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"

	"configManager/models"
)

// FindKeypairByCustomerHandlerFunc turns a function with the right signature into a find keypair by customer handler
type FindKeypairByCustomerHandlerFunc func(FindKeypairByCustomerParams, *models.Customer) middleware.Responder

// Handle executing the request and returning a response
func (fn FindKeypairByCustomerHandlerFunc) Handle(params FindKeypairByCustomerParams, principal *models.Customer) middleware.Responder {
	return fn(params, principal)
}

// FindKeypairByCustomerHandler interface for that can handle valid find keypair by customer params
type FindKeypairByCustomerHandler interface {
	Handle(FindKeypairByCustomerParams, *models.Customer) middleware.Responder
}

// NewFindKeypairByCustomer creates a new http.Handler for the find keypair by customer operation
func NewFindKeypairByCustomer(ctx *middleware.Context, handler FindKeypairByCustomerHandler) *FindKeypairByCustomer {
	return &FindKeypairByCustomer{Context: ctx, Handler: handler}
}

/*FindKeypairByCustomer swagger:route GET /keypairs keypair findKeypairByCustomer

Finds Keypair by customer

*/
type FindKeypairByCustomer struct {
	Context *middleware.Context
	Handler FindKeypairByCustomerHandler
}

func (o *FindKeypairByCustomer) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewFindKeypairByCustomerParams()

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
