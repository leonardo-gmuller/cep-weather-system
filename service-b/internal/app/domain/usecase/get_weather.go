package usecase

import (
	"context"

	"github.com/leonardo-gmuller/cep-weather-system/service-b/internal/app/domain/dto"
	"github.com/leonardo-gmuller/cep-weather-system/service-b/internal/app/telemetry"
	"go.opentelemetry.io/otel/codes"
)

type WeatherResponse struct {
	TempC float64
	TempF float64
	TempK float64
}

func (u *UseCase) GetWeather(ctx context.Context, address dto.Address) (*WeatherResponse, error) {
	ctxWeather, spanWeather := telemetry.StartClientSpan(ctx, "Call Weather API")
	defer spanWeather.End()
	data, err := u.WeatherGateway.GetWeatherByCity(ctxWeather, address.City, address.UF)
	if err != nil {
		spanWeather.RecordError(err)
		spanWeather.SetStatus(codes.Error, "failed to call Weather API")
		return nil, err
	}

	c := data.TempC
	return &WeatherResponse{
		TempC: c,
		TempF: c*1.8 + 32,
		TempK: c + 273,
	}, nil
}
