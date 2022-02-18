package main

import (
	"bytes"
	"fmt"
	"go/format"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/tylermmorton/gocodeshift/pkg/loader"
	"github.com/tylermmorton/gocodeshift/pkg/openapi"
	"github.com/tylermmorton/gocodeshift/pkg/parser"
	"github.com/urfave/cli/v2"
)

const (
	ANNOTATION_OPT = "@knit"
	ANNOTATION_BEG = "@+knit"
	ANNOTATION_END = "@!knit"

	BEG_PATTERN = `(.*)@\+knit(.*)`
	END_PATTERN = `(.*)@!knit(.*)`
)

type Option = string

const (
	Input    Option = "input"
	Loader   Option = "loader"
	Variable Option = "variable"
	Template Option = "template"
)

func main() {
	(&cli.App{
		Name:  "boom",
		Usage: "make an explosive entrance",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "glob",
				Value: "*.go",
				Usage: "glob pattern of files to search for",
			},
		},
		Action: func(c *cli.Context) error {
			pattern := c.String("glob")
			matches, err := filepath.Glob(pattern)
			if err != nil {
				return errors.Wrap(err, "failed to match any files")
			}
			for _, match := range matches {
				log.Printf("Knitting File: " + match)
				err = shiftFile(match)
				if err != nil {
					log.Printf("failed to knit file: %+v", err)
					return errors.Wrap(err, "failed to knit file")
				}
			}

			return nil
		},
	}).Run(os.Args)
}

type ExecutionContext struct {
	Options  []string
	Input    string
	Loader   loader.SchemaLoader
	Template *template.Template
}

func parseOptions(block string) (*ExecutionContext, error) {
	e := &ExecutionContext{}
	opts, err := parser.Parse(block)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse block")
	}
	for _, opt := range opts {
		switch opt.Type {
		case Input:
			_, err := os.Stat(opt.Value)
			if err != nil {
			} else {
				e.Input = opt.Value
			}
		case Loader:
			loaderType := opt.Value
			if loaderType == "openapi" {
				e.Loader = &openapi.Loader{}
			} else if loaderType == "yml" ||
				loaderType == "yaml" {
				e.Loader = &loader.YamlLoader{}
			}
		case Template:
			tmpl := template.New("knit")
			_, err := os.Stat(opt.Value)
			if err != nil {
				if len(opt.Literal) == 0 {
					return nil, fmt.Errorf("template %s is not a file and did not find inline template definition", opt.Value)
				}

				tmpl, err = tmpl.Parse(opt.Literal)
				if err != nil {
					return nil, errors.Wrap(err, "failed to parse inline template")
				}
			} else {
				// template is a file, load from disk
				file, err := os.ReadFile(opt.Value)
				if err != nil {
					return nil, errors.Wrap(err, "failed to load template file")
				}

				tmpl, err = tmpl.Parse(string(file))
				if err != nil {
					return nil, errors.Wrap(err, "failed to read template from file")
				}
			}
			e.Template = tmpl
		default:
			return nil, fmt.Errorf("unknown option type: %s", opt.Type)
		}
	}
	return e, nil
}

func process(block string) (string, error) {
	opts, err := parseOptions(block)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse knit options")
	}

	dat, err := opts.Loader.LoadFromFile(opts.Input)
	if err != nil {
		return "", errors.Wrap(err, "loader failed to load file")
	}

	tpl := &bytes.Buffer{}
	err = opts.Template.Execute(tpl, dat)
	if err != nil {
		return "", errors.Wrap(err, "failed to execute template")
	}

	return tpl.String(), nil
}

func shiftFile(location string) error {
	inFile, err := os.ReadFile(location)
	if err != nil {
		return errors.Wrap(err, "failed to read file")
	}

	b := strings.Builder{}

	// SplitAfter includes all text before and including the annotation
	blocks := strings.SplitAfter(string(inFile), ANNOTATION_END)
	fmt.Printf("blocks: %d", len(blocks))

	for _, block := range blocks {
		// Look for a begin annotation. If there isn't one in
		// this code block, just write the block and continue
		begin := strings.Index(block, ANNOTATION_BEG)
		if begin == -1 {
			b.WriteString(block)
			continue
		}

		// If there is a begin annotation, append all the
		// text before it
		b.WriteString(block[:begin+len(ANNOTATION_BEG)])
		b.WriteString("\n")

		// Process the annotated block and append the codegen result
		generated, err := process(block)
		if err != nil {
			return errors.Wrap(err, "failed to generate knit code block")
		}
		b.WriteString(generated)

		re := regexp.MustCompile(END_PATTERN)
		endAnnotation := re.FindString(block)
		if len(endAnnotation) == 0 {
			return fmt.Errorf("failed to match end annotation")
		}

		b.WriteString(endAnnotation)
	}

	newFile, err := format.Source([]byte(b.String()))
	if err != nil {
		return errors.Wrap(err, "failed to format result")
	}

	err = os.WriteFile(location, newFile, os.ModeExclusive)
	if err != nil {
		return errors.Wrap(err, "failed to write result to file")
	}
	return nil
}
