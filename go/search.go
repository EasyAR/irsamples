package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
	"gopkg.in/vmihailenco/msgpack.v2"
)

const (
	host      = "localhost:8080"
	appKey    = "test_app_key"
	appSecret = "test_app_secret"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	params := map[string]string{}
	signParam(params, appKey, appSecret)

	jsonData, err := json.Marshal(params)
	check(err)

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/tunnels/", host), bytes.NewBuffer(jsonData))
	check(err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	check(err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	check(err)

	var jso map[string]interface{}
	if err = json.Unmarshal(body, &jso); err != nil {
		log.Fatal(err)
	}
	host := jso["host"]
	port := jso["port"]
	// host := "localhost"
	// port := "8080"
	tunnel := jso["result"].(map[string]interface{})["tunnel"]
	wsURL := fmt.Sprintf("ws://%s:%s/services/recognize/%s", host, port, tunnel)
	originURL := fmt.Sprintf("http://%s/", host)

	data, err := ioutil.ReadFile("../java/test_search_image.jpg")
	check(err)

	searchParams := map[string]interface{}{
		"image": data,
		"egg":   "spam",
	}
	msgData, err := msgpack.Marshal(searchParams)
	check(err)

	ws, err := websocket.Dial(wsURL, "", originURL)
	check(err)
	defer ws.Close()

	err = websocket.Message.Send(ws, msgData)
	check(err)

	var respData []byte
	err = websocket.Message.Receive(ws, &respData)
	check(err)

	fmt.Println(string(respData))
}
