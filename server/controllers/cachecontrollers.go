package controllers

import (
	"lrucachesystem/handlers"

	"github.com/gin-gonic/gin"
)

type SetRequestBody struct {
	Key   string `json:"key"`
	Value string `json:"value"`
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
	c.Bind(&request)
	handlers.HandleSetcache(request.Key, request.Value)
	c.JSON(200, gin.H{"STATUS": "Cache Entered"})

}
