package wechat

type Config struct {
	AppId          string `json:"app_id" mapstructure:"app_id"`
	AppSecret      string `json:"app_secret" mapstructure:"app_secret"`
	IsProd         bool   `json:"is_prod" mapstructure:"is_prod"`
	Token          string `json:"token" mapstructure:"token"`
	EncodingAESKey string `json:"encoding_aes_key" mapstructure:"encoding_aes_key"`
}

type PayConfig struct {
	AppId           string `json:"app_id" mapstructure:"app_id"`
	AppKey          string `json:"app_key" mapstructure:"app_key"`
	MchId           string `json:"mch_id" mapstructure:"mch_id"`
	IsProd          bool   `json:"is_prod" mapstructure:"is_prod"`
	CertFilePath    string `json:"cert_file_path" mapstructure:"cert_file_path"`
	KeyFilePath     string `json:"key_file_path" mapstructure:"key_file_path"`
	Pkcs12FilePath  string `json:"pkcs_12_file_path" mapstructure:"pkcs_12_file_path"`
	NotifyUrl       string `json:"notify_url" mapstructure:"notify_url"`
	RefundNotifyUrl string `json:"refund_notify_url" mapstructure:"refund_notify_url"`
}
