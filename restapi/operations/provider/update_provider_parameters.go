// Code generated by go-swagger; DO NOT EDIT.

package provider

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"

	"configManager/models"
)

// NewUpdateProviderParams creates a new UpdateProviderParams object
// with the default values initialized.
func NewUpdateProviderParams() UpdateProviderParams {
	var ()
	return UpdateProviderParams{}
}

// UpdateProviderParams contains all the bound params for the update provider operation
// typically these are obtained from a http.Request
//
// swagger:parameters updateProvider
type UpdateProviderParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Provider object that needs to be updated
	  Required: true
	  In: body
	*/
	Body *models.Provider
	/*ID of cell that will be used
	  Required: true
	  In: path
	*/
	CellID int64
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls
func (o *UpdateProviderParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error
	o.HTTPRequest = r

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body models.Provider
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("body", "body"))
			} else {
				res = append(res, errors.NewParseError("body", "body", "", err))
			}

		} else {
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.Body = &body
			}
		}

	} else {
		res = append(res, errors.Required("body", "body"))
	}

	rCellID, rhkCellID, _ := route.Params.GetOK("cell_id")
	if err := o.bindCellID(rCellID, rhkCellID, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *UpdateProviderParams) bindCellID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("cell_id", "path", "int64", raw)
	}
	o.CellID = value

	return nil
}
