package provider

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"../../../models"
)

// UpdateProviderOKCode is the HTTP code returned for type UpdateProviderOK
const UpdateProviderOKCode int = 200

/*UpdateProviderOK successful operation

swagger:response updateProviderOK
*/
type UpdateProviderOK struct {
}

// NewUpdateProviderOK creates UpdateProviderOK with default headers values
func NewUpdateProviderOK() *UpdateProviderOK {
	return &UpdateProviderOK{}
}

// WriteResponse to the client
func (o *UpdateProviderOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
}

// UpdateProviderNotFoundCode is the HTTP code returned for type UpdateProviderNotFound
const UpdateProviderNotFoundCode int = 404

/*UpdateProviderNotFound Does not exists

swagger:response updateProviderNotFound
*/
type UpdateProviderNotFound struct {

	/*
	  In: Body
	*/
	Payload models.APIResponse `json:"body,omitempty"`
}

// NewUpdateProviderNotFound creates UpdateProviderNotFound with default headers values
func NewUpdateProviderNotFound() *UpdateProviderNotFound {
	return &UpdateProviderNotFound{}
}

// WithPayload adds the payload to the update provider not found response
func (o *UpdateProviderNotFound) WithPayload(payload models.APIResponse) *UpdateProviderNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update provider not found response
func (o *UpdateProviderNotFound) SetPayload(payload models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateProviderNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// UpdateProviderConflictCode is the HTTP code returned for type UpdateProviderConflict
const UpdateProviderConflictCode int = 409

/*UpdateProviderConflict Already exists

swagger:response updateProviderConflict
*/
type UpdateProviderConflict struct {

	/*
	  In: Body
	*/
	Payload models.APIResponse `json:"body,omitempty"`
}

// NewUpdateProviderConflict creates UpdateProviderConflict with default headers values
func NewUpdateProviderConflict() *UpdateProviderConflict {
	return &UpdateProviderConflict{}
}

// WithPayload adds the payload to the update provider conflict response
func (o *UpdateProviderConflict) WithPayload(payload models.APIResponse) *UpdateProviderConflict {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update provider conflict response
func (o *UpdateProviderConflict) SetPayload(payload models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateProviderConflict) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(409)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// UpdateProviderInternalServerErrorCode is the HTTP code returned for type UpdateProviderInternalServerError
const UpdateProviderInternalServerErrorCode int = 500

/*UpdateProviderInternalServerError Internal error

swagger:response updateProviderInternalServerError
*/
type UpdateProviderInternalServerError struct {

	/*
	  In: Body
	*/
	Payload models.APIResponse `json:"body,omitempty"`
}

// NewUpdateProviderInternalServerError creates UpdateProviderInternalServerError with default headers values
func NewUpdateProviderInternalServerError() *UpdateProviderInternalServerError {
	return &UpdateProviderInternalServerError{}
}

// WithPayload adds the payload to the update provider internal server error response
func (o *UpdateProviderInternalServerError) WithPayload(payload models.APIResponse) *UpdateProviderInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update provider internal server error response
func (o *UpdateProviderInternalServerError) SetPayload(payload models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateProviderInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}
