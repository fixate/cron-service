package main

import (
	"log"
	"os"

	"github.com/fixate/cron-service/cron"
	mfst "github.com/fixate/cron-service/manifest"

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
		cli.BoolFlag{
			Name:  "ensure-topics-created",
			Usage: "Create pubsub topics if they don't exist",
		},
		cli.StringFlag{
			Name:   "project-id",
			Usage:  "project id",
			EnvVar: "CRON_GOOG_PROJECT_ID",
		},

		cli.StringFlag{
			Name:   "m, manifest-path",
			Usage:  "cron.yaml manifest file",
			EnvVar: "CRON_MANIFEST_PATH",
		},

		cli.StringFlag{
			Name:   "pubsub-emulator-host",
			Usage:  "Use the pubsub emulator",
			EnvVar: "PUBSUB_EMULATOR_HOST",
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

	err, crn := cron.New(c, manifest)
	if err != nil {
		return err
	}
	crn.Run()
	return nil
}
