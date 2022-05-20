package health

type HealthCheckResponse struct {
	Error error `json:"-"`
}
