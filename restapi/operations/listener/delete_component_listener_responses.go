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

package listener

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"configManager/models"
)

// DeleteComponentListenerOKCode is the HTTP code returned for type DeleteComponentListenerOK
const DeleteComponentListenerOKCode int = 200

/*DeleteComponentListenerOK successful operation

swagger:response deleteComponentListenerOK
*/
type DeleteComponentListenerOK struct {
}

// NewDeleteComponentListenerOK creates DeleteComponentListenerOK with default headers values
func NewDeleteComponentListenerOK() *DeleteComponentListenerOK {
	return &DeleteComponentListenerOK{}
}

// WriteResponse to the client
func (o *DeleteComponentListenerOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(200)
}

// DeleteComponentListenerBadRequestCode is the HTTP code returned for type DeleteComponentListenerBadRequest
const DeleteComponentListenerBadRequestCode int = 400

/*DeleteComponentListenerBadRequest Invalid cell id, component id or listener id

swagger:response deleteComponentListenerBadRequest
*/
type DeleteComponentListenerBadRequest struct {

	/*
	  In: Body
	*/
	Payload models.APIResponse `json:"body,omitempty"`
}

// NewDeleteComponentListenerBadRequest creates DeleteComponentListenerBadRequest with default headers values
func NewDeleteComponentListenerBadRequest() *DeleteComponentListenerBadRequest {
	return &DeleteComponentListenerBadRequest{}
}

// WithPayload adds the payload to the delete component listener bad request response
func (o *DeleteComponentListenerBadRequest) WithPayload(payload models.APIResponse) *DeleteComponentListenerBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete component listener bad request response
func (o *DeleteComponentListenerBadRequest) SetPayload(payload models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteComponentListenerBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// DeleteComponentListenerNotFoundCode is the HTTP code returned for type DeleteComponentListenerNotFound
const DeleteComponentListenerNotFoundCode int = 404

/*DeleteComponentListenerNotFound Hostgroup does not exists

swagger:response deleteComponentListenerNotFound
*/
type DeleteComponentListenerNotFound struct {
}

// NewDeleteComponentListenerNotFound creates DeleteComponentListenerNotFound with default headers values
func NewDeleteComponentListenerNotFound() *DeleteComponentListenerNotFound {
	return &DeleteComponentListenerNotFound{}
}

// WriteResponse to the client
func (o *DeleteComponentListenerNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}

// DeleteComponentListenerInternalServerErrorCode is the HTTP code returned for type DeleteComponentListenerInternalServerError
const DeleteComponentListenerInternalServerErrorCode int = 500

/*DeleteComponentListenerInternalServerError Internal error

swagger:response deleteComponentListenerInternalServerError
*/
type DeleteComponentListenerInternalServerError struct {

	/*
	  In: Body
	*/
	Payload models.APIResponse `json:"body,omitempty"`
}

// NewDeleteComponentListenerInternalServerError creates DeleteComponentListenerInternalServerError with default headers values
func NewDeleteComponentListenerInternalServerError() *DeleteComponentListenerInternalServerError {
	return &DeleteComponentListenerInternalServerError{}
}

// WithPayload adds the payload to the delete component listener internal server error response
func (o *DeleteComponentListenerInternalServerError) WithPayload(payload models.APIResponse) *DeleteComponentListenerInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete component listener internal server error response
func (o *DeleteComponentListenerInternalServerError) SetPayload(payload models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteComponentListenerInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}