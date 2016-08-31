package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	host      = "http://localhost:8888"
	appKey    = "test_app_key"
	appSecret = "test_app_secret"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	data, err := ioutil.ReadFile("../java/test_target_image.jpg")
	check(err)

	imgData := base64.StdEncoding.EncodeToString(data)
	params := map[string]string{
		"image": imgData,
	}
	signParam(params, appKey, appSecret)

	jsonData, err := json.Marshal(params)
	check(err)

	req, err := http.NewRequest("POST", host+"/similar/", bytes.NewBuffer(jsonData))
	check(err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	check(err)
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, err := ioutil.ReadAll(resp.Body)
	check(err)
	fmt.Println("response Body:", string(body))
}
