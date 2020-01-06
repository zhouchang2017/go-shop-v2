package fields

import "github.com/gin-gonic/gin"

// 省市区选择器
type AreaCascader struct {
	*Field `inline`
}

func NewAreaCascader(name string, fieldName string, opts ...FieldOption) *AreaCascader {
	var options = []FieldOption{
		SetPrefixComponent(true),
		SetShowOnIndex(true),
		SetShowOnDetail(true),
		SetShowOnCreation(true),
		SetShowOnUpdate(true),
		WithComponent("area-cascader-field"),
		SetTextAlign("left"),
	}
	options = append(options, opts...)
	area := &AreaCascader{Field: NewField(name, fieldName, options...)}
	area.Value = &areaCascaderValue{}
	return area
}

type areaCascaderValue struct {
	Province string `json:"province"`
	City     string `json:"city"`
	Areas    string `json:"areas"`
}

// 赋值
func (this *AreaCascader) Resolve(ctx *gin.Context, model interface{}) {
	if this.resolveForDisplay != nil {
		this.Value = this.resolveForDisplay(ctx, model)
		return
	}
	addressModels := this.resolveAttribute(ctx, model)
	value := &areaCascaderValue{
		Province: getValueByField(addressModels, "Province").(string),
		City:     getValueByField(addressModels, "City").(string),
		Areas:    getValueByField(addressModels, "Areas").(string),
	}
	this.Value = value
}
