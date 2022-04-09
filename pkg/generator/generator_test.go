package generator

import (
	"os"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
)

var (
	inputFileMissing    = "./testdata/does_not_exist.json"
	templateFileMissing = "./testdata/does_not_exist.tmpl"

	inputFileJson    = "./testdata/inputs/golden.json"
	inputFileYaml    = "./testdata/inputs/golden.yml"
	inputFileGraphql = "./testdata/inputs/golden.graphql"

	inputTmplFileGolden        = "./testdata/templates/golden.tmpl"
	inputTmplFileGoldenGraphql = "./testdata/templates/golden_graphql.tmpl"
)

func fromFile(t *testing.T, filepath string) string {
	byt, err := os.ReadFile(filepath)
	if err != nil {
		t.Fatalf("failed to read file in test. %v", err)
	}
	return string(byt)
}

func Test_Generate(t *testing.T) {
	type input = Generator

	type want struct {
		err        bool
		errMessage string
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			name: "handles missing loader type",
			input: &generator{
				InputLiteral:    fromFile(t, inputFileJson),
				TemplateLiteral: fromFile(t, inputTmplFileGolden),
			},
			want: want{
				err:        true,
				errMessage: "missing loader type",
			},
		},
		{
			name: "handles missing input",
			input: &generator{
				LoaderType:      "json",
				TemplateLiteral: fromFile(t, inputTmplFileGolden),
			},
			want: want{
				err:        true,
				errMessage: "missing input",
			},
		},
		{
			name: "handles missing template",
			input: &generator{
				LoaderType:   "json",
				InputLiteral: fromFile(t, inputFileJson),
			},
			want: want{
				err:        true,
				errMessage: "missing template",
			},
		},
		{
			name: "handles failure to load input file",
			input: &generator{
				LoaderType: "json",
				InputFile:  &inputFileMissing,
			},
			want: want{
				err:        true,
				errMessage: "failed to load input file",
			},
		},
		{
			name: "handles failure to load template file",
			input: &generator{
				LoaderType:   "json",
				InputLiteral: fromFile(t, inputFileJson),
				TemplateFile: &templateFileMissing,
			},
			want: want{
				err:        true,
				errMessage: "failed to load template file",
			},
		},
		{
			name: "handles failure to create loader",
			input: &generator{
				LoaderType:   "bad",
				InputFile:    &inputFileJson,
				TemplateFile: &inputTmplFileGolden,
			},
			want: want{
				err:        true,
				errMessage: "failed to create loader",
			},
		},
		{
			name: "handles json input literal",
			input: &generator{
				LoaderType:      "json",
				InputLiteral:    fromFile(t, inputFileJson),
				TemplateLiteral: fromFile(t, inputTmplFileGolden),
			},
			want: want{},
		},
		{
			name: "handles json input file",
			input: &generator{
				LoaderType:   "json",
				InputFile:    &inputFileJson,
				TemplateFile: &inputTmplFileGolden,
			},
			want: want{},
		},
		{
			name: "handles yaml input literal",
			input: &generator{
				LoaderType:      "yaml",
				InputLiteral:    fromFile(t, inputFileYaml),
				TemplateLiteral: fromFile(t, inputTmplFileGolden),
			},
			want: want{},
		},
		{
			name: "handles yaml input file",
			input: &generator{
				LoaderType:   "yaml",
				InputFile:    &inputFileYaml,
				TemplateFile: &inputTmplFileGolden,
			},
			want: want{},
		},
		{
			name: "handles graphql input literal",
			input: &generator{
				LoaderType:      "graphql",
				InputLiteral:    fromFile(t, inputFileGraphql),
				TemplateLiteral: fromFile(t, inputTmplFileGoldenGraphql),
			},
			want: want{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			codegen, err := c.input.Generate()
			if c.want.err {
				assert.Errorf(t, err, c.want.errMessage)
			} else {
				cupaloy.SnapshotT(t, codegen)
			}
		})
	}
}
