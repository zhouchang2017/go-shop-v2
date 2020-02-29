package qiniu

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

type Image string

func NewImage(url string) Image {
	split := strings.Split(url, fmt.Sprintf("%s/", GetQiniu().domain))
	join := strings.Join(split, "")
	return Image(join)
}

var stringReg = regexp.MustCompile(`^\"(.*)\"$`)

func (i *Image) UnmarshalJSON(data []byte) error {
	all := bytes.ReplaceAll(data, []byte("\""), []byte(""))

	var buffer *bytes.Buffer

	buffer = bytes.NewBuffer(all)

	dataString := buffer.String()
	split := strings.Split(dataString, fmt.Sprintf("%s/", GetQiniu().domain))
	//split := strings.Split(dataString, fmt.Sprintf("%s/", "http://q5q1efml2.bkt.clouddn.com"))
	join := strings.Join(split, "")
	*i = Image(join)
	return nil
}

func (i Image) MarshalJSON() ([]byte, error) {
	if string(i) == "" {
		return bytes.NewBufferString("null").Bytes(), nil
	}
	bufferString := bytes.NewBufferString(`"`)
	bufferString.WriteString(i.Src())
	bufferString.WriteString(`"`)
	return bufferString.Bytes(), nil
}

var isUrl = regexp.MustCompile(`^https?:\/\/`)

func (i Image) Src() string {
	if string(i) == "" {
		return ""
	}
	if isUrl.MatchString(string(i)) {
		return string(i)
	}
	if GetQiniu() == nil {
		return string(i)
	}
	return fmt.Sprintf("%s/%s", GetQiniu().domain, string(i))
}

func (i Image) IsUrl() bool {
	return isUrl.MatchString(string(i))
}
