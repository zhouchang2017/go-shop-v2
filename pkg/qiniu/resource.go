package qiniu

import (
	"fmt"
)

type ImageInfo struct {
	Format string `json:"format"`
	Height int64  `json:"height"`
	Width  int64  `json:"width"`
	Size   int64  `json:"size"`
}

type ImageAve struct {
	RGB string `json:"RGB"`
}

// 静态资源结构
type Resource struct {
	Key       string     `json:"key"`                          // 文件保存在空间中的资源名
	Name      string     `json:"name"`                         // 原始文件名
	Bucket    string     `json:"bucket"`                       // 目标空间名
	Domain    string     `json:"domain"`                       // 储存域名
	Drive     string     `json:"drive"`                        // 驱动
	MimeType  string     `json:"mime_type" bson:"mime_type"`   // 文件类型
	Ext       string     `json:"ext"`                          // 文件扩展名
	ImageInfo *ImageInfo `json:"image_info" bson:"image_info"` // 图片信息
	ImageAve  *ImageAve  `json:"image_ave" bson:"image_ave"`   // 图片主色调
}

func (this Resource) GetKey() string {
	return this.Key
}

func (this Resource) GetName() string {
	return this.Name
}

func (this Resource) GetDrive() string {
	return this.Drive
}

// 预览地址
func (this Resource) PreviewUrl() string {
	return fmt.Sprintf("%s/%s", this.Domain, this.Key)
}
