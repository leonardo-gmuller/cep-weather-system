package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/leonardo-gmuller/cep-weather-system/service-a/internal/app/gateway/client"
	"github.com/leonardo-gmuller/cep-weather-system/service-a/internal/app/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type CEPInput struct {
	Cep string `json:"cep"`
}

const validateCepPattern = "/validate-cep"

func (h *Handler) ValidateCepSetup(router chi.Router) {
	router.Route(validateCepPattern, func(r chi.Router) {
		r.Post("/", h.ValidateCep)
	})
}

func (h *Handler) ValidateCep(rw http.ResponseWriter, r *http.Request) {
	ctx, span := telemetry.StartServerSpan(r.Context(), "Handle CEP Request")
	defer span.End()
	var input CEPInput
	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, &input)

	span.SetAttributes(
		attribute.String("cep", input.Cep),
	)

	valid, err := h.useCase.ValidateCEP(r.Context(), input.Cep)
	if err != nil {
		http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
		span.RecordError(err)
		span.SetStatus(codes.Error, "internal server error")
		return
	}

	if !valid {
		http.Error(rw, "invalid zipcode", http.StatusBadRequest)
		span.RecordError(err)
		span.SetStatus(codes.Error, "invalid zipcode")
		return
	}

	ctxClient, spanClient := telemetry.StartClientSpan(ctx, "Call Service B API")
	defer spanClient.End()

	client := client.NewServiceBClient(&h.cfg)
	if client == nil {
		spanClient.RecordError(err)
		spanClient.SetStatus(codes.Error, "failed to create service B client")
		http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	resp, err := client.GetWeatherByCep(ctxClient, input.Cep)
	if err != nil {
		http.Error(rw, "Failed to fetch weather data", http.StatusInternalServerError)
		spanClient.RecordError(err)
		spanClient.SetStatus(codes.Error, "failed to fetch weather data")
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(resp)
}
