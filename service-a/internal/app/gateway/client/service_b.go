package client

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/leonardo-gmuller/cep-weather-system/service-a/internal/app/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type ServiceBClientInterface interface {
	GetWeatherByCep(cep string) ([]byte, error)
}

type ServiceBClient struct {
	url string
}

func NewServiceBClient(config *config.Config) *ServiceBClient {
	return &ServiceBClient{
		url: config.ServiceB.URL,
	}
}

func (c *ServiceBClient) GetWeatherByCep(ctx context.Context, cep string) ([]byte, error) {
	fmt.Printf(c.url+"/weather/%s\n", cep)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url+"/weather/"+cep, nil)
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
