package main

import (
	"fmt"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {

	var opt Option

	app := cli.NewApp()
	//生成命令行模式
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
			Value:       "./output/",
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
		var metaData *ServiceMetaData
		err := genMgr.Run(&opt, metaData)
		if err != nil {
			fmt.Sprintf("code generator failed, err:%v\n", err)
			return err
		}
		fmt.Println("code generator succ")
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
