package sample

type GetSampleProxyModel struct {
	Id         int    `validate:"required"`
	SampleName string `validate:"required"`
}
