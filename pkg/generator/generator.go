package generator

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/knitcodegen/knit/pkg/loader"
	"github.com/knitcodegen/knit/pkg/parser"
	"github.com/pkg/errors"
)

type Generator interface {
	// Execute runs knit and returns the generated code block
	Generate() (string, error)
}

type generator struct {
	// Options represents the options used to construct this executor
	Options []*parser.Option
	// Input represents an input file or input literal type
	Input string
	// InputLiteral represents the input literal, if one is defined
	InputLiteral string
	// Loader represents the input loader for this executor
	Loader loader.SchemaLoader

	// Templater represents the template engine
	Templater *template.Template
}

// FromOpts constructs a Generator from parser.Options
func FromOpts(options []*parser.Option) (Generator, error) {
	g := &generator{
		Options: options,
	}

	for _, opt := range options {
		switch opt.Type {
		case parser.Input:
			if len(opt.Literal) != 0 {
				g.Input = opt.Value
				g.InputLiteral = opt.Literal
			} else {
				byt, err := os.ReadFile(opt.Value)
				if err != nil {
					return nil, errors.Wrap(err, "failed to read input file")
				}

				g.Input = string(byt)
			}
		case parser.Loader:
			// toss the error here... there's a chance
			// the loader type is defined by the input
			// literal
			loader, _ := createLoader(opt.Value)
			g.Loader = loader
		case parser.Template:
			if len(opt.Literal) != 0 {
				tmpl, err := template.New("knit").Parse(opt.Literal)
				if err != nil {
					return nil, err
				}

				g.Templater = tmpl
			} else {
				byt, err := os.ReadFile(opt.Value)
				if err != nil {
					return nil, errors.Wrap(err, "failed to load template file")
				}

				tmpl, err := template.New("knit").Parse(string(byt))
				if err != nil {
					return nil, errors.Wrap(err, "failed to parse template file")
				}

				g.Templater = tmpl
			}
		}
	}

	// If no loader option was specified, see if we can
	// construct the loader from an input literal type
	if g.Loader == nil {
		if len(g.InputLiteral) != 0 {
			loader, err := createLoader(g.Input)
			if err != nil {
				return nil, err
			}

			g.Loader = loader
		} else {
			return nil, errors.New("no loader or input literal was specified")
		}
	}

	return g, nil
}

func createLoader(loaderType string) (loader.SchemaLoader, error) {
	switch loaderType {
	case "yml":
		fallthrough
	case "yaml":
		return &loader.YamlLoader{}, nil
	case "json":
		return &loader.JsonLoader{}, nil
	case "graphql":
		return &loader.GraphqlLoader{}, nil
	case "openapi3":
		return &loader.OpenAPI3Loader{}, nil
	}
	return nil, fmt.Errorf("undefined loader type %s", loaderType)
}

func (e *generator) Generate() (string, error) {
	var err error
	var data interface{}

	if len(e.InputLiteral) != 0 {
		data, err = e.Loader.LoadFromData([]byte(e.InputLiteral))
		if err != nil {
			return "", err
		}
	} else {
		data, err = e.Loader.LoadFromData([]byte(e.Input))
		if err != nil {
			return "", err
		}
	}

	tpl := &bytes.Buffer{}
	err = e.Templater.Execute(tpl, data)
	if err != nil {
		return "", err
	}

	return tpl.String(), nil
}
