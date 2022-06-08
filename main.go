package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Config struct {
	FileId string
	Token  string
	NodeId string
}

func loadConfig() Config {
	data, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Println(err)
	}
	var config Config
	//读取的数据为json格式，需要进行解码
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Println(err)
	}
	return config
}

func getFileJson(config Config) string {
	url := "https://api.figma.com/v1/files/" + config.FileId
	if len(config.NodeId) > 0 {
		url = url + "/nodes?ids=" + config.NodeId
	}
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("x-figma-token", config.Token)
	client := http.Client{
		// Timeout: time.Duration(5 * time.Second),
	}
	res, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return ""
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(body)
}

func saveResult(content string, config Config) {
	var temp map[string]interface{}
	json.Unmarshal([]byte(content), &temp)
	name := temp["name"].(string)
	if len(config.NodeId) > 0 {
		name = name + ":" + config.NodeId
	}
	name = name + ".json"
	create, err := os.Create(name)
	if err != nil {
		log.Println("cretre error", err)
		return
	}
	_, err = create.Write([]byte(content))
	if err != nil {
		log.Println("write error", err)
		return
	}

	// 用后关闭
	defer create.Close()

}
func main() {
	config := loadConfig()
	log.Println(config)
	result := getFileJson(config)
	saveResult(result, config)

}
