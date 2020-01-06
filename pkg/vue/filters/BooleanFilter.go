package filters

type BooleanFilter struct {
	*AbstractSelectFilter
}

func NewBooleanFilter() *BooleanFilter {
	return &BooleanFilter{
		AbstractSelectFilter: NewAbstractSelectFilter(),
	}
}

func (this BooleanFilter) Component() string  {
	return "boolean-filter"
}