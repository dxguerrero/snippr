package middleware

import (
	"fmt"
    //"net/http"
	"strings"

    "github.com/gin-gonic/gin"
    "github.com/dgrijalva/jwt-go"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(401, gin.H{
                "message": "Unable to get Token.",
            })
            return
        }

        tokenString := strings.Split(authHeader, " ")[1]
        fmt.Println(tokenString)

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            // Check the signing method
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return []byte("your_secret_key"), nil
        })

        if err != nil {
            c.JSON(401, gin.H{
                "message": err.Error(),
            })
            return
        }

        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            userID := claims["user_id"] 
            c.Set("user_id", userID)
            c.Next()
        } else {
            c.JSON(404, gin.H{
                "message": "Unable to get user id.",
            })
        }
	}}