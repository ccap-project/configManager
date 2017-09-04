package providertype

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetProviderTypeByIDParams creates a new GetProviderTypeByIDParams object
// with the default values initialized.
func NewGetProviderTypeByIDParams() GetProviderTypeByIDParams {
	var ()
	return GetProviderTypeByIDParams{}
}

// GetProviderTypeByIDParams contains all the bound params for the get provider type by Id operation
// typically these are obtained from a http.Request
//
// swagger:parameters getProviderTypeById
type GetProviderTypeByIDParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request

	/*ID of providertype to return
	  Required: true
	  In: path
	*/
	ProvidertypeID int64
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls
func (o *GetProviderTypeByIDParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error
	o.HTTPRequest = r

	rProvidertypeID, rhkProvidertypeID, _ := route.Params.GetOK("providertype_id")
	if err := o.bindProvidertypeID(rProvidertypeID, rhkProvidertypeID, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetProviderTypeByIDParams) bindProvidertypeID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("providertype_id", "path", "int64", raw)
	}
	o.ProvidertypeID = value

	return nil
}
