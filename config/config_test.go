package config

import (
	"testing"
)

func TestNewConfig(t *testing.T) {
	//// get env file normal without test func
	//envPath, pathErr := utils.GetFilePath(1, ".env")
	//assert.NoError(t, pathErr)
	//// open file
	//file, openErr := os.Open(envPath)
	//assert.NoError(t, openErr)
	//defer file.Close()
	//// decode json
	//decoder := json.NewDecoder(file)
	//config := config{}
	//decodeErr := decoder.Decode(&config)
	//assert.NoError(t, decodeErr)
	//
	//// test func
	//testConfig := NewConfig()
	//assert.NotNil(t, testConfig)
	//assert.NotNil(t, testConfig.MongoCfg)
	//assert.NotNil(t, testConfig.MysqlCfg)
	//assert.NotNil(t, testConfig.QiniuCfg)
}