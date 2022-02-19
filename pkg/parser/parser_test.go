package parser_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tylermmorton/knit/pkg/parser"
)

const (
	TEST_FILENAME = "./test.yml"
)

func setupEnvironmentVars(t *testing.T) {
	err := os.Setenv("TEST_FILENAME", TEST_FILENAME)
	assert.NoError(t, err)
}

func Test_Options(t *testing.T) {
	setupEnvironmentVars(t)

	type input = string

	type want struct {
		opts       []*parser.Option
		err        bool
		errMessage string
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			name:  "handles empty input block",
			input: "",
			want: want{
				opts: []*parser.Option{},
			},
		},
		{
			name:  "parses option successfully",
			input: "@knit input test.yml",
			want: want{
				opts: []*parser.Option{
					{
						Type:  "input",
						Value: "test.yml",
					},
				},
			},
		},
		{
			name:  "parses option with literal",
			input: "@knit input yml`hello world`",
			want: want{
				opts: []*parser.Option{
					{
						Type:    "input",
						Value:   "yml",
						Literal: "hello world",
					},
				},
			},
		},
		{
			name:  "handles escaped backtick in literal",
			input: "@knit input yml`hello\\`world\\``",
			want: want{
				opts: []*parser.Option{
					{
						Type:    "input",
						Value:   "yml",
						Literal: "hello`world`",
					},
				},
			},
		},
		{
			name:  "expands environment variable in option value",
			input: "@knit input $TEST_FILENAME",
			want: want{
				opts: []*parser.Option{
					{
						Type:  "input",
						Value: TEST_FILENAME,
					},
				},
			},
		},
	}

	for _, c := range cases {
		opts, err := parser.Options(c.input)
		if c.want.err {
			assert.Errorf(t, err, c.want.errMessage)
		}
		assert.Equal(t, c.want.opts, opts)
	}
}

func Test_EndAnnotation(t *testing.T) {
	type input = string

	type want struct {
		match      string
		err        bool
		errMessage string
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			name:  "parses successfully",
			input: "// @!knit **abc",
			want: want{
				match: "// @!knit **abc",
			},
		},
		{
			name:  "handles empty input",
			input: "",
			want: want{
				err:        true,
				errMessage: "did not match end annotation",
			},
		},
	}
	for _, c := range cases {
		match, err := parser.EndAnnotation(c.input)
		if c.want.err {
			assert.Error(t, err, c.want.errMessage)
		}
		assert.Equal(t, c.want.match, match)
	}
}
