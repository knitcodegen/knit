package openapi

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/pkg/errors"
)

type Loader struct{}

func (l *Loader) LoadFromFile(location string) (interface{}, error) {
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile(location)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load from file")
	}
	err = doc.Validate(loader.Context)
	if err != nil {
		return nil, errors.Wrap(err, "failed to validate openapi v3 spec")
	}
	return doc, nil
}
