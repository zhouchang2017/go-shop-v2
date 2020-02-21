package fields

import (
	"context"
	"go-shop-v2/pkg/qiniu"
	"log"
)

// 默认上传地址
var DefaultFileUploadAction string

// 文件类型字段
// https://element.eleme.cn/#/zh-CN/component/upload
// 目前存储驱动仅用七牛，暂未抽象接口
type File struct {
	*Field             `inline`
	Drive              string  `json:"drive"`
	IsMultiple         bool    `json:"multiple"`             // 是否支持多选文件
	Accept             *string `json:"accept"`               // 接受上传的文件类型,多类型逗号分割
	Drag               bool    `json:"drag"`                 // 是否开启拖拽上传
	ShouldShowPreview  bool    `json:"should_show_preview"`  // 是否支持预览
	ShouldShowDownload bool    `json:"should_show_download"` // 是否支持下载
	Token              string  `json:"token"`                // 上传token
	Action             string  `json:"action"`               // 上传地址
	ListType           string  `json:"list_type"`            // 文件列表的类型 text/picture/picture-card
	LimitMax           int64   `json:"limit"`                // 最大允许上传个数
	LimitMaxSize       *int64  `json:"limit_max_size"`       // 允许上传文件最大体积
	Type               string  `json:"type"`                 // 文件字段默认
}

func NewFileField(name string, fieldName string, opts ...FieldOption) *File {
	var options = []FieldOption{
		SetPrefixComponent(true),
		SetShowOnIndex(false),
		SetShowOnDetail(true),
		SetShowOnCreation(true),
		SetShowOnUpdate(true),
		WithComponent("file-field"),
		SetTextAlign("left"),
	}
	options = append(options, opts...)
	return &File{Field: NewField(name, fieldName, options...),
		Drive:    "qiniu",
		ListType: "text",
		Type:     "file",
		LimitMax: 1,
	}
}

// 在创建页、更新页响应JSON格式返还给前端时调用的一个钩子
func (this *File) Call() {
	// 目前就只有七牛、暂未抽象。写死了哈

	var token string
	var err error

	switch this.Type {
	case "file":
		token, err = qiniu.GetQiniu().FileToken(context.Background())
	default:
		token, err = qiniu.GetQiniu().ImageToken(context.Background())
	}
	if err != nil {
		log.Printf("get qiniu token error:%s\n", err)
	}
	if this.Action == "" {
		this.Action = DefaultFileUploadAction
	}
	this.Token = token
}

func (this *File) Multiple() *File {
	this.IsMultiple = true
	return this
}

func (this *File) WithDrag() *File {
	this.Drag = true
	return this
}

func (this *File) WithAccept(accept string) *File {
	this.Accept = &accept
	return this
}

func (this *File) OnlyAcceptImage() *File {
	accept := "image/*"
	this.Accept = &accept
	return this
}

func (this *File) ShowDownload() *File {
	this.ShouldShowDownload = true
	return this
}

func (this *File) HiddenDownload() *File {
	this.ShouldShowDownload = false
	return this
}

func (this *File) WithAction(action string) *File {
	this.Action = action
	return this
}

func (this *File) WithListTypeText() *File {
	this.ListType = "text"
	return this
}

func (this *File) WithListTypePictureCard() *File {
	this.ListType = "picture-card"
	return this
}

func (this *File) Limit(max int64) *File {
	this.LimitMax = max
	return this
}

func (this *File) Size(size int64) *File {
	this.LimitMaxSize = &size
	return this
}
