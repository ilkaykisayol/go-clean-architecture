package sample

import (
	"encoding/json"
	"go-clean-architecture/internal/data/pubsub/receiver"
	"go-clean-architecture/internal/service/sample"

	"go-clean-architecture/internal/util/cacher"
	"go-clean-architecture/internal/util/env"
	"go-clean-architecture/internal/util/logger"
	"go-clean-architecture/internal/util/validator"
)

type ISampleReceiverHandler interface {
	Handle(ch chan error, model *receiver.ReceiverHandlerModel)
}

type SampleReceiverHandler struct {
	environment   env.IEnvironment
	loggr         logger.ILogger
	validatr      validator.IValidator
	cachr         cacher.ICacher
	sampleService sample.ISampleService
}

// NewSampleReceiverHandler
// Returns a new SampleReceiverHandler
func NewSampleReceiverHandler(environment env.IEnvironment, loggr logger.ILogger, validatr validator.IValidator, cachr cacher.ICacher, sampleService sample.ISampleService) ISampleReceiverHandler {
	handler := SampleReceiverHandler{
		environment: environment,
		loggr:       loggr,
		validatr:    validatr,
		cachr:       cachr,
	}

	if sampleService != nil {
		handler.sampleService = sampleService
	} else {
		handler.sampleService = sample.NewSampleService(environment, loggr, validatr, cachr, nil, nil, nil, nil)
	}

	return &handler
}

// Handle
// Process events & messages and handle necessary logic/business actions.
func (h *SampleReceiverHandler) Handle(ch chan error, model *receiver.ReceiverHandlerModel) {
	err := h.validatr.ValidateStruct(model)
	if err != nil {
		ch <- err
		return
	}

	var handlerModel SampleReceiverHandlerModel
	err = json.Unmarshal(model.Data, &handlerModel)
	if err != nil {
		ch <- err
		return
	}

	err = h.validatr.ValidateStruct(handlerModel)
	if err != nil {
		ch <- err
		return
	}

	sampleCh := make(chan *sample.UpdateSampleServiceResponse)
	defer close(sampleCh)
	go h.sampleService.UpdateSample(sampleCh, &sample.UpdateSampleServiceModel{
		SampleId:     handlerModel.SampleId,
		SampleStatus: handlerModel.SampleStatus,
		ModifiedBy:   handlerModel.ModifiedBy,
	})

	response := <-sampleCh

	if response.Error != nil {
		ch <- response.Error
		return
	}

	ch <- nil
}
