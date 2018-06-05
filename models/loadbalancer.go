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
	"encoding/json"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Loadbalancer loadbalancer
// swagger:model Loadbalancer
type Loadbalancer struct {

	// algorithm
	// Required: true
	Algorithm *string `json:"algorithm"`

	// connection drain
	ConnectionDrain string `json:"connection_drain,omitempty"`

	// connection idle timeout
	ConnectionIDLETimeout int64 `json:"connection_idle_timeout,omitempty"`

	// id
	ID ULID `json:"id,omitempty"`

	// members
	Members []string `json:"members"`

	// name
	// Required: true
	Name *string `json:"name"`

	// network
	Network []string `json:"network"`

	// port
	// Required: true
	Port *int64 `json:"port"`

	// protocol
	// Required: true
	Protocol *string `json:"protocol"`

	// securitygroups
	Securitygroups []string `json:"securitygroups"`

	// type
	Type string `json:"type,omitempty"`
}

// Validate validates this loadbalancer
func (m *Loadbalancer) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAlgorithm(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateID(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateMembers(formats); err != nil {
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

	if err := m.validatePort(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateProtocol(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateSecuritygroups(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Loadbalancer) validateAlgorithm(formats strfmt.Registry) error {

	if err := validate.Required("algorithm", "body", m.Algorithm); err != nil {
		return err
	}

	return nil
}

func (m *Loadbalancer) validateID(formats strfmt.Registry) error {

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

func (m *Loadbalancer) validateMembers(formats strfmt.Registry) error {

	if swag.IsZero(m.Members) { // not required
		return nil
	}

	return nil
}

func (m *Loadbalancer) validateName(formats strfmt.Registry) error {

	if err := validate.Required("name", "body", m.Name); err != nil {
		return err
	}

	return nil
}

func (m *Loadbalancer) validateNetwork(formats strfmt.Registry) error {

	if swag.IsZero(m.Network) { // not required
		return nil
	}

	return nil
}

func (m *Loadbalancer) validatePort(formats strfmt.Registry) error {

	if err := validate.Required("port", "body", m.Port); err != nil {
		return err
	}

	return nil
}

var loadbalancerTypeProtocolPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["tcp","http","https","tls"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		loadbalancerTypeProtocolPropEnum = append(loadbalancerTypeProtocolPropEnum, v)
	}
}

const (
	// LoadbalancerProtocolTCP captures enum value "tcp"
	LoadbalancerProtocolTCP string = "tcp"
	// LoadbalancerProtocolHTTP captures enum value "http"
	LoadbalancerProtocolHTTP string = "http"
	// LoadbalancerProtocolHTTPS captures enum value "https"
	LoadbalancerProtocolHTTPS string = "https"
	// LoadbalancerProtocolTLS captures enum value "tls"
	LoadbalancerProtocolTLS string = "tls"
)

// prop value enum
func (m *Loadbalancer) validateProtocolEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, loadbalancerTypeProtocolPropEnum); err != nil {
		return err
	}
	return nil
}

func (m *Loadbalancer) validateProtocol(formats strfmt.Registry) error {

	if err := validate.Required("protocol", "body", m.Protocol); err != nil {
		return err
	}

	// value enum
	if err := m.validateProtocolEnum("protocol", "body", *m.Protocol); err != nil {
		return err
	}

	return nil
}

func (m *Loadbalancer) validateSecuritygroups(formats strfmt.Registry) error {

	if swag.IsZero(m.Securitygroups) { // not required
		return nil
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Loadbalancer) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Loadbalancer) UnmarshalBinary(b []byte) error {
	var res Loadbalancer
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
