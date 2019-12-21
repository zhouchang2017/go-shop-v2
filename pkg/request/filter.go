package request

import (
	"encoding/base64"
	"encoding/json"
	"go-shop-v2/pkg/utils"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type FilterDirection int32

const (
	Filter_ASC  FilterDirection = 1  // 升序
	Filter_DESC FilterDirection = -1 // 降序
)

func (f FilterDirection) String() string {
	if f == 1 {
		return "asc"
	}
	return "desc"
}

type Filter struct {
	Search         string          `json:"search" form:"search"`
	WithTrashed    bool            `json:"with_trashed" form:"with_trashed"`
	CreatedAt      []time.Time     `json:"created_at" form:"created_at[]"`
	UpdatedAt      []time.Time     `json:"updated_at" form:"updated_at[]"`
	Page           int64           `json:"page" form:"page"`
	PerPage        int64           `json:"per_page" form:"per_page"`
	OrderBy        string          `json:"order_by" form:"order_by"`
	OrderDirection FilterDirection `json:"order_direction" form:"order_direction"`
	Options        string          `json:"options" form:"options"`
	WithAll        bool            `form:"with_all"`
	searchField    string
	Only           string `json:"only" form:"only"`
	Hidden         string `json:"hidden" form:"hidden"`
	only           []string
	hidden         []string
	query          map[string]interface{}
}

func (this *Filter) GetQuery() map[string]interface{}  {
	return this.query
}

func (this *Filter) AppendFilter(key string, value interface{}) {
	if this.query == nil {
		this.query = map[string]interface{}{}
	}
	this.query[key] = value
}

func (this *Filter) GetPerPage() int64 {
	if this.PerPage == 0 {
		return 15
	}
	return this.PerPage
}

func (this *Filter) SetSearchField(field string) {
	this.searchField = field
}

func (this *Filter) GetSearchField() string {
	if this.searchField == "" {
		return "name"
	}
	return this.searchField
}

func (this *Filter) AddOnly(key string) {
	for _, item := range this.only {
		if exists, _ := utils.InArray(key, item); !exists {
			this.only = append(this.only, key)
		}
	}
}

func (this *Filter) GetOnly() (res []string) {
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

func (this *Filter) AddHidden(key string) {
	for _, item := range this.hidden {
		if exists, _ := utils.InArray(key, item); !exists {
			this.hidden = append(this.hidden, key)
		}
	}
}

func (this *Filter) GetHidden() (res []string) {
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

func (this *Filter) GetPage() int64 {

	if this.WithAll {
		return -1
	}

	if this.Page <= 0 {
		return 1
	}
	return this.Page
}

func (this *Filter) GetFilterOptions() (opts map[string]string) {
	// javascript btoa(encodeURIComponent(JSON.stringify( options )))
	opts = map[string]string{}
	if this.Options == "" {
		return
	}

	// decodeBase64
	n, err := base64.StdEncoding.DecodeString(this.Options)
	if err != nil {
		return
	}

	st, err := url.PathUnescape(string(n))
	if err != nil {
		return
	}

	var options map[string]interface{}
	err = json.Unmarshal([]byte(st), &options)
	if err != nil {
		return
	}

	for key, value := range options {
		opts[key] = toString(value)
	}

	return opts
}

func toString(v interface{}) string {
	switch v.(type) {
	case float64:
		return strconv.Itoa(int(v.(float64)))
	case float32:
		return strconv.Itoa(int(v.(float32)))
	case string:
		return v.(string)
	case []interface{}:
		arr := []string{}
		for _, v := range v.([]interface{}) {
			arr = append(arr, toString(v))
		}
		return strings.Join(arr, ",")
	default:
		return ""
	}
}

func (this *Filter) GetWithTrashed() bool {
	return this.WithTrashed
}