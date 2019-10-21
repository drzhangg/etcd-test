package main

import "github.com/urfave/cli"

func main() {

	var opt Option

	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "f",
			Usage:       "idl filename",
			Value:       "./test.proto",
			Destination: &opt.Proto3Filename,
		},
	}
}
