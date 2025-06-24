package config

import (
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Environment string

type Config struct {
	Environment Environment `required:"true" envconfig:"ENVIRONMENT"`

	App    App
	Server Server

	ServiceB ServiceB
	Otel     Otel
}

type App struct {
	Name                    string        `required:"true" envconfig:"SERVICE_A_NAME"`
	ID                      string        `required:"true" envconfig:"SERVICE_A_ID"`
	GracefulShutdownTimeout time.Duration `required:"true" envconfig:"SERVICE_A_GRACEFUL_SHUTDOWN_TIMEOUT"`
}

type Server struct {
	Address      string        `required:"true" envconfig:"SERVER_A_ADDRESS"`
	ReadTimeout  time.Duration `required:"true" envconfig:"SERVER_A_READ_TIMEOUT"`
	WriteTimeout time.Duration `required:"true" envconfig:"SERVER_A_WRITE_TIMEOUT"`
}

type ServiceB struct {
	URL string `required:"true" envconfig:"SERVICE_B_URL"`
}

type Otel struct {
	CollectorEndpoint string        `required:"true" envconfig:"OTEL_COLLECTOR_ENDPOINT"`
	ExporterTimeout   time.Duration `required:"true" envconfig:"OTEL_EXPORTER_TIMEOUT"`

	// The ratio of samples sent by TraceID. See more on TraceIDRatioBased.
	// NOTE: The sampling in production is always 1% (100:1). So just values lesser than 1% will make an effect.
	SamplingRatio    float64 `required:"true" envconfig:"OTEL_SAMPLING_RATIO"` // 0.01 is a 100:1 ratio
	ServiceName      string  `required:"true" envconfig:"OTEL_SERVICE_A_NAME"`
	ServiceNamespace string  `required:"true" envconfig:"OTEL_SERVICE_NAMESPACE"`
}

func New() (Config, error) {
	const operation = "Config.New"

	var cfg Config

	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using system environment variables")
	}

	err := envconfig.Process("", &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("%s -> %w", operation, err)
	}

	return cfg, nil
}
