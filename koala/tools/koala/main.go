package main

import (
	"github.com/urfave/cli"
	"log"
	"os"
)

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
		cli.StringFlag{
			Name:        "o",
			Usage:       "output directory",
			Destination: &opt.Output,
		},
		cli.BoolFlag{
			Name:        "c",
			Usage:       "generate grpc client code",
			Destination: &opt.GenClientCode,
		},
		cli.BoolFlag{
			Name:        "s",
			Usage:       "generate grpc server code",
			Destination: &opt.GenServerCode,
		},
	}

	app.Action = func(c *cli.Context) error {
		//name := "someone"
		//if c.NArg() > 0 {
		//	name = c.Args()[0]
		//}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
