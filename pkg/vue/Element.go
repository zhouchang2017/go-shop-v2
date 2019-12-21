package vue

type Element struct {
	Meta         MetaItems `json:"meta"`
	Component    string    `json:"component"`
	OnlyOnDetail bool      `json:"onlyOnDetail"`
}

func (m *Element) WithMeta(key string, value interface{}) {
	m.Meta = append(m.Meta, &metaItem{key, value})
}