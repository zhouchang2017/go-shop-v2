package qiniu

type Config struct {
	Drive          string `json:"drive"`
	QiniuDomain    string `json:"qiniu_domain"`
	QiniuAccessKey string `json:"qiniu_access_key"`
	QiniuSecretKey string `json:"qiniu_secret_key"`
	Bucket         string `json:"bucket"`
}
