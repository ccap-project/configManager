// Code generated by go-swagger; DO NOT EDIT.

package host

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"
)

// NewFindCellHostsParams creates a new FindCellHostsParams object
// with the default values initialized.
func NewFindCellHostsParams() FindCellHostsParams {
	var ()
	return FindCellHostsParams{}
}

// FindCellHostsParams contains all the bound params for the find cell hosts operation
// typically these are obtained from a http.Request
//
// swagger:parameters findCellHosts
type FindCellHostsParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*ID of cell that will be used
	  Required: true
	  In: path
	*/
	CellID int64
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls
func (o *FindCellHostsParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error
	o.HTTPRequest = r

	rCellID, rhkCellID, _ := route.Params.GetOK("cell_id")
	if err := o.bindCellID(rCellID, rhkCellID, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *FindCellHostsParams) bindCellID(rawData []string, hasKey bool, formats strfmt.Registry) error {
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
