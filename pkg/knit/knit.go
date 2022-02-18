package knit

import (
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/tylermmorton/gocodeshift/pkg/executor"
	"github.com/tylermmorton/gocodeshift/pkg/parser"
)

// process executes knit against the given text block and
// returns the generated code as a result
func process(block string) (string, error) {
	opts, err := parser.Parse(block)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse knit options")
	}

	exec, err := executor.FromOpts(opts)
	if err != nil {
		return "", errors.Wrap(err, "failed to setup execution context")
	}

	return exec.Execute()
}

func KnitFile(filepath string) error {
	inputFile, err := os.ReadFile(filepath)
	if err != nil {
		return errors.Wrap(err, "failed to load file")
	}

	b := strings.Builder{}

	// Split the file into blocks based on the ending annotation.
	blocks := strings.SplitAfter(string(inputFile), parser.ANNOTATION_END)
	for _, block := range blocks {
		// Look for a begin annotation. If there isn't one in
		// this code block, just write the block and continue
		begin := strings.Index(block, parser.ANNOTATION_BEG)
		if begin == -1 {
			b.WriteString(block)
			continue
		}

		// If there is a begin annotation, append all the
		// text before it
		b.WriteString(block[:begin+len(parser.ANNOTATION_BEG)])
		b.WriteString("\n")

		opts, err := parser.Parse(block)
		if err != nil {
			return errors.Wrap(err, "failed to parse knit options")
		}

		exec, err := executor.FromOpts(opts)
		if err != nil {
			return errors.Wrap(err, "failed to setup execution context")
		}

		// Run the current block through the knit processor
		codegen, err := exec.Execute()
		if err != nil {
			return errors.Wrap(err, "failed to generate knit code block")
		}

		b.WriteString(codegen)

		// Find the end annotation and insert it back into the block
		endAnnotation, err := parser.ParseCodegenEnd(block)
		if err != nil {
			return errors.Wrap(err, "failed to parse codegen end annotation")
		}

		b.WriteString(endAnnotation)
	}

	err = os.WriteFile(filepath, []byte(b.String()), os.ModeExclusive)
	if err != nil {
		return errors.Wrap(err, "failed to write result to file")
	}

	return nil
}
