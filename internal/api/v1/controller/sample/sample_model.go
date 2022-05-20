package sample

import (
	"go-clean-architecture/internal/data/pubsub/publisher/sample"
)

type AddSampleModel struct {
	SampleName *string `json:"SampleName"`
	SampleType *string `json:"SampleType"`
	SampleCode *int    `json:"SampleCode"`
}

type UpdateSampleModel struct {
	SampleStatus int    `json:"SampleStatus" validate:"required,gte=0"`
	ModifiedBy   string `json:"ModifiedBy" validate:"required"`
}

type PublishPubSubMessageModel struct {
	Count   int                         `json:"Count" validate:"required"`
	Message sample.SamplePublisherModel `json:"Message" validate:"required"`
}

type PostSampleXmlModel struct {
	SampleName string `json:"SampleName"`
	SampleType string `json:"SampleType"`
	SampleCode *int   `json:"SampleCode"`
}
