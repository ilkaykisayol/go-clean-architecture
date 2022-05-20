package sample

type SampleReceiverHandlerModel struct {
	SampleId     int    `validate:"required,gte=0"`
	SampleStatus int    `validate:"required,gte=0"`
	ModifiedBy   string `validate:"required"`
}
