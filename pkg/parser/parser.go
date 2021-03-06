package parser

import (
	"os"
	"strings"

	"github.com/dlclark/regexp2"
	"github.com/pkg/errors"
)

const (
	ANNOTATION_OPT = "@knit"
	ANNOTATION_BEG = "@+knit"
	ANNOTATION_END = "@!knit"

	BEG_PATTERN    = `(.*)@\+knit(.*)`
	END_PATTERN    = `(.*)@!knit(.*)`
	OPTION_PATTERN = "@knit.(\\w*).([^`\\n]*)(?:`((.|\\n)*?(?<!\\\\))`)?"

	GROUP_OPTION_TYPE    = 1
	GROUP_OPTION_VALUE   = 2
	GROUP_OPTION_LITERAL = 3
)

// Option represents options read through the regex parser.
type Option struct {
	Type    string
	Value   string
	Literal string
}

// Options uses regexp to match all parser.OPTION_PATTERN in the given input
// text. Returns a slice of *Option representing the parsed results.
// This function will match everything including the option type, option value
// and any option literal defined.
func Options(input string) ([]*Option, error) {
	re2 := regexp2.MustCompile(OPTION_PATTERN, regexp2.None)

	opts := make([]*Option, 0)
	for m, err := re2.FindStringMatch(input); m != nil; m, err = re2.FindNextMatch(m) {
		if err != nil {
			return nil, errors.Wrap(err, "failed to find next knit option")
		}

		opt := &Option{
			Type:    m.GroupByNumber(GROUP_OPTION_TYPE).String(),
			Value:   m.GroupByNumber(GROUP_OPTION_VALUE).String(),
			Literal: m.GroupByNumber(GROUP_OPTION_LITERAL).String(),
		}

		if len(opt.Value) != 0 {
			opt.Value = os.ExpandEnv(opt.Value)
		}

		if len(opt.Literal) != 0 {
			opt.Literal = replaceEscaped(opt.Literal)
		}

		opts = append(opts, opt)
	}

	return opts, nil
}

func replaceEscaped(str string) string {
	return strings.ReplaceAll(str, "\\`", "`")
}

// BeginAnnotation uses regexp to match parser.BEG_PATTERN in the given
// input text. Returns an error if no match was found
func BeginAnnotation(input string) (string, error) {
	re := regexp2.MustCompile(BEG_PATTERN, regexp2.None)
	m, err := re.FindStringMatch(input)
	if err != nil {
		return "", err
	}
	if m == nil {
		return "", errors.New("did not match begin annotation")
	}

	return m.String(), nil
}

// EndAnnotation uses regexp to match parser.END_PATTERN in the given
// input text. Returns an error if no match was found
func EndAnnotation(input string) (string, error) {
	re := regexp2.MustCompile(END_PATTERN, regexp2.None)
	m, err := re.FindStringMatch(input)
	if err != nil {
		return "", err
	}
	if m == nil {
		return "", errors.New("did not match end annotation")
	}

	return m.String(), nil
}
