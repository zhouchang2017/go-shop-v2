package fields

import (
	"go-shop-v2/pkg/vue/contracts"
	"go-shop-v2/pkg/vue/helper"
)

type Relations struct {
	*Field       `inline`
	ResourceName string `json:"resource_name"`
}

func NewRelationsField(resource contracts.ResourceRelations, fieldName string, opts ...FieldOption) *Relations {
	var fieldOptions = []FieldOption{
		SetPrefixComponent(true),
		SetShowOnIndex(false),
		SetShowOnDetail(true),
		SetShowOnCreation(true),
		SetShowOnUpdate(true),
		WithComponent("relations-field"),
		SetTextAlign("left"),
	}
	fieldOptions = append(fieldOptions, opts...)
	return &Relations{
		Field:        NewField(resource.Title(), fieldName, fieldOptions...),
		ResourceName: helper.ResourceUriKey(resource),
	}
}

func (this *Relations) Searchable() *Relations {
	this.WithMeta("searchable", true)
	return this
}

func (this *Relations) Multiple() *Relations {
	this.WithMeta("multiple", true)
	return this
}

func (this *Relations) MultipleLimit(num int64) *Relations {
	this.WithMeta("multiple_limit", num)
	return this
}

func (this *Relations) WithName(name string) *Relations {
	this.Name = name
	return this
}
