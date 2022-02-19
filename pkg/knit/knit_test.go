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
				knit: &knit{},
				file: "./testdata/golden.go",
			},
			want: want{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.input.knit.ProcessFile(c.input.file)
			if c.want.err {
				assert.Errorf(t, err, c.want.errMessage)
			}
			cupaloy.SnapshotT(t, fromFile(t, c.input.file))
		})
	}
}
