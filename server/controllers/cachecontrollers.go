package controllers

import (
	"fmt"
	"lrucachesystem/handlers"

	"github.com/gin-gonic/gin"
)

type SetRequestBody struct {
	Key        string `json:"key"`
	Value      string `json:"value"`
	Expiration int    `json:"expiration"` // Corrected the struct tag
}

func Getcache(c *gin.Context) {
	key := c.Param("key")
	data, flag := handlers.HandleGetcache(key)

	switch flag {
	case true:
		c.JSON(200, gin.H{
			"STATUS": "DATA FOUND",
			"DATA":   data,
		})
	case false:
		c.JSON(200, gin.H{
			"STATUS": "DATA NOT FOUND",
			"DATA":   "",
		})
	}
}

func Setcache(c *gin.Context) {
	var request SetRequestBody
	if err := c.Bind(&request); err != nil {
		c.JSON(400, gin.H{"STATUS": "Invalid request"})
		return
	}
	fmt.Println(request)
	handlers.HandleSetcache(request.Key, request.Value, request.Expiration)
	c.JSON(200, gin.H{"STATUS": "Cache Entered"})
}
