package http

import (
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"go-shop-v2/pkg/utils"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"testing"
)

func TestTaobao_Api(t *testing.T) {
	api, err := url.Parse("https://acs.m.taobao.com/h5/mtop.taobao.detail.getdetail/6.0")
	if err != nil {
		t.Fatal(err)
	}

	query := url.Values{}

	data := map[string]string{
		"itemNumId": "600740693156",
	}

	bytes, _ := json.Marshal(data)

	query.Add("data", string(bytes))

	api.RawQuery = query.Encode()

	resp, err := http.Get(api.String())

	if err != nil {
		t.Fatal(err)
	}

	all, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%s", all)
}

func TestIndexController_Index(t *testing.T) {
	desc := "<p align=\"center\"><img src=\"//img.alicdn.com/imgextra/i1/844323786/O1CN01mcRswc1dq23UdR3yw_!!844323786.jpg\" align=\"middle\" size=\"750x750\"><img src=\"//img.alicdn.com/imgextra/i3/844323786/O1CN01zLu8gX1dq23bcYE5U_!!844323786.jpg\" align=\"middle\" size=\"750x750\"><img src=\"//img.alicdn.com/imgextra/i2/844323786/O1CN01eZsbIq1dq23cFfmkN_!!844323786.jpg\" align=\"middle\" size=\"750x750\" /></p>"
	compile, err := regexp.Compile("src=\"//img")
	if err != nil {
		t.Fatal(err)
	}
	s := compile.ReplaceAllString(desc, "src=\"https//img")
	spew.Dump(s)
}

func TestIndexController_TaobaoDetail(t *testing.T) {
	desc := "<p>&nbsp;</p> <p>&nbsp;</p> <p style=\"text-align:center;\"><img src=\"https://img.alicdn.com/imgextra/i3/2616970884/O1CN01Zp75Vq1IOugGap1ID_!!2616970884.jpg\" align=\"absmiddle\" size=\"790x450\" /></p> <p><img src=\"https://img.alicdn.com/imgextra/i2/2616970884/O1CN01Dxk5ZI1IOugbTsKSS_!!2616970884.jpg\" align=\"absmiddle\"><img src=\"https://img.alicdn.com/imgextra/i3/2616970884/O1CN01V1xlL71IOugOAH051_!!2616970884.jpg\" align=\"absmiddle\" size=\"790x14975\"><img src=\"https://img.alicdn.com/imgextra/i2/2616970884/O1CN01Vaq7K81IOuhd3Tsca_!!2616970884.jpg\" align=\"absmiddle\"><img src=\"https://img.alicdn.com/imgextra/i3/2616970884/O1CN01pOWuWh1IOuhZM3HhB_!!2616970884.jpg\" align=\"absmiddle\" /></p>"
	var imgRE = regexp.MustCompile(`<img[^>]+\bsrc=["']([^"']+)["']`)
	submatch := imgRE.FindAllStringSubmatch(desc, -1)

	images := sync.Map{}


	for _, match := range submatch {
		if len(match) >= 2 {
			images.Store(match[1],utils.RandomString(32))
		}
	}

	images.Range(func(key, value interface{}) bool {
		spew.Dump(key)
		spew.Dump(value)
		desc = strings.ReplaceAll(desc,key.(string),value.(string))
		return true
	})
	//for k,v:=range images {
	//	desc = strings.ReplaceAll(desc,k,v)
	//}

	spew.Dump(desc)

}

func TestTaobaoShortUrl(t *testing.T)  {
	uri:= "https://m.tb.cn/h.VX7dGft?sm=3d0625"
	resp, err := http.Get(uri)
	if err!=nil {
		t.Fatal(err)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err!=nil {
		t.Fatal(err)
	}
	spew.Dump(string(bytes))

}
