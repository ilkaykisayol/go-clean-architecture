package sample

import (
	"go-clean-architecture/internal/data/pubsub/publisher/sample"
)

type GetSampleServiceModel struct {
	Id         int    `validate:"required,gte=0"`
	SampleName string `validate:"required"`
}

type UpdateSampleServiceModel struct {
	SampleId     int    `validate:"required,gte=0"`
	SampleStatus int    `validate:"required,gte=0"`
	ModifiedBy   string `validate:"required"`
}

type PublishPubSubMessageServiceModel struct {
	Count   int                         `validate:"required,gte=0,lte=50"`
	Message sample.SamplePublisherModel `validate:"required"`
}

type PostSampleXmlServiceModel struct {
	SampleName string `validate:"required"`
	SampleType string `validate:"required"`
	SampleCode *int
}
