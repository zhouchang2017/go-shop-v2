package fields

import (
	"context"
	"go-shop-v2/pkg/qiniu"
	"log"
)

type RichText struct {
	*Field
	qiniu bool
	Token string `json:"token"`
	Action string `json:"action"`
}

func NewRichTextField(name string, fieldName string, opts ...FieldOption) *RichText {
	var options = []FieldOption{
		SetPrefixComponent(true),
		SetShowOnIndex(false),
		SetShowOnDetail(true),
		SetShowOnCreation(true),
		SetShowOnUpdate(true),
		WithComponent("rich-text-field"),
		SetTextAlign("left"),
		SetAsHtml(true),
	}
	options = append(options,opts...)

	return &RichText{Field: NewField(name, fieldName, options...)}
}

// 在创建页、更新页响应JSON格式返还给前端时调用的一个钩子
func (this *RichText) Call() {
	if this.qiniu {
		// 目前就只有七牛、暂未抽象。写死了哈
		token, err := qiniu.Token(context.Background())
		if err != nil {
			log.Printf("get qiniu token error:%s\n", err)
		}
		if this.Action == "" {
			this.Action = DefaultFileUploadAction
		}
		this.Token = token
	}
}

func (this *RichText) UseQiniu() *RichText  {
	this.qiniu = true
	return this
}