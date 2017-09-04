package providertype

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"../../../models"
)

// AddProviderTypeOKCode is the HTTP code returned for type AddProviderTypeOK
const AddProviderTypeOKCode int = 200

/*AddProviderTypeOK Already exists

swagger:response addProviderTypeOK
*/
type AddProviderTypeOK struct {

	/*
	  In: Body
	*/
	Payload models.APIResponse `json:"body,omitempty"`
}

// NewAddProviderTypeOK creates AddProviderTypeOK with default headers values
func NewAddProviderTypeOK() *AddProviderTypeOK {
	return &AddProviderTypeOK{}
}

// WithPayload adds the payload to the add provider type o k response
func (o *AddProviderTypeOK) WithPayload(payload models.APIResponse) *AddProviderTypeOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add provider type o k response
func (o *AddProviderTypeOK) SetPayload(payload models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddProviderTypeOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// AddProviderTypeCreatedCode is the HTTP code returned for type AddProviderTypeCreated
const AddProviderTypeCreatedCode int = 201

/*AddProviderTypeCreated Created

swagger:response addProviderTypeCreated
*/
type AddProviderTypeCreated struct {

	/*
	  In: Body
	*/
	Payload models.ProviderTypeID `json:"body,omitempty"`
}

// NewAddProviderTypeCreated creates AddProviderTypeCreated with default headers values
func NewAddProviderTypeCreated() *AddProviderTypeCreated {
	return &AddProviderTypeCreated{}
}

// WithPayload adds the payload to the add provider type created response
func (o *AddProviderTypeCreated) WithPayload(payload models.ProviderTypeID) *AddProviderTypeCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add provider type created response
func (o *AddProviderTypeCreated) SetPayload(payload models.ProviderTypeID) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddProviderTypeCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// AddProviderTypeMethodNotAllowedCode is the HTTP code returned for type AddProviderTypeMethodNotAllowed
const AddProviderTypeMethodNotAllowedCode int = 405

/*AddProviderTypeMethodNotAllowed Invalid input

swagger:response addProviderTypeMethodNotAllowed
*/
type AddProviderTypeMethodNotAllowed struct {
}

// NewAddProviderTypeMethodNotAllowed creates AddProviderTypeMethodNotAllowed with default headers values
func NewAddProviderTypeMethodNotAllowed() *AddProviderTypeMethodNotAllowed {
	return &AddProviderTypeMethodNotAllowed{}
}

// WriteResponse to the client
func (o *AddProviderTypeMethodNotAllowed) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(405)
}

// AddProviderTypeInternalServerErrorCode is the HTTP code returned for type AddProviderTypeInternalServerError
const AddProviderTypeInternalServerErrorCode int = 500

/*AddProviderTypeInternalServerError Internal error

swagger:response addProviderTypeInternalServerError
*/
type AddProviderTypeInternalServerError struct {

	/*
	  In: Body
	*/
	Payload models.APIResponse `json:"body,omitempty"`
}

// NewAddProviderTypeInternalServerError creates AddProviderTypeInternalServerError with default headers values
func NewAddProviderTypeInternalServerError() *AddProviderTypeInternalServerError {
	return &AddProviderTypeInternalServerError{}
}

// WithPayload adds the payload to the add provider type internal server error response
func (o *AddProviderTypeInternalServerError) WithPayload(payload models.APIResponse) *AddProviderTypeInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add provider type internal server error response
func (o *AddProviderTypeInternalServerError) SetPayload(payload models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddProviderTypeInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}
