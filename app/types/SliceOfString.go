package types

import (
	"encoding/json"
	"strings"
)

type SliceOfString string

func (s SliceOfString) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Split())
}

func (s SliceOfString) Split() []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(string(s), ",")
}

func (s SliceOfString) Make(data []string) SliceOfString  {
	return SliceOfString(strings.Join(data, ","))
}
