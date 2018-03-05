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
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"configManager/models"
)

// DeleteLoadbalancerRelationshipOKCode is the HTTP code returned for type DeleteLoadbalancerRelationshipOK
const DeleteLoadbalancerRelationshipOKCode int = 200

/*DeleteLoadbalancerRelationshipOK successful operation

swagger:response deleteLoadbalancerRelationshipOK
*/
type DeleteLoadbalancerRelationshipOK struct {
}

// NewDeleteLoadbalancerRelationshipOK creates DeleteLoadbalancerRelationshipOK with default headers values
func NewDeleteLoadbalancerRelationshipOK() *DeleteLoadbalancerRelationshipOK {
	return &DeleteLoadbalancerRelationshipOK{}
}

// WriteResponse to the client
func (o *DeleteLoadbalancerRelationshipOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(200)
}

// DeleteLoadbalancerRelationshipNotFoundCode is the HTTP code returned for type DeleteLoadbalancerRelationshipNotFound
const DeleteLoadbalancerRelationshipNotFoundCode int = 404

/*DeleteLoadbalancerRelationshipNotFound loadbalancer or entity not found

swagger:response deleteLoadbalancerRelationshipNotFound
*/
type DeleteLoadbalancerRelationshipNotFound struct {

	/*
	  In: Body
	*/
	Payload models.APIResponse `json:"body,omitempty"`
}

// NewDeleteLoadbalancerRelationshipNotFound creates DeleteLoadbalancerRelationshipNotFound with default headers values
func NewDeleteLoadbalancerRelationshipNotFound() *DeleteLoadbalancerRelationshipNotFound {
	return &DeleteLoadbalancerRelationshipNotFound{}
}

// WithPayload adds the payload to the delete loadbalancer relationship not found response
func (o *DeleteLoadbalancerRelationshipNotFound) WithPayload(payload models.APIResponse) *DeleteLoadbalancerRelationshipNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete loadbalancer relationship not found response
func (o *DeleteLoadbalancerRelationshipNotFound) SetPayload(payload models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteLoadbalancerRelationshipNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// DeleteLoadbalancerRelationshipInternalServerErrorCode is the HTTP code returned for type DeleteLoadbalancerRelationshipInternalServerError
const DeleteLoadbalancerRelationshipInternalServerErrorCode int = 500

/*DeleteLoadbalancerRelationshipInternalServerError Internal error

swagger:response deleteLoadbalancerRelationshipInternalServerError
*/
type DeleteLoadbalancerRelationshipInternalServerError struct {

	/*
	  In: Body
	*/
	Payload models.APIResponse `json:"body,omitempty"`
}

// NewDeleteLoadbalancerRelationshipInternalServerError creates DeleteLoadbalancerRelationshipInternalServerError with default headers values
func NewDeleteLoadbalancerRelationshipInternalServerError() *DeleteLoadbalancerRelationshipInternalServerError {
	return &DeleteLoadbalancerRelationshipInternalServerError{}
}

// WithPayload adds the payload to the delete loadbalancer relationship internal server error response
func (o *DeleteLoadbalancerRelationshipInternalServerError) WithPayload(payload models.APIResponse) *DeleteLoadbalancerRelationshipInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete loadbalancer relationship internal server error response
func (o *DeleteLoadbalancerRelationshipInternalServerError) SetPayload(payload models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteLoadbalancerRelationshipInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}
