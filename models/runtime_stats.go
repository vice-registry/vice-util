// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// RuntimeStats runtime stats
// swagger:model RuntimeStats

type RuntimeStats struct {

	// export worker
	ExportWorker int64 `json:"exportWorker,omitempty"`

	// exports pending
	ExportsPending int64 `json:"exportsPending,omitempty"`

	// import worker
	ImportWorker int64 `json:"importWorker,omitempty"`

	// imports pending
	ImportsPending int64 `json:"importsPending,omitempty"`

	// store pending
	StorePending int64 `json:"storePending,omitempty"`

	// store worker
	StoreWorker int64 `json:"storeWorker,omitempty"`
}

/* polymorph RuntimeStats exportWorker false */

/* polymorph RuntimeStats exportsPending false */

/* polymorph RuntimeStats importWorker false */

/* polymorph RuntimeStats importsPending false */

/* polymorph RuntimeStats storePending false */

/* polymorph RuntimeStats storeWorker false */

// Validate validates this runtime stats
func (m *RuntimeStats) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *RuntimeStats) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *RuntimeStats) UnmarshalBinary(b []byte) error {
	var res RuntimeStats
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}