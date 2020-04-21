package fields

// 货币字段
// https://dm4t2.github.io/vue-currency-input/config/#component
type Currency struct {
	*Field
}

func NewCurrencyField(name string, fieldName string, opts ...FieldOption) *Currency {
	var options = []FieldOption{
		SetPrefixComponent(true),
		SetShowOnIndex(true),
		SetShowOnDetail(true),
		SetShowOnCreation(true),
		SetShowOnUpdate(true),
		WithComponent("currency-field"),
		SetTextAlign("left"),
		SetNullValue(0),
	}
	options = append(options, opts...)

	return &Currency{Field: NewField(name, fieldName, options...)}
}

// A ISO 4217 currency code (for example USD or EUR). Default is CNY
// https://en.wikipedia.org/wiki/ISO_4217
func (this *Currency) Currency(code string) *Currency {
	this.WithMeta("currency", code)
	return this
}

// A BCP 47 language tag (for example en or de-DE).
// Default is zh-CN (use the runtime's default locale).
// https://tools.ietf.org/html/bcp47
func (this *Currency) Locale(locale string) *Currency {
	this.WithMeta("locale", locale)
	return this
}

// The number of displayed decimal digits.
// Default is 2 (use the currency's default).
// 单位换算，默认2位，100 => 1
func (this *Currency) Precision(precision int8) *Currency {
	this.WithMeta("precision", precision)
	return this
}

// 值是否为数字类型
// Whether the number value should be handled as integer instead of float value.
// Default is false.
func (this *Currency) ValueAsInteger(ok bool) *Currency {
	this.WithMeta("valueAsInteger", ok)
	return this
}

// 扩展 css class
func (this *Currency) ExtraClass(class string) *Currency {
	this.WithMeta("extraClass", class)
	return this
}

// 扩展 el-input props
func (this *Currency) ExtraProps(props map[string]interface{}) *Currency {
	this.WithMeta("extraProps", props)
	return this
}
