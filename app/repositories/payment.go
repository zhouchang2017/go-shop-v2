package repositories

import "go-shop-v2/pkg/repository"

type PaymentRep struct {
	repository.IRepository
}

func NewPaymentRep(IRepository repository.IRepository) *PaymentRep {
	return &PaymentRep{IRepository: IRepository}
}
