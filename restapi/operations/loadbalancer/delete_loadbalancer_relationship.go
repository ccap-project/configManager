// Code generated by go-swagger; DO NOT EDIT.

// Copyright (c) 2016, 2017, 2018 Alexandre Biancalana <ale@biancalanas.net>.
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//     * Redistributions of source code must retain the above copyright
//       notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above copyright
//       notice, this list of conditions and the following disclaimer in the
//       documentation and/or other materials provided with the distribution.
//     * Neither the name of the <organization> nor the
//       names of its contributors may be used to endorse or promote products
//       derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL <COPYRIGHT HOLDER> BE LIABLE FOR ANY
// DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//

package loadbalancer

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"

	"configManager/models"
)

// DeleteLoadbalancerRelationshipHandlerFunc turns a function with the right signature into a delete loadbalancer relationship handler
type DeleteLoadbalancerRelationshipHandlerFunc func(DeleteLoadbalancerRelationshipParams, *models.Customer) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteLoadbalancerRelationshipHandlerFunc) Handle(params DeleteLoadbalancerRelationshipParams, principal *models.Customer) middleware.Responder {
	return fn(params, principal)
}

// DeleteLoadbalancerRelationshipHandler interface for that can handle valid delete loadbalancer relationship params
type DeleteLoadbalancerRelationshipHandler interface {
	Handle(DeleteLoadbalancerRelationshipParams, *models.Customer) middleware.Responder
}

// NewDeleteLoadbalancerRelationship creates a new http.Handler for the delete loadbalancer relationship operation
func NewDeleteLoadbalancerRelationship(ctx *middleware.Context, handler DeleteLoadbalancerRelationshipHandler) *DeleteLoadbalancerRelationship {
	return &DeleteLoadbalancerRelationship{Context: ctx, Handler: handler}
}

/*DeleteLoadbalancerRelationship swagger:route DELETE /cell/{cell_id}/loadbalancer/{loadbalancer_id}/connect_to/{listener_id} loadbalancer deleteLoadbalancerRelationship

delete loadbalancer relationship

*/
type DeleteLoadbalancerRelationship struct {
	Context *middleware.Context
	Handler DeleteLoadbalancerRelationshipHandler
}

func (o *DeleteLoadbalancerRelationship) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewDeleteLoadbalancerRelationshipParams()

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
