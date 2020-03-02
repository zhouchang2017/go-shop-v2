package repositories

import "go-shop-v2/pkg/repository"

type AddressRep struct {
	repository.IRepository
}

func NewAddressRep(IRepository repository.IRepository) *AddressRep {
	return &AddressRep{IRepository: IRepository}
}
