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

package provider

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"configManager/models"
)

// AddProviderCreatedCode is the HTTP code returned for type AddProviderCreated
const AddProviderCreatedCode int = 201

/*AddProviderCreated Created

swagger:response addProviderCreated
*/
type AddProviderCreated struct {

	/*
	  In: Body
	*/
	Payload models.ULID `json:"body,omitempty"`
}

// NewAddProviderCreated creates AddProviderCreated with default headers values
func NewAddProviderCreated() *AddProviderCreated {
	return &AddProviderCreated{}
}

// WithPayload adds the payload to the add provider created response
func (o *AddProviderCreated) WithPayload(payload models.ULID) *AddProviderCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add provider created response
func (o *AddProviderCreated) SetPayload(payload models.ULID) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddProviderCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// AddProviderConflictCode is the HTTP code returned for type AddProviderConflict
const AddProviderConflictCode int = 409

/*AddProviderConflict Already exists

swagger:response addProviderConflict
*/
type AddProviderConflict struct {

	/*
	  In: Body
	*/
	Payload models.APIResponse `json:"body,omitempty"`
}

// NewAddProviderConflict creates AddProviderConflict with default headers values
func NewAddProviderConflict() *AddProviderConflict {
	return &AddProviderConflict{}
}

// WithPayload adds the payload to the add provider conflict response
func (o *AddProviderConflict) WithPayload(payload models.APIResponse) *AddProviderConflict {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add provider conflict response
func (o *AddProviderConflict) SetPayload(payload models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddProviderConflict) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(409)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// AddProviderInternalServerErrorCode is the HTTP code returned for type AddProviderInternalServerError
const AddProviderInternalServerErrorCode int = 500

/*AddProviderInternalServerError Internal error

swagger:response addProviderInternalServerError
*/
type AddProviderInternalServerError struct {

	/*
	  In: Body
	*/
	Payload models.APIResponse `json:"body,omitempty"`
}

// NewAddProviderInternalServerError creates AddProviderInternalServerError with default headers values
func NewAddProviderInternalServerError() *AddProviderInternalServerError {
	return &AddProviderInternalServerError{}
}

// WithPayload adds the payload to the add provider internal server error response
func (o *AddProviderInternalServerError) WithPayload(payload models.APIResponse) *AddProviderInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add provider internal server error response
func (o *AddProviderInternalServerError) SetPayload(payload models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddProviderInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}
