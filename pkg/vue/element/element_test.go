package element

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestNewElement(t *testing.T) {
	element := NewElement()
	element.WithMeta("name","张")
	bytes, _ := json.Marshal(element)
	fmt.Printf("%s",bytes)
}
