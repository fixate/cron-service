package cron

import (
	"fmt"
	"log"

	mfst "github.com/fixate/cron-server/manifest"
	"github.com/fixate/cron-server/pubsub"
	"github.com/fixate/cron-server/request"

	"github.com/robfig/cron"
	"github.com/urfave/cli"
)

type Cron struct {
	cli      *cli.Context
	cron     *cron.Cron
	manifest *mfst.CronManifest
}

func New(c *cli.Context, manifest mfst.CronManifest) (error, *Cron) {
	gocron := cron.New()
	crn := Cron{c, gocron, &manifest}

	if err := crn.setupTasks(); err != nil {
		return err, nil
	}

	return nil, &crn
}

func (c *Cron) Run() {
	log.Println("Cron process started")
	c.cron.Run()
}

type ActionProvider interface {
	Setup() error
	Handler() func()
}

func (c *Cron) getProviderForTask(task *mfst.CronTaskDef) ActionProvider {
	if task.PubSub != nil {
		return pubsub.NewProvider(c.cli, task)
	}

	if task.Request != nil {
		return request.NewProvider(c.cli, task)
	}

	return nil
}

func (crn *Cron) setupTasks() error {
	for _, task := range *crn.manifest {
		if !task.Enabled {
			log.Printf("SKIPPING task '%s'. It is not enabled.\n", task.Description)
			continue
		}
		log.Printf("Adding task '%s'.\n", task.Description)

		provider := crn.getProviderForTask(&task)
		if provider == nil {
			return fmt.Errorf("Invalid manifest. Specify pubsub or request for cron task '%s'\n", task.Description)
		}

		if err := provider.Setup(); err != nil {
			return err
		}

		crn.cron.AddFunc(task.Schedule, provider.Handler())
	}
	return nil
}

func (c *Cron) newRequestHandler(task mfst.CronTaskDef) func() {
	return func() {
		log.Printf("[REQUEST] Task start: '%s'\n", task.Description)
		// TODO
	}
}
