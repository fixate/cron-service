package main

import (
	"log"
	"os"

	"github.com/fixate/cron-service/cron"
	mfst "github.com/fixate/cron-service/manifest"

	"github.com/urfave/cli"
	//"net/http"
	//_ "net/http/pprof"
)

const version string = "0.2.0"

func main() {
	//go func() {
	//log.Println(http.ListenAndServe("localhost:6060", nil))
	//}()
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
			Name:   "credentials-file",
			Usage:  "Service account or refresh token JSON credentials file",
			EnvVar: "CRON_CREDENTIALS_FILE",
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

	err, crn := cron.New(c, manifest)
	if err != nil {
		return err
	}
	crn.Run()
	return nil
}
