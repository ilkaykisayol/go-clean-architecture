package receiver

type ReceiverHandlerModel struct {
	MessageId  string            `validate:"required"`
	Data       []byte            `validate:"required"`
	Attributes map[string]string `validate:"required"`
}
