package fields

// 地图字段
type Map struct {
	*Field `inline`
}

var DefaultMapLocation *MapValue

func NewMapField(name string, fieldName string, opts ...FieldOption) *Map {
	var options = []FieldOption{
		SetPrefixComponent(true),
		SetShowOnIndex(false),
		SetShowOnDetail(true),
		SetShowOnCreation(true),
		SetShowOnUpdate(true),
		WithComponent("map-field"),
		SetTextAlign("left"),
	}
	options = append(options, opts...)
	fieldComponent := &Map{Field: NewField(name, fieldName, options...)}
	if DefaultMapLocation != nil {
		fieldComponent.Value = DefaultMapLocation
	} else {
		fieldComponent.Value = &MapValue{}
	}

	return fieldComponent
}

type MapValue struct {
	Lng float64 `json:"lng"`
	Lat float64 `json:"lat"`
}
