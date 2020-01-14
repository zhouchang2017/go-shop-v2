package fields

// 头像字段
type Avatar struct {
	*Image
	Round string `json:"round"`
}

func NewAvatar(name string, fieldName string, opts ...FieldOption) *Avatar {
	avatar := &Avatar{
		Image:NewImageField(name,fieldName,opts...),
	}
	avatar.Image.Size(1)
	avatar.IsMultiple = false
	return avatar
}

// https://tailwindcss.com/docs/border-radius/#app
// 圆角
func (this *Avatar) RoundedFull() *Avatar {
	this.Round = "rounded-full"
	return this
}

func (this *Avatar) RoundedLg() *Avatar {
	this.Round = "rounded-lg"
	return this
}

func (this *Avatar) RoundedSm() *Avatar {
	this.Round = "rounded-sm"
	return this
}

func (this *Avatar) Rounded() *Avatar  {
	this.Round = "rounded"
	return this
}

// 设置最大上传尺寸
func (this *Avatar) Size(size int64) *Avatar {
	this.LimitMaxSize = &size
	return this
}