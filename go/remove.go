package main

import (
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

const targetID = "00ed20c3-53ea-4cdc-a5ed-4766ce3adb60"

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	params := map[string]string{}
	signParam(params, appKey, appSecret)

	req, err := http.NewRequest("DELETE", host+"/target/"+targetID, nil)
	check(err)

	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

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
