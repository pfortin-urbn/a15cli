package models

type AppConfig struct {
	ListenPort      string `envconfig:"LISTEN_PORT"`
	ServerDirectory string `envconfig:"SERVER_DIRECTORY"`

	MongoHost string `envconfig:"MONGO_HOST"`
	MongoPort string `envconfig:"MONGO_PORT"`

	ApiServices map[string]string `envconfig:"API_SERVICES"`

	PayloadHmacKey string `envconfig:"PAYLOAD_HMAC_KEY"`
	PayloadKey     string `envconfig:"PAYLOAD_KEY"`
	SigningSecret  string `envconfig:"SIGNING_SECRET"`

	FulfillmentQueueURL       string `envconfig:"FULFILLMENT_QUEUE_URL"`
	FulfillmentErrorsQueueURL string `envconfig:"FULFILLMENT_ERROR_QUEUE_URL"`
}

var Config AppConfig
