package charts

type Value struct {
	*Charts
}

func NewValue() *Value {
	value:= &Value{
		NewCharts(),
	}
	value.SetGrid()
	return value
}

func (this *Value) Component() string {
	return "charts-value"
}


