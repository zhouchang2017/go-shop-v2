package vue

import (
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/auth"
)

type Element interface {
	AuthorizedTo(ctx *gin.Context, user auth.Authenticatable) bool
}

type BasicElement struct {
	Meta            MetaItems `json:"meta"`
	PrefixComponent bool      `json:"prefix_component"`
	Component       string    `json:"component"`
	authorizedTo    func(ctx *gin.Context, user auth.Authenticatable) bool
}

func AuthorizedCallback(cb func(ctx *gin.Context, user auth.Authenticatable) bool) FieldOption {
	return func(field interface{}) {
		basicField, err := resolveBasicField(field)
		if err == nil {
			basicField.AuthorizedCallback(cb)
		}
	}
}

func SetPrefixComponent(ok bool) FieldOption {
	return func(field interface{}) {
		basicField, err := resolveBasicField(field)
		if err == nil {
			basicField.SetPrefixComponent(ok)
		}
	}
}

func WithMeta(key string, value interface{}) FieldOption {
	return func(field interface{}) {
		basicField, err := resolveBasicField(field)
		if err == nil {
			basicField.WithMeta(key, value)
		}
	}
}

func SetComponent(component string) FieldOption {
	return func(field interface{}) {
		basicField, err := resolveBasicField(field)
		if err == nil {
			basicField.SetComponent(component)
		}
	}
}

func (m *BasicElement) AuthorizedCallback(cb func(ctx *gin.Context, user auth.Authenticatable) bool) {
	m.authorizedTo = cb
}

func (m *BasicElement) AuthorizedTo(ctx *gin.Context, user auth.Authenticatable) bool {
	if m.authorizedTo != nil {
		return m.authorizedTo(ctx, user)
	}
	return true
}

func (m *BasicElement) SetPrefixComponent(ok bool) {
	m.PrefixComponent = ok
}

func (m *BasicElement) WithMeta(key string, value interface{}) {
	m.Meta = append(m.Meta, &metaItem{key, value})
}

func (m *BasicElement) SetComponent(component string) {
	m.Component = component
}
