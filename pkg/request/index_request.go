package request

import (
	"encoding/base64"
	"encoding/json"
	"go-shop-v2/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"net/url"
	"strings"
)

// 列表页请求通用参数
type IndexRequest struct {
	searchField    string
	Search         string          `json:"search" form:"search"`                   // 搜索关键词
	Trashed        bool            `json:"trashed" form:"trashed"`                 // 软删除
	Page           int64           `json:"page" form:"page"`                       // 当前页
	PerPage        int64           `json:"per_page" form:"per_page"`               // 分页步长
	OrderBy        string          `json:"order_by" form:"order_by"`               // 排序字段
	OrderDirection FilterDirection `json:"order_direction" form:"order_direction"` // 方向
	Only           string          `json:"only" form:"only"`                       // 仅显示字段
	Hidden         string          `json:"hidden" form:"hidden"`                   // 忽略字段
	only           []string
	hidden         []string
	Filters        Filters                `json:"filters" form:"filters"` // 筛选
	query          map[string]interface{} // 自定义搜索
	projection     map[string]interface{} // 自定义projection
}

func (this *IndexRequest) SetSearchField(field string) {
	this.searchField = field
}

func (this *IndexRequest) Query() map[string]interface{} {
	return this.query
}

func (this *IndexRequest) GetSearchField() string {
	return this.searchField
}

func (this *IndexRequest) AppendFilter(key string, value interface{}) {
	if this.query == nil {
		this.query = map[string]interface{}{}
	}
	this.query[key] = value
}

func (this *IndexRequest) GetPage() int64 {
	if (this.Page == 0) {
		this.Page = 1
	}
	return this.Page
}

func (this *IndexRequest) GetPerPage() int64 {
	if (this.PerPage == 0) {
		this.PerPage = 15
	}
	return this.PerPage
}

func (this *IndexRequest) Sort() (bson.M, bool) {
	if this.OrderBy != "" {
		negative := 1
		if this.OrderDirection == Filter_DESC {
			negative = -1
		}
		return bson.M{this.OrderBy: negative}, true
	}
	return bson.M{"_id": -1}, true
}

func (this *IndexRequest) AppendProjection(key string,value interface{}) {
	if this.projection == nil {
		this.projection = map[string]interface{}{}
	}
	this.projection[key] = value
}

func (this *IndexRequest) Projection() (bson.M, bool) {
	res := bson.M{}
	for _, field := range this.GetOnly() {
		res[field] = 1
	}
	for _, field := range this.getHidden() {
		res[field] = 0
	}

	if this.projection != nil {
		for key, value := range this.projection {
			res[key] = value
		}
	}

	return res, len(res) > 0
}

func (this *IndexRequest) AddOnly(key string) {
	for _, item := range this.only {
		if exists, _ := utils.InArray(key, item); !exists {
			this.only = append(this.only, key)
		}
	}
}

func (this *IndexRequest) AddHidden(key string) {
	for _, item := range this.hidden {
		if exists, _ := utils.InArray(key, item); !exists {
			this.hidden = append(this.hidden, key)
		}
	}
}

func (this *IndexRequest) GetOnly() (res []string) {
	if this.Only == "" {
		return this.only
	}
	res = this.only
	for _, field := range strings.Split(this.Only, ",") {
		if exists, _ := utils.InArray(field, res); !exists {
			res = append(res, field)
		}
	}
	return
}

func (this *IndexRequest) getHidden() (res []string) {
	if this.Hidden == "" {
		return this.hidden
	}
	res = this.hidden
	for _, field := range strings.Split(this.Hidden, ",") {
		if exists, _ := utils.InArray(field, res); !exists {
			res = append(res, field)
		}
	}
	return
}

type Filters string

func (f Filters) Unmarshal() (filters map[string]interface{}) {
	// javascript btoa(encodeURIComponent(JSON.stringify( options )))
	filters = map[string]interface{}{}
	if f == "" {
		return
	}

	// decodeBase64
	n, err := base64.StdEncoding.DecodeString(string(f))

	if err != nil {
		return
	}

	st, err := url.PathUnescape(string(n))
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(st), &filters)
	if err != nil {
		return
	}

	return
}

func (f Filters) Decode() (filters interface{}) {

	if f == "" {
		return nil
	}

	// decodeBase64
	n, err := base64.StdEncoding.DecodeString(string(f))
	if err != nil {
		return
	}

	st, err := url.PathUnescape(string(n))
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(st), &filters)
	if err != nil {
		return
	}

	return
}
