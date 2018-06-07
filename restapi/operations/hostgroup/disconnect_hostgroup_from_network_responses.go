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

// DisconnectHostgroupFromNetworkOKCode is the HTTP code returned for type DisconnectHostgroupFromNetworkOK
const DisconnectHostgroupFromNetworkOKCode int = 200

/*DisconnectHostgroupFromNetworkOK successful operation

swagger:response disconnectHostgroupFromNetworkOK
*/
type DisconnectHostgroupFromNetworkOK struct {
}

// NewDisconnectHostgroupFromNetworkOK creates DisconnectHostgroupFromNetworkOK with default headers values
func NewDisconnectHostgroupFromNetworkOK() *DisconnectHostgroupFromNetworkOK {
	return &DisconnectHostgroupFromNetworkOK{}
}

// WriteResponse to the client
func (o *DisconnectHostgroupFromNetworkOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(200)
}

// DisconnectHostgroupFromNetworkBadRequestCode is the HTTP code returned for type DisconnectHostgroupFromNetworkBadRequest
const DisconnectHostgroupFromNetworkBadRequestCode int = 400

/*DisconnectHostgroupFromNetworkBadRequest Invalid cell id, component id or hostgroup id

swagger:response disconnectHostgroupFromNetworkBadRequest
*/
type DisconnectHostgroupFromNetworkBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewDisconnectHostgroupFromNetworkBadRequest creates DisconnectHostgroupFromNetworkBadRequest with default headers values
func NewDisconnectHostgroupFromNetworkBadRequest() *DisconnectHostgroupFromNetworkBadRequest {
	return &DisconnectHostgroupFromNetworkBadRequest{}
}

// WithPayload adds the payload to the disconnect hostgroup from network bad request response
func (o *DisconnectHostgroupFromNetworkBadRequest) WithPayload(payload *models.APIResponse) *DisconnectHostgroupFromNetworkBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the disconnect hostgroup from network bad request response
func (o *DisconnectHostgroupFromNetworkBadRequest) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DisconnectHostgroupFromNetworkBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// DisconnectHostgroupFromNetworkNotFoundCode is the HTTP code returned for type DisconnectHostgroupFromNetworkNotFound
const DisconnectHostgroupFromNetworkNotFoundCode int = 404

/*DisconnectHostgroupFromNetworkNotFound Hostgroup does not exists

swagger:response disconnectHostgroupFromNetworkNotFound
*/
type DisconnectHostgroupFromNetworkNotFound struct {
}

// NewDisconnectHostgroupFromNetworkNotFound creates DisconnectHostgroupFromNetworkNotFound with default headers values
func NewDisconnectHostgroupFromNetworkNotFound() *DisconnectHostgroupFromNetworkNotFound {
	return &DisconnectHostgroupFromNetworkNotFound{}
}

// WriteResponse to the client
func (o *DisconnectHostgroupFromNetworkNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}

// DisconnectHostgroupFromNetworkInternalServerErrorCode is the HTTP code returned for type DisconnectHostgroupFromNetworkInternalServerError
const DisconnectHostgroupFromNetworkInternalServerErrorCode int = 500

/*DisconnectHostgroupFromNetworkInternalServerError Internal error

swagger:response disconnectHostgroupFromNetworkInternalServerError
*/
type DisconnectHostgroupFromNetworkInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewDisconnectHostgroupFromNetworkInternalServerError creates DisconnectHostgroupFromNetworkInternalServerError with default headers values
func NewDisconnectHostgroupFromNetworkInternalServerError() *DisconnectHostgroupFromNetworkInternalServerError {
	return &DisconnectHostgroupFromNetworkInternalServerError{}
}

// WithPayload adds the payload to the disconnect hostgroup from network internal server error response
func (o *DisconnectHostgroupFromNetworkInternalServerError) WithPayload(payload *models.APIResponse) *DisconnectHostgroupFromNetworkInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the disconnect hostgroup from network internal server error response
func (o *DisconnectHostgroupFromNetworkInternalServerError) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DisconnectHostgroupFromNetworkInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}