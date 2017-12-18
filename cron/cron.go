package cron

import (
	"errors"
	"fmt"

	mfst "github.com/fixate/cron-server/manifest"
	"github.com/robfig/cron"
)

func Run(manifest mfst.CronManifest) error {
	c := cron.New()
	if err := setupTasks(c, manifest); err != nil {
		return err
	}
	fmt.Println("Cron process started")
	c.Run()
	return nil
}

func setupTasks(c *cron.Cron, manifest mfst.CronManifest) error {
	for _, task := range manifest {
		fmt.Printf("Adding task '%s'.\n", task.Description)
		err, fn := getHandleFunc(task)
		if err != nil {
			return err
		}
		c.AddFunc(task.Schedule, fn)
	}
	return nil
}

func getHandleFunc(task mfst.CronTaskDef) (error, func()) {
	if task.PubSub != nil {
		return nil, newPubSubHandler(task)
	}

	if task.Request != nil {
		return nil, newRequestHandler(task)
	}

	return errors.New("Invalid manifest. Specify pubsub or request for cron task"), nil
}

func newPubSubHandler(task mfst.CronTaskDef) func() {
	return func() {
		fmt.Printf("[PUBSUB] Task start: '%s'\n", task.Description)
		pubsub := task.PubSub
		fmt.Printf("[PUBSUB] Publishing topic: '%s'\n", pubsub.Topic)

		// TODO
	}
}

func newRequestHandler(task mfst.CronTaskDef) func() {
	return func() {
		fmt.Printf("[REQUEST] Task start: '%s'\n", task.Description)
		// TODO
	}
}
