package tb

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-shop-v2/app/models"
	err2 "go-shop-v2/pkg/err"
	"go-shop-v2/pkg/qiniu"
	"go-shop-v2/pkg/utils"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"sync"
)

type taobaoResponse struct {
	Data *tbResponseBody `json:"data"`
}

type tbResponseBody struct {
	//ApiStack []tbApiStackItem       `json:"apiStack"`
	Item                  *tbResponseBodyItem    `json:"item"`
	MockData              string                 `json:"mockData,omitempty"`
	Data                  map[string]interface{} `json:"data,omitempty"`
	Props                 map[string]interface{} `json:"props"`
	PropsCut              string                 `json:"propsCut"`
	SkuBase               *skuBase               `json:"skuBase"`
	productOptions        models.ProductOptions
	optionValuesMap       map[string]*models.OptionValue
	productOptionsRunOnce sync.Once
}

func (this *tbResponseBody) parseMockData() {
	// 解析 MockData
	if this.MockData != "" && this.Data == nil {
		var data map[string]interface{}
		err := json.Unmarshal([]byte(this.MockData), &data)

		if err == nil {
			this.Data = data
		}
	}
}

func (this *tbResponseBody) GetProductPrice() uint64 {
	this.parseMockData()

	if this.Data != nil {
		if price, err := utils.MapGet(this.Data, "skuCore.sku2info.0.price.priceMoney"); err == nil {
			return uint64(price.(float64))
		}
	}

	return 0
}

func (this *tbResponseBody) GetItemPrice(itemId string) uint64 {
	this.parseMockData()
	if this.Data != nil {
		key := fmt.Sprintf("skuCore.sku2info.%s.price.priceMoney", itemId)
		if price, err := utils.MapGet(this.Data, key); err == nil {
			return uint64(price.(float64))
		}
	}

	return 0

}

func (this *tbResponseBody) GetProductOptions() []*models.ProductOption {

	this.productOptionsRunOnce.Do(func() {
		if this.SkuBase != nil {
			this.optionValuesMap = map[string]*models.OptionValue{}
			for _, prop := range this.SkuBase.Props {
				option := models.NewProductOption(prop.Name)
				var hasImage bool
				for _, value := range prop.Values {
					image := value.Image
					if image != "" {
						if parse, err := url.Parse(image); err == nil {
							hasImage = true
							parse.Scheme = "https"
							image = parse.String()
						}
					}

					optionValue := option.NewValue(value.Name).SetImage(image)

					option.AddValues(optionValue)

					this.optionValuesMap[fmt.Sprintf("%s:%s", prop.Pid, value.Vid)] = optionValue
				}
				if hasImage {
					option.Image = true
				}
				this.productOptions = append(this.productOptions, option)
			}
			sort.Sort(this.productOptions)
		}
	})
	return this.productOptions

}

func (this *tbResponseBody) GetOptionValueByPropPath(path string) (optionValue *models.OptionValue, err error) {
	this.GetProductOptions()
	if value, ok := this.optionValuesMap[path]; ok {
		optionValue = value
		return
	}
	err = fmt.Errorf("OptionValue prop path = %s , not found!!", path)
	return
}

func (this *tbResponseBody) GetAttributes() []*models.ProductAttribute {
	// 基本属性组
	var productAttributes []*models.ProductAttribute

	if this.Props != nil {
		basicAttrs, err := utils.MapGet(this.Props, "groupProps.0.基本信息")
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
	}
	return productAttributes
}

func (this *tbResponseBody) GetName() string {
	if this.Item != nil {
		return this.Item.Title
	}
	return ""
}

func (this *tbResponseBody) GetImages() []string {
	var images []string
	if this.Item != nil {
		for _, image := range this.Item.Images {
			if parse, err := url.Parse(image); err == nil {
				parse.Scheme = "https"
				images = append(images, parse.String())
			}
		}
	}
	return images
}

func (this *tbResponseBody) GetItems() (items []*models.Item, err error) {
	if this.SkuBase != nil {
		for _, sku := range this.SkuBase.Skus {
			item := &models.Item{}

			item.Price = this.GetItemPrice(sku.SkuId)
			propPaths := strings.Split(sku.PropPath, ";")
			var optionValues []*models.OptionValue
			for _, path := range propPaths {
				optionValue, err := this.GetOptionValueByPropPath(path)
				if err != nil {
					return nil, err
				}
				optionValues = append(optionValues, optionValue)
			}

			// 排序
			sortValues := models.SortOptionValues{
				Values:  optionValues,
				Options: this.GetProductOptions(),
			}

			sort.Sort(sortValues)
			item.OptionValues = sortValues.Values
			items = append(items, item)
		}
	}
	return
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
	Name  string `json:"name"`
	Vid   string `json:"vid"`
	Image string `json:"image,omitempty"`
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

		product.Name = res.Data.GetName()
		// 图集
		for _, image := range res.Data.GetImages() {
			product.Images = append(product.Images,qiniu.NewImage(image))
		}
		// 描述
		if s, err := this.Description(id); err == nil {
			product.Description = imgReg.ReplaceAllString(s, "src=\"https://img")
		}

		// 基本属性组
		product.Attributes = res.Data.GetAttributes()
		// 销售属性
		product.Options = res.Data.GetProductOptions()
		// 价格
		product.Price = res.Data.GetProductPrice()
		// 变体
		if items, err := res.Data.GetItems(); err == nil {
			product.Items = items
		}

	}

	return product, nil
}
