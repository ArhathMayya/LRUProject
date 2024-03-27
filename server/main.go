package main

import (
	"lrucachesystem/controllers"
	"lrucachesystem/handlers"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

//I'm using gin library instead of standard http library since I'm more confortable with gin.

func main() {
	r := gin.Default()
	r.Use(cors.Default())
	//cors will allow the request from the react fron end.

	go handlers.CleanupExpiredEntries()
	//cleaning of cache will run in separate thread.

	r.GET("/get/:key", controllers.Getcache)
	r.POST("/set", controllers.Setcache)

	r.Run(":8000") //Server will start running at port 8000

}
