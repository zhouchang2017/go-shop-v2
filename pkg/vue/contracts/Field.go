package contracts

import "github.com/gin-gonic/gin"



// 字段接口
type Field interface {
	Element
	ShowOnIndex() bool
	ShowOnDetail() bool
	ShowOnCreation() bool
	ShowOnUpdate() bool
	Resolve(ctx *gin.Context, model interface{})
	SetPanel(name string)
	GetPanel() string
	GetRules() []FieldRule
	GetAttribute() string
	Call(model interface{})
	// Fill(ctx *gin.Context, data map[string]interface{}, model interface{})
}

// 组件自定义列表页对应组件
type CustomIndexFieldComponent interface {
	IndexComponent()
}
// 组件自定义详情页对应组件
type CustomDetailFieldComponent interface {
	DetailComponent()
}
// 组件自定义创建页对应组件
type CustomCreationFieldComponent interface {
	CreationComponent()
}
// 组件自定义更新页对应组件
type CustomUpdateFieldComponent interface {
	UpdateComponent()
}
type FieldRule interface {
	GetRule() string
	GetMessage() string
}
