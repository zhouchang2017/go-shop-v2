package qiniu

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	err2 "go-shop-v2/pkg/err"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sync"
)

var once sync.Once
var instance *Qiniu

func GetQiniu() *Qiniu {
	return instance
}

type Qiniu struct {
	accessKey        string
	secretKey        string
	bucket           string
	domain           string
	fileUploadAction string
}

func NewQiniu(config Config) *Qiniu {
	once.Do(func() {
		instance = &Qiniu{
			accessKey:        config.QiniuAccessKey,
			secretKey:        config.QiniuSecretKey,
			bucket:           config.Bucket,
			domain:           config.QiniuDomain,
			fileUploadAction: config.FileUploadAction,
		}
	})
	return instance
}

func (this *Qiniu) mac() *qbox.Mac {
	return qbox.NewMac(this.accessKey, this.secretKey)
}

func (this *Qiniu) FileUploadAction() string {
	return this.fileUploadAction
}

func (this *Qiniu) BucketManager() *storage.BucketManager {
	config := storage.Config{}
	return storage.NewBucketManager(this.mac(), &config)
}

func (this Qiniu)Domain() string  {
	return this.domain
}

func (this Qiniu) Name() string {
	return "qiniu"
}

// 图片上传token
func (this *Qiniu) ImageToken(ctx context.Context) (token string, err error) {
	putPolicy := storage.PutPolicy{
		Scope: this.bucket,
	}
	putPolicy.Expires = 7200 //示例2小时有效期
	return putPolicy.UploadToken(this.mac()), nil
}

// 文件上传token
func (this *Qiniu) FileToken(ctx context.Context) (token string, err error) {
	body := fmt.Sprintf(`{"key":"$(key)","name":"$(fname)","bucket":"$(bucket)","mime_type":"$(mimeType)","ext":"$(ext)","drive":"%s","domain":"%s"}`, this.Name(), this.domain)
	putPolicy := storage.PutPolicy{
		Scope:      this.bucket,
		ReturnBody: body,
	}
	putPolicy.Expires = 7200 //示例2小时有效期

	return putPolicy.UploadToken(this.mac()), nil
}

func (this *Qiniu) Token(ctx context.Context) (token string, err error) {
	body := fmt.Sprintf(`{"key":"$(key)","name":"$(fname)","bucket":"$(bucket)","mime_type":"$(mimeType)","ext":"$(ext)","drive":"%s","domain":"%s"}`, this.Name(), this.domain)
	putPolicy := storage.PutPolicy{
		Scope:      this.bucket,
		ReturnBody: body,
	}
	putPolicy.Expires = 7200 //示例2小时有效期

	return putPolicy.UploadToken(this.mac()), nil
}

// 辅助函数
func Token(ctx context.Context) (token string, err error) {
	return instance.Token(ctx)
}

func (this *Qiniu) PutByUrl(ctx context.Context, url string, key string) (res Image, err error) {
	fetchRet, err := this.BucketManager().Fetch(url, this.bucket, key)
	if err != nil {
		return res, err
	}
	return Image(fetchRet.Key), nil
}

func (this *Qiniu) Put(ctx context.Context, key string, file *os.File) (res *Resource, err error) {
	formUploader := storage.NewFormUploader(nil)

	res = &Resource{}

	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:drive": this.Name(),
		},
	}

	defer file.Close()
	token, err := this.Token(ctx)
	if err != nil {
		return
	}
	err = formUploader.PutFileWithoutKey(ctx, res, token, file.Name(), &putExtra)

	defer os.Remove(file.Name())
	return
}

func (this *Qiniu) GetResourceURL(ctx context.Context, key string) (url string, err error) {
	return this.getFileUrl(key), nil
}

func (this *Qiniu) getFileUrl(key string) string {
	return storage.MakePublicURL(this.domain, key)
}

func (this *Qiniu) Get(resouce *Resource) (f *os.File, err error) {
	return this.get(resouce.Key)
}

func (this *Qiniu) get(key string) (f *os.File, err error) {
	publicAccessURL := this.getFileUrl(key)
	parse, err := url.Parse(publicAccessURL)
	if err != nil {
		return nil, err
	}
	parse.Scheme = "http"
	//http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, err := http.Get(parse.String())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	file, err := ioutil.TempFile("", "tmpfile")
	if err != nil {
		return nil, err
	}
	defer os.Remove(file.Name())
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	_, err = file.Write(bytes)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (this *Qiniu) Delete(ctx context.Context, key string) error {
	return this.BucketManager().Delete(this.bucket, key)
}

func (this *Qiniu) HttpHandle(router gin.IRouter) {
	router.GET("/qiniu", func(c *gin.Context) {
		tokenType := c.Query("type")
		if tokenType == "" {
			tokenType = "image"
		}
		var token string
		var err error

		switch tokenType {
		case "file":
			token, err = instance.FileToken(context.Background())
		default:
			token, err = instance.ImageToken(context.Background())
		}
		if err != nil {
			err2.ErrorEncoder(nil, err, c.Writer)
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
	})
}
