package awspub

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/foodora/go-ranger/pubsub"
)

// publisher will accept AWS configuration and an SNS topic name
// and it will emit any publish events to it.
type publisher struct {
	sns    snsiface.SNSAPI
	topic  string
	Logger pubsub.Logger
}

// NewPublisher will initiate the SNS client.
func NewPublisher(cfg SNSConfig) (pubsub.Publisher, error) {
	p := &publisher{}
	p.Logger = pubsub.DefaultLogger

	if cfg.Topic == "" {
		return p, errors.New("SNS topic name is required")
	}
	p.topic = cfg.Topic

	if cfg.Region == nil {
		return p, errors.New("SNS region is required")
	}

	sess, err := session.NewSession()
	if err != nil {
		return p, err
	}

	p.sns = sns.New(sess, &aws.Config{
		Region:   cfg.Region,
		Endpoint: cfg.Endpoint, //optional
	})
	return p, nil
}

// Publish send the message to the SNS topic.
// The key will be used as the SNS message subject which is optional.
func (p *publisher) Publish(ctx context.Context, key string, m string) error {
	msg := &sns.PublishInput{
		TopicArn: &p.topic,
		Subject:  &key, //optional
		Message:  aws.String(m),
	}

	_, err := p.sns.Publish(msg)
	return err
}
