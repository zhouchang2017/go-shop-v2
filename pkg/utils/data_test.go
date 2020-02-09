package utils

import (
	"github.com/davecgh/go-spew/spew"
	"testing"
)

func TestDataGet(t *testing.T) {
	data := map[string]interface{}{
		"props": map[string]interface{}{
			"groupProps": []map[string]interface{}{
				{
					"基本信息": []map[string]interface{}{
						{"品牌": "Lacoste/拉科斯特"},
						{
							"功能": "轻质",
						},
						{
							"闭合方式": "系带",
						},
						{
							"尺码": "6,6.5,7.5,8,9",
						},
						{
							"图案": "拼色",
						},
						{
							"风格": "运动",
						},
						{
							"流行元素": "车缝线",
						},
						{
							"鞋跟高": "中跟(3-5cm)",
						},
						{
							"颜色分类": "蓝色 092,白色 147",
						},
						{
							"货号": "M0023R M2",
						},
						{
							"季节": "春秋",
						},
						{
							"鞋头款式": "尖头",
						},
						{
							"场合": "日常",
						},
						{
							"跟底款式": "平跟",
						},
						{
							"鞋面内里材质": "混合材质",
						},
						{
							"鞋制作工艺": "胶粘鞋",
						},
						{
							"鞋面材质": "涤沦",
						},
						{
							"款式": "运动休闲鞋",
						},
					},
				},
			},
		},
	}



	get,err := MapGet(data, "props.groupProps.0.基本信息.0.品牌")
	if err!=nil {
		t.Fatal(err)
	}
	spew.Dump(get)
}
