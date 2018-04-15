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

package customer

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"configManager/models"
)

// DeleteCustomerOKCode is the HTTP code returned for type DeleteCustomerOK
const DeleteCustomerOKCode int = 200

/*DeleteCustomerOK successful operation

swagger:response deleteCustomerOK
*/
type DeleteCustomerOK struct {
}

// NewDeleteCustomerOK creates DeleteCustomerOK with default headers values
func NewDeleteCustomerOK() *DeleteCustomerOK {
	return &DeleteCustomerOK{}
}

// WriteResponse to the client
func (o *DeleteCustomerOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(200)
}

// DeleteCustomerBadRequestCode is the HTTP code returned for type DeleteCustomerBadRequest
const DeleteCustomerBadRequestCode int = 400

/*DeleteCustomerBadRequest Invalid ID supplied

swagger:response deleteCustomerBadRequest
*/
type DeleteCustomerBadRequest struct {
}

// NewDeleteCustomerBadRequest creates DeleteCustomerBadRequest with default headers values
func NewDeleteCustomerBadRequest() *DeleteCustomerBadRequest {
	return &DeleteCustomerBadRequest{}
}

// WriteResponse to the client
func (o *DeleteCustomerBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(400)
}

// DeleteCustomerInternalServerErrorCode is the HTTP code returned for type DeleteCustomerInternalServerError
const DeleteCustomerInternalServerErrorCode int = 500

/*DeleteCustomerInternalServerError Internal error

swagger:response deleteCustomerInternalServerError
*/
type DeleteCustomerInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewDeleteCustomerInternalServerError creates DeleteCustomerInternalServerError with default headers values
func NewDeleteCustomerInternalServerError() *DeleteCustomerInternalServerError {
	return &DeleteCustomerInternalServerError{}
}

// WithPayload adds the payload to the delete customer internal server error response
func (o *DeleteCustomerInternalServerError) WithPayload(payload *models.APIResponse) *DeleteCustomerInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete customer internal server error response
func (o *DeleteCustomerInternalServerError) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteCustomerInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
