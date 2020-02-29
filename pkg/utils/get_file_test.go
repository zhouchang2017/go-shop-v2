package utils

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func TestGetFilePath(t *testing.T) {
	// init
	projectDir, err := os.Getwd()
	assert.NoError(t, err)
	pathArr := strings.Split(projectDir, "/")
	// test case
	situations := []struct {
		name      string
		input     func() (int, string)
		expectRes string
		expectErr error
	}{
		{
			"not exist file path",
			func() (int, string) {
				return 0, "notexistdir/notexistfile"
			},
			"",
			errors.New("not exist file"),
		},
		{
			"normal case",
			func() (int, string) {
				return 2, ".config.example"
			},
			".config.example",
			nil,
		},
	}
	// test
	for _, situation := range situations {
		// get test input
		absRootBackTimes, absRootTargetPath := situation.input()
		// get prefix root path of expect result
		projectRoot := strings.Join(pathArr[:len(pathArr)-absRootBackTimes], "/")
		// do func
		expectRes := fmt.Sprintf("%s/%s", projectRoot, situation.expectRes)
		res, err := GetFilePath(absRootBackTimes, absRootTargetPath)
		// test result
		if situation.expectErr != nil {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, expectRes, res)
		}
	}
}
