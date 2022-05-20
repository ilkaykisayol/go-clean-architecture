package sample

type GetSampleProxyResponse struct {
	Error      error `json:"-"`
	Id         int
	SampleName string
}
