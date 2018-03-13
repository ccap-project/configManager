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

package providertype

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/validate"

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
	HTTPRequest *http.Request `json:"-"`

	/*ProvidertypeID
	  Required: true
	  Max Length: 26
	  Min Length: 26
	  Pattern: ^[0123456789ABCDEFGHJKMNPQRSTVWXYZ]{26}$
	  In: path
	*/
	ProvidertypeID string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls
func (o *GetProviderTypeByIDParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error
	o.HTTPRequest = r

	rProvidertypeID, rhkProvidertypeID, _ := route.Params.GetOK("providertypeId")
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

	o.ProvidertypeID = raw

	if err := o.validateProvidertypeID(formats); err != nil {
		return err
	}

	return nil
}

func (o *GetProviderTypeByIDParams) validateProvidertypeID(formats strfmt.Registry) error {

	if err := validate.MinLength("providertypeId", "path", o.ProvidertypeID, 26); err != nil {
		return err
	}

	if err := validate.MaxLength("providertypeId", "path", o.ProvidertypeID, 26); err != nil {
		return err
	}

	if err := validate.Pattern("providertypeId", "path", o.ProvidertypeID, `^[0123456789ABCDEFGHJKMNPQRSTVWXYZ]{26}$`); err != nil {
		return err
	}

	return nil
}
