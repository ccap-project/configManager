package keypair

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetKeypairByIDParams creates a new GetKeypairByIDParams object
// with the default values initialized.
func NewGetKeypairByIDParams() GetKeypairByIDParams {
	var ()
	return GetKeypairByIDParams{}
}

// GetKeypairByIDParams contains all the bound params for the get keypair by Id operation
// typically these are obtained from a http.Request
//
// swagger:parameters getKeypairById
type GetKeypairByIDParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request

	/*ID of keypair to return
	  Required: true
	  In: path
	*/
	KeypairID int64
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls
func (o *GetKeypairByIDParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error
	o.HTTPRequest = r

	rKeypairID, rhkKeypairID, _ := route.Params.GetOK("keypair_id")
	if err := o.bindKeypairID(rKeypairID, rhkKeypairID, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetKeypairByIDParams) bindKeypairID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("keypair_id", "path", "int64", raw)
	}
	o.KeypairID = value

	return nil
}
