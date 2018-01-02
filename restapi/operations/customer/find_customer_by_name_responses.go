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

// FindCustomerByNameOKCode is the HTTP code returned for type FindCustomerByNameOK
const FindCustomerByNameOKCode int = 200

/*FindCustomerByNameOK successful operation

swagger:response findCustomerByNameOK
*/
type FindCustomerByNameOK struct {

	/*
	  In: Body
	*/
	Payload *models.Customer `json:"body,omitempty"`
}

// NewFindCustomerByNameOK creates FindCustomerByNameOK with default headers values
func NewFindCustomerByNameOK() *FindCustomerByNameOK {
	return &FindCustomerByNameOK{}
}

// WithPayload adds the payload to the find customer by name o k response
func (o *FindCustomerByNameOK) WithPayload(payload *models.Customer) *FindCustomerByNameOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the find customer by name o k response
func (o *FindCustomerByNameOK) SetPayload(payload *models.Customer) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *FindCustomerByNameOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// FindCustomerByNameBadRequestCode is the HTTP code returned for type FindCustomerByNameBadRequest
const FindCustomerByNameBadRequestCode int = 400

/*FindCustomerByNameBadRequest Invalid customer name

swagger:response findCustomerByNameBadRequest
*/
type FindCustomerByNameBadRequest struct {
}

// NewFindCustomerByNameBadRequest creates FindCustomerByNameBadRequest with default headers values
func NewFindCustomerByNameBadRequest() *FindCustomerByNameBadRequest {
	return &FindCustomerByNameBadRequest{}
}

// WriteResponse to the client
func (o *FindCustomerByNameBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(400)
}

// FindCustomerByNameNotFoundCode is the HTTP code returned for type FindCustomerByNameNotFound
const FindCustomerByNameNotFoundCode int = 404

/*FindCustomerByNameNotFound Customer not found

swagger:response findCustomerByNameNotFound
*/
type FindCustomerByNameNotFound struct {
}

// NewFindCustomerByNameNotFound creates FindCustomerByNameNotFound with default headers values
func NewFindCustomerByNameNotFound() *FindCustomerByNameNotFound {
	return &FindCustomerByNameNotFound{}
}

// WriteResponse to the client
func (o *FindCustomerByNameNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}
