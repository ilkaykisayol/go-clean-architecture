package sample

type GetSampleDbResponse struct {
	Error      error `json:"-"`
	Id         int
	SampleName string
}
