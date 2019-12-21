package request

import (
	"github.com/davecgh/go-spew/spew"
	"testing"
)

func TestFilters_Unmarshal(t *testing.T) {
	filters := Filters("")
	strings := filters.Unmarshal()
	spew.Dump(strings)
}
