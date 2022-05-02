package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

var apiUrl = "https://zhblogs.ohyee.cc/api/blogs?size=-1&status=0"

func main() {
	if res, err := http.Get(apiUrl); err == nil {
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(res.Body)
		var responseBody map[string]interface{}
		if err = json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
			log.Print(err)
		}
		// success
		isResValid(responseBody)
		total := responseBody["data"].(map[string]interface{})["total"].(float64)
		log.Printf("get %.0f blogs", total)
		parseData(responseBody["data"].(map[string]interface{})["blogs"].([]interface{}))
	}
}

func isResValid(res map[string]interface{}) {
	if res["success"] == false {
		log.Fatalf("%s", res["message"])
	}
	if res["data"] == nil {
		log.Fatal("res.data is nil")
	}
	if res["data"].(map[string]interface{})["blogs"] == nil || res["data"].(map[string]interface{})["total"] == nil {
		log.Fatal("res.data.blogs or res.data.total is nil")
	}
	if len(res["data"].(map[string]interface{})["blogs"].([]interface{})) == 0 {
		log.Fatal("res.data.blogs is empty")
	}
}

func parseData(blogs []interface{}) {
	if data, err := json.Marshal(blogs); err != nil {
		log.Fatalf("error when marshal blogs: %s", err)
	} else {
		// write to ./tmp/data.json
		if err = ioutil.WriteFile("./data.json", data, 0644); err != nil {
			log.Fatalf("error when write data.json: %s", err)
		}
	}
}
