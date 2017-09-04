package provider

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// DeleteProviderBadRequestCode is the HTTP code returned for type DeleteProviderBadRequest
const DeleteProviderBadRequestCode int = 400

/*DeleteProviderBadRequest Invalid ID supplied

swagger:response deleteProviderBadRequest
*/
type DeleteProviderBadRequest struct {
}

// NewDeleteProviderBadRequest creates DeleteProviderBadRequest with default headers values
func NewDeleteProviderBadRequest() *DeleteProviderBadRequest {
	return &DeleteProviderBadRequest{}
}

// WriteResponse to the client
func (o *DeleteProviderBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
}

// DeleteProviderNotFoundCode is the HTTP code returned for type DeleteProviderNotFound
const DeleteProviderNotFoundCode int = 404

/*DeleteProviderNotFound Provider not found

swagger:response deleteProviderNotFound
*/
type DeleteProviderNotFound struct {
}

// NewDeleteProviderNotFound creates DeleteProviderNotFound with default headers values
func NewDeleteProviderNotFound() *DeleteProviderNotFound {
	return &DeleteProviderNotFound{}
}

// WriteResponse to the client
func (o *DeleteProviderNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
}
