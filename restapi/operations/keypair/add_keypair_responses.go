// Code generated by go-swagger; DO NOT EDIT.

package keypair

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"configManager/models"
)

// AddKeypairCreatedCode is the HTTP code returned for type AddKeypairCreated
const AddKeypairCreatedCode int = 201

/*AddKeypairCreated Created

swagger:response addKeypairCreated
*/
type AddKeypairCreated struct {

	/*
	  In: Body
	*/
	Payload int64 `json:"body,omitempty"`
}

// NewAddKeypairCreated creates AddKeypairCreated with default headers values
func NewAddKeypairCreated() *AddKeypairCreated {
	return &AddKeypairCreated{}
}

// WithPayload adds the payload to the add keypair created response
func (o *AddKeypairCreated) WithPayload(payload int64) *AddKeypairCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add keypair created response
func (o *AddKeypairCreated) SetPayload(payload int64) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddKeypairCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// AddKeypairMethodNotAllowedCode is the HTTP code returned for type AddKeypairMethodNotAllowed
const AddKeypairMethodNotAllowedCode int = 405

/*AddKeypairMethodNotAllowed Invalid input

swagger:response addKeypairMethodNotAllowed
*/
type AddKeypairMethodNotAllowed struct {

	/*
	  In: Body
	*/
	Payload models.APIResponse `json:"body,omitempty"`
}

// NewAddKeypairMethodNotAllowed creates AddKeypairMethodNotAllowed with default headers values
func NewAddKeypairMethodNotAllowed() *AddKeypairMethodNotAllowed {
	return &AddKeypairMethodNotAllowed{}
}

// WithPayload adds the payload to the add keypair method not allowed response
func (o *AddKeypairMethodNotAllowed) WithPayload(payload models.APIResponse) *AddKeypairMethodNotAllowed {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add keypair method not allowed response
func (o *AddKeypairMethodNotAllowed) SetPayload(payload models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddKeypairMethodNotAllowed) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(405)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// AddKeypairConflictCode is the HTTP code returned for type AddKeypairConflict
const AddKeypairConflictCode int = 409

/*AddKeypairConflict Already exists

swagger:response addKeypairConflict
*/
type AddKeypairConflict struct {

	/*
	  In: Body
	*/
	Payload models.APIResponse `json:"body,omitempty"`
}

// NewAddKeypairConflict creates AddKeypairConflict with default headers values
func NewAddKeypairConflict() *AddKeypairConflict {
	return &AddKeypairConflict{}
}

// WithPayload adds the payload to the add keypair conflict response
func (o *AddKeypairConflict) WithPayload(payload models.APIResponse) *AddKeypairConflict {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add keypair conflict response
func (o *AddKeypairConflict) SetPayload(payload models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddKeypairConflict) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(409)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// AddKeypairInternalServerErrorCode is the HTTP code returned for type AddKeypairInternalServerError
const AddKeypairInternalServerErrorCode int = 500

/*AddKeypairInternalServerError Internal error

swagger:response addKeypairInternalServerError
*/
type AddKeypairInternalServerError struct {

	/*
	  In: Body
	*/
	Payload models.APIResponse `json:"body,omitempty"`
}

// NewAddKeypairInternalServerError creates AddKeypairInternalServerError with default headers values
func NewAddKeypairInternalServerError() *AddKeypairInternalServerError {
	return &AddKeypairInternalServerError{}
}

// WithPayload adds the payload to the add keypair internal server error response
func (o *AddKeypairInternalServerError) WithPayload(payload models.APIResponse) *AddKeypairInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the add keypair internal server error response
func (o *AddKeypairInternalServerError) SetPayload(payload models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *AddKeypairInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}
