package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/leonardo-gmuller/cep-weather-system/service-b/internal/app/gateway/api/handler/schema"
	"github.com/leonardo-gmuller/cep-weather-system/service-b/internal/app/telemetry"
	"github.com/leonardo-gmuller/go-weather/app/domain/usecase"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
)

const weatherPattern = "/weather"

func (h *Handler) WeatherSetup(router chi.Router) {
	router.Route(weatherPattern, func(r chi.Router) {
		r.Get("/{cep}", h.GetWeather)
	})
}

func (h *Handler) GetWeather(w http.ResponseWriter, r *http.Request) {
	ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
	ctx, span := telemetry.StartServerSpan(ctx, "Handle Weather Request")
	defer span.End()

	cep := chi.URLParam(r, "cep")

	span.SetAttributes(attribute.String("cep", cep))

	address, err := h.useCase.GetAddress(ctx, cep)
	if err != nil {
		var status int
		switch err {
		case usecase.ErrInvalidZipcode:
			status = http.StatusUnprocessableEntity
		case usecase.ErrNotFound:
			status = http.StatusNotFound
		default:
			status = http.StatusInternalServerError
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return

	}

	weather, err := h.useCase.GetWeather(ctx, address.Address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := schema.WeatherResponse{
		TempC: weather.TempC,
		TempF: weather.TempF,
		TempK: weather.TempK,
		City:  address.Address.City,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
