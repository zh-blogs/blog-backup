package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var apiUrl = "https://zhblogs.ohyee.cc/api/blogs?size=-1&status=0&hash=1tYgiQWZKDjIihy"
var maxNum = 20
var dataDir = "./database"

func main() {
	log.Print("getting blogs...")
	if res, err := http.Get(apiUrl); err == nil {
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(res.Body)
		var responseBody map[string]interface{}
		if err = json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
			log.Fatal(err)
		}
		// success
		isResValid(responseBody)
		total := responseBody["data"].(map[string]interface{})["total"].(float64)
		log.Printf("get %.0f blogs", total)
		parseData(responseBody["data"].(map[string]interface{})["blogs"].([]interface{}))
	} else {
		log.Fatalf("httpGet Err %s", err)
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
		//check if database dir exists
		makeDataDir()
		// if more than %maxNum% files, delete the oldest one
		deleteOldestFile()
		// write to ./tmp/data.json
		writeToFile(data)
	}
}

func writeToFile(data []byte) {
	timeNow := time.Now().Format("01-02-15-04")
	path := dataDir + "/data-" + timeNow + ".json"

	if err := ioutil.WriteFile(path, data, 0644); err == nil {
		log.Printf("write to %s", path)
	} else {
		log.Fatalf("error when write to %s: %s", path, err)
	}
}

func deleteOldestFile() {
	if list, err := ioutil.ReadDir("database"); err == nil {
		for _, file := range list {
			if getDataFileNum() < maxNum {
				break
			}
			if err := os.Remove(dataDir + "/" + file.Name()); err == nil {
				log.Printf("delete %s", file.Name())
			} else {
				log.Printf("error when delete file: %s", err)
			}
		}
	} else {
		log.Fatalf("error when read database dir: %s", err)
	}
}

func getDataFileNum() int {
	if list, err := ioutil.ReadDir(dataDir); err == nil {
		return len(list)
	} else {
		log.Printf("error when read database dir: %s", err)
	}
	return 0
}

func makeDataDir() {
	if _, err := ioutil.ReadDir(dataDir); err != nil {
		if err := os.Mkdir(dataDir, os.ModePerm); err != nil {
			log.Fatalf("error when create database dir: %s", err)
		}
	}
}
