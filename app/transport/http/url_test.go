package http

import (
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"testing"
)

func TestTaobao_Api(t *testing.T) {
	api, err := url.Parse("https://acs.m.taobao.com/h5/mtop.taobao.detail.getdetail/6.0")
	if err!=nil {
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

	if err!=nil {
		t.Fatal(err)
	}

	all, err := ioutil.ReadAll(resp.Body)

	if err!=nil {
		t.Fatal(err)
	}

	t.Logf("%s",all)
}

func TestIndexController_Index(t *testing.T) {
	desc:="<p align=\"center\"><img src=\"//img.alicdn.com/imgextra/i1/844323786/O1CN01mcRswc1dq23UdR3yw_!!844323786.jpg\" align=\"middle\" size=\"750x750\"><img src=\"//img.alicdn.com/imgextra/i3/844323786/O1CN01zLu8gX1dq23bcYE5U_!!844323786.jpg\" align=\"middle\" size=\"750x750\"><img src=\"//img.alicdn.com/imgextra/i2/844323786/O1CN01eZsbIq1dq23cFfmkN_!!844323786.jpg\" align=\"middle\" size=\"750x750\" /></p>"
	compile, err := regexp.Compile("src=\"//img")
	if err!=nil {
		t.Fatal(err)
	}
	s := compile.ReplaceAllString(desc, "src=\"https//img")
	spew.Dump(s)
}