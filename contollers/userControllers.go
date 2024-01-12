package controllers

import (
	"net/http"

	"github.com/dxguerrero/snippr/utils"
	"github.com/gin-gonic/gin"
)

type User struct {
	ID int `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var users []User

func GetUsers(c *gin.Context) {
	// langauge query parameter
	username := c.Query("username")

	filteredUsers := make([]User, 0)
	for _, u := range users {
		if u.Username == username || username == "" {
			filteredUsers = append(filteredUsers, u)
		}
	}

	c.IndentedJSON(http.StatusOK, filteredUsers)
}

func PostUser(c *gin.Context) {
	var newUser User

	if err := c.BindJSON(&newUser); err != nil {
		return
	}

	hashedPassword, err := utils.HashPassword(newUser.Password)
	if err != nil {
		panic(err)
	}

	newUser.Password = hashedPassword

	users = append(users, newUser)
	c.IndentedJSON(http.StatusCreated, newUser)
}