package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/leonardo-gmuller/cep-weather-system/service-b/internal/app/config"
)

type Weather struct {
	TempC float64
}

type WeatherGateway interface {
	GetWeatherByCity(ctx context.Context, city, uf string) (*Weather, error)
}

func NewWeatherGateway(config *config.Config) *weatherClient {
	return &weatherClient{
		apiKey: config.WeatherAPI.APIKey,
		url:    config.WeatherAPI.URL,
	}
}

type weatherClient struct {
	apiKey string
	url    string
}

type weatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

func (c *weatherClient) GetWeatherByCity(ctx context.Context, city, uf string) (*Weather, error) {
	location := fmt.Sprintf("%s,%s", city, uf)
	escapedLocation := url.QueryEscape(location)

	fullURL := fmt.Sprintf("%s?key=%s&q=%s", c.url, c.apiKey, escapedLocation)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching weather: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making weather request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather api error: %s", resp.Status)
	}

	var data weatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("error decoding weather json: %w", err)
	}

	return &Weather{
		TempC: data.Current.TempC,
	}, nil
}
