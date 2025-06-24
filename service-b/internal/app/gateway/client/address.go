package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const urlViaCep = "https://viacep.com.br/ws/%s/json"

type ViaCepResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
	Erro        bool   `json:"erro,omitempty"`
}

type AddressGateway interface {
	GetAddressByCEP(ctx context.Context, cep string) (*ViaCepResponse, error)
}

func NewAddressGateway() AddressGateway {
	return &addressGateway{}
}

type addressGateway struct{}

func (a *addressGateway) GetAddressByCEP(ctx context.Context, cep string) (*ViaCepResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(urlViaCep, cep), nil)
	if err != nil {
		return &ViaCepResponse{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return &ViaCepResponse{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return &ViaCepResponse{}, err
	}

	var data ViaCepResponse

	err = json.Unmarshal(body, &data)

	if err != nil {
		return &ViaCepResponse{}, err
	}

	if data.Erro {
		return &ViaCepResponse{}, fmt.Errorf("zipcode not found")
	}

	return &data, nil
}
