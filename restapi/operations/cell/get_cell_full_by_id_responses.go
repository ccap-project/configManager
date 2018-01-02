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

package cell

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"configManager/models"
)

// GetCellFullByIDOKCode is the HTTP code returned for type GetCellFullByIDOK
const GetCellFullByIDOKCode int = 200

/*GetCellFullByIDOK successful operation

swagger:response getCellFullByIdOK
*/
type GetCellFullByIDOK struct {

	/*
	  In: Body
	*/
	Payload *models.FullCell `json:"body,omitempty"`
}

// NewGetCellFullByIDOK creates GetCellFullByIDOK with default headers values
func NewGetCellFullByIDOK() *GetCellFullByIDOK {
	return &GetCellFullByIDOK{}
}

// WithPayload adds the payload to the get cell full by Id o k response
func (o *GetCellFullByIDOK) WithPayload(payload *models.FullCell) *GetCellFullByIDOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get cell full by Id o k response
func (o *GetCellFullByIDOK) SetPayload(payload *models.FullCell) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetCellFullByIDOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetCellFullByIDNotFoundCode is the HTTP code returned for type GetCellFullByIDNotFound
const GetCellFullByIDNotFoundCode int = 404

/*GetCellFullByIDNotFound Cell not found

swagger:response getCellFullByIdNotFound
*/
type GetCellFullByIDNotFound struct {
}

// NewGetCellFullByIDNotFound creates GetCellFullByIDNotFound with default headers values
func NewGetCellFullByIDNotFound() *GetCellFullByIDNotFound {
	return &GetCellFullByIDNotFound{}
}

// WriteResponse to the client
func (o *GetCellFullByIDNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}

// GetCellFullByIDMethodNotAllowedCode is the HTTP code returned for type GetCellFullByIDMethodNotAllowed
const GetCellFullByIDMethodNotAllowedCode int = 405

/*GetCellFullByIDMethodNotAllowed Invalid input

swagger:response getCellFullByIdMethodNotAllowed
*/
type GetCellFullByIDMethodNotAllowed struct {

	/*
	  In: Body
	*/
	Payload models.APIResponse `json:"body,omitempty"`
}

// NewGetCellFullByIDMethodNotAllowed creates GetCellFullByIDMethodNotAllowed with default headers values
func NewGetCellFullByIDMethodNotAllowed() *GetCellFullByIDMethodNotAllowed {
	return &GetCellFullByIDMethodNotAllowed{}
}

// WithPayload adds the payload to the get cell full by Id method not allowed response
func (o *GetCellFullByIDMethodNotAllowed) WithPayload(payload models.APIResponse) *GetCellFullByIDMethodNotAllowed {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get cell full by Id method not allowed response
func (o *GetCellFullByIDMethodNotAllowed) SetPayload(payload models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetCellFullByIDMethodNotAllowed) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(405)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// GetCellFullByIDInternalServerErrorCode is the HTTP code returned for type GetCellFullByIDInternalServerError
const GetCellFullByIDInternalServerErrorCode int = 500

/*GetCellFullByIDInternalServerError Internal error

swagger:response getCellFullByIdInternalServerError
*/
type GetCellFullByIDInternalServerError struct {

	/*
	  In: Body
	*/
	Payload models.APIResponse `json:"body,omitempty"`
}

// NewGetCellFullByIDInternalServerError creates GetCellFullByIDInternalServerError with default headers values
func NewGetCellFullByIDInternalServerError() *GetCellFullByIDInternalServerError {
	return &GetCellFullByIDInternalServerError{}
}

// WithPayload adds the payload to the get cell full by Id internal server error response
func (o *GetCellFullByIDInternalServerError) WithPayload(payload models.APIResponse) *GetCellFullByIDInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get cell full by Id internal server error response
func (o *GetCellFullByIDInternalServerError) SetPayload(payload models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetCellFullByIDInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}
