package qiniu

import (
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"testing"
)

func TestImage_MarshalJSON(t *testing.T) {
	image := Image("ezzNCEuFFJP5aAHeKAo9RqV1uCvylMmg")

	bytes, e := json.Marshal(image)
	if e != nil {
		t.Fatal(e)
	}

	spew.Dump(string(bytes))
}

func TestImage_UnmarshalJSON(t *testing.T) {
	data := []string{
		`"http://q5q1efml2.bkt.clouddn.com/ezzNCEuFFJP5aAHeKAo9RqV1uCvylMmg"`,
		`"https://img.alicdn.com/imgextra/i1/844323786/O1CN01TmUjkE1dq23evKe22_!!844323786.jpg"`,
		`"ezzNCEuFFJP5aAHeKAo9RqV1uCvylMmg"`,
	}

	type request struct {
		Image Image `json:"image"`
	}

	//req:=&request{Image:"http://q5q1efml2.bkt.clouddn.com/ezzNCEuFFJP5aAHeKAo9RqV1uCvylMmg"}

	for _, i := range data {
		var image Image
		err := json.Unmarshal([]byte(i), &image)
		if err != nil {
			t.Fatal(err)
		}
		spew.Dump(image)
	}
}
