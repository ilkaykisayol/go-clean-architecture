package sample

import (
	sampleDb "go-clean-architecture/internal/data/database/sample"
	sampleProxy "go-clean-architecture/internal/data/proxy/sample"
	samplePubSubPublisher "go-clean-architecture/internal/data/pubsub/publisher/sample"
	"go-clean-architecture/internal/util/cacher"
	"go-clean-architecture/internal/util/env"
	"go-clean-architecture/internal/util/logger"
	"go-clean-architecture/internal/util/validator"

	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type SampleServiceTestSuite struct {
	suite.Suite
	sampleService                  ISampleService
	mockEnvironment                *env.MockIEnvironment
	mockLogger                     *logger.MockILogger
	mockValidator                  *validator.MockIValidator
	mockCacher                     *cacher.MockICacher
	mockSampleProxy                *sampleProxy.MockISampleProxy
	mockSampleDb                   *sampleDb.MockISampleDb
	mockSamplePubSubPublisher      *samplePubSubPublisher.MockISamplePublisher
	mockSamplePubSubProtoPublisher *samplePubSubPublisher.MockISampleProtoPublisher
	mockSampleXmlProxy             *sampleProxy.MockISampleXmlProxy
}

// Run suite.
func TestSampleService(t *testing.T) {
	suite.Run(t, new(SampleServiceTestSuite))
}

// Runs before each test in the suite.
func (s *SampleServiceTestSuite) SetupTest() {
	s.T().Log("Setup")

	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.mockEnvironment = env.NewMockIEnvironment(ctrl)
	s.mockLogger = logger.NewMockILogger(ctrl)
	s.mockValidator = validator.NewMockIValidator(ctrl)
	s.mockCacher = cacher.NewMockICacher(ctrl)
	s.mockSampleProxy = sampleProxy.NewMockISampleProxy(ctrl)
	s.mockSampleDb = sampleDb.NewMockISampleDb(ctrl)
	s.mockSamplePubSubPublisher = samplePubSubPublisher.NewMockISamplePublisher(ctrl)
	s.mockSamplePubSubProtoPublisher = samplePubSubPublisher.NewMockISampleProtoPublisher(ctrl)
	s.mockSampleXmlProxy = sampleProxy.NewMockISampleXmlProxy(ctrl)

	s.sampleService = NewSampleService(s.mockEnvironment, s.mockLogger, s.mockValidator, s.mockCacher, s.mockSampleProxy, s.mockSampleDb, s.mockSamplePubSubPublisher, s.mockSamplePubSubProtoPublisher, s.mockSampleXmlProxy)
}

// Runs after each test in the suite.
func (s *SampleServiceTestSuite) TearDownTest() {
	s.T().Log("Teardown")
}

func (s *SampleServiceTestSuite) TestGetGoogle_ModelHasValidationError_ReturnsError() {
	// Given
	id := 99
	sampleName := "test_get_google_service"

	s.mockValidator.
		EXPECT().
		ValidateStruct(gomock.Eq(&GetSampleServiceModel{Id: id, SampleName: sampleName})).
		Return(errors.New("validation error"))

	// When
	ch := make(chan *GetSampleServiceResponse)
	defer close(ch)
	go s.sampleService.GetGoogle(ch, &GetSampleServiceModel{
		Id:         id,
		SampleName: sampleName,
	})
	response := <-ch

	// Then
	s.Error(response.Error)
}

func (s *SampleServiceTestSuite) TestGetGoogle_ProxyReturnedError_ReturnsError() {
	// Given
	id := 99
	sampleName := "test_get_google_service"

	s.mockValidator.
		EXPECT().
		ValidateStruct(gomock.Eq(&GetSampleServiceModel{Id: id, SampleName: sampleName})).
		Return(nil)

	s.mockSampleProxy.
		EXPECT().
		GetGoogle(gomock.Any(), gomock.Eq(&sampleProxy.GetSampleProxyModel{Id: id, SampleName: sampleName})).
		DoAndReturn(func(ch chan *sampleProxy.GetSampleProxyResponse, model *sampleProxy.GetSampleProxyModel) {
			ch <- &sampleProxy.GetSampleProxyResponse{Error: errors.New("proxy has returned error")}
			return
		})

	// When
	ch := make(chan *GetSampleServiceResponse)
	defer close(ch)
	go s.sampleService.GetGoogle(ch, &GetSampleServiceModel{
		Id:         id,
		SampleName: sampleName,
	})
	response := <-ch

	// Then
	s.Error(response.Error)
}

func (s *SampleServiceTestSuite) TestGetGoogle_HappyPath_Success() {
	// Given
	id := 99
	sampleName := "test_get_google_service"

	s.mockValidator.
		EXPECT().
		ValidateStruct(gomock.Eq(&GetSampleServiceModel{Id: id, SampleName: sampleName})).
		Return(nil)
	s.mockSampleProxy.
		EXPECT().
		GetGoogle(gomock.Any(), gomock.Eq(&sampleProxy.GetSampleProxyModel{Id: id, SampleName: sampleName})).
		DoAndReturn(func(ch chan *sampleProxy.GetSampleProxyResponse, model *sampleProxy.GetSampleProxyModel) {
			ch <- &sampleProxy.GetSampleProxyResponse{Id: id, SampleName: sampleName}
			return
		})

	// When
	ch := make(chan *GetSampleServiceResponse)
	defer close(ch)
	go s.sampleService.GetGoogle(ch, &GetSampleServiceModel{
		Id:         id,
		SampleName: sampleName,
	})
	response := <-ch

	// Then
	s.Equal(id, response.Id)
	s.Equal(sampleName, response.SampleName)
}
