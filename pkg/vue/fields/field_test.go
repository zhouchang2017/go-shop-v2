package fields

import (
	"encoding/json"
	"testing"
)

func TestNewField(t *testing.T) {
	field := NewField("姓名", "Name")
	bytes, _ := json.Marshal(field)

	t.Logf("%s",bytes)
}
