package generator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/knitcodegen/knit/pkg/loader"
	"github.com/knitcodegen/knit/pkg/parser"

	"github.com/pkg/errors"
)

type Generator interface {
	// Validate ensures the generator is properly configured
	Validate() error
	// Generate runs the code generator and returns the generated code block
	Generate() (string, error)
}

type generator struct {
	// Options represent the parser options used to construct this generator
	Options []*parser.Option
	// LoaderType is the string representation of the loader type
	LoaderType string
	// InputFile is the fully resolved path to the input file, if provided
	InputFile *string
	// InputLiteral is the provided input literal OR the loaded InputFile
	InputLiteral string
	// TemplateFile is the fully resolved path to the template file, if provided
	TemplateFile *string
	// TemplateLiteral is the provided template literal OR the loaded TemplateFile
	TemplateLiteral string
}

type OptionType = string

const (
	Input    OptionType = "input"
	Loader   OptionType = "loader"
	Template OptionType = "template"
)

func New(opts ...*parser.Option) (Generator, error) {
	gen := &generator{
		Options: opts,
	}

	for _, opt := range opts {
		switch opt.Type {
		case Input:
			if len(opt.Literal) != 0 {
				gen.InputFile = nil
				gen.InputLiteral = opt.Literal

				if len(opt.Value) != 0 {
					gen.LoaderType = opt.Value
				} else {
					return nil, errors.New("failed to determine loader type from input literal")
				}
			} else {
				path, err := filepath.Abs(opt.Value)
				if err != nil {
					return nil, errors.Wrap(err, "failed to resolve absolute path to input file")
				}

				gen.InputFile = &path
			}
		case Loader:
			gen.LoaderType = opt.Value
		case Template:
			if len(opt.Literal) != 0 {
				gen.TemplateFile = nil
				gen.TemplateLiteral = opt.Literal
			} else {
				path, err := filepath.Abs(opt.Value)
				if err != nil {
					return nil, errors.Wrap(err, "failed to resolve absolute path to template file")
				}

				gen.TemplateFile = &path
			}
		}
	}

	return gen, nil
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

func (gen *generator) Validate() error {
	if len(gen.LoaderType) == 0 {
		return errors.New("missing loader type")
	}

	if len(gen.InputLiteral) == 0 {
		return errors.New("missing input")
	}

	if len(gen.TemplateLiteral) == 0 {
		return errors.New("missing template")
	}

	return nil
}

func (gen *generator) Generate() (string, error) {
	if gen.InputFile != nil {
		byt, err := os.ReadFile(*gen.InputFile)
		if err != nil {
			return "", errors.Wrap(err, "failed to load input file")
		}

		gen.InputLiteral = string(byt)
	}

	if gen.TemplateFile != nil {
		byt, err := os.ReadFile(*gen.TemplateFile)
		if err != nil {
			return "", errors.Wrap(err, "failed to load template file")
		}

		gen.TemplateLiteral = string(byt)
	}

	err := gen.Validate()
	if err != nil {
		return "", errors.Wrap(err, "failed to validate generator configuration")
	}

	loader, err := createLoader(gen.LoaderType)
	if err != nil {
		return "", errors.Wrap(err, "failed to create loader")
	}

	data, err := loader.LoadFromData([]byte(gen.InputLiteral))
	if err != nil {
		return "", errors.Wrap(err, "failed to decode input into schema object")
	}

	tmpl, err := template.
		New("knit").
		Funcs(sprig.TxtFuncMap()).
		Parse(gen.TemplateLiteral)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse template")
	}

	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, data)
	if err != nil {
		return "", errors.Wrap(err, "failed to execute template")
	}

	return buf.String(), nil
}
