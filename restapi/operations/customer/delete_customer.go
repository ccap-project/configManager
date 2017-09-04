package customer

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// DeleteCustomerHandlerFunc turns a function with the right signature into a delete customer handler
type DeleteCustomerHandlerFunc func(DeleteCustomerParams) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteCustomerHandlerFunc) Handle(params DeleteCustomerParams) middleware.Responder {
	return fn(params)
}

// DeleteCustomerHandler interface for that can handle valid delete customer params
type DeleteCustomerHandler interface {
	Handle(DeleteCustomerParams) middleware.Responder
}

// NewDeleteCustomer creates a new http.Handler for the delete customer operation
func NewDeleteCustomer(ctx *middleware.Context, handler DeleteCustomerHandler) *DeleteCustomer {
	return &DeleteCustomer{Context: ctx, Handler: handler}
}

/*DeleteCustomer swagger:route DELETE /customer/{customerId} customer deleteCustomer

Deletes a customer

*/
type DeleteCustomer struct {
	Context *middleware.Context
	Handler DeleteCustomerHandler
}

func (o *DeleteCustomer) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewDeleteCustomerParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
