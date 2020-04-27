package qiniu

type Config struct {
	Drive            string `json:"drive" mapstructure:"drive"`
	QiniuDomain      string `json:"qiniu_domain" mapstructure:"qiniu_domain"`
	QiniuAccessKey   string `json:"qiniu_access_key" mapstructure:"qiniu_access_key"`
	QiniuSecretKey   string `json:"qiniu_secret_key" mapstructure:"qiniu_secret_key"`
	Bucket           string `json:"bucket" mapstructure:"bucket"`
	FileUploadAction string `json:"file_upload_action" mapstructure:"file_upload_action"`
}
