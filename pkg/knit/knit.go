package knit

import (
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/tylermmorton/gocodeshift/pkg/generator"
	"github.com/tylermmorton/gocodeshift/pkg/parser"
)

type Config struct {
}

type Knit interface {
	ProcessText(text string) error
	ProcessFile(filepath string) error
}

type knit struct {
	cfg *Config
}

func New() Knit {
	return &knit{}
}

// ProcessFile reads and parses knit options from file
// then executes all configured codegen templates
func (k *knit) ProcessFile(filepath string) error {
	file, err := os.ReadFile(filepath)
	if err != nil {
		return errors.Wrap(err, "failed to load file")
	}
	return k.ProcessText(string(file))
}

// ProcessText parses knit options and executes
// all configured codegen templates
func (k *knit) ProcessText(text string) error {
	b := strings.Builder{}

	// Split the file into blocks based on the ending annotation.
	blocks := strings.SplitAfter(text, parser.ANNOTATION_END)
	for _, block := range blocks {
		// Look for a begin annotation. If there isn't one in
		// this code block, just write the block and continue
		begin := strings.Index(block, parser.ANNOTATION_BEG)
		if begin == -1 {
			b.WriteString(block)
			continue
		}

		// If there is a begin annotation, append all the
		// text before it plus the annotation
		b.WriteString(block[:begin+len(parser.ANNOTATION_BEG)])
		b.WriteString("\n")

		opts, err := parser.Options(block)
		if err != nil {
			return errors.Wrap(err, "failed to parse knit options")
		}

		generator, err := generator.FromOpts(opts)
		if err != nil {
			return errors.Wrap(err, "failed to setup generator context")
		}

		codegen, err := generator.Generate()
		if err != nil {
			return errors.Wrap(err, "failed to generate knit code block")
		}

		b.WriteString(codegen)

		// Find the end annotation and insert it back into the block
		endAnnotation, err := parser.EndAnnotation(block)
		if err != nil {
			return errors.Wrap(err, "failed to parse codegen end annotation")
		}

		b.WriteString(endAnnotation)
	}
	return nil
}
