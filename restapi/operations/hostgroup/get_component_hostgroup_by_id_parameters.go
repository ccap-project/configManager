// Code generated by go-swagger; DO NOT EDIT.

// Copyright (c) 2016, 2017, 2018 Alexandre Biancalana <ale@biancalanas.net>.
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//     * Redistributions of source code must retain the above copyright
//       notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above copyright
//       notice, this list of conditions and the following disclaimer in the
//       documentation and/or other materials provided with the distribution.
//     * Neither the name of the <organization> nor the
//       names of its contributors may be used to endorse or promote products
//       derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL <COPYRIGHT HOLDER> BE LIABLE FOR ANY
// DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//

package hostgroup

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"

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
	HTTPRequest *http.Request `json:"-"`

	/*Cell ID
	  Required: true
	  Max Length: 26
	  Min Length: 26
	  Pattern: ^[0123456789ABCDEFGHJKMNPQRSTVWXYZ]{26}$
	  In: path
	*/
	CellID string
	/*Component ID
	  Required: true
	  Max Length: 26
	  Min Length: 26
	  Pattern: ^[0123456789ABCDEFGHJKMNPQRSTVWXYZ]{26}$
	  In: path
	*/
	ComponentID string
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

	rCellID, rhkCellID, _ := route.Params.GetOK("cellId")
	if err := o.bindCellID(rCellID, rhkCellID, route.Formats); err != nil {
		res = append(res, err)
	}

	rComponentID, rhkComponentID, _ := route.Params.GetOK("componentId")
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

	o.CellID = raw

	if err := o.validateCellID(formats); err != nil {
		return err
	}

	return nil
}

func (o *GetComponentHostgroupByIDParams) validateCellID(formats strfmt.Registry) error {

	if err := validate.MinLength("cellId", "path", o.CellID, 26); err != nil {
		return err
	}

	if err := validate.MaxLength("cellId", "path", o.CellID, 26); err != nil {
		return err
	}

	if err := validate.Pattern("cellId", "path", o.CellID, `^[0123456789ABCDEFGHJKMNPQRSTVWXYZ]{26}$`); err != nil {
		return err
	}

	return nil
}

func (o *GetComponentHostgroupByIDParams) bindComponentID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	o.ComponentID = raw

	if err := o.validateComponentID(formats); err != nil {
		return err
	}

	return nil
}

func (o *GetComponentHostgroupByIDParams) validateComponentID(formats strfmt.Registry) error {

	if err := validate.MinLength("componentId", "path", o.ComponentID, 26); err != nil {
		return err
	}

	if err := validate.MaxLength("componentId", "path", o.ComponentID, 26); err != nil {
		return err
	}

	if err := validate.Pattern("componentId", "path", o.ComponentID, `^[0123456789ABCDEFGHJKMNPQRSTVWXYZ]{26}$`); err != nil {
		return err
	}

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
