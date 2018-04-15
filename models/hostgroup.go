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

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Hostgroup hostgroup
// swagger:model Hostgroup
type Hostgroup struct {

	// bootstrap command
	BootstrapCommand string `json:"bootstrap_command,omitempty"`

	// component
	Component string `json:"component,omitempty"`

	// count
	// Required: true
	Count *int64 `json:"count"`

	// flavor
	// Required: true
	Flavor *string `json:"flavor"`

	// id
	ID ULID `json:"id,omitempty"`

	// image
	// Required: true
	Image *string `json:"image"`

	// listeners
	Listeners HostgroupListeners `json:"listeners"`

	// name
	// Required: true
	Name *string `json:"name"`

	// network
	// Required: true
	Network *string `json:"network"`

	// order
	Order *int64 `json:"order,omitempty"`

	// roles
	Roles HostgroupRoles `json:"roles"`

	// securitygroups
	Securitygroups []string `json:"securitygroups"`

	// username
	// Required: true
	Username *string `json:"username"`
}

// Validate validates this hostgroup
func (m *Hostgroup) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCount(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateFlavor(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateID(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateImage(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateName(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateNetwork(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateSecuritygroups(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateUsername(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Hostgroup) validateCount(formats strfmt.Registry) error {

	if err := validate.Required("count", "body", m.Count); err != nil {
		return err
	}

	return nil
}

func (m *Hostgroup) validateFlavor(formats strfmt.Registry) error {

	if err := validate.Required("flavor", "body", m.Flavor); err != nil {
		return err
	}

	return nil
}

func (m *Hostgroup) validateID(formats strfmt.Registry) error {

	if swag.IsZero(m.ID) { // not required
		return nil
	}

	if err := m.ID.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("id")
		}
		return err
	}

	return nil
}

func (m *Hostgroup) validateImage(formats strfmt.Registry) error {

	if err := validate.Required("image", "body", m.Image); err != nil {
		return err
	}

	return nil
}

func (m *Hostgroup) validateName(formats strfmt.Registry) error {

	if err := validate.Required("name", "body", m.Name); err != nil {
		return err
	}

	return nil
}

func (m *Hostgroup) validateNetwork(formats strfmt.Registry) error {

	if err := validate.Required("network", "body", m.Network); err != nil {
		return err
	}

	return nil
}

func (m *Hostgroup) validateSecuritygroups(formats strfmt.Registry) error {

	if swag.IsZero(m.Securitygroups) { // not required
		return nil
	}

	return nil
}

func (m *Hostgroup) validateUsername(formats strfmt.Registry) error {

	if err := validate.Required("username", "body", m.Username); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Hostgroup) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Hostgroup) UnmarshalBinary(b []byte) error {
	var res Hostgroup
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
