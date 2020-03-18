package wechat

type Config struct {
	AppId     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

type PayConfig struct {
	AppId     string `json:"app_id"`
	AppKey    string `json:"app_key"`
	MchId     string `json:"mch_id"`
	CertPath  string `json:"cert_path"`
	KeyPath   string `json:"key_path"`
	NotifyUrl string `json:"notify_url"`
}
