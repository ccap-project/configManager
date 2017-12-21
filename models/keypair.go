// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Keypair keypair
// swagger:model Keypair
type Keypair struct {

	// id
	ID int64 `json:"id,omitempty"`

	// name
	// Required: true
	Name *string `json:"name"`

	// public key
	// Required: true
	PublicKey *string `json:"public_key"`
}

// Validate validates this keypair
func (m *Keypair) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateName(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validatePublicKey(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Keypair) validateName(formats strfmt.Registry) error {

	if err := validate.Required("name", "body", m.Name); err != nil {
		return err
	}

	return nil
}

func (m *Keypair) validatePublicKey(formats strfmt.Registry) error {

	if err := validate.Required("public_key", "body", m.PublicKey); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Keypair) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Keypair) UnmarshalBinary(b []byte) error {
	var res Keypair
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
