package controllers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dxguerrero/snippr/utils"
	"github.com/gin-gonic/gin"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var users []User

func GetUser(c *gin.Context) {
	// username and password query parameters
	username := c.Query("username")
	password := c.Query("password")

	var filteredUser User
	for _, u := range users {
		fmt.Println(u.Username)
		if u.Username == username {
			if utils.CheckPasswordHash(u.Password, password) {
				filteredUser = u
			} else {
				c.JSON(404, gin.H{
					"message": "Passwords do not match.",
				})
			}
		} else {
			c.JSON(404, gin.H{
				"message": "No user with that username.",
			})
		}
	}
	c.IndentedJSON(http.StatusOK, filteredUser)
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

func Login(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Authorization header is missing",
		})
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Basic" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid authorization header format",
		})
		return
	}

	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Failed to decode authorization header",
		})
		return
	}

	credentials := strings.Split(string(decoded), ":")
	if len(credentials) != 2 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials format",
		})
		return
	}

	username := credentials[0]
	password := credentials[1]

	var user User
	for _, u := range users {
		if u.Username == username {
			user = u
			break
		}
	}

	if user.Username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid Credentials",
		})
		return
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid Credentials",
		})
		return	
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id" : user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte("your_secret_key"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}
