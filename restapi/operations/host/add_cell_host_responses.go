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

package host

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"configManager/models"
)

// AddCellHostCreatedCode is the HTTP code returned for type AddCellHostCreated
const AddCellHostCreatedCode int = 201

/*AddCellHostCreated Created

swagger:response addCellHostCreated
*/
type AddCellHostCreated struct {

	/*
	  In: Body
	*/
	Payload models.ULID `json:"body,omitempty"`
}

// NewAddCellHostCreated creates AddCellHostCreated with default headers values
func NewAddCellHostCreated() *AddCellHostCreated {
	return &AddCellHostCreated{}
}

// WithPayload adds the payload to the add cell host created response
func (o *AddCellHostCreated) WithPayload(payload models.ULID) *AddCellHostCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add cell host created response
func (o *AddCellHostCreated) SetPayload(payload models.ULID) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddCellHostCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// AddCellHostMethodNotAllowedCode is the HTTP code returned for type AddCellHostMethodNotAllowed
const AddCellHostMethodNotAllowedCode int = 405

/*AddCellHostMethodNotAllowed Invalid input

swagger:response addCellHostMethodNotAllowed
*/
type AddCellHostMethodNotAllowed struct {

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewAddCellHostMethodNotAllowed creates AddCellHostMethodNotAllowed with default headers values
func NewAddCellHostMethodNotAllowed() *AddCellHostMethodNotAllowed {
	return &AddCellHostMethodNotAllowed{}
}

// WithPayload adds the payload to the add cell host method not allowed response
func (o *AddCellHostMethodNotAllowed) WithPayload(payload *models.APIResponse) *AddCellHostMethodNotAllowed {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add cell host method not allowed response
func (o *AddCellHostMethodNotAllowed) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddCellHostMethodNotAllowed) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(405)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// AddCellHostConflictCode is the HTTP code returned for type AddCellHostConflict
const AddCellHostConflictCode int = 409

/*AddCellHostConflict Already exists

swagger:response addCellHostConflict
*/
type AddCellHostConflict struct {

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewAddCellHostConflict creates AddCellHostConflict with default headers values
func NewAddCellHostConflict() *AddCellHostConflict {
	return &AddCellHostConflict{}
}

// WithPayload adds the payload to the add cell host conflict response
func (o *AddCellHostConflict) WithPayload(payload *models.APIResponse) *AddCellHostConflict {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add cell host conflict response
func (o *AddCellHostConflict) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddCellHostConflict) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(409)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// AddCellHostInternalServerErrorCode is the HTTP code returned for type AddCellHostInternalServerError
const AddCellHostInternalServerErrorCode int = 500

/*AddCellHostInternalServerError Internal error

swagger:response addCellHostInternalServerError
*/
type AddCellHostInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewAddCellHostInternalServerError creates AddCellHostInternalServerError with default headers values
func NewAddCellHostInternalServerError() *AddCellHostInternalServerError {
	return &AddCellHostInternalServerError{}
}

// WithPayload adds the payload to the add cell host internal server error response
func (o *AddCellHostInternalServerError) WithPayload(payload *models.APIResponse) *AddCellHostInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add cell host internal server error response
func (o *AddCellHostInternalServerError) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddCellHostInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
