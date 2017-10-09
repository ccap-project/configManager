package host

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"../../../models"
)

// FindCellHostsOKCode is the HTTP code returned for type FindCellHostsOK
const FindCellHostsOKCode int = 200

/*FindCellHostsOK successful operation

swagger:response findCellHostsOK
*/
type FindCellHostsOK struct {

	/*
	  In: Body
	*/
	Payload []*models.Host `json:"body,omitempty"`
}

// NewFindCellHostsOK creates FindCellHostsOK with default headers values
func NewFindCellHostsOK() *FindCellHostsOK {
	return &FindCellHostsOK{}
}

// WithPayload adds the payload to the find cell hosts o k response
func (o *FindCellHostsOK) WithPayload(payload []*models.Host) *FindCellHostsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the find cell hosts o k response
func (o *FindCellHostsOK) SetPayload(payload []*models.Host) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *FindCellHostsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		payload = make([]*models.Host, 0, 50)
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// FindCellHostsBadRequestCode is the HTTP code returned for type FindCellHostsBadRequest
const FindCellHostsBadRequestCode int = 400

/*FindCellHostsBadRequest Invalid cell id

swagger:response findCellHostsBadRequest
*/
type FindCellHostsBadRequest struct {
}

// NewFindCellHostsBadRequest creates FindCellHostsBadRequest with default headers values
func NewFindCellHostsBadRequest() *FindCellHostsBadRequest {
	return &FindCellHostsBadRequest{}
}

// WriteResponse to the client
func (o *FindCellHostsBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
}

// FindCellHostsNotFoundCode is the HTTP code returned for type FindCellHostsNotFound
const FindCellHostsNotFoundCode int = 404

/*FindCellHostsNotFound component not found

swagger:response findCellHostsNotFound
*/
type FindCellHostsNotFound struct {
}

// NewFindCellHostsNotFound creates FindCellHostsNotFound with default headers values
func NewFindCellHostsNotFound() *FindCellHostsNotFound {
	return &FindCellHostsNotFound{}
}

// WriteResponse to the client
func (o *FindCellHostsNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
}

// FindCellHostsInternalServerErrorCode is the HTTP code returned for type FindCellHostsInternalServerError
const FindCellHostsInternalServerErrorCode int = 500

/*FindCellHostsInternalServerError Internal error

swagger:response findCellHostsInternalServerError
*/
type FindCellHostsInternalServerError struct {

	/*
	  In: Body
	*/
	Payload models.APIResponse `json:"body,omitempty"`
}

// NewFindCellHostsInternalServerError creates FindCellHostsInternalServerError with default headers values
func NewFindCellHostsInternalServerError() *FindCellHostsInternalServerError {
	return &FindCellHostsInternalServerError{}
}

// WithPayload adds the payload to the find cell hosts internal server error response
func (o *FindCellHostsInternalServerError) WithPayload(payload models.APIResponse) *FindCellHostsInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the find cell hosts internal server error response
func (o *FindCellHostsInternalServerError) SetPayload(payload models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *FindCellHostsInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}
