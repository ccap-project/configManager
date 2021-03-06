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

package component

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"configManager/models"
)

// FindCellComponentsOKCode is the HTTP code returned for type FindCellComponentsOK
const FindCellComponentsOKCode int = 200

/*FindCellComponentsOK successful operation

swagger:response findCellComponentsOK
*/
type FindCellComponentsOK struct {

	/*
	  In: Body
	*/
	Payload models.FindCellComponentsOKBody `json:"body,omitempty"`
}

// NewFindCellComponentsOK creates FindCellComponentsOK with default headers values
func NewFindCellComponentsOK() *FindCellComponentsOK {
	return &FindCellComponentsOK{}
}

// WithPayload adds the payload to the find cell components o k response
func (o *FindCellComponentsOK) WithPayload(payload models.FindCellComponentsOKBody) *FindCellComponentsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the find cell components o k response
func (o *FindCellComponentsOK) SetPayload(payload models.FindCellComponentsOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *FindCellComponentsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		payload = make(models.FindCellComponentsOKBody, 0, 50)
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// FindCellComponentsBadRequestCode is the HTTP code returned for type FindCellComponentsBadRequest
const FindCellComponentsBadRequestCode int = 400

/*FindCellComponentsBadRequest Invalid cell id

swagger:response findCellComponentsBadRequest
*/
type FindCellComponentsBadRequest struct {
}

// NewFindCellComponentsBadRequest creates FindCellComponentsBadRequest with default headers values
func NewFindCellComponentsBadRequest() *FindCellComponentsBadRequest {
	return &FindCellComponentsBadRequest{}
}

// WriteResponse to the client
func (o *FindCellComponentsBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(400)
}

// FindCellComponentsNotFoundCode is the HTTP code returned for type FindCellComponentsNotFound
const FindCellComponentsNotFoundCode int = 404

/*FindCellComponentsNotFound component not found

swagger:response findCellComponentsNotFound
*/
type FindCellComponentsNotFound struct {
}

// NewFindCellComponentsNotFound creates FindCellComponentsNotFound with default headers values
func NewFindCellComponentsNotFound() *FindCellComponentsNotFound {
	return &FindCellComponentsNotFound{}
}

// WriteResponse to the client
func (o *FindCellComponentsNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}

// FindCellComponentsInternalServerErrorCode is the HTTP code returned for type FindCellComponentsInternalServerError
const FindCellComponentsInternalServerErrorCode int = 500

/*FindCellComponentsInternalServerError Internal error

swagger:response findCellComponentsInternalServerError
*/
type FindCellComponentsInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewFindCellComponentsInternalServerError creates FindCellComponentsInternalServerError with default headers values
func NewFindCellComponentsInternalServerError() *FindCellComponentsInternalServerError {
	return &FindCellComponentsInternalServerError{}
}

// WithPayload adds the payload to the find cell components internal server error response
func (o *FindCellComponentsInternalServerError) WithPayload(payload *models.APIResponse) *FindCellComponentsInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the find cell components internal server error response
func (o *FindCellComponentsInternalServerError) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *FindCellComponentsInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
