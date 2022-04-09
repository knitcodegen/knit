package main

import (
	"fmt"
	"log"
	"os"

	"github.com/knitcodegen/knit/pkg/generator"
	"github.com/knitcodegen/knit/pkg/knit"
	"github.com/knitcodegen/knit/pkg/parser"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

// ldflags set by goreleaser
var (
	version = "dev"
	commit  = "none"
)

func main() {
	(&cli.App{
		Name:    "knit",
		Usage:   "language & schema agnostic code generation toolkit",
		Version: fmt.Sprintf("%s\n%s", version, commit),

		UsageText: "DEFAULT: knit ./**/*.gen.go\n\t COMMAND: knit [global options] command [command options] [arguments...]",
		// Default Action
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "format",
				Usage:   "Enable auto-formatting of .go source files",
				Aliases: []string{"f"},
				Value:   true,
			},
			&cli.BoolFlag{
				Name:  "verbose",
				Usage: "Enable verbose logging",
				Value: false,
			},
			&cli.BoolFlag{
				Name:  "parallel",
				Usage: "Enable parallel file processing",
				Value: true,
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				return errors.New("at least one argument is required")
			}

			files := c.Args().Slice()

			k := knit.New(&knit.Config{
				Format:   c.Bool("format"),
				Verbose:  c.Bool("verbose"),
				Parallel: c.Bool("parallel"),
			})

			k.ProcessFiles(files, func(res knit.ProcessResult) {
				if res.Error != nil {
					log.Printf("knit failed to process file: %s\n%+v", res.File, res.Error)
				} else {
					log.Printf("knit processed file successfully: %s", res.File)
				}
			})

			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "generate",
				Usage:   "Runs the knit code generator using the specified options",
				Aliases: []string{"gen"},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Required: false,
						Aliases:  []string{"l"},
						Name:     "loader",
						Value:    "",
						Usage:    "loader type for input file",
					},
					&cli.PathFlag{
						Aliases:  []string{"i"},
						Required: true,
						Name:     "input",
						Value:    "",
						Usage:    "input file",
					},
					&cli.PathFlag{
						Aliases:  []string{"t"},
						Required: true,
						Name:     "template",
						Value:    "",
						Usage:    "template file",
					},
				},
				Action: func(c *cli.Context) error {
					opts := []*parser.Option{
						{
							Type:  "loader",
							Value: c.String("loader"),
						},
						{
							Type:  "input",
							Value: c.Path("input"),
						},
						{
							Type:  "template",
							Value: c.Path("template"),
						},
					}

					gen, err := generator.New(opts...)
					if err != nil {
						return err
					}

					codegen, err := gen.Generate()
					if err != nil {
						return err
					}

					_, err = c.App.Writer.Write([]byte(codegen))
					if err != nil {
						return err
					}

					return nil
				},
			},
		},
	}).Run(os.Args)
}
