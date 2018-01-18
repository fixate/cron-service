package pubsub

import (
	"golang.org/x/net/context"
	"log"

	//"cloud.google.com/go/iam"
	"cloud.google.com/go/pubsub"

	mfst "github.com/fixate/cron-service/manifest"
)

type pubSubClient struct {
	client *pubsub.Client
}

func NewClient(projectId string) (error, *pubSubClient) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectId)
	if err != nil {
		return err, nil
	}

	return nil, &pubSubClient{client}
}

func (p *pubSubClient) EnsureTopic(topicName string) (err error, topic *pubsub.Topic) {
	ctx := context.Background()
	topic, err = p.client.CreateTopic(ctx, topicName)
	if err != nil {
		topic = p.client.Topic(topicName)
		err = nil
	}
	return
}

func (p *pubSubClient) Publish(ps *mfst.PubSubDef) (error, string) {
	t := p.client.Topic(ps.Topic)

	log.Printf("[PUBSUB] Publishing topic: '%s'\n", t)
	ctx := context.Background()
	data := []byte(ps.Message)
	if len(data) == 0 {
		data = []byte{0}
	}
	result := t.Publish(ctx, &pubsub.Message{
		Data:       data,
		Attributes: ps.Attributes,
	})

	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		return err, id
	}

	return nil, id
}
