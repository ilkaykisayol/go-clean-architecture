package sample

import (
	"encoding/json"
	"go-clean-architecture/internal/data/database/sample"
	sample_proxy "go-clean-architecture/internal/data/proxy/sample"
	sample_publisher "go-clean-architecture/internal/data/pubsub/publisher/sample"
	"time"

	"go-clean-architecture/internal/data/pubsub/publisher"
	"go-clean-architecture/internal/util/cacher"
	"go-clean-architecture/internal/util/env"
	"go-clean-architecture/internal/util/logger"
	"go-clean-architecture/internal/util/validator"
)

type ISampleService interface {
	GetGoogle(ch chan *GetSampleServiceResponse, model *GetSampleServiceModel)
	GetDatabase(ch chan *GetSampleServiceResponse, model *GetSampleServiceModel)
	GetCache(ch chan *GetSampleServiceResponse, model *GetSampleServiceModel)
	PublishPubSubMessage(ch chan *PublishPubSubMessageServiceResponse, model *PublishPubSubMessageServiceModel)
	UpdateSample(ch chan *UpdateSampleServiceResponse, model *UpdateSampleServiceModel)
	PostSampleXml(ch chan *PostSampleXmlServiceResponse, model *PostSampleXmlServiceModel)
}

type SampleService struct {
	environment     env.IEnvironment
	loggr           logger.ILogger
	validatr        validator.IValidator
	cachr           cacher.ICacher
	sampleProxy     sample_proxy.ISampleProxy
	sampleDb        sample.ISampleDb
	samplePublisher sample_publisher.ISamplePublisher
	sampleXmlProxy  sample_proxy.ISampleXmlProxy
}

// NewSampleService
// Returns a new SampleService.
func NewSampleService(
	environment env.IEnvironment,
	loggr logger.ILogger,
	validatr validator.IValidator,
	cachr cacher.ICacher,
	sampleProxy sample_proxy.ISampleProxy, sampleDb sample.ISampleDb, samplePublisher sample_publisher.ISamplePublisher, sampleXmlProxy sample_proxy.ISampleXmlProxy) ISampleService {
	service := SampleService{
		environment: environment,
		loggr:       loggr,
		validatr:    validatr,
		cachr:       cachr,
	}

	if sampleProxy != nil {
		service.sampleProxy = sampleProxy
	} else {
		service.sampleProxy = sample_proxy.NewSampleProxy(environment, loggr, validatr, cachr)
	}

	if sampleXmlProxy != nil {
		service.sampleXmlProxy = sampleXmlProxy
	} else {
		service.sampleXmlProxy = sample_proxy.NewSampleXmlProxy(environment, loggr, validatr, cachr)
	}

	if sampleDb != nil {
		service.sampleDb = sampleDb
	} else {
		service.sampleDb = sample.NewSampleDb(environment, loggr, validatr, cachr)
	}

	if samplePublisher != nil {
		service.samplePublisher = samplePublisher
	} else {
		service.samplePublisher = sample_publisher.NewSamplePublisher(environment, loggr, validatr, cachr)
	}

	return &service
}

// GetGoogle
// Gets a sample service response from Google via proxy.
func (s *SampleService) GetGoogle(ch chan *GetSampleServiceResponse, model *GetSampleServiceModel) {
	modelErr := s.validatr.ValidateStruct(model)
	if modelErr != nil {
		ch <- &GetSampleServiceResponse{Error: modelErr}
		return
	}

	getSampleCh := make(chan *sample_proxy.GetSampleProxyResponse)
	defer close(getSampleCh)
	go s.sampleProxy.GetGoogle(getSampleCh, &sample_proxy.GetSampleProxyModel{
		Id:         model.Id,
		SampleName: model.SampleName,
	})
	sampleProxyResponse := <-getSampleCh
	if sampleProxyResponse.Error != nil {
		ch <- &GetSampleServiceResponse{Error: sampleProxyResponse.Error}
		return
	}

	ch <- &GetSampleServiceResponse{
		Id:         sampleProxyResponse.Id,
		SampleName: sampleProxyResponse.SampleName,
	}
}

