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

package loadbalancer

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"configManager/models"
)

// AddLoadbalancerCreatedCode is the HTTP code returned for type AddLoadbalancerCreated
const AddLoadbalancerCreatedCode int = 201

/*AddLoadbalancerCreated Created

swagger:response addLoadbalancerCreated
*/
type AddLoadbalancerCreated struct {

	/*
	  In: Body
	*/
	Payload models.ULID `json:"body,omitempty"`
}

// NewAddLoadbalancerCreated creates AddLoadbalancerCreated with default headers values
func NewAddLoadbalancerCreated() *AddLoadbalancerCreated {
	return &AddLoadbalancerCreated{}
}

// WithPayload adds the payload to the add loadbalancer created response
func (o *AddLoadbalancerCreated) WithPayload(payload models.ULID) *AddLoadbalancerCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add loadbalancer created response
func (o *AddLoadbalancerCreated) SetPayload(payload models.ULID) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddLoadbalancerCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// AddLoadbalancerMethodNotAllowedCode is the HTTP code returned for type AddLoadbalancerMethodNotAllowed
const AddLoadbalancerMethodNotAllowedCode int = 405

/*AddLoadbalancerMethodNotAllowed Invalid input

swagger:response addLoadbalancerMethodNotAllowed
*/
type AddLoadbalancerMethodNotAllowed struct {

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewAddLoadbalancerMethodNotAllowed creates AddLoadbalancerMethodNotAllowed with default headers values
func NewAddLoadbalancerMethodNotAllowed() *AddLoadbalancerMethodNotAllowed {
	return &AddLoadbalancerMethodNotAllowed{}
}

// WithPayload adds the payload to the add loadbalancer method not allowed response
func (o *AddLoadbalancerMethodNotAllowed) WithPayload(payload *models.APIResponse) *AddLoadbalancerMethodNotAllowed {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add loadbalancer method not allowed response
func (o *AddLoadbalancerMethodNotAllowed) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddLoadbalancerMethodNotAllowed) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(405)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// AddLoadbalancerConflictCode is the HTTP code returned for type AddLoadbalancerConflict
const AddLoadbalancerConflictCode int = 409

/*AddLoadbalancerConflict Already exists

swagger:response addLoadbalancerConflict
*/
type AddLoadbalancerConflict struct {

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewAddLoadbalancerConflict creates AddLoadbalancerConflict with default headers values
func NewAddLoadbalancerConflict() *AddLoadbalancerConflict {
	return &AddLoadbalancerConflict{}
}

// WithPayload adds the payload to the add loadbalancer conflict response
func (o *AddLoadbalancerConflict) WithPayload(payload *models.APIResponse) *AddLoadbalancerConflict {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add loadbalancer conflict response
func (o *AddLoadbalancerConflict) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddLoadbalancerConflict) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(409)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// AddLoadbalancerInternalServerErrorCode is the HTTP code returned for type AddLoadbalancerInternalServerError
const AddLoadbalancerInternalServerErrorCode int = 500

/*AddLoadbalancerInternalServerError Internal error

swagger:response addLoadbalancerInternalServerError
*/
type AddLoadbalancerInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewAddLoadbalancerInternalServerError creates AddLoadbalancerInternalServerError with default headers values
func NewAddLoadbalancerInternalServerError() *AddLoadbalancerInternalServerError {
	return &AddLoadbalancerInternalServerError{}
}

// WithPayload adds the payload to the add loadbalancer internal server error response
func (o *AddLoadbalancerInternalServerError) WithPayload(payload *models.APIResponse) *AddLoadbalancerInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add loadbalancer internal server error response
func (o *AddLoadbalancerInternalServerError) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddLoadbalancerInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
