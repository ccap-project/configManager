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

// GetKeypairByIDOKCode is the HTTP code returned for type GetKeypairByIDOK
const GetKeypairByIDOKCode int = 200

/*GetKeypairByIDOK successful operation

swagger:response getKeypairByIdOK
*/
type GetKeypairByIDOK struct {

	/*
	  In: Body
	*/
	Payload *models.Keypair `json:"body,omitempty"`
}

// NewGetKeypairByIDOK creates GetKeypairByIDOK with default headers values
func NewGetKeypairByIDOK() *GetKeypairByIDOK {
	return &GetKeypairByIDOK{}
}

// WithPayload adds the payload to the get keypair by Id o k response
func (o *GetKeypairByIDOK) WithPayload(payload *models.Keypair) *GetKeypairByIDOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get keypair by Id o k response
func (o *GetKeypairByIDOK) SetPayload(payload *models.Keypair) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetKeypairByIDOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetKeypairByIDBadRequestCode is the HTTP code returned for type GetKeypairByIDBadRequest
const GetKeypairByIDBadRequestCode int = 400

/*GetKeypairByIDBadRequest Invalid ID supplied

swagger:response getKeypairByIdBadRequest
*/
type GetKeypairByIDBadRequest struct {
}

// NewGetKeypairByIDBadRequest creates GetKeypairByIDBadRequest with default headers values
func NewGetKeypairByIDBadRequest() *GetKeypairByIDBadRequest {
	return &GetKeypairByIDBadRequest{}
}

// WriteResponse to the client
func (o *GetKeypairByIDBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(400)
}

// GetKeypairByIDNotFoundCode is the HTTP code returned for type GetKeypairByIDNotFound
const GetKeypairByIDNotFoundCode int = 404

/*GetKeypairByIDNotFound Keypair not found

swagger:response getKeypairByIdNotFound
*/
type GetKeypairByIDNotFound struct {
}

// NewGetKeypairByIDNotFound creates GetKeypairByIDNotFound with default headers values
func NewGetKeypairByIDNotFound() *GetKeypairByIDNotFound {
	return &GetKeypairByIDNotFound{}
}

// WriteResponse to the client
func (o *GetKeypairByIDNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}

// GetKeypairByIDInternalServerErrorCode is the HTTP code returned for type GetKeypairByIDInternalServerError
const GetKeypairByIDInternalServerErrorCode int = 500

/*GetKeypairByIDInternalServerError InternalError

swagger:response getKeypairByIdInternalServerError
*/
type GetKeypairByIDInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewGetKeypairByIDInternalServerError creates GetKeypairByIDInternalServerError with default headers values
func NewGetKeypairByIDInternalServerError() *GetKeypairByIDInternalServerError {
	return &GetKeypairByIDInternalServerError{}
}

// WithPayload adds the payload to the get keypair by Id internal server error response
func (o *GetKeypairByIDInternalServerError) WithPayload(payload *models.APIResponse) *GetKeypairByIDInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get keypair by Id internal server error response
func (o *GetKeypairByIDInternalServerError) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetKeypairByIDInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
