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

package router

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"configManager/models"
)

// GetCellRouterOKCode is the HTTP code returned for type GetCellRouterOK
const GetCellRouterOKCode int = 200

/*GetCellRouterOK successful operation

swagger:response getCellRouterOK
*/
type GetCellRouterOK struct {

	/*
	  In: Body
	*/
	Payload *models.Router `json:"body,omitempty"`
}

// NewGetCellRouterOK creates GetCellRouterOK with default headers values
func NewGetCellRouterOK() *GetCellRouterOK {
	return &GetCellRouterOK{}
}

// WithPayload adds the payload to the get cell router o k response
func (o *GetCellRouterOK) WithPayload(payload *models.Router) *GetCellRouterOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get cell router o k response
func (o *GetCellRouterOK) SetPayload(payload *models.Router) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetCellRouterOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetCellRouterBadRequestCode is the HTTP code returned for type GetCellRouterBadRequest
const GetCellRouterBadRequestCode int = 400

/*GetCellRouterBadRequest Invalid cell id or router id

swagger:response getCellRouterBadRequest
*/
type GetCellRouterBadRequest struct {
}

// NewGetCellRouterBadRequest creates GetCellRouterBadRequest with default headers values
func NewGetCellRouterBadRequest() *GetCellRouterBadRequest {
	return &GetCellRouterBadRequest{}
}

// WriteResponse to the client
func (o *GetCellRouterBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(400)
}

// GetCellRouterNotFoundCode is the HTTP code returned for type GetCellRouterNotFound
const GetCellRouterNotFoundCode int = 404

/*GetCellRouterNotFound router not found

swagger:response getCellRouterNotFound
*/
type GetCellRouterNotFound struct {
}

// NewGetCellRouterNotFound creates GetCellRouterNotFound with default headers values
func NewGetCellRouterNotFound() *GetCellRouterNotFound {
	return &GetCellRouterNotFound{}
}

// WriteResponse to the client
func (o *GetCellRouterNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}

// GetCellRouterInternalServerErrorCode is the HTTP code returned for type GetCellRouterInternalServerError
const GetCellRouterInternalServerErrorCode int = 500

/*GetCellRouterInternalServerError Internal error

swagger:response getCellRouterInternalServerError
*/
type GetCellRouterInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewGetCellRouterInternalServerError creates GetCellRouterInternalServerError with default headers values
func NewGetCellRouterInternalServerError() *GetCellRouterInternalServerError {
	return &GetCellRouterInternalServerError{}
}

// WithPayload adds the payload to the get cell router internal server error response
func (o *GetCellRouterInternalServerError) WithPayload(payload *models.APIResponse) *GetCellRouterInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get cell router internal server error response
func (o *GetCellRouterInternalServerError) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetCellRouterInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
