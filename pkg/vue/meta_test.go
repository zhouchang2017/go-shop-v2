package vue

import (
	"encoding/json"
	"testing"
)

func TestMeta_WithMeta(t *testing.T) {
	element := NewField("ID","id")
	element.Help("主键值")
	element.Component = "product-index"
	element.WithMeta("price",10)
	bytes, e := json.Marshal(element)
	if e!=nil {
		t.Fatal(e)
	}
	t.Logf("%s",bytes)
}

func TestElement_WithMeta(t *testing.T) {
	router := &Router{
		Path: "products",
		Name: "product.index",
	}

	i := &Router{
		Path: "products/:id",
		Name: "product.detail",
	}

	i.WithMeta("icon","www")

	router.AddChild(i)
	bytes, e := json.Marshal(router)
	if e!=nil {
		t.Fatal(e)
	}
	t.Logf("%s",bytes)
}