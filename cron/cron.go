package cron

import (
	"errors"
	"log"

	mfst "github.com/fixate/cron-server/manifest"
	"github.com/fixate/cron-server/pubsub"
	"github.com/robfig/cron"
	"github.com/urfave/cli"
)

type Cron struct {
	*cli.Context
	cron     *cron.Cron
	manifest *mfst.CronManifest
}

func New(c *cli.Context, manifest mfst.CronManifest) (error, *Cron) {
	gocron := cron.New()
	crn := Cron{c, gocron, &manifest}

	if err := crn.setupTasks(); err != nil {
		return err, nil
	}

	if c.Bool("ensure-topics-created") {
		if err := crn.EnsureTopics(); err != nil {
			return err, nil
		}
	}
	return nil, &crn
}

func (c *Cron) Run() {
	log.Println("Cron process started")
	c.cron.Run()
}

func (c *Cron) EnsureTopics() error {
	projectId := c.String("project-id")
	_, client := pubsub.NewClient(projectId)
	for _, task := range *c.manifest {
		if task.PubSub != nil {
			if err, _ := client.CreateTopic(task.PubSub.Topic); err != nil {
				return err
			}
		}
	}
	return nil
}

func (crn *Cron) setupTasks() error {
	for _, task := range *crn.manifest {
		log.Printf("Adding task '%s'.\n", task.Description)
		err, fn := crn.getHandleFunc(task)
		if err != nil {
			return err
		}
		crn.cron.AddFunc(task.Schedule, fn)
	}
	return nil
}

func (crn *Cron) getHandleFunc(task mfst.CronTaskDef) (error, func()) {
	if task.PubSub != nil {
		return nil, crn.newPubSubHandler(task)
	}

	if task.Request != nil {
		return nil, crn.newRequestHandler(task)
	}

	return errors.New("Invalid manifest. Specify pubsub or request for cron task"), nil
}

func (c *Cron) newPubSubHandler(task mfst.CronTaskDef) func() {
	projectId := c.String("project-id")
	return func() {
		log.Printf("[PUBSUB] Task start: '%s'\n", task.Description)
		ps := task.PubSub
		err, client := pubsub.NewClient(projectId)

		if err != nil {
			log.Fatal(err)
		}

		log.Printf("[PUBSUB] Publishing topic: '%s'\n", ps.Topic)
		err, id := client.Publish(ps.Topic, ps.Message)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Published a message; msg ID: %v\n", id)
	}
}

func (c *Cron) newRequestHandler(task mfst.CronTaskDef) func() {
	return func() {
		log.Printf("[REQUEST] Task start: '%s'\n", task.Description)
		// TODO
	}
}
