package cell

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"../../../models"
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
