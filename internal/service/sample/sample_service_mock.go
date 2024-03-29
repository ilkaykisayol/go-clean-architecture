// Code generated by MockGen. DO NOT EDIT.
// Source: internal\service\sample\sample_service.go

// Package sample is a generated GoMock package.
package sample

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockISampleService is a mock of ISampleService interface.
type MockISampleService struct {
	ctrl     *gomock.Controller
	recorder *MockISampleServiceMockRecorder
}

// MockISampleServiceMockRecorder is the mock recorder for MockISampleService.
type MockISampleServiceMockRecorder struct {
	mock *MockISampleService
}

// NewMockISampleService creates a new mock instance.
func NewMockISampleService(ctrl *gomock.Controller) *MockISampleService {
	mock := &MockISampleService{ctrl: ctrl}
	mock.recorder = &MockISampleServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockISampleService) EXPECT() *MockISampleServiceMockRecorder {
	return m.recorder
}

// GetCache mocks base method.
func (m *MockISampleService) GetCache(ch chan *GetSampleServiceResponse, model *GetSampleServiceModel) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetCache", ch, model)
}

// GetCache indicates an expected call of GetCache.
func (mr *MockISampleServiceMockRecorder) GetCache(ch, model interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCache", reflect.TypeOf((*MockISampleService)(nil).GetCache), ch, model)
}

// GetDatabase mocks base method.
func (m *MockISampleService) GetDatabase(ch chan *GetSampleServiceResponse, model *GetSampleServiceModel) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetDatabase", ch, model)
}

// GetDatabase indicates an expected call of GetDatabase.
func (mr *MockISampleServiceMockRecorder) GetDatabase(ch, model interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDatabase", reflect.TypeOf((*MockISampleService)(nil).GetDatabase), ch, model)
}

// GetGoogle mocks base method.
func (m *MockISampleService) GetGoogle(ch chan *GetSampleServiceResponse, model *GetSampleServiceModel) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetGoogle", ch, model)
}

// GetGoogle indicates an expected call of GetGoogle.
func (mr *MockISampleServiceMockRecorder) GetGoogle(ch, model interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGoogle", reflect.TypeOf((*MockISampleService)(nil).GetGoogle), ch, model)
}

// PostSampleXml mocks base method.
func (m *MockISampleService) PostSampleXml(ch chan *PostSampleXmlServiceResponse, model *PostSampleXmlServiceModel) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "PostSampleXml", ch, model)
}

// PostSampleXml indicates an expected call of PostSampleXml.
func (mr *MockISampleServiceMockRecorder) PostSampleXml(ch, model interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostSampleXml", reflect.TypeOf((*MockISampleService)(nil).PostSampleXml), ch, model)
}

// PublishPubSubMessage mocks base method.
func (m *MockISampleService) PublishPubSubMessage(ch chan *PublishPubSubMessageServiceResponse, model *PublishPubSubMessageServiceModel) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "PublishPubSubMessage", ch, model)
}

// PublishPubSubMessage indicates an expected call of PublishPubSubMessage.
func (mr *MockISampleServiceMockRecorder) PublishPubSubMessage(ch, model interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublishPubSubMessage", reflect.TypeOf((*MockISampleService)(nil).PublishPubSubMessage), ch, model)
}

// UpdateSample mocks base method.
func (m *MockISampleService) UpdateSample(ch chan *UpdateSampleServiceResponse, model *UpdateSampleServiceModel) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateSample", ch, model)
}

// UpdateSample indicates an expected call of UpdateSample.
func (mr *MockISampleServiceMockRecorder) UpdateSample(ch, model interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSample", reflect.TypeOf((*MockISampleService)(nil).UpdateSample), ch, model)
}
