// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// Credentials credentials
// swagger:model Credentials

type Credentials struct {

	// location
	Location string `json:"location,omitempty"`

	// password
	Password string `json:"password,omitempty"`

	// specifics
	Specifics interface{} `json:"specifics,omitempty"`

	// username
	Username string `json:"username,omitempty"`
}

/* polymorph Credentials location false */

/* polymorph Credentials password false */

/* polymorph Credentials specifics false */

/* polymorph Credentials username false */

// Validate validates this credentials
func (m *Credentials) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *Credentials) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Credentials) UnmarshalBinary(b []byte) error {
	var res Credentials
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}