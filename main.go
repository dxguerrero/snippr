package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	//"github.com/gin-gonic/gin"
)

type snippet struct {
	ID int `json:"id"`
	Language string `json:"language"`
	Code string `json:"code"`
}

var snippets []snippet  

func main() {
	readFile()
	fmt.Println(snippets)
}

func readFile() []snippet {
	file, err := os.Open("/home/danxguerrero/snippr/seedData.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(data, &snippets)
	if err != nil {
		log.Fatal(err)
	}
	
	return snippets
}