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
)

// FullCell full cell
// swagger:model FullCell
type FullCell struct {

	// components
	Components FullCellComponents `json:"components"`

	// customer name
	CustomerName string `json:"customer_name,omitempty"`

	// keypair
	Keypair *Keypair `json:"keypair,omitempty"`

	// loadbalancers
	Loadbalancers FullCellLoadbalancers `json:"loadbalancers"`

	// name
	Name string `json:"name,omitempty"`

	// provider
	Provider *Provider `json:"provider,omitempty"`
}

// Validate validates this full cell
func (m *FullCell) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateKeypair(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateProvider(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *FullCell) validateKeypair(formats strfmt.Registry) error {

	if swag.IsZero(m.Keypair) { // not required
		return nil
	}

	if m.Keypair != nil {

		if err := m.Keypair.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("keypair")
			}
			return err
		}
	}

	return nil
}

func (m *FullCell) validateProvider(formats strfmt.Registry) error {

	if swag.IsZero(m.Provider) { // not required
		return nil
	}

	if m.Provider != nil {

		if err := m.Provider.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("provider")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *FullCell) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *FullCell) UnmarshalBinary(b []byte) error {
	var res FullCell
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
