package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/leonardo-gmuller/go-weather/app/domain/usecase"
	"github.com/leonardo-gmuller/go-weather/app/gateway/client"
)

type MockAddressGateway struct {
	GetAddressByCEPFunc func(string) (*client.ViaCepResponse, error)
}

func (m *MockAddressGateway) GetAddressByCEP(cep string) (*client.ViaCepResponse, error) {
	return m.GetAddressByCEPFunc(cep)
}

func TestGetAddress_InvalidZipcode(t *testing.T) {
	u := &usecase.UseCase{}
	_, err := u.GetAddress(context.Background(), "123")
	if !errors.Is(err, usecase.ErrInvalidZipcode) {
		t.Errorf("expected ErrInvalidZipcode, got %v", err)
	}
}

func TestGetAddress_NotFound(t *testing.T) {
	mockGateway := &MockAddressGateway{
		GetAddressByCEPFunc: func(cep string) (*client.ViaCepResponse, error) {
			return nil, errors.New("not found")
		},
	}
	u := &usecase.UseCase{AddressGateway: mockGateway}
	_, err := u.GetAddress(context.Background(), "12345678")
	if !errors.Is(err, usecase.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestGetAddress_Success(t *testing.T) {
	mockGateway := &MockAddressGateway{
		GetAddressByCEPFunc: func(cep string) (*client.ViaCepResponse, error) {
			return &client.ViaCepResponse{
				Localidade: "Sao Paulo",
				Uf:         "SP",
			}, nil
		},
	}
	u := &usecase.UseCase{AddressGateway: mockGateway}
	resp, err := u.GetAddress(context.Background(), "12345678")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Address.City != "Sao Paulo" || resp.Address.UF != "SP" {
		t.Errorf("unexpected address: %+v", resp.Address)
	}
}
