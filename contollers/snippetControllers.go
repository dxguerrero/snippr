package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/dxguerrero/snippr/utils"
	"github.com/gin-gonic/gin"
)

type snippet struct {
	ID       int    `json:"id"`
	Language string `json:"language"`
	Code     string `json:"code"`
}

var snippets []snippet

func ReadFile() []snippet {
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

func GetSnippets(c *gin.Context) {
	// langauge query parameter
	language := c.Query("language")

	filteredSnippets := make([]snippet, 0)
	for _, s := range snippets {
		if s.Language == language || language == "" {
			filteredSnippets = append(filteredSnippets, s)
		}
	}

	c.IndentedJSON(http.StatusOK, filteredSnippets)
}

func PostSnippet(c *gin.Context) {
	var newSnippet snippet

	if err := c.BindJSON(&newSnippet); err != nil {
		return
	}

	encryptedCode, err := utils.GetAESEncrypted(newSnippet.Code)
	newSnippet.Code = encryptedCode
	fmt.Println(encryptedCode)

	if err != nil {
		fmt.Println("Error during encryption", err)
	}

	snippets = append(snippets, newSnippet)
	c.IndentedJSON(http.StatusCreated, newSnippet)
}

func GetSnippetByID(c *gin.Context) {
	id := c.Param("id")

	for _, s := range snippets {
		if fmt.Sprintf("%d", s.ID) == id {
			decryptedCode, err := utils.GetAESDecrypted(s.Code)
			if err != nil {
				panic(err)
			}
			s.Code = string(decryptedCode)
			c.IndentedJSON(http.StatusOK, s)
			return
		}
	}
}
