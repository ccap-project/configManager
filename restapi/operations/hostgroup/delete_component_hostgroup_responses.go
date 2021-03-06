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

package hostgroup

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"configManager/models"
)

// DeleteComponentHostgroupOKCode is the HTTP code returned for type DeleteComponentHostgroupOK
const DeleteComponentHostgroupOKCode int = 200

/*DeleteComponentHostgroupOK successful operation

swagger:response deleteComponentHostgroupOK
*/
type DeleteComponentHostgroupOK struct {
}

// NewDeleteComponentHostgroupOK creates DeleteComponentHostgroupOK with default headers values
func NewDeleteComponentHostgroupOK() *DeleteComponentHostgroupOK {
	return &DeleteComponentHostgroupOK{}
}

// WriteResponse to the client
func (o *DeleteComponentHostgroupOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(200)
}

// DeleteComponentHostgroupBadRequestCode is the HTTP code returned for type DeleteComponentHostgroupBadRequest
const DeleteComponentHostgroupBadRequestCode int = 400

/*DeleteComponentHostgroupBadRequest Invalid cell id, component id or hostgroup id

swagger:response deleteComponentHostgroupBadRequest
*/
type DeleteComponentHostgroupBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewDeleteComponentHostgroupBadRequest creates DeleteComponentHostgroupBadRequest with default headers values
func NewDeleteComponentHostgroupBadRequest() *DeleteComponentHostgroupBadRequest {
	return &DeleteComponentHostgroupBadRequest{}
}

// WithPayload adds the payload to the delete component hostgroup bad request response
func (o *DeleteComponentHostgroupBadRequest) WithPayload(payload *models.APIResponse) *DeleteComponentHostgroupBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete component hostgroup bad request response
func (o *DeleteComponentHostgroupBadRequest) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteComponentHostgroupBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// DeleteComponentHostgroupNotFoundCode is the HTTP code returned for type DeleteComponentHostgroupNotFound
const DeleteComponentHostgroupNotFoundCode int = 404

/*DeleteComponentHostgroupNotFound Hostgroup does not exists

swagger:response deleteComponentHostgroupNotFound
*/
type DeleteComponentHostgroupNotFound struct {
}

// NewDeleteComponentHostgroupNotFound creates DeleteComponentHostgroupNotFound with default headers values
func NewDeleteComponentHostgroupNotFound() *DeleteComponentHostgroupNotFound {
	return &DeleteComponentHostgroupNotFound{}
}

// WriteResponse to the client
func (o *DeleteComponentHostgroupNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}

// DeleteComponentHostgroupInternalServerErrorCode is the HTTP code returned for type DeleteComponentHostgroupInternalServerError
const DeleteComponentHostgroupInternalServerErrorCode int = 500

/*DeleteComponentHostgroupInternalServerError Internal error

swagger:response deleteComponentHostgroupInternalServerError
*/
type DeleteComponentHostgroupInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewDeleteComponentHostgroupInternalServerError creates DeleteComponentHostgroupInternalServerError with default headers values
func NewDeleteComponentHostgroupInternalServerError() *DeleteComponentHostgroupInternalServerError {
	return &DeleteComponentHostgroupInternalServerError{}
}

// WithPayload adds the payload to the delete component hostgroup internal server error response
func (o *DeleteComponentHostgroupInternalServerError) WithPayload(payload *models.APIResponse) *DeleteComponentHostgroupInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete component hostgroup internal server error response
func (o *DeleteComponentHostgroupInternalServerError) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteComponentHostgroupInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
