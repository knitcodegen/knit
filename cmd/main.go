package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/knitgo/knit/pkg/generator"
	"github.com/knitgo/knit/pkg/knit"
	"github.com/knitgo/knit/pkg/parser"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func main() {
	(&cli.App{
		Name:  "knit",
		Usage: "find and execute knit generators in specified file glob",
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				return errors.New("at least one argument is required")
			}

			pattern := c.Args().First()
			matches, err := filepath.Glob(pattern)
			if err != nil {
				return errors.Wrap(err, "failed to match any files")
			}

			k := knit.New()
			for _, match := range matches {
				log.Printf("Knitting file: " + match)
				err := k.ProcessFile(match)
				if err != nil {
					log.Printf("Failed to execute knit against file: %+v", err)
					return errors.Wrap(err, "failed to knit file")
				}
			}

			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "generate",
				Aliases: []string{"gen"},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Required: false,
						Name:     "loader",
						Value:    "",
						Usage:    "loader type for input file",
					},
					&cli.PathFlag{
						Required: true,
						Name:     "input",
						Value:    "",
						Usage:    "input file",
					},
					&cli.PathFlag{
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
