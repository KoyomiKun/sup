package main

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/KoyomiKun/sup/pkg/config"
)

var (
	configPath string
)

func main() {
	app := &cli.App{
		Name:        "sup",
		Description: "Keep process in live && rotate logs",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config-path",
				Aliases:     []string{"c"},
				Usage:       "Config Path.",
				Destination: &configPath,
			},
		},
		Action: func(c *cli.Context) error {
			// init config
			config := config.NewConfig(configPath)

			if c.NArg() == 0 {
				// init server

			} else {
				// init client

			}
			return nil
		},
	}
	app.Run(os.Args)
}
