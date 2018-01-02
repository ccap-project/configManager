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

// GetCellByIDOKCode is the HTTP code returned for type GetCellByIDOK
const GetCellByIDOKCode int = 200

/*GetCellByIDOK successful operation

swagger:response getCellByIdOK
*/
type GetCellByIDOK struct {

	/*
	  In: Body
	*/
	Payload *models.Cell `json:"body,omitempty"`
}

// NewGetCellByIDOK creates GetCellByIDOK with default headers values
func NewGetCellByIDOK() *GetCellByIDOK {
	return &GetCellByIDOK{}
}

// WithPayload adds the payload to the get cell by Id o k response
func (o *GetCellByIDOK) WithPayload(payload *models.Cell) *GetCellByIDOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get cell by Id o k response
func (o *GetCellByIDOK) SetPayload(payload *models.Cell) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetCellByIDOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetCellByIDNotFoundCode is the HTTP code returned for type GetCellByIDNotFound
const GetCellByIDNotFoundCode int = 404

/*GetCellByIDNotFound Cell not found

swagger:response getCellByIdNotFound
*/
type GetCellByIDNotFound struct {
}

// NewGetCellByIDNotFound creates GetCellByIDNotFound with default headers values
func NewGetCellByIDNotFound() *GetCellByIDNotFound {
	return &GetCellByIDNotFound{}
}

// WriteResponse to the client
func (o *GetCellByIDNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}

// GetCellByIDMethodNotAllowedCode is the HTTP code returned for type GetCellByIDMethodNotAllowed
const GetCellByIDMethodNotAllowedCode int = 405

/*GetCellByIDMethodNotAllowed Invalid input

swagger:response getCellByIdMethodNotAllowed
*/
type GetCellByIDMethodNotAllowed struct {

	/*
	  In: Body
	*/
	Payload models.APIResponse `json:"body,omitempty"`
}

// NewGetCellByIDMethodNotAllowed creates GetCellByIDMethodNotAllowed with default headers values
func NewGetCellByIDMethodNotAllowed() *GetCellByIDMethodNotAllowed {
	return &GetCellByIDMethodNotAllowed{}
}

// WithPayload adds the payload to the get cell by Id method not allowed response
func (o *GetCellByIDMethodNotAllowed) WithPayload(payload models.APIResponse) *GetCellByIDMethodNotAllowed {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get cell by Id method not allowed response
func (o *GetCellByIDMethodNotAllowed) SetPayload(payload models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetCellByIDMethodNotAllowed) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(405)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// GetCellByIDInternalServerErrorCode is the HTTP code returned for type GetCellByIDInternalServerError
const GetCellByIDInternalServerErrorCode int = 500

/*GetCellByIDInternalServerError Internal error

swagger:response getCellByIdInternalServerError
*/
type GetCellByIDInternalServerError struct {

	/*
	  In: Body
	*/
	Payload models.APIResponse `json:"body,omitempty"`
}

// NewGetCellByIDInternalServerError creates GetCellByIDInternalServerError with default headers values
func NewGetCellByIDInternalServerError() *GetCellByIDInternalServerError {
	return &GetCellByIDInternalServerError{}
}

// WithPayload adds the payload to the get cell by Id internal server error response
func (o *GetCellByIDInternalServerError) WithPayload(payload models.APIResponse) *GetCellByIDInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get cell by Id internal server error response
func (o *GetCellByIDInternalServerError) SetPayload(payload models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetCellByIDInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}
