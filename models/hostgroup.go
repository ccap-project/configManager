package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"strconv"

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
	ID int64 `json:"id,omitempty"`

	// image
	// Required: true
	Image *string `json:"image"`

	// name
	// Required: true
	Name *string `json:"name"`

	// network
	// Required: true
	Network *string `json:"network"`

	// roles
	Roles []*Role `json:"roles"`

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

	if err := m.validateRoles(formats); err != nil {
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

func (m *Hostgroup) validateRoles(formats strfmt.Registry) error {

	if swag.IsZero(m.Roles) { // not required
		return nil
	}

	for i := 0; i < len(m.Roles); i++ {

		if swag.IsZero(m.Roles[i]) { // not required
			continue
		}

		if m.Roles[i] != nil {

			if err := m.Roles[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("roles" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

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
