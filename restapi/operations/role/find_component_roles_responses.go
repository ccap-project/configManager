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

package role

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"configManager/models"
)

// FindComponentRolesOKCode is the HTTP code returned for type FindComponentRolesOK
const FindComponentRolesOKCode int = 200

/*FindComponentRolesOK successful operation

swagger:response findComponentRolesOK
*/
type FindComponentRolesOK struct {

	/*
	  In: Body
	*/
	Payload models.FindComponentRolesOKBody `json:"body,omitempty"`
}

// NewFindComponentRolesOK creates FindComponentRolesOK with default headers values
func NewFindComponentRolesOK() *FindComponentRolesOK {
	return &FindComponentRolesOK{}
}

// WithPayload adds the payload to the find component roles o k response
func (o *FindComponentRolesOK) WithPayload(payload models.FindComponentRolesOKBody) *FindComponentRolesOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the find component roles o k response
func (o *FindComponentRolesOK) SetPayload(payload models.FindComponentRolesOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *FindComponentRolesOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		payload = make(models.FindComponentRolesOKBody, 0, 50)
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// FindComponentRolesBadRequestCode is the HTTP code returned for type FindComponentRolesBadRequest
const FindComponentRolesBadRequestCode int = 400

/*FindComponentRolesBadRequest Invalid cell id or role id

swagger:response findComponentRolesBadRequest
*/
type FindComponentRolesBadRequest struct {

	/*
	  In: Body
	*/
	Payload models.APIResponse `json:"body,omitempty"`
}

// NewFindComponentRolesBadRequest creates FindComponentRolesBadRequest with default headers values
func NewFindComponentRolesBadRequest() *FindComponentRolesBadRequest {
	return &FindComponentRolesBadRequest{}
}

// WithPayload adds the payload to the find component roles bad request response
func (o *FindComponentRolesBadRequest) WithPayload(payload models.APIResponse) *FindComponentRolesBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the find component roles bad request response
func (o *FindComponentRolesBadRequest) SetPayload(payload models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *FindComponentRolesBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// FindComponentRolesInternalServerErrorCode is the HTTP code returned for type FindComponentRolesInternalServerError
const FindComponentRolesInternalServerErrorCode int = 500

/*FindComponentRolesInternalServerError Internal error

swagger:response findComponentRolesInternalServerError
*/
type FindComponentRolesInternalServerError struct {

	/*
	  In: Body
	*/
	Payload models.APIResponse `json:"body,omitempty"`
}

// NewFindComponentRolesInternalServerError creates FindComponentRolesInternalServerError with default headers values
func NewFindComponentRolesInternalServerError() *FindComponentRolesInternalServerError {
	return &FindComponentRolesInternalServerError{}
}

// WithPayload adds the payload to the find component roles internal server error response
func (o *FindComponentRolesInternalServerError) WithPayload(payload models.APIResponse) *FindComponentRolesInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the find component roles internal server error response
func (o *FindComponentRolesInternalServerError) SetPayload(payload models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *FindComponentRolesInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}
