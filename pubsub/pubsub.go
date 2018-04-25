package pubsub

import (
	"log"

	"github.com/urfave/cli"

	mfst "github.com/fixate/cron-service/manifest"
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

func (p *PubSubProvider) Name() string {
	return "PUBSUB"
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
	credentialsFile := p.cli.String("credentials-file")
	var client *pubSubClient
	var err error
	log.Printf("[%s] Creating Client for project %s.\n", p.Name(), projectId)
	if err, client = NewClient(projectId, credentialsFile); err != nil {
		return err
	}
	log.Printf("[%s] New Client created.\n", p.Name())

	p.client = client

	if p.cli.Bool("ensure-topics-created") {
		if err := p.ensureTopics(); err != nil {
			return err
		}
	}
	log.Printf("[%s] Setup complete.\n", p.Name())
	return nil
}

func (p *PubSubProvider) Handler() func() {
	var task mfst.CronTaskDef = *p.Task
	ps := task.PubSub
	topic := p.client.Topic(ps.Topic)
	return func() {
		log.Printf("[PUBSUB] Task start: '%s'\n", task.Description)

		err, id := p.client.Publish(topic, ps)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Published a message; msg ID: %v\n", id)
	}
}
