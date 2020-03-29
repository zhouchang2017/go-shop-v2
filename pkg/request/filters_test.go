package request

import (
	"encoding/base64"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"net/url"
	"testing"
)

func TestFilters_Unmarshal(t *testing.T) {
	filters := Filters("JTVCJTdCJTIya2V5JTIyJTNBJTIyc2hvcHMlMjIlMkMlMjJ2YWx1ZSUyMiUzQSU1QiUyMjVkZmQ4YWFkN2FhYjEyMDg4YTI2YmFmNiUyMiUyQyUyMjVlMDIyNTk4NGQ1NzBiYzk1ODNiN2NlMyUyMiU1RCU3RCU1RA==")
	strings := filters.Unmarshal()
	spew.Dump(strings)

	filter := map[string]interface{}{
		"status": []int{0},
	}

	marshal, err := json.Marshal(filter)
	if err != nil {
		panic(err)
	}

	escape := url.PathEscape(string(marshal))

	encodeToString := base64.StdEncoding.EncodeToString([]byte(escape))
	spew.Dump(encodeToString)

	f := Filters(encodeToString)

	spew.Dump(f.Unmarshal())
}
