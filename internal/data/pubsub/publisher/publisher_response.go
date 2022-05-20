package publisher

type PublisherResponse struct {
	Error     error `json:"-"`
	MessageId *string
}
