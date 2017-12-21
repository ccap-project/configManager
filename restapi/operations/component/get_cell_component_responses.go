// Code generated by go-swagger; DO NOT EDIT.

package component

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"configManager/models"
)

// GetCellComponentOKCode is the HTTP code returned for type GetCellComponentOK
const GetCellComponentOKCode int = 200

/*GetCellComponentOK successful operation

swagger:response getCellComponentOK
*/
type GetCellComponentOK struct {

	/*
	  In: Body
	*/
	Payload *models.Component `json:"body,omitempty"`
}

// NewGetCellComponentOK creates GetCellComponentOK with default headers values
func NewGetCellComponentOK() *GetCellComponentOK {
	return &GetCellComponentOK{}
}

// WithPayload adds the payload to the get cell component o k response
func (o *GetCellComponentOK) WithPayload(payload *models.Component) *GetCellComponentOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get cell component o k response
func (o *GetCellComponentOK) SetPayload(payload *models.Component) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetCellComponentOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetCellComponentBadRequestCode is the HTTP code returned for type GetCellComponentBadRequest
const GetCellComponentBadRequestCode int = 400

/*GetCellComponentBadRequest Invalid cell id or component id

swagger:response getCellComponentBadRequest
*/
type GetCellComponentBadRequest struct {
}

// NewGetCellComponentBadRequest creates GetCellComponentBadRequest with default headers values
func NewGetCellComponentBadRequest() *GetCellComponentBadRequest {
	return &GetCellComponentBadRequest{}
}

// WriteResponse to the client
func (o *GetCellComponentBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(400)
}

// GetCellComponentNotFoundCode is the HTTP code returned for type GetCellComponentNotFound
const GetCellComponentNotFoundCode int = 404

/*GetCellComponentNotFound component not found

swagger:response getCellComponentNotFound
*/
type GetCellComponentNotFound struct {
}

// NewGetCellComponentNotFound creates GetCellComponentNotFound with default headers values
func NewGetCellComponentNotFound() *GetCellComponentNotFound {
	return &GetCellComponentNotFound{}
}

// WriteResponse to the client
func (o *GetCellComponentNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}

// GetCellComponentInternalServerErrorCode is the HTTP code returned for type GetCellComponentInternalServerError
const GetCellComponentInternalServerErrorCode int = 500

/*GetCellComponentInternalServerError Internal error

swagger:response getCellComponentInternalServerError
*/
type GetCellComponentInternalServerError struct {

	/*
	  In: Body
	*/
	Payload models.APIResponse `json:"body,omitempty"`
}

// NewGetCellComponentInternalServerError creates GetCellComponentInternalServerError with default headers values
func NewGetCellComponentInternalServerError() *GetCellComponentInternalServerError {
	return &GetCellComponentInternalServerError{}
}

// WithPayload adds the payload to the get cell component internal server error response
func (o *GetCellComponentInternalServerError) WithPayload(payload models.APIResponse) *GetCellComponentInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get cell component internal server error response
func (o *GetCellComponentInternalServerError) SetPayload(payload models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetCellComponentInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}
