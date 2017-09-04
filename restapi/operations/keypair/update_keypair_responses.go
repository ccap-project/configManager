package keypair

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// UpdateKeypairBadRequestCode is the HTTP code returned for type UpdateKeypairBadRequest
const UpdateKeypairBadRequestCode int = 400

/*UpdateKeypairBadRequest Invalid ID supplied

swagger:response updateKeypairBadRequest
*/
type UpdateKeypairBadRequest struct {
}

// NewUpdateKeypairBadRequest creates UpdateKeypairBadRequest with default headers values
func NewUpdateKeypairBadRequest() *UpdateKeypairBadRequest {
	return &UpdateKeypairBadRequest{}
}

// WriteResponse to the client
func (o *UpdateKeypairBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
}

// UpdateKeypairNotFoundCode is the HTTP code returned for type UpdateKeypairNotFound
const UpdateKeypairNotFoundCode int = 404

/*UpdateKeypairNotFound Keypair not found

swagger:response updateKeypairNotFound
*/
type UpdateKeypairNotFound struct {
}

// NewUpdateKeypairNotFound creates UpdateKeypairNotFound with default headers values
func NewUpdateKeypairNotFound() *UpdateKeypairNotFound {
	return &UpdateKeypairNotFound{}
}

// WriteResponse to the client
func (o *UpdateKeypairNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
}

// UpdateKeypairMethodNotAllowedCode is the HTTP code returned for type UpdateKeypairMethodNotAllowed
const UpdateKeypairMethodNotAllowedCode int = 405

/*UpdateKeypairMethodNotAllowed Validation exception

swagger:response updateKeypairMethodNotAllowed
*/
type UpdateKeypairMethodNotAllowed struct {
}

// NewUpdateKeypairMethodNotAllowed creates UpdateKeypairMethodNotAllowed with default headers values
func NewUpdateKeypairMethodNotAllowed() *UpdateKeypairMethodNotAllowed {
	return &UpdateKeypairMethodNotAllowed{}
}

// WriteResponse to the client
func (o *UpdateKeypairMethodNotAllowed) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(405)
}
