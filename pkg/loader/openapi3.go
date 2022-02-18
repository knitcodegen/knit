package loader

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/pkg/errors"
)

type OpenAPI3Loader struct{}

func (l *OpenAPI3Loader) LoadFromData(data []byte) (interface{}, error) {
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load from data")
	}
	err = doc.Validate(loader.Context)
	if err != nil {
		return nil, errors.Wrap(err, "failed to validate openapi v3 spec")
	}
	return doc, nil
}

func (l *OpenAPI3Loader) LoadFromFile(location string) (interface{}, error) {
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
