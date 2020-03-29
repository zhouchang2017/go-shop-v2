package app

import (
	"github.com/davecgh/go-spew/spew"
	"testing"
)

func TestConfig_URI(t *testing.T) {
	config := Config{
		Username: "root",
		Password: "12345678",
	}

	spew.Dump(config.URI())
}
