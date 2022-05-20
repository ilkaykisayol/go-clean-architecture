package sample

type GetSampleDbModel struct {
	Id         int    `validate:"required"`
	SampleName string `validate:"required"`
}