// GetDatabase
// Gets a sample service response from postgresql database.
func (s *SampleService) GetDatabase(ch chan *GetSampleServiceResponse, model *GetSampleServiceModel) {
	modelErr := s.validatr.ValidateStruct(model)
	if modelErr != nil {
		ch <- &GetSampleServiceResponse{Error: modelErr}
		return
	}

	getSampleCh := make(chan *sample.GetSampleDbResponse)
	defer close(getSampleCh)
	go s.sampleDb.GetSample(getSampleCh, &sample.GetSampleDbModel{
		Id:         model.Id,
		SampleName: model.SampleName,
	})

	sampleDbResponse := <-getSampleCh
	if sampleDbResponse.Error != nil {
		ch <- &GetSampleServiceResponse{Error: sampleDbResponse.Error}
		return
	}

	ch <- &GetSampleServiceResponse{
		Id:         sampleDbResponse.Id,
		SampleName: sampleDbResponse.SampleName,
	}
}

// GetCache
// Gets a sample service response from redis cache.
func (s *SampleService) GetCache(ch chan *GetSampleServiceResponse, model *GetSampleServiceModel) {
	modelErr := s.validatr.ValidateStruct(model)
	if modelErr != nil {
		ch <- &GetSampleServiceResponse{Error: modelErr}
		return
	}

	var response GetSampleServiceResponse
	cacheValue := s.cachr.Get("dummy_cache_key")
	if cacheValue != nil {
		err := json.Unmarshal([]byte(*cacheValue), &response)
		if err != nil {
			s.loggr.Panic("Panicked while unmarshalling json to type.")
		}

		ch <- &response
		return
	}

	response = GetSampleServiceResponse{
		Id:         1,
		SampleName: "Cached new response here!",
	}

	s.cachr.Set("dummy_cache_key", response, 30*time.Second)
	ch <- &response
}

// PublishPubSubMessage
// Publishes a pub sub sample message to topic.
func (s *SampleService) PublishPubSubMessage(ch chan *PublishPubSubMessageServiceResponse, model *PublishPubSubMessageServiceModel) {
	modelErr := s.validatr.ValidateStruct(model)
	if modelErr != nil {
		ch <- &PublishPubSubMessageServiceResponse{Error: modelErr}
		return
	}

	publisherCh := make(chan publisher.PublisherResponse, model.Count)
	defer close(publisherCh)
	for i := 0; i < model.Count; i++ {
		go s.samplePublisher.Publish(publisherCh, &model.Message, map[string]string{publisher.DummyAttribute: "overriding dummy attribute value here"})
	}

	// Await all responses from routines simultaneously.
	var messageIds []string
	for i := 0; i < model.Count; i++ {
		response := <-publisherCh
		if response.Error != nil {
			ch <- &PublishPubSubMessageServiceResponse{Error: response.Error}
			return
		}
		messageIds = append(messageIds, *response.MessageId)
	}

	ch <- &PublishPubSubMessageServiceResponse{
		MessageIds:   messageIds,
		IsSuccessful: true,
	}
}

// UpdateSample
// Updates a sample with details provided by google pub sub message or via http endpoint.
func (s *SampleService) UpdateSample(ch chan *UpdateSampleServiceResponse, model *UpdateSampleServiceModel) {
	modelErr := s.validatr.ValidateStruct(model)
	if modelErr != nil {
		ch <- &UpdateSampleServiceResponse{Error: modelErr}
		return
	}

	// Doing some dummy stuff.
	// Launching rocket to the sky..
	// Calculating...

	response := UpdateSampleServiceResponse{IsSuccessful: true}
	ch <- &response
}

// PostSampleXml
// Post sample xml
func (s *SampleService) PostSampleXml(ch chan *PostSampleXmlServiceResponse, model *PostSampleXmlServiceModel) {
	modelErr := s.validatr.ValidateStruct(model)
	if modelErr != nil {
		ch <- &PostSampleXmlServiceResponse{Error: modelErr}
		return
	}

	getSampleCh := make(chan *sample_proxy.PostSampleXmlProxyResponse)
	defer close(getSampleCh)

	go s.sampleXmlProxy.PostSampleXml(getSampleCh, &sample_proxy.PostSampleXmlProxyModel{
		SampleName: model.SampleName,
		SampleType: model.SampleType,
		SampleCode: model.SampleCode,
	})

	sampleProxyResponse := <-getSampleCh
	if sampleProxyResponse.Error != nil {
		ch <- &PostSampleXmlServiceResponse{Error: sampleProxyResponse.Error}
		return
	}

	ch <- &PostSampleXmlServiceResponse{}
}
