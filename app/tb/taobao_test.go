package tb

import (
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"net/http"
	"regexp"
	"testing"
)

func TestTaobaoSdkService_Detail(t *testing.T) {
	service := &TaobaoSdkService{}

	data, err := service.Detail("600740693156")
	if err != nil {
		t.Fatal(err)
	}

	spew.Dump(data)
}

func TestIsUrlReg(t *testing.T) {
	type item struct {
		url  string
		isOK bool
	}
	var isUrl = regexp.MustCompile(`^https?:\/\/`)
	data := []item{
		{url: "http://q5q1efml2.bkt.clouddn.com/ezzNCEuFFJP5aAHeKAo9RqV1uCvylMmg", isOK: true},
		{url: "httpss://q5q1efml2.bkt.clouddn.com/ezzNCEuFFJP5aAHeKAo9RqV1uCvylMmg", isOK: false},
		{url: "https://img.alicdn.com/imgextra/i1/844323786/O1CN01TmUjkE1dq23evKe22_!!844323786.jpg", isOK: true},
		{url: "ezzNCEuFFJP5aAHeKAo9RqV1uCvylMmg", isOK: false},
	}

	for _, i := range data {
		matchString := isUrl.MatchString(i.url)
		if matchString != i.isOK {
			t.Fatal("err", i, matchString)
		}
	}
}

func TestKuaiDi100Api(t *testing.T) {
	// https://www.kuaidi100.com/query?type=zhongtong&postid=73123917441103
	service := &TaobaoSdkService{}
	build, err := service.build("https://www.kuaidi100.com/query", map[string]interface{}{
		"type":   "shunfeng",
		"postid": "SF1019151340747",
		"temp":   "0.9780776625299978",
	})
	if err != nil {
		t.Fatal(err)
	}
	get, err := http.Get(build)
	if err != nil {
		t.Fatal(err)
	}
	//all, err := ioutil.ReadAll(get.Body)
	//if err != nil {
	//	t.Fatal(err)
	//}
	data := map[string]interface{}{}
	err = json.NewDecoder(get.Body).Decode(&data)
	if err != nil {
		t.Fatal(err)
	}
	spew.Dump(data)
}
