package auth

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"go-shop-v2/pkg/utils"
	"reflect"
)

type GormUserProvider struct {
	model Authenticatable
	db    *gorm.DB
}

func NewGormUserProvider(model Authenticatable) *GormUserProvider {
	return &GormUserProvider{model: model}
}

func (this *GormUserProvider) RetrieveById(identifier interface{}) (Authenticatable, error) {
	model := this.createModel()
	first := this.newModelQuery().Where(fmt.Sprintf("%s = ?", this.model.GetAuthIdentifierName()), identifier).First(model)
	return model, first.Error
}

func (this *GormUserProvider) RetrieveByCredentials(credentials map[string]string) (Authenticatable, error) {
	if _, ok := credentials["password"]; !ok && credentials == nil {
		return nil, errors.New("credentials is empty!")
	}
	model := this.createModel()
	query := this.newModelQuery()

	for key, value := range credentials {
		if key == "password" {
			continue
		}
		query.Where(fmt.Sprintf("%s = ?", key), value)
	}

	first := query.First(model)
	return model, first.Error
}

func (this *GormUserProvider) ValidateCredentials(user Authenticatable, credentials map[string]string) bool {
	if password, ok := credentials["password"]; ok {
		if utils.Compare(user.GetAuthPassword(), password) != nil {
			return false
		}
		return true
	}
	return false
}

func (this *GormUserProvider) newModelQuery() *gorm.DB {
	return this.db.Model(this.model)
}

func (this *GormUserProvider) createModel() Authenticatable {
	t := reflect.ValueOf(this.model).Elem().Type()
	return reflect.New(t).Interface().(Authenticatable)
}
