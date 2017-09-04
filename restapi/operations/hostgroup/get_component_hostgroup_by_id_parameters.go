package hostgroup

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetComponentHostgroupByIDParams creates a new GetComponentHostgroupByIDParams object
// with the default values initialized.
func NewGetComponentHostgroupByIDParams() GetComponentHostgroupByIDParams {
	var ()
	return GetComponentHostgroupByIDParams{}
}

// GetComponentHostgroupByIDParams contains all the bound params for the get component hostgroup by ID operation
// typically these are obtained from a http.Request
//
// swagger:parameters getComponentHostgroupByID
type GetComponentHostgroupByIDParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request

	/*ID of cell that will be used
	  Required: true
	  In: path
	*/
	CellID int64
	/*ID of component that will be used
	  Required: true
	  In: path
	*/
	ComponentID int64
	/*ID of hostgroup that will be used
	  Required: true
	  In: path
	*/
	HostgroupID int64
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls
func (o *GetComponentHostgroupByIDParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error
	o.HTTPRequest = r

	rCellID, rhkCellID, _ := route.Params.GetOK("cell_id")
	if err := o.bindCellID(rCellID, rhkCellID, route.Formats); err != nil {
		res = append(res, err)
	}

	rComponentID, rhkComponentID, _ := route.Params.GetOK("component_id")
	if err := o.bindComponentID(rComponentID, rhkComponentID, route.Formats); err != nil {
		res = append(res, err)
	}

	rHostgroupID, rhkHostgroupID, _ := route.Params.GetOK("hostgroup_id")
	if err := o.bindHostgroupID(rHostgroupID, rhkHostgroupID, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetComponentHostgroupByIDParams) bindCellID(rawData []string, hasKey bool, formats strfmt.Registry) error {
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

func (o *GetComponentHostgroupByIDParams) bindComponentID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("component_id", "path", "int64", raw)
	}
	o.ComponentID = value

	return nil
}

func (o *GetComponentHostgroupByIDParams) bindHostgroupID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("hostgroup_id", "path", "int64", raw)
	}
	o.HostgroupID = value

	return nil
}
