package request

import (
	"github.com/davecgh/go-spew/spew"
	"testing"
)

func TestFilters_Unmarshal(t *testing.T) {
	filters := Filters("JTVCJTdCJTIya2V5JTIyJTNBJTIyc2hvcHMlMjIlMkMlMjJ2YWx1ZSUyMiUzQSU1QiUyMjVkZmQ4YWFkN2FhYjEyMDg4YTI2YmFmNiUyMiUyQyUyMjVlMDIyNTk4NGQ1NzBiYzk1ODNiN2NlMyUyMiU1RCU3RCU1RA==")
	strings := filters.Unmarshal()
	spew.Dump(strings)
}
