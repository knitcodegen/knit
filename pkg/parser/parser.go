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

type Options = []*Option

func Parse(input string) (Options, error) {
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
	str = strings.ReplaceAll(str, "\\`", "`")

	return str
}
