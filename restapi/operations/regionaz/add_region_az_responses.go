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

package regionaz

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"configManager/models"
)

// AddRegionAZOKCode is the HTTP code returned for type AddRegionAZOK
const AddRegionAZOKCode int = 200

/*AddRegionAZOK Already exists

swagger:response addRegionAZOK
*/
type AddRegionAZOK struct {

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewAddRegionAZOK creates AddRegionAZOK with default headers values
func NewAddRegionAZOK() *AddRegionAZOK {
	return &AddRegionAZOK{}
}

// WithPayload adds the payload to the add region a z o k response
func (o *AddRegionAZOK) WithPayload(payload *models.APIResponse) *AddRegionAZOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add region a z o k response
func (o *AddRegionAZOK) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddRegionAZOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// AddRegionAZCreatedCode is the HTTP code returned for type AddRegionAZCreated
const AddRegionAZCreatedCode int = 201

/*AddRegionAZCreated Created

swagger:response addRegionAZCreated
*/
type AddRegionAZCreated struct {

	/*
	  In: Body
	*/
	Payload models.ULID `json:"body,omitempty"`
}

// NewAddRegionAZCreated creates AddRegionAZCreated with default headers values
func NewAddRegionAZCreated() *AddRegionAZCreated {
	return &AddRegionAZCreated{}
}

// WithPayload adds the payload to the add region a z created response
func (o *AddRegionAZCreated) WithPayload(payload models.ULID) *AddRegionAZCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add region a z created response
func (o *AddRegionAZCreated) SetPayload(payload models.ULID) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddRegionAZCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// AddRegionAZMethodNotAllowedCode is the HTTP code returned for type AddRegionAZMethodNotAllowed
const AddRegionAZMethodNotAllowedCode int = 405

/*AddRegionAZMethodNotAllowed Invalid input

swagger:response addRegionAZMethodNotAllowed
*/
type AddRegionAZMethodNotAllowed struct {
}

// NewAddRegionAZMethodNotAllowed creates AddRegionAZMethodNotAllowed with default headers values
func NewAddRegionAZMethodNotAllowed() *AddRegionAZMethodNotAllowed {
	return &AddRegionAZMethodNotAllowed{}
}

// WriteResponse to the client
func (o *AddRegionAZMethodNotAllowed) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(405)
}

// AddRegionAZInternalServerErrorCode is the HTTP code returned for type AddRegionAZInternalServerError
const AddRegionAZInternalServerErrorCode int = 500

/*AddRegionAZInternalServerError Internal error

swagger:response addRegionAZInternalServerError
*/
type AddRegionAZInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewAddRegionAZInternalServerError creates AddRegionAZInternalServerError with default headers values
func NewAddRegionAZInternalServerError() *AddRegionAZInternalServerError {
	return &AddRegionAZInternalServerError{}
}

// WithPayload adds the payload to the add region a z internal server error response
func (o *AddRegionAZInternalServerError) WithPayload(payload *models.APIResponse) *AddRegionAZInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add region a z internal server error response
func (o *AddRegionAZInternalServerError) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddRegionAZInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
