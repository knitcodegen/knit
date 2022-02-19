package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/tylermmorton/gocodeshift/pkg/knit"
	"github.com/urfave/cli/v2"
)

func main() {
	(&cli.App{
		Name:  "knit",
		Usage: "execute knit code generators in specified files",
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				return errors.New("no arguments specified")
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
	}).Run(os.Args)
}
