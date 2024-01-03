package main

import (
	"fmt"
	"bytes"
	"io/ioutil"
	"net/http"
	"encoding/base64"
	"gopkg.in/yaml.v2"
)

type Config struct {
	ServerAddress string `yaml:"server_address"`
	Port          int    `yaml:"port"`
}

func main() {
	configData, err := ioutil.ReadFile("config.yml")
	if err != nil {
		fmt.Println("读取配置文件失败:", err)
		return
	}

	var config Config
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		fmt.Println("解析配置文件失败:", err)
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "读取请求内容失败", http.StatusInternalServerError)
			return
		}

		encodedContent := base64.StdEncoding.EncodeToString(content)

		encodedBytes := []byte(encodedContent)

		_, err = http.Post(config.ServerAddress, "text/plain", bytes.NewBuffer(encodedBytes))
		if err != nil {
			http.Error(w, "发送请求失败", http.StatusInternalServerError)
			return
		}


		fmt.Fprintf(w, "已桥接请求")
	})

	addr := fmt.Sprintf(":%d", config.Port)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("启动HTTP服务失败:", err)
		return
	}
}
