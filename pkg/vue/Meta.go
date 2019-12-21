package vue

import "encoding/json"

type Metable interface {
	WithMeta(key string, value interface{})
}

type MetaItems []*metaItem

func (m MetaItems) MarshalJSON() ([]byte, error) {
	maps := map[string]interface{}{}
	for _, item := range m {
		maps[item.Key] = item.Value
	}
	return json.Marshal(maps)
}


type metaItem struct {
	Key   string
	Value interface{}
}

