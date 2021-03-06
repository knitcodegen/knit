package knit

import (
	"os"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
)

func fromFile(t *testing.T, filepath string) string {
	byt, err := os.ReadFile(filepath)
	if err != nil {
		t.Fatalf("failed to read file in test. %v", err)
	}
	return string(byt)
}

func Test_Knit(t *testing.T) {
	type input = struct {
		knit Knit
		file string
	}

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
			name: "knits file successfully",
			input: input{
				knit: &knit{
					cfg: &Config{},
				},
				file: "./testdata/golden.go",
			},
			want: want{},
		},
		{
			name: "handles empty file",
			input: input{
				knit: &knit{
					cfg: &Config{},
				},
				file: "./testdata/empty",
			},
			want: want{},
		},
		{
			name: "handles only end annotation",
			input: input{
				knit: &knit{
					cfg: &Config{},
				},
				file: "./testdata/only_end_annotation",
			},
			want: want{},
		},
		{
			name: "handles no options configured",
			input: input{
				knit: &knit{
					cfg: &Config{},
				},
				file: "./testdata/no_options",
			},
			want: want{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res := c.input.knit.ProcessFile(c.input.file)
			if c.want.err {
				assert.Errorf(t, res.Error, c.want.errMessage)
			}
			cupaloy.SnapshotT(t, fromFile(t, c.input.file))
		})
	}
}
