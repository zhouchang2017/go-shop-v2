package go_shop_v2

import (
	"encoding/json"
	"fmt"
	conf "go-shop-v2/config"
	"os"
	"path"
	"strings"
)

func TestInit() {
	projectName := "go-shop-v2"
	configFileName := ".config"
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	split := strings.Split(dir, projectName)

	configFilePath := path.Join(split[0], projectName, configFileName)

	// open file
	file, openErr := os.Open(configFilePath)
	if openErr != nil {
		panic(fmt.Sprintf("open config file failed caused of %s", openErr.Error()))
	}
	// decode json
	if err := json.NewDecoder(file).Decode(&conf.Config); err != nil {
		panic(fmt.Sprintf("decode config file failed caused of %s", err))
	}

	defer file.Close()
}
