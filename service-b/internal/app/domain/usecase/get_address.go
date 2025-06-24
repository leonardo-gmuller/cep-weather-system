package usecase

import (
	"context"
	"errors"

	"github.com/leonardo-gmuller/cep-weather-system/service-b/internal/app/domain/dto"
	"github.com/leonardo-gmuller/cep-weather-system/service-b/internal/app/telemetry"
	"go.opentelemetry.io/otel/codes"
)

var (
	ErrInvalidZipcode = errors.New("invalid zipcode")
	ErrNotFound       = errors.New("can not find zipcode")
)

type AddressResponse struct {
	Address dto.Address
}

func (u *UseCase) GetAddress(ctx context.Context, zipcode string) (*AddressResponse, error) {
	if len(zipcode) != 8 {
		return nil, ErrInvalidZipcode
	}

	ctxViaCEP, spanViaCEP := telemetry.StartClientSpan(ctx, "Call ViaCEP API")
	defer spanViaCEP.End()

	address, err := u.AddressGateway.GetAddressByCEP(ctxViaCEP, zipcode)
	if err != nil {
		spanViaCEP.RecordError(err)
		spanViaCEP.SetStatus(codes.Error, "failed to call ViaCEP API")
		return nil, ErrNotFound
	}

	return &AddressResponse{Address: dto.Address{
		City: address.Localidade,
		UF:   address.Uf,
	}}, nil
}
