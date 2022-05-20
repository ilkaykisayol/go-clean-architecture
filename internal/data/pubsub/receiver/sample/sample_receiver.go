package sample

import (
	"context"
	"encoding/base64"
	receiver2 "go-clean-architecture/internal/data/pubsub/receiver"
	"strconv"
	"time"

	"google.golang.org/api/option"

	"go-clean-architecture/internal/util/cacher"
	"go-clean-architecture/internal/util/env"
	"go-clean-architecture/internal/util/logger"
	"go-clean-architecture/internal/util/validator"

	"cloud.google.com/go/pubsub"
	"go.uber.org/zap"
)

type ISampleReceiver interface {
	InitReceivers(count int)
}

type SampleReceiver struct {
	environment           env.IEnvironment
	loggr                 logger.ILogger
	validatr              validator.IValidator
	cachr                 cacher.ICacher
	projectId             string
	subscriptionId        string
	timeout               time.Duration
	defaultAttributes     map[string]string
	messageCacheKeyPrefix string
	handler               ISampleReceiverHandler
	receiverName          string
}

// NewSampleReceiver
// Returns a new SampleReceiver.
func NewSampleReceiver(environment env.IEnvironment, loggr logger.ILogger, validatr validator.IValidator, cachr cacher.ICacher, handler ISampleReceiverHandler) ISampleReceiver {
	projectId := environment.Get(env.SampleReceiverProjectId)
	subscriptionId := environment.Get(env.SampleReceiverSubscriptionId)
	receiverName := "SampleReceiver"

	receiver := SampleReceiver{
		environment: environment,
		loggr: loggr.With(
			zap.String("receiverName", receiverName),
			zap.String("projectId", projectId),
			zap.String("subscriptionId", subscriptionId),
		),
		validatr:              validatr,
		cachr:                 cachr,
		projectId:             projectId,
		subscriptionId:        subscriptionId,
		timeout:               time.Second * 5,
		defaultAttributes:     map[string]string{receiver2.DummyAttribute: "dummy_attribute_value"},
		messageCacheKeyPrefix: "sample-receiver-message-id",
		receiverName:          receiverName,
	}

	if handler != nil {
		receiver.handler = handler
	} else {
		receiver.handler = NewSampleReceiverHandler(environment, loggr, validatr, cachr, nil)
	}

	return &receiver
}

// InitReceivers
// Initializes multiple receivers that listen given topic for new messages.
func (r *SampleReceiver) InitReceivers(count int) {
	r.loggr.Info(r.receiverName + " Initializing receivers.")
	for i := 0; i < count; i++ {
		go r.receive()
	}
	r.loggr.Info(r.receiverName + " Initialized " + strconv.Itoa(count) + " receivers.")
}

func (r *SampleReceiver) receive() {
	// Register another receiver if one of them fails.
	defer func() {
		if rec := recover(); rec != nil {
			r.loggr.Error(r.receiverName+" Recovered the panic. Trying to receive again.", zap.Any("panic", rec))
			r.receive()
		}
	}()

	saJson, err := base64.StdEncoding.DecodeString(r.environment.Get(env.SampleReceiverSaJson))
	if err != nil {
		r.loggr.Panic(r.receiverName + " Panicked while SampleReceiverSaJson decoding base64.")
	}

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, r.projectId, option.WithCredentialsJSON(saJson))
	if err != nil {
		r.loggr.Panic(r.receiverName + " Panicked while creating pub sub client.")
	}
	defer client.Close()

	subscription := client.Subscription(r.subscriptionId)
	err = subscription.Receive(ctx, r.eventHandler)
	if err != nil {
		r.loggr.Panic(r.receiverName + " Panicked while receiving messages from subscription.")
	}
}

func (r *SampleReceiver) eventHandler(ctx context.Context, msg *pubsub.Message) {
	defer func() {
		if rec := recover(); rec != nil {
			r.loggr.Error(r.receiverName+" "+msg.ID+" ID message is panicked during execution in event handler.",
				zap.String("messageId", msg.ID),
				zap.String("data", string(msg.Data)),
				zap.Any("attributes", msg.Attributes),
				zap.Any("panic", rec),
			)
		}
	}()

	// Cache message id for two days to prevent duplication.
	existingMessageId := r.cachr.Get(r.getMessageCacheKey(msg.ID))
	if existingMessageId != nil {
		r.loggr.Error(r.receiverName+" "+msg.ID+" ID message is duplicate.",
			zap.String("messageId", msg.ID),
			zap.String("data", string(msg.Data)),
			zap.Any("attributes", msg.Attributes),
		)

		msg.Ack()
		return
	}

	ch := make(chan error)
	defer close(ch)
	go r.handler.Handle(ch, &receiver2.ReceiverHandlerModel{
		MessageId:  msg.ID,
		Data:       msg.Data,
		Attributes: msg.Attributes,
	})

	err := <-ch
	if err != nil {
		r.loggr.Error(r.receiverName+" "+msg.ID+" ID message is failed to process.",
			zap.String("messageId", msg.ID),
			zap.String("data", string(msg.Data)),
			zap.Any("attributes", msg.Attributes),
			zap.Error(err),
		)
		return
	}

	// Success
	msg.Ack()
	r.cachr.Set(r.getMessageCacheKey(msg.ID), "", time.Hour*24*2)
	r.loggr.Info(r.receiverName+" "+msg.ID+" ID message is processed successfully.",
		zap.String("messageId", msg.ID),
		zap.String("data", string(msg.Data)),
		zap.Any("attributes", msg.Attributes),
	)
}

// Returns a cache key for a specific message id.
func (r *SampleReceiver) getMessageCacheKey(messageId string) string {
	return r.messageCacheKeyPrefix + ":" + messageId
}
