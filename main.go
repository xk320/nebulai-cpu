package main

import (
	"encoding/json"
	"io/ioutil"
	"nebulai-cpu/apis"
	"nebulai-cpu/logger"
)

type AccountConfig struct {
	Token    string `json:"token"`
	JwtToken string `json:"jwt_token"`
	Proxy    string `json:"proxy"`
}

type MultiConfig struct {
	GpuEnabled bool            `json:"gpu_enabled"`
	Accounts   []AccountConfig `json:"accounts"`
}

func loadMultiConfig(filename string) (*MultiConfig, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var config MultiConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func main() {
	config, err := loadMultiConfig("config.json")
	if err != nil {
		logger.LogError("无法读取配置文件: %v", err)
		return
	}
	for i, acc := range config.Accounts {
		go func(idx int, acc AccountConfig) {
			token := acc.Token
			jwtToken := acc.JwtToken
			proxy := acc.Proxy
			logger.LogInfo("账号%d 启动，代理: %s", idx+1, proxy)
			apis.RunAccountTask(token, jwtToken, proxy, idx+1, config.GpuEnabled)
		}(i, acc)
	}
	select {} // 阻塞主线程
}
