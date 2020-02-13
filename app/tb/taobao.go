package tb

import (
	"encoding/json"
	"errors"
	"go-shop-v2/app/models"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/utils"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
)

type taobaoResponse struct {
	Data *tbResponseBody `json:"data"`
}

type tbResponseBody struct {
	//ApiStack []tbApiStackItem       `json:"apiStack"`
	Item     *tbResponseBodyItem    `json:"item"`
	MockData string                 `json:"mockData,omitempty"`
	Data     map[string]interface{} `json:"data,omitempty"`
	Props    map[string]interface{} `json:"props"`
	PropsCut string                 `json:"propsCut"`
	SkuBase  *skuBase               `json:"skuBase"`
}

type tbApiStackItem map[string]interface{}

type skuBase struct {
	Props []*tbSkuBaseItem `json:"props"`
	Skus  []*tbSkuBaseSku  `json:"skus"`
}

type tbSkuBaseSku struct {
	SkuId    string `json:"skuId"`
	PropPath string `json:"propPath"`
}

type tbSkuBaseItem struct {
	Name   string                `json:"name"`
	Pid    string                `json:"pid"`
	Values []*tbSkuBaseItemValue `json:"values"`
}

type tbSkuBaseItemValue struct {
	Name string `json:"name"`
	Vid  string `json:"vid"`
}

type tbResponseBodyItem struct {
	BrandValueId    string   `json:"brandValueId"`
	CartUrl         string   `json:"cartUrl"`
	CategoryId      string   `json:"categoryId"`
	CommentCount    string   `json:"commentCount"`
	H5moduleDescUrl string   `json:"h5moduleDescUrl"`
	Images          []string `json:"images"`
	ItemId          string   `json:"itemId"`
	Title           string   `json:"title"`
}

const DETAIL_API = "https://acs.m.taobao.com/h5/mtop.taobao.detail.getdetail/6.0"
const DESCRIPTION_API = "https://h5api.m.taobao.com/h5/mtop.taobao.detail.getdesc/6.0"

type TaobaoSdkService struct {
}

func (this *TaobaoSdkService) build(uri string, query map[string]interface{}) (string, error) {
	api, _ := url.Parse(uri)
	q := url.Values{}
	for key, value := range query {
		typeOf := reflect.TypeOf(value)
		var v string
		switch typeOf.Kind() {
		case reflect.String:
			v = value.(string)
		case reflect.Map:
			bytes, err := json.Marshal(value)
			if err != nil {
				return "", err
			}
			v = string(bytes)
		default:
			return "", errors.New("value 必须为 string 或者 map")
		}
		q.Add(key, v)
	}
	api.RawQuery = q.Encode()

	return api.String(), nil
}

func (this *TaobaoSdkService) Description(id string) (data string, err error) {
	uri, err := this.build(DESCRIPTION_API, map[string]interface{}{
		"jsv":     "2.5.1",
		"appKey":  "12574478",
		"t":       "1581237506638",
		"sign":    "281f62c619862ac35bd5592056b57a7d",
		"api":     "mtop.taobao.detail.getdesc",
		"v":       "6.0",
		"timeout": "20000",
		"data": map[string]string{
			"id":   id,
			"type": "1",
			"f":    "TB1Duv9sKH2gK0jSZJn8quT1Fla",
		},
	})

	if err != nil {
		return "", err
	}

	resp, err := http.Get(uri)
	if err != nil {
		return "", err
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var res map[string]interface{}
	err = json.Unmarshal(bytes, &res)
	if err != nil {
		return "", err
	}
	v, err := utils.MapGet(res, "data.pcDescContent")
	if err != nil {
		return "", err
	}
	return v.(string), nil
}

var imgReg, _ = regexp.Compile("src=\"//img")

func (this *TaobaoSdkService) Detail(id string) (data *models.Product, err error) {
	api, err := this.build(DETAIL_API, map[string]interface{}{
		"data": map[string]string{
			"itemNumId": id,
		},
	})
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(api)

	if err != nil {
		return nil, err
	}

	all, err := ioutil.ReadAll(resp.Body)

	var res taobaoResponse

	err = json.Unmarshal(all, &res)

	if err != nil {
		return
	}

	product := &models.Product{}

	if res.Data != nil {
		if res.Data.Item == nil {
			return nil, err2.NewFromCode(404).F("该商品不存在")
		}
		if res.Data.MockData != "" {
			var data map[string]interface{}
			if err := json.Unmarshal([]byte(res.Data.MockData), &data); err != nil {
				return nil, err
			}
			res.Data.MockData = ""
			res.Data.Data = data
		}

		if res.Data.Item != nil {
			product.Name = res.Data.Item.Title

			if s, err := this.Description(id); err == nil {
				product.Description = imgReg.ReplaceAllString(s, "src=\"https://img")
			}

			//if res.Data.Item.H5moduleDescUrl != "" {
			//	parse, err := url.Parse(res.Data.Item.H5moduleDescUrl)
			//	if err != nil {
			//		panic(err)
			//	}
			//
			//	parse.Scheme = "https"
			//
			//	if response, err := http.Get(parse.String()); err == nil {
			//		if bytes, err := ioutil.ReadAll(response.Body); err == nil {
			//
			//			spew.Dump(string(bytes))
			//		}
			//
			//	}
			//
			//}

		}

		// 基本属性组
		var productAttributes []*models.ProductAttribute

		basicAttrs, err := utils.MapGet(res.Data.Props, "groupProps.0.基本信息")
		if err == nil {
			if attrs, ok := basicAttrs.([]interface{}); ok {
				for _, attr := range attrs {
					if item, ok := attr.(map[string]interface{}); ok {
						for k, v := range item {
							productAttributes = append(productAttributes, &models.ProductAttribute{
								Name:  k,
								Value: v.(string),
							})
						}
					}

				}
			}

		}
		product.Attributes = productAttributes

		var options []*models.ProductOption
		if res.Data.SkuBase != nil {
			for _, prop := range res.Data.SkuBase.Props {
				option := models.NewProductOption(prop.Name)
				for _, value := range prop.Values {
					option.AddValues(option.NewValue(value.Name, value.Vid))
				}

				options = append(options, option)

			}
		}

		product.Options = options

		if res.Data.Data != nil {
			if price, err := utils.MapGet(res.Data.Data, "skuCore.sku2info.0.price.priceMoney"); err == nil {
				product.Price = int64(price.(float64))
			}
		}

		//qn := qiniu.GetQiniu()
		var images []string
		for _, image := range res.Data.Item.Images {
			if parse, err := url.Parse(image); err == nil {
				parse.Scheme = "https"
				images = append(images, parse.String())
			}
		}
		product.WithMeta("images", images)
	}

	return product, nil
}
