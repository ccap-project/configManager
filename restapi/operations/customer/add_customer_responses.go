// Code generated by go-swagger; DO NOT EDIT.

package customer

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"configManager/models"
)

// AddCustomerCreatedCode is the HTTP code returned for type AddCustomerCreated
const AddCustomerCreatedCode int = 201

/*AddCustomerCreated Created

swagger:response addCustomerCreated
*/
type AddCustomerCreated struct {

	/*
	  In: Body
	*/
	Payload int64 `json:"body,omitempty"`
}

// NewAddCustomerCreated creates AddCustomerCreated with default headers values
func NewAddCustomerCreated() *AddCustomerCreated {
	return &AddCustomerCreated{}
}

// WithPayload adds the payload to the add customer created response
func (o *AddCustomerCreated) WithPayload(payload int64) *AddCustomerCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add customer created response
func (o *AddCustomerCreated) SetPayload(payload int64) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddCustomerCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// AddCustomerBadRequestCode is the HTTP code returned for type AddCustomerBadRequest
const AddCustomerBadRequestCode int = 400

/*AddCustomerBadRequest Invalid input

swagger:response addCustomerBadRequest
*/
type AddCustomerBadRequest struct {
}

// NewAddCustomerBadRequest creates AddCustomerBadRequest with default headers values
func NewAddCustomerBadRequest() *AddCustomerBadRequest {
	return &AddCustomerBadRequest{}
}

// WriteResponse to the client
func (o *AddCustomerBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(400)
}

// AddCustomerConflictCode is the HTTP code returned for type AddCustomerConflict
const AddCustomerConflictCode int = 409

/*AddCustomerConflict Already exists

swagger:response addCustomerConflict
*/
type AddCustomerConflict struct {

	/*
	  In: Body
	*/
	Payload models.APIResponse `json:"body,omitempty"`
}

// NewAddCustomerConflict creates AddCustomerConflict with default headers values
func NewAddCustomerConflict() *AddCustomerConflict {
	return &AddCustomerConflict{}
}

// WithPayload adds the payload to the add customer conflict response
func (o *AddCustomerConflict) WithPayload(payload models.APIResponse) *AddCustomerConflict {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add customer conflict response
func (o *AddCustomerConflict) SetPayload(payload models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddCustomerConflict) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(409)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// AddCustomerInternalServerErrorCode is the HTTP code returned for type AddCustomerInternalServerError
const AddCustomerInternalServerErrorCode int = 500

/*AddCustomerInternalServerError Internal error

swagger:response addCustomerInternalServerError
*/
type AddCustomerInternalServerError struct {

	/*
	  In: Body
	*/
	Payload models.APIResponse `json:"body,omitempty"`
}

// NewAddCustomerInternalServerError creates AddCustomerInternalServerError with default headers values
func NewAddCustomerInternalServerError() *AddCustomerInternalServerError {
	return &AddCustomerInternalServerError{}
}

// WithPayload adds the payload to the add customer internal server error response
func (o *AddCustomerInternalServerError) WithPayload(payload models.APIResponse) *AddCustomerInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add customer internal server error response
func (o *AddCustomerInternalServerError) SetPayload(payload models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddCustomerInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}
