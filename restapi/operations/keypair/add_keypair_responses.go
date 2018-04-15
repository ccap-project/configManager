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

package keypair

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"configManager/models"
)

// AddKeypairCreatedCode is the HTTP code returned for type AddKeypairCreated
const AddKeypairCreatedCode int = 201

/*AddKeypairCreated Created

swagger:response addKeypairCreated
*/
type AddKeypairCreated struct {

	/*
	  In: Body
	*/
	Payload models.ULID `json:"body,omitempty"`
}

// NewAddKeypairCreated creates AddKeypairCreated with default headers values
func NewAddKeypairCreated() *AddKeypairCreated {
	return &AddKeypairCreated{}
}

// WithPayload adds the payload to the add keypair created response
func (o *AddKeypairCreated) WithPayload(payload models.ULID) *AddKeypairCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add keypair created response
func (o *AddKeypairCreated) SetPayload(payload models.ULID) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddKeypairCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// AddKeypairMethodNotAllowedCode is the HTTP code returned for type AddKeypairMethodNotAllowed
const AddKeypairMethodNotAllowedCode int = 405

/*AddKeypairMethodNotAllowed Invalid input

swagger:response addKeypairMethodNotAllowed
*/
type AddKeypairMethodNotAllowed struct {

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewAddKeypairMethodNotAllowed creates AddKeypairMethodNotAllowed with default headers values
func NewAddKeypairMethodNotAllowed() *AddKeypairMethodNotAllowed {
	return &AddKeypairMethodNotAllowed{}
}

// WithPayload adds the payload to the add keypair method not allowed response
func (o *AddKeypairMethodNotAllowed) WithPayload(payload *models.APIResponse) *AddKeypairMethodNotAllowed {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add keypair method not allowed response
func (o *AddKeypairMethodNotAllowed) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddKeypairMethodNotAllowed) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(405)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// AddKeypairConflictCode is the HTTP code returned for type AddKeypairConflict
const AddKeypairConflictCode int = 409

/*AddKeypairConflict Already exists

swagger:response addKeypairConflict
*/
type AddKeypairConflict struct {

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewAddKeypairConflict creates AddKeypairConflict with default headers values
func NewAddKeypairConflict() *AddKeypairConflict {
	return &AddKeypairConflict{}
}

// WithPayload adds the payload to the add keypair conflict response
func (o *AddKeypairConflict) WithPayload(payload *models.APIResponse) *AddKeypairConflict {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add keypair conflict response
func (o *AddKeypairConflict) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddKeypairConflict) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(409)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// AddKeypairInternalServerErrorCode is the HTTP code returned for type AddKeypairInternalServerError
const AddKeypairInternalServerErrorCode int = 500

/*AddKeypairInternalServerError Internal error

swagger:response addKeypairInternalServerError
*/
type AddKeypairInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewAddKeypairInternalServerError creates AddKeypairInternalServerError with default headers values
func NewAddKeypairInternalServerError() *AddKeypairInternalServerError {
	return &AddKeypairInternalServerError{}
}

// WithPayload adds the payload to the add keypair internal server error response
func (o *AddKeypairInternalServerError) WithPayload(payload *models.APIResponse) *AddKeypairInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add keypair internal server error response
func (o *AddKeypairInternalServerError) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddKeypairInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
