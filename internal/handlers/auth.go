package handlers

import (
	"log"
	"sica/internal/repositories"
	"sica/pkg/bcrypt"
	"sica/pkg/jwt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Credentials struct {
	Password string `json:"password"`
}

type Token struct {
	Refresh string `json:"refresh"`
}

type NewCredentials struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func Login(c *gin.Context) {

	var credentials Credentials
	err := c.BindJSON(&credentials)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{
			"error": "invalid json format",
		})
		return
	}

	r := repositories.NewAuthRepository()
	user, err := r.Get(1)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": "error getting user",
		})
		return
	}
	err = bcrypt.CheckPassword(user.Password, credentials.Password)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{
			"error": "invalid credentials",
		})
		return
	}

	access, err := jwt.CreateToken(strconv.Itoa(1), "access", 5)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": "error creating access token",
		})
		return
	}

	refresh, err := jwt.CreateToken(strconv.Itoa(1), "refresh", 1440)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": "error creating refresh token",
		})
		return
	}

	c.JSON(200, gin.H{
		"access":  access,
		"refresh": refresh,
	})
}

func Refresh(c *gin.Context) {

	var token Token
	err := c.BindJSON(&token)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{
			"error": "invalid json format",
		})
		return
	}

	claims, err := jwt.ValidateToken(token.Refresh)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{
			"error": "invalid token",
		})
		return
	}

	if claims.Type != "refresh" {
		log.Println("not a refresh token")
		c.JSON(400, gin.H{
			"error": "invalid token",
		})
		return
	}

	access, err := jwt.CreateToken(strconv.Itoa(1), "access", 5)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": "error creating access token",
		})
		return
	}

	c.JSON(200, gin.H{
		"access":  access,
		"refresh": token.Refresh,
	})

}

func ChangePassword(c *gin.Context) {

	var newCredentials NewCredentials
	err := c.BindJSON(&newCredentials)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{
			"error": "invalid json format",
		})
		return
	}

	r := repositories.NewAuthRepository()
	user, err := r.Get(1)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": "error getting user",
		})
		return
	}
	err = bcrypt.CheckPassword(user.Password, newCredentials.OldPassword)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{
			"error": "invalid credentials",
		})
		return
	}

	newPassword, err := bcrypt.HashPassword(newCredentials.NewPassword)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": "error hashing new password",
		})
		return
	}

	newValues := map[string]interface{}{
		"password": newPassword,
	}

	user, err = r.Update(1, newValues)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": "error updating user password",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "password changed",
	})
}