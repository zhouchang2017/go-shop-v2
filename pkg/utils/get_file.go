package utils

import (
	"fmt"
	"os"
	"strings"
)

func GetFilePath(absRootBackTimes int, absRootTargetPath string) (string, error) {
	// get pwd dir
	pwdDir, pwdErr := os.Getwd()
	if pwdErr != nil {
		return "", pwdErr
	}
	// contact path
	pathArr := strings.Split(pwdDir, "/")
	if absRootBackTimes > 0 && len(pathArr) <= absRootBackTimes {
		return "", fmt.Errorf("not enough path length %s", pwdDir)
	}
	pathArr = pathArr[:len(pathArr) - absRootBackTimes]
	filePath := fmt.Sprintf("%s/%s", strings.Join(pathArr, "/"), absRootTargetPath)
	// check exist
	if _, fileExist := os.Stat(filePath); fileExist != nil {
		return "", fileExist
	}
	// return
	return filePath, nil
}
