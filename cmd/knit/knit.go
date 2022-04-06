package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

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
			var wg sync.WaitGroup
			var parallel = c.Bool("parallel")

			if c.NArg() == 0 {
				return errors.New("at least one argument is required")
			}
			files := c.Args().Slice()

			k := knit.New(&knit.Config{
				Format:  c.Bool("format"),
				Verbose: c.Bool("verbose"),
			})

			for _, file := range files {
				if parallel {
					wg.Add(1)
					go knitWorker(c, &wg, k, file)
				} else {
					knitWorker(c, nil, k, file)
				}
			}

			wg.Wait()

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

func knitWorker(c *cli.Context, wg *sync.WaitGroup, k knit.Knit, file string) {
	startTime := time.Now()

	if wg != nil {
		defer wg.Done()
	}

	modified, err := k.ProcessFile(file)
	if c.Bool("verbose") {
		if err != nil {
			log.Printf("(%s) knit failed to process file: %s\n%+v",
				time.Since(startTime),
				file,
				err,
			)
		}

		if modified {
			log.Printf("(%s) knit modified file %s",
				time.Since(startTime),
				file,
			)
		}
	}
}
