package main

import (
	"bytes"
	"fmt"
	"go/format"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/tylermmorton/gocodeshift/pkg/loader"
	"github.com/tylermmorton/gocodeshift/pkg/openapi"
	"github.com/urfave/cli/v2"
)

const (
	ANNOTATION_BEG = "// @knit"
	ANNOTATION_END = "// @!knit"
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
				shiftFile(match)
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

func parse(block string) (*shifter, error) {
	s := &shifter{}
	lines := strings.Split(block, "\n")
	for _, l := range lines {
		// if the current line does not contain an annotation
		// then skip it
		begin := strings.Index(l, ANNOTATION_BEG)
		if begin == -1 {
			continue
		}

		// skip past the annotation and parse the options
		begin += len(ANNOTATION_BEG) + 1
		rawOpts := l[begin:]
		// save raw option text for use in output
		s.Options = append(s.Options, rawOpts)
		// expand environment variables in options
		rawOpts = os.ExpandEnv(rawOpts)

		split := strings.Split(rawOpts, " ")
		if split[0] == "input" {
			s.Input = split[1]
		} else if split[0] == "loader" {
			if split[1] == "openapi" {
				s.Loader = &openapi.Loader{}
			}
		} else if split[0] == "template" {
			file, err := os.ReadFile(split[1])
			if err != nil {
				return nil, errors.Wrap(err, "failed to load template file")
			}
			tmpl, err := template.New("shift").Parse(string(file))
			if err != nil {
				return nil, errors.Wrap(err, "failed to parse template file")
			}
			s.Template = tmpl
		} else {
			log.Printf("unknown option: %s", split[0])
			continue
		}
	}
	return s, nil
}

func shift(block string) (string, error) {
	b := strings.Builder{}

	s, err := parse(block)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse options")
	}

	dat, err := s.Loader.LoadFromFile(s.Input)
	if err != nil {
		return "", errors.Wrap(err, "loader failed to load file")
	}

	tpl := &bytes.Buffer{}
	err = s.Template.Execute(tpl, dat)
	if err != nil {
		return "", errors.Wrap(err, "failed to execute template")
	}

	for _, opt := range s.Options {
		b.WriteString(fmt.Sprintf("%s %s\n", ANNOTATION_BEG, opt))
	}
	b.Write(tpl.Bytes())
	b.WriteString(fmt.Sprintf("%s\n", ANNOTATION_END))

	return b.String(), nil
}

func shiftFile(location string) error {
	inFile, err := os.ReadFile(location)
	if err != nil {
		return errors.Wrap(err, "failed to read file")
	}

	b := strings.Builder{}

	// SplitAfter includes all text before and including the annotation
	blocks := strings.SplitAfter(string(inFile), ANNOTATION_END)

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
		b.WriteString(block[:begin])

		// Now what remains is the begin/end annotations and
		// everything in between
		// Process the annotation block and append the result
		shifted, err := shift(block[begin:])
		if err != nil {
			return errors.Wrap(err, "failed to shift block")
		}
		b.WriteString(shifted)
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
