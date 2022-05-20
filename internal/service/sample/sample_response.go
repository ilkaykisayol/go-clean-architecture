package sample

type GetSampleServiceResponse struct {
	Error      error `json:"-"`
	Id         int
	SampleName string
}

type UpdateSampleServiceResponse struct {
	Error        error `json:"-"`
	IsSuccessful bool
}

type PublishPubSubMessageServiceResponse struct {
	Error        error `json:"-"`
	IsSuccessful bool
	MessageIds   []string
}

type PostSampleXmlServiceResponse struct {
	Error     error `json:"-"`
	IsSuccess bool
	Message   string
}
