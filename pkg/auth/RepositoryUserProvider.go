package auth

import (
	"context"
	"errors"
	"go-shop-v2/pkg/repository"
	"go-shop-v2/pkg/utils"
)

type RepositoryUserProvider struct {
	rep repository.IRepository
}

func NewRepositoryUserProvider(rep repository.IRepository) *RepositoryUserProvider {
	return &RepositoryUserProvider{rep: rep}
}

func (this *RepositoryUserProvider) RetrieveById(identifier interface{}) (Authenticatable, error) {
	result := <-this.rep.FindById(context.Background(), identifier.(string))
	if result.Error != nil {
		return nil, result.Error
	}
	return result.Result.(Authenticatable), nil
}

func (this *RepositoryUserProvider) RetrieveByCredentials(credentials map[string]string) (Authenticatable, error) {
	if _, ok := credentials["password"]; !ok && credentials == nil {
		return nil, errors.New("credentials is empty!")
	}

	query := map[string]interface{}{}
	for key, value := range credentials {
		if key == "password" {
			continue
		}
		query[key] = value
	}
	result := <-this.rep.FindOne(context.Background(), query)
	if result.Error != nil {
		return nil, result.Error
	}
	return result.Result.(Authenticatable), nil
}

func (this *RepositoryUserProvider) ValidateCredentials(user Authenticatable, credentials map[string]string) bool {
	if password, ok := credentials["password"]; ok {
		if utils.Compare(user.GetAuthPassword(), password) != nil {
			return false
		}
		return true
	}
	return false
}
