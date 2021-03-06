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

// DeleteComponentRoleOKCode is the HTTP code returned for type DeleteComponentRoleOK
const DeleteComponentRoleOKCode int = 200

/*DeleteComponentRoleOK successful operation

swagger:response deleteComponentRoleOK
*/
type DeleteComponentRoleOK struct {
}

// NewDeleteComponentRoleOK creates DeleteComponentRoleOK with default headers values
func NewDeleteComponentRoleOK() *DeleteComponentRoleOK {
	return &DeleteComponentRoleOK{}
}

// WriteResponse to the client
func (o *DeleteComponentRoleOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(200)
}

// DeleteComponentRoleBadRequestCode is the HTTP code returned for type DeleteComponentRoleBadRequest
const DeleteComponentRoleBadRequestCode int = 400

/*DeleteComponentRoleBadRequest Invalid cell id or role id

swagger:response deleteComponentRoleBadRequest
*/
type DeleteComponentRoleBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewDeleteComponentRoleBadRequest creates DeleteComponentRoleBadRequest with default headers values
func NewDeleteComponentRoleBadRequest() *DeleteComponentRoleBadRequest {
	return &DeleteComponentRoleBadRequest{}
}

// WithPayload adds the payload to the delete component role bad request response
func (o *DeleteComponentRoleBadRequest) WithPayload(payload *models.APIResponse) *DeleteComponentRoleBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete component role bad request response
func (o *DeleteComponentRoleBadRequest) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteComponentRoleBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// DeleteComponentRoleNotFoundCode is the HTTP code returned for type DeleteComponentRoleNotFound
const DeleteComponentRoleNotFoundCode int = 404

/*DeleteComponentRoleNotFound Role does not exists

swagger:response deleteComponentRoleNotFound
*/
type DeleteComponentRoleNotFound struct {
}

// NewDeleteComponentRoleNotFound creates DeleteComponentRoleNotFound with default headers values
func NewDeleteComponentRoleNotFound() *DeleteComponentRoleNotFound {
	return &DeleteComponentRoleNotFound{}
}

// WriteResponse to the client
func (o *DeleteComponentRoleNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}

// DeleteComponentRoleInternalServerErrorCode is the HTTP code returned for type DeleteComponentRoleInternalServerError
const DeleteComponentRoleInternalServerErrorCode int = 500

/*DeleteComponentRoleInternalServerError Internal error

swagger:response deleteComponentRoleInternalServerError
*/
type DeleteComponentRoleInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewDeleteComponentRoleInternalServerError creates DeleteComponentRoleInternalServerError with default headers values
func NewDeleteComponentRoleInternalServerError() *DeleteComponentRoleInternalServerError {
	return &DeleteComponentRoleInternalServerError{}
}

// WithPayload adds the payload to the delete component role internal server error response
func (o *DeleteComponentRoleInternalServerError) WithPayload(payload *models.APIResponse) *DeleteComponentRoleInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete component role internal server error response
func (o *DeleteComponentRoleInternalServerError) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteComponentRoleInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
