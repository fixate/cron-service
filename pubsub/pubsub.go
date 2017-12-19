package pubsub

import (
	"log"

	mfst "github.com/fixate/cron-service/manifest"
	"github.com/urfave/cli"
)

type PubSubProvider struct {
	cli    *cli.Context
	client *pubSubClient

	Task *mfst.CronTaskDef
}

func NewProvider(cli *cli.Context, task *mfst.CronTaskDef) *PubSubProvider {
	return &PubSubProvider{
		cli:  cli,
		Task: task,
	}
}

func (p *PubSubProvider) ensureTopics() error {
	task := p.Task
	if err, _ := p.client.EnsureTopic(task.PubSub.Topic); err != nil {
		return err
	}
	return nil
}

func (p *PubSubProvider) Setup() error {
	projectId := p.cli.String("project-id")
	var client *pubSubClient
	var err error
	if err, client = NewClient(projectId); err != nil {
		return err
	}

	p.client = client

	if err := p.ensureTopics(); err != nil {
		return err
	}
	return nil
}

func (p *PubSubProvider) Handler() func() {
	var task mfst.CronTaskDef = *p.Task
	ps := task.PubSub
	return func() {
		log.Printf("[PUBSUB] Task start: '%s'\n", task.Description)

		log.Printf("[PUBSUB] Publishing topic: '%s'\n", ps.Topic)
		err, id := p.client.Publish(ps.Topic, ps.Message)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Published a message; msg ID: %v\n", id)
	}
}
