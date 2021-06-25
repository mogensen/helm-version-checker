package app

type config struct {
	MetricsPort int `env:"METRICS_PORT" envDefault:"8080"`
	WebPort     int `env:"WEB_PORT" envDefault:"8081"`

	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
}
