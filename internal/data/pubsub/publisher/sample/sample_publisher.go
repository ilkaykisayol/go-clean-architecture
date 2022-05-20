package sample

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"go-clean-architecture/internal/data/pubsub/publisher"
	"google.golang.org/api/option"
	"time"

	"go-clean-architecture/internal/util/cacher"
	"go-clean-architecture/internal/util/env"
	"go-clean-architecture/internal/util/logger"
	"go-clean-architecture/internal/util/validator"

	"cloud.google.com/go/pubsub"
	"go.uber.org/zap"
)

type ISamplePublisher interface {
	Publish(ch chan publisher.PublisherResponse, message *SamplePublisherModel, attributeOverrides map[string]string)
}

type SamplePublisher struct {
	environment       env.IEnvironment
	loggr             logger.ILogger
	validatr          validator.IValidator
	cachr             cacher.ICacher
	projectId         string
	topicId           string
	timeout           time.Duration
	defaultAttributes map[string]string
}

// NewSamplePublisher
// Returns a new SamplePublisher.
func NewSamplePublisher(environment env.IEnvironment, loggr logger.ILogger, validatr validator.IValidator, cachr cacher.ICacher) ISamplePublisher {
	projectId := environment.Get(env.SamplePublisherProjectId)
	topicId := environment.Get(env.SamplePublisherTopicId)

	publisher := SamplePublisher{
		environment: environment,
		loggr: loggr.With(
			zap.String("projectId", projectId),
			zap.String("topicId", topicId),
		),
		validatr:          validatr,
		cachr:             cachr,
		projectId:         projectId,
		topicId:           topicId,
		timeout:           time.Second * 5,
		defaultAttributes: map[string]string{publisher.DummyAttribute: "dummy_attribute_value"},
	}

	return &publisher
}

func (p *SamplePublisher) Publish(ch chan publisher.PublisherResponse, model *SamplePublisherModel, attributeOverrides map[string]string) {
	err := p.validatr.ValidateStruct(model)
	if err != nil {
		p.loggr.Error(err.Error())
		ch <- publisher.PublisherResponse{Error: err}
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	saJson, err := base64.StdEncoding.DecodeString(p.environment.Get(env.SamplePublisherSaJson))
	if err != nil {
		p.loggr.Panic("Panicked while decoding SamplePublisherSaJson base64.")
	}

	client, err := pubsub.NewClient(ctx, p.projectId, option.WithCredentialsJSON(saJson))
	if err != nil {
		p.loggr.Error(err.Error())
		ch <- publisher.PublisherResponse{Error: err}
		return
	}
	defer client.Close()

	bytes, err := json.Marshal(model)
	if err != nil {
		p.loggr.Error(err.Error())
		ch <- publisher.PublisherResponse{Error: err}
		return
	}

	attributes := p.overrideAttributes(attributeOverrides)

	topic := client.Topic(p.topicId)
	result := topic.Publish(ctx, &pubsub.Message{
		Data:       bytes,
		Attributes: attributes,
	})

	messageId, err := result.Get(ctx)
	if err != nil {
		p.loggr.Error(err.Error())
		ch <- publisher.PublisherResponse{Error: err}
		return
	}

	p.loggr.Info(messageId+" ID message is published successfully.", zap.String("messageId", messageId))
	ch <- publisher.PublisherResponse{MessageId: &messageId}
}

func (p *SamplePublisher) overrideAttributes(overrides map[string]string) map[string]string {
	attr := make(map[string]string)

	for key, value := range p.defaultAttributes {
		attr[key] = value
	}

	for key, value := range overrides {
		attr[key] = value
	}

	return attr
}
