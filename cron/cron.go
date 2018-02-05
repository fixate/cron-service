package cron

import (
	"fmt"
	"log"

	mfst "github.com/fixate/cron-service/manifest"
	"github.com/fixate/cron-service/pubsub"
	"github.com/fixate/cron-service/request"

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
	log.Println("Cron process starting")
	c.cron.Run()
}

type ActionProvider interface {
	Name() string
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

		provider := crn.getProviderForTask(&task)
		if provider == nil {
			return fmt.Errorf("Invalid manifest. Specify pubsub or request for cron task '%s'\n", task.Description)
		}
		log.Printf("[%s] Adding task '%s'.\n", provider.Name(), task.Description)

		if err := provider.Setup(); err != nil {
			log.Fatal(err)
			return err
		}

		log.Printf("[%s] Schedule '%s'.\n", provider.Name(), task.Schedule)

		log.Printf("[%s] debug fireOnStart=%t.\n", provider.Name(), task.FireOnStart)
		handler := provider.Handler()
		if task.FireOnStart {
			log.Printf("[%s] Fire On Start.\n", provider.Name())
			go handler()
		}

		crn.cron.AddFunc(task.Schedule, handler)
	}
	return nil
}
