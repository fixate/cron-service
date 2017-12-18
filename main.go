package main

import (
	"log"
	"os"

	"github.com/fixate/cron-server/cron"
	mfst "github.com/fixate/cron-server/manifest"

	"github.com/urfave/cli"
)

const version string = "1.0.0"

func main() {
	app := cli.NewApp()
	app.Name = "Cron service"
	app.Version = version
	app.Usage = "Add a cron yaml file to trigger time-based events"
	app.EnableBashCompletion = true
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "key-file",
			Usage:  "optional keyfile to use for pub sub",
			EnvVar: "CRON_KEYFILE",
		},
		cli.StringFlag{
			Name:   "m, manifest-path",
			Usage:  "cron.yaml manifest file",
			EnvVar: "CRON_MANIFEST_PATH",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Println("Fatal Error")
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	manifest, err := mfst.Load(c.String("manifest-path"))
	if err != nil {
		return err
	}

	return cron.Run(manifest)
}
