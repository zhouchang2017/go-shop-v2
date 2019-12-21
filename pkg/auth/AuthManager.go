package auth

import (
	"fmt"
	"sync"
)

var Auth *AuthManager
var once sync.Once

const CtxUserKey = "user"

type AuthManager struct {
	guards map[string]func() StatefulGuard // guard工厂方法
}

func (this *AuthManager) Register(factory func() StatefulGuard) {
	this.guards[factory().Name()] = factory
}

func (this *AuthManager) Guard(name string) (StatefulGuard, error) {
	if guardFactory, ok := this.guards[name]; ok {
		return guardFactory(), nil
	}
	return nil, fmt.Errorf("Auth guard [%s] is not defined.", name)
}


func NewAuth() *AuthManager {
	once.Do(func() {
		Auth = &AuthManager{
			guards: map[string]func() StatefulGuard{},
		}
	})
	return Auth
}
