package fields

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-shop-v2/pkg/vue/contracts"
	"reflect"
)

type FieldOption func(field interface{})

func resolveBasicField(field interface{}) (*Field, error) {
	if basicField, ok := field.(*Field); ok {
		return basicField, nil
	}
	if reflect.ValueOf(field).Kind() == reflect.Ptr {
		elem := reflect.ValueOf(field).Elem()
		for i := 0; i < elem.NumField(); i++ {
			value := elem.Field(i)
			if value.IsValid() && value.Type() == reflect.ValueOf(&Field{}).Type() {
				return value.Interface().(*Field), nil
			}
		}
	}
	return nil, fmt.Errorf("basic field not found in %+v\n", field)
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

func SetSortable(sort bool) FieldOption {
	return func(field interface{}) {
		basicField, err := resolveBasicField(field)
		if err == nil {
			basicField.Sortable = sort
		}
	}
}

func SetResolveForDisplay(cb func(ctx *gin.Context, model interface{}) interface{}) FieldOption {
	return func(field interface{}) {
		basicField, err := resolveBasicField(field)
		if err == nil {
			basicField.resolveForDisplay = cb
		}
	}
}

func SetRules(rules []*FieldRule) FieldOption {
	return func(field interface{}) {
		basicField, err := resolveBasicField(field)
		if err == nil {
			var r []contracts.FieldRule
			for _, rule := range rules {
				r = append(r, rule)
			}
			basicField.Rules = r
		}
	}
}

func SetAttribute(attr string) FieldOption {
	return func(field interface{}) {
		basicField, err := resolveBasicField(field)
		if err == nil {
			basicField.Attribute = attr
		}
	}
}

func SetNullable(nullable bool) FieldOption {
	return func(field interface{}) {
		basicField, err := resolveBasicField(field)
		if err == nil {
			basicField.Nullable = nullable
		}
	}
}

func SetValue(value interface{}) FieldOption {
	return func(field interface{}) {
		basicField, err := resolveBasicField(field)
		if err == nil {
			basicField.Value = value
		}
	}
}

func SetNullValue(value interface{}) FieldOption {
	return func(field interface{}) {
		basicField, err := resolveBasicField(field)
		if err == nil {
			basicField.NullValue = value
		}
	}
}

func SetTextAlign(align string) FieldOption {
	return func(field interface{}) {
		basicField, err := resolveBasicField(field)
		if err == nil {
			basicField.TextAlign = align
		}
	}
}

func SetAsHtml(flag bool) FieldOption {
	return func(field interface{}) {
		basicField, err := resolveBasicField(field)
		if err == nil {
			basicField.AsHtml = flag
		}
	}
}

func SetShowOnIndex(show bool) FieldOption {
	return func(field interface{}) {
		basicField, err := resolveBasicField(field)
		if err == nil {
			basicField.showOnIndex = show
		}
	}
}

func SetShowOnDetail(show bool) FieldOption {
	return func(field interface{}) {
		basicField, err := resolveBasicField(field)
		if err == nil {
			basicField.showOnDetail = show
		}
	}
}

func SetShowOnUpdate(show bool) FieldOption {
	return func(field interface{}) {
		basicField, err := resolveBasicField(field)
		if err == nil {
			basicField.showOnUpdate = show
		}
	}
}

func SetShowOnCreation(show bool) FieldOption {
	return func(field interface{}) {
		basicField, err := resolveBasicField(field)
		if err == nil {
			basicField.showOnCreation = show
		}
	}
}

func OnlyOnIndex() FieldOption {
	return func(field interface{}) {
		basicField, err := resolveBasicField(field)
		if err == nil {
			basicField.OnlyOnIndex()
		}
	}
}

func OnlyOnDetail() FieldOption {
	return func(field interface{}) {
		basicField, err := resolveBasicField(field)
		if err == nil {
			basicField.OnlyOnDetail()
		}
	}
}

func OnlyOnForm() FieldOption {
	return func(field interface{}) {
		basicField, err := resolveBasicField(field)
		if err == nil {
			basicField.OnlyOnForm()
		}
	}
}

func ExceptOnForms() FieldOption {
	return func(field interface{}) {
		basicField, err := resolveBasicField(field)
		if err == nil {
			basicField.ExceptOnForms()
		}
	}
}
