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
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				return errors.New("at least one argument is required")
			}
			files := c.Args().Slice()

			k := knit.New()
			for _, file := range files {
				err := k.ProcessFile(file)
				if err != nil {
					log.Fatalf("knit failed to process file: %s\n%+v", file, err)
					return err
				}
			}

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
					gen, err := generator.FromOpts([]*parser.Option{
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
					})
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
