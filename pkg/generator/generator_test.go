package generator

import (
	"os"
	"testing"
	"text/template"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/knitcodegen/knit/pkg/loader"
	"github.com/stretchr/testify/assert"
)

func fromFile(t *testing.T, filepath string) string {
	byt, err := os.ReadFile(filepath)
	if err != nil {
		t.Fatalf("failed to read file in test. %v", err)
	}
	return string(byt)
}

func Test_Generate(t *testing.T) {
	goldenTemplate, err := template.New("golden").
		Parse(fromFile(t, "./testdata/templates/golden.tmpl"))
	assert.NoError(t, err)

	graphqlTemplate, err := template.New("golden_graphql").
		Parse(fromFile(t, "./testdata/templates/golden_graphql.tmpl"))
	assert.NoError(t, err)

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
			name: "handles json input",
			input: &generator{
				Input:     fromFile(t, "./testdata/inputs/golden.json"),
				Loader:    &loader.JsonLoader{},
				Templater: goldenTemplate,
			},
			want: want{},
		},
		{
			name: "handles yaml input",
			input: &generator{
				Input:     fromFile(t, "./testdata/inputs/golden.yml"),
				Loader:    &loader.YamlLoader{},
				Templater: goldenTemplate,
			},
			want: want{},
		},
		{
			name: "handles graphql input",
			input: &generator{
				Input:     fromFile(t, "./testdata/inputs/golden.graphql"),
				Loader:    &loader.GraphqlLoader{},
				Templater: graphqlTemplate,
			},
			want: want{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			codegen, err := c.input.Generate()
			if c.want.err {
				assert.Errorf(t, err, c.want.errMessage)
			}
			cupaloy.SnapshotT(t, codegen)
		})
	}
}
