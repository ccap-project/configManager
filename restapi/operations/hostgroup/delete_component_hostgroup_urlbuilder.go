// Code generated by go-swagger; DO NOT EDIT.

package hostgroup

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"errors"
	"net/url"
	golangswaggerpaths "path"
	"strings"

	"github.com/go-openapi/swag"
)

// DeleteComponentHostgroupURL generates an URL for the delete component hostgroup operation
type DeleteComponentHostgroupURL struct {
	CellID      int64
	ComponentID int64
	HostgroupID int64

	_basePath string
	// avoid unkeyed usage
	_ struct{}
}

// WithBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *DeleteComponentHostgroupURL) WithBasePath(bp string) *DeleteComponentHostgroupURL {
	o.SetBasePath(bp)
	return o
}

// SetBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *DeleteComponentHostgroupURL) SetBasePath(bp string) {
	o._basePath = bp
}

// Build a url path and query string
func (o *DeleteComponentHostgroupURL) Build() (*url.URL, error) {
	var result url.URL

	var _path = "/cell/{cell_id}/component/{component_id}/hostgroup/{hostgroup_id}"

	cellID := swag.FormatInt64(o.CellID)
	if cellID != "" {
		_path = strings.Replace(_path, "{cell_id}", cellID, -1)
	} else {
		return nil, errors.New("CellID is required on DeleteComponentHostgroupURL")
	}
	componentID := swag.FormatInt64(o.ComponentID)
	if componentID != "" {
		_path = strings.Replace(_path, "{component_id}", componentID, -1)
	} else {
		return nil, errors.New("ComponentID is required on DeleteComponentHostgroupURL")
	}
	hostgroupID := swag.FormatInt64(o.HostgroupID)
	if hostgroupID != "" {
		_path = strings.Replace(_path, "{hostgroup_id}", hostgroupID, -1)
	} else {
		return nil, errors.New("HostgroupID is required on DeleteComponentHostgroupURL")
	}
	_basePath := o._basePath
	if _basePath == "" {
		_basePath = "/v1"
	}
	result.Path = golangswaggerpaths.Join(_basePath, _path)

	return &result, nil
}

// Must is a helper function to panic when the url builder returns an error
func (o *DeleteComponentHostgroupURL) Must(u *url.URL, err error) *url.URL {
	if err != nil {
		panic(err)
	}
	if u == nil {
		panic("url can't be nil")
	}
	return u
}

// String returns the string representation of the path with query string
func (o *DeleteComponentHostgroupURL) String() string {
	return o.Must(o.Build()).String()
}

// BuildFull builds a full url with scheme, host, path and query string
func (o *DeleteComponentHostgroupURL) BuildFull(scheme, host string) (*url.URL, error) {
	if scheme == "" {
		return nil, errors.New("scheme is required for a full url on DeleteComponentHostgroupURL")
	}
	if host == "" {
		return nil, errors.New("host is required for a full url on DeleteComponentHostgroupURL")
	}

	base, err := o.Build()
	if err != nil {
		return nil, err
	}

	base.Scheme = scheme
	base.Host = host
	return base, nil
}

// StringFull returns the string representation of a complete url
func (o *DeleteComponentHostgroupURL) StringFull(scheme, host string) string {
	return o.Must(o.BuildFull(scheme, host)).String()
}
