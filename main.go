package main

import (
	"github.com/lpisces/marketboy/cmds/boot"
	"gopkg.in/urfave/cli.v1"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "marketboy"
	app.Usage = "tools for market maker"

	app.Commands = []cli.Command{
		{
			Name:    "boot",
			Aliases: []string{"b"},
			Usage:   "cmd demo",
			Action:  boot.Run,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "debug, d",
					Usage: "debug switch",
				},
				cli.StringFlag{
					Name:  "config, c",
					Usage: "load config file",
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
