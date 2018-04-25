package pubsub

import (
	"golang.org/x/net/context"
	"log"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"

	mfst "github.com/fixate/cron-service/manifest"
)

type pubSubClient struct {
	context context.Context
	client  *pubsub.Client
}

func NewClient(projectId string, credentialsFile string) (error, *pubSubClient) {
	ctx := context.Background()
	var client *pubsub.Client
	var err error
	if len(credentialsFile) > 0 {
		clientOptions := option.WithCredentialsFile(credentialsFile)
		client, err = pubsub.NewClient(ctx, projectId, clientOptions)

		if err != nil {
			return err, nil
		}

		return nil, &pubSubClient{ctx, client}
	}

	client, err = pubsub.NewClient(ctx, projectId)

	if err != nil {
		return err, nil
	}

	return nil, &pubSubClient{ctx, client}
}

func (p *pubSubClient) EnsureTopic(topicName string) (err error, topic *pubsub.Topic) {
	topic, err = p.client.CreateTopic(p.context, topicName)
	if err != nil {
		log.Printf("[PUBSUB CLIENT] '%s' topic exists.\n", topicName)
		topic = p.client.Topic(topicName)
		err = nil
	} else {
		log.Printf("[PUBSUB CLIENT] '%s' topic created.\n", topicName)
	}
	return
}

func (p *pubSubClient) Topic(topicName string) *pubsub.Topic {
	return p.client.Topic(topicName)
}

func (p *pubSubClient) Publish(topic *pubsub.Topic, ps *mfst.PubSubDef) (error, string) {
	log.Printf("[PUBSUB CLIENT] Publishing message to topic: '%s'\n", topic)
	data := []byte(ps.Message)
	if len(data) == 0 {
		data = []byte{0}
	}
	result := topic.Publish(p.context, &pubsub.Message{
		Data:       data,
		Attributes: ps.Attributes,
	})

	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(p.context)
	if err != nil {
		return err, id
	}

	return nil, id
}
