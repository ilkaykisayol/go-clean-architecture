package sample

type SamplePublisherModel struct {
	SampleId     int    `json:"SampleId" validate:"required,gte=0"`
	SampleStatus int    `json:"SampleStatus" validate:"required,gte=0"`
	ModifiedBy   string `json:"ModifiedBy" validate:"required"`
}
