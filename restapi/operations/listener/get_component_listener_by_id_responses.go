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

package listener

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"configManager/models"
)

// GetComponentListenerByIDOKCode is the HTTP code returned for type GetComponentListenerByIDOK
const GetComponentListenerByIDOKCode int = 200

/*GetComponentListenerByIDOK successful operation

swagger:response getComponentListenerByIdOK
*/
type GetComponentListenerByIDOK struct {

	/*
	  In: Body
	*/
	Payload *models.Listener `json:"body,omitempty"`
}

// NewGetComponentListenerByIDOK creates GetComponentListenerByIDOK with default headers values
func NewGetComponentListenerByIDOK() *GetComponentListenerByIDOK {
	return &GetComponentListenerByIDOK{}
}

// WithPayload adds the payload to the get component listener by Id o k response
func (o *GetComponentListenerByIDOK) WithPayload(payload *models.Listener) *GetComponentListenerByIDOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get component listener by Id o k response
func (o *GetComponentListenerByIDOK) SetPayload(payload *models.Listener) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetComponentListenerByIDOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetComponentListenerByIDBadRequestCode is the HTTP code returned for type GetComponentListenerByIDBadRequest
const GetComponentListenerByIDBadRequestCode int = 400

/*GetComponentListenerByIDBadRequest Invalid cell id or component id

swagger:response getComponentListenerByIdBadRequest
*/
type GetComponentListenerByIDBadRequest struct {

	/*
	  In: Body
	*/
	Payload models.APIResponse `json:"body,omitempty"`
}

// NewGetComponentListenerByIDBadRequest creates GetComponentListenerByIDBadRequest with default headers values
func NewGetComponentListenerByIDBadRequest() *GetComponentListenerByIDBadRequest {
	return &GetComponentListenerByIDBadRequest{}
}

// WithPayload adds the payload to the get component listener by Id bad request response
func (o *GetComponentListenerByIDBadRequest) WithPayload(payload models.APIResponse) *GetComponentListenerByIDBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get component listener by Id bad request response
func (o *GetComponentListenerByIDBadRequest) SetPayload(payload models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetComponentListenerByIDBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// GetComponentListenerByIDNotFoundCode is the HTTP code returned for type GetComponentListenerByIDNotFound
const GetComponentListenerByIDNotFoundCode int = 404

/*GetComponentListenerByIDNotFound hostgroup does not exists

swagger:response getComponentListenerByIdNotFound
*/
type GetComponentListenerByIDNotFound struct {
}

// NewGetComponentListenerByIDNotFound creates GetComponentListenerByIDNotFound with default headers values
func NewGetComponentListenerByIDNotFound() *GetComponentListenerByIDNotFound {
	return &GetComponentListenerByIDNotFound{}
}

// WriteResponse to the client
func (o *GetComponentListenerByIDNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}

// GetComponentListenerByIDInternalServerErrorCode is the HTTP code returned for type GetComponentListenerByIDInternalServerError
const GetComponentListenerByIDInternalServerErrorCode int = 500

/*GetComponentListenerByIDInternalServerError Internal error

swagger:response getComponentListenerByIdInternalServerError
*/
type GetComponentListenerByIDInternalServerError struct {

	/*
	  In: Body
	*/
	Payload models.APIResponse `json:"body,omitempty"`
}

// NewGetComponentListenerByIDInternalServerError creates GetComponentListenerByIDInternalServerError with default headers values
func NewGetComponentListenerByIDInternalServerError() *GetComponentListenerByIDInternalServerError {
	return &GetComponentListenerByIDInternalServerError{}
}

// WithPayload adds the payload to the get component listener by Id internal server error response
func (o *GetComponentListenerByIDInternalServerError) WithPayload(payload models.APIResponse) *GetComponentListenerByIDInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get component listener by Id internal server error response
func (o *GetComponentListenerByIDInternalServerError) SetPayload(payload models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetComponentListenerByIDInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}