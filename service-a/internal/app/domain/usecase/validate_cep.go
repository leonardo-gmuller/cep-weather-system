package usecase

import "context"

func (u *UseCase) ValidateCEP(ctx context.Context, cep string) (bool, error) {
	if len(cep) != 8 {
		return false, nil
	}

	for _, char := range cep {
		if char < '0' || char > '9' {
			return false, nil
		}
	}

	return true, nil
}
