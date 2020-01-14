package fields

// 图片字段
type Image struct {
	*File
	Round string `json:"round"`
}

func NewImageField(name string, fieldName string, opts ...FieldOption) *Image {
	var options = []FieldOption{
		SetShowOnIndex(true),
	}
	options = append(options, opts...)
	image := &Image{
		File:NewFileField(name,fieldName,options...),
	}
	image.Type = "image"
	// 默认仅允许图片类型文件上传
	image.OnlyAcceptImage()
	// 默认样式为图片卡
	image.WithListTypePictureCard()
	// 默认开启预览
	image.ShowPreview()
	// 默认开启下载
	image.ShowDownload()
	return image
}


// https://tailwindcss.com/docs/border-radius/#app
// 圆角
func (this *Image) RoundedFull() *Image {
	this.Round = "rounded-full"
	return this
}

func (this *Image) RoundedLg() *Image {
	this.Round = "rounded-lg"
	return this
}

func (this *Image) RoundedSm() *Image {
	this.Round = "rounded-sm"
	return this
}

func (this *Image) Rounded() *Image  {
	this.Round = "rounded"
	return this
}

// 重写File
func (this *Image) Multiple() *Image {
	this.IsMultiple = true
	return this
}

// 允许上传类型
func (this *Image) WithAccept(accept string) *Image {
	this.Accept = &accept
	return this
}

// 开启预览
func (this *Image) ShowPreview() *Image {
	this.ShouldShowPreview = true
	return this
}

// 禁用预览
func (this *Image) HiddenPreview() *Image {
	this.ShouldShowPreview = true
	return this
}

// 开启下载
func (this *Image) ShowDownload() *Image {
	this.ShouldShowDownload = true
	return this
}

// 禁用下载
func (this *Image) HiddenDownload() *Image {
	this.ShouldShowDownload = false
	return this
}

// 设置上传地址
func (this *Image) WithAction(action string) *Image {
	this.Action = action
	return this
}

// 显示格式
func (this *Image) WithListTypePictureCard() *Image {
	this.ListType = "picture-card"
	return this
}

// 设置最大上传张数
func (this *Image) Limit(max int64) *Image {
	this.LimitMax = max
	return this
}

// 设置最大上传尺寸
func (this *Image) Size(size int64) *Image {
	this.LimitMaxSize = &size
	return this
}
