package env

// AppEnviroment Application
const (
	AppEnvironment = "APP_ENVIRONMENT"
	AppName        = "APP_NAME"
	AppHost        = "APP_HOST"
)

// Jwt
const JwtSecret = "JWT_SECRET"

// Database
const PostgresqlConnectionString = "POSTGRESQL_CONNECTION_STRING"

// Redis
const (
	RedisAddress  = "REDIS_ADDRESS"
	RedisPassword = "REDIS_PASSWORD"
)

// PubSub
const (
	SamplePublisherSaJson        = "SAMPLE_PUBLISHER_SA_JSON"
	SampleReceiverSaJson         = "SAMPLE_RECEIVER_SA_JSON"
	SamplePublisherProjectId     = "SAMPLE_PUBLISHER_PROJECT_ID"
	SamplePublisherTopicId       = "SAMPLE_PUBLISHER_TOPIC_ID"
	SampleReceiverProjectId      = "SAMPLE_RECEIVER_PROJECT_ID"
	SampleReceiverSubscriptionId = "SAMPLE_RECEIVER_SUBSCRIPTION_ID"
)

// Proxy
const (
	SampleProxyUrl     = "SAMPLE_PROXY_URL"
	SampleProxyTimeout = "SAMPLE_PROXY_TIMEOUT"
)
