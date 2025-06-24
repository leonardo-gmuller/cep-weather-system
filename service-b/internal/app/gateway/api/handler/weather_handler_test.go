package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/leonardo-gmuller/cep-weather-system/service-b/internal/app/config"
	"github.com/leonardo-gmuller/cep-weather-system/service-b/internal/app/domain/dto"
	"github.com/leonardo-gmuller/cep-weather-system/service-b/internal/app/domain/usecase"
	"github.com/leonardo-gmuller/cep-weather-system/service-b/internal/app/gateway/api/handler"

	"github.com/leonardo-gmuller/cep-weather-system/service-b/internal/app/tests/mocks"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// --- Handler Setup ---

func TestGetWeather_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockUseCase(ctrl)
	mockUseCase.EXPECT().GetAddress(gomock.Any(), "12345678").Return(&usecase.AddressResponse{Address: dto.Address{}}, nil).Times(1)
	mockUseCase.EXPECT().GetWeather(gomock.Any(), usecase.AddressResponse{}.Address).Return(&usecase.WeatherResponse{}, nil).Times(1)

	req := httptest.NewRequest(http.MethodGet, "/weather/12345678", nil)
	respWriter := httptest.NewRecorder()

	router := chi.NewRouter()
	cfg := config.Config{}
	handler.RegisterRoutes(router, cfg, mockUseCase)

	router.ServeHTTP(respWriter, req)

	if respWriter.Code != http.StatusOK {
		t.Fatalf("expected status code 200, got %d", respWriter.Code)
	}
	var response map[string]float64
	if err := json.NewDecoder(respWriter.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	expectedWeather := usecase.WeatherResponse{}
	if response["temp_C"] != expectedWeather.TempC {
		t.Errorf("expected temp_C %f, got %f", expectedWeather.TempC, response["temp_C"])
	}
	if response["temp_F"] != expectedWeather.TempF {
		t.Errorf("expected temp_F %f, got %f", expectedWeather.TempF, response["temp_F"])
	}
	if response["temp_K"] != expectedWeather.TempK {
		t.Errorf("expected temp_K %f, got %f", expectedWeather.TempK, response["temp_K"])
	}
}
func TestGetWeather_InvalidZipcode(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockUseCase(ctrl)
	mockUseCase.EXPECT().
		GetAddress(gomock.Any(), "invalid").
		Return(nil, usecase.ErrInvalidZipcode).
		Times(1)

	req := httptest.NewRequest(http.MethodGet, "/weather/invalid", nil)
	respWriter := httptest.NewRecorder()

	router := chi.NewRouter()
	cfg := config.Config{}
	handler.RegisterRoutes(router, cfg, mockUseCase)

	router.ServeHTTP(respWriter, req)

	if respWriter.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status code 422, got %d", respWriter.Code)
	}
}

func TestGetWeather_NotFound(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockUseCase(ctrl)
	mockUseCase.EXPECT().
		GetAddress(gomock.Any(), "00000000").
		Return(nil, usecase.ErrNotFound).
		Times(1)

	req := httptest.NewRequest(http.MethodGet, "/weather/00000000", nil)
	respWriter := httptest.NewRecorder()

	router := chi.NewRouter()
	cfg := config.Config{}
	handler.RegisterRoutes(router, cfg, mockUseCase)

	router.ServeHTTP(respWriter, req)

	if respWriter.Code != http.StatusNotFound {
		t.Fatalf("expected status code 404, got %d", respWriter.Code)
	}
}

func TestGetWeather_InternalServerErrorOnGetAddress(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockUseCase(ctrl)
	mockUseCase.EXPECT().
		GetAddress(gomock.Any(), "99999999").
		Return(nil, assert.AnError).
		Times(1)

	req := httptest.NewRequest(http.MethodGet, "/weather/99999999", nil)
	respWriter := httptest.NewRecorder()

	router := chi.NewRouter()
	cfg := config.Config{}
	handler.RegisterRoutes(router, cfg, mockUseCase)

	router.ServeHTTP(respWriter, req)

	if respWriter.Code != http.StatusInternalServerError {
		t.Fatalf("expected status code 500, got %d", respWriter.Code)
	}
}

func TestGetWeather_InternalServerErrorOnGetWeather(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockUseCase(ctrl)
	mockUseCase.EXPECT().
		GetAddress(gomock.Any(), "12341234").
		Return(&usecase.AddressResponse{Address: dto.Address{}}, nil).
		Times(1)
	mockUseCase.EXPECT().
		GetWeather(gomock.Any(), gomock.Any()).
		Return(nil, assert.AnError).
		Times(1)

	req := httptest.NewRequest(http.MethodGet, "/weather/12341234", nil)
	respWriter := httptest.NewRecorder()

	router := chi.NewRouter()
	cfg := config.Config{}
	handler.RegisterRoutes(router, cfg, mockUseCase)

	router.ServeHTTP(respWriter, req)

	if respWriter.Code != http.StatusInternalServerError {
		t.Fatalf("expected status code 500, got %d", respWriter.Code)
	}
}
