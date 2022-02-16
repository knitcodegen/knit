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
	"github.com/urfave/cli/v2"
)

const (
	ANNOTATION_OPT = "@knit"
	ANNOTATION_BEG = "@+knit"
	ANNOTATION_END = "@!knit"

	BEG_PATTERN    = `(.*)@\+knit(.*)`
	END_PATTERN    = `(.*)@!knit(.*)`
	OPTION_PATTERN = "@knit[^\\n](\\w*)[^\\n]([^`\\n]*)(?:(`(.|\\n)*`))?"
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
				log.Printf("shifting file: " + match)
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

type shifter struct {
	Options  []string
	Input    string
	Loader   loader.SchemaLoader
	Template *template.Template
}

func parseOptions(block string) (*shifter, error) {
	s := &shifter{}
	re := regexp.MustCompile(OPTION_PATTERN)

	matches := re.FindAllStringSubmatch(block, -1)
	if matches == nil {
		return nil, errors.New("failed to match any knit options")
	}
	//log.Printf("%+v", matches)
	for _, match := range matches {
		s.Options = append(s.Options, match[0])

		// expand env vars on the options value
		match[2] = os.ExpandEnv(match[2])

		optionType := match[1]
		switch optionType {
		case Input:
			s.Input = match[2]
		case Loader:
			loaderType := match[2]
			if loaderType == "openapi" {
				s.Loader = &openapi.Loader{}
			} else if loaderType == "yml" ||
				loaderType == "yaml" {
				s.Loader = &loader.YamlLoader{}
			}
		case Template:
			tmpl := template.New("knit")
			_, err := os.Stat(match[2])
			if err != nil {
				if len(match[3]) == 0 {
					return nil, fmt.Errorf("template %s is not a file and did not find inline template definition", match[2])
				}

				str := strings.TrimSpace(match[3])
				// remove the first and last characters (`backticks`)
				str = str[1 : len(match[3])-1]

				tmpl, err = tmpl.Parse(str)
				if err != nil {
					return nil, errors.Wrap(err, "failed to parse inline template")
				}
			} else {
				// template is a file, load from disk
				file, err := os.ReadFile(match[2])
				if err != nil {
					return nil, errors.Wrap(err, "failed to load template file")
				}

				tmpl, err = tmpl.Parse(string(file))
				if err != nil {
					return nil, errors.Wrap(err, "failed to read template from file")
				}
			}
			s.Template = tmpl
		default:
			return nil, fmt.Errorf("unknown option type: %s", match[0])
		}
	}
	return s, nil
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
