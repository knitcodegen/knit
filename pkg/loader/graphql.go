package loader

import (
	"github.com/pkg/errors"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/parser"
)

type GraphqlLoader struct {
}

func (g *GraphqlLoader) LoadFromData(data []byte) (interface{}, error) {
	doc, err := parser.ParseSchema(&ast.Source{
		Input: string(data),
		Name:  "spec",
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse schema")
	}
	return doc, nil
}
