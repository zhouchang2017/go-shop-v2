package element

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/auth"
)

type Element struct {
	ElementPrefixComponent bool                   `json:"prefix_component"`
	ElementComponent       string                 `json:"component"`
	ElementMeta            map[string]interface{} `json:"meta"`
	authorizedTo           func(ctx *gin.Context, user auth.Authenticatable) bool
}

func NewElement() *Element {
	return &Element{ElementMeta: map[string]interface{}{}}
}

func (m *Element) AuthorizedCallback(cb func(ctx *gin.Context, user auth.Authenticatable) bool) {
	m.authorizedTo = cb
}

func (m *Element) AuthorizedTo(ctx *gin.Context, user auth.Authenticatable) bool {
	if m.authorizedTo != nil {
		return m.authorizedTo(ctx, user)
	}
	return true
}

func (m *Element) SetPrefixComponent(ok bool) {
	m.ElementPrefixComponent = ok
}

func (m *Element) WithMeta(key string, value interface{}) {
	m.ElementMeta[key] = value
}

func (m Element) Component() string {
	return m.ElementComponent
}

func (m Element) PrefixComponent() bool {
	return m.ElementPrefixComponent
}

func (m Element) Meta() map[string]interface{} {
	return m.ElementMeta
}

func (m *Element) WithComponent(component string) {
	m.ElementComponent = component
}
