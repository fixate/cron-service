package pubsub

import (
	"encoding/base64"
	"golang.org/x/net/context"

	//"cloud.google.com/go/iam"
	"cloud.google.com/go/pubsub"
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

func (p *pubSubClient) Publish(topic string, msg string) (error, string) {
	t := p.client.Topic(topic)
	ctx := context.Background()
	encodedMsg := make([]byte, base64.StdEncoding.EncodedLen(len(msg)))
	base64.StdEncoding.Encode(encodedMsg, []byte(msg))
	result := t.Publish(ctx, &pubsub.Message{
		Data: encodedMsg,
	})

	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		return err, id
	}

	return nil, id
}
