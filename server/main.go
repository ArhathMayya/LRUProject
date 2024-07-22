package main

import (
	"fmt"
	"lrucachesystem/controllers"
	"lrucachesystem/handlers"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)

func handleConnections(c *gin.Context) {
	w := c.Writer
	r := c.Request
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ws.Close()

	clients[ws] = true

	for {
		var msg map[string]string
		err := ws.ReadJSON(&msg)
		if err != nil {
			fmt.Println(err)
			delete(clients, ws)
			break
		}
	}
}

func handleMessages() {
	for {
		msg := <-handlers.Broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				fmt.Println(err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	// cors will allow the request from the React frontend.

	// config := cors.DefaultConfig()
	// config.AllowOrigins = []string{"http://localhost:3000"}
	// r.Use(cors.New(config))
	go handlers.CleanupExpiredEntries()
	// cleaning of cache will run in a separate thread.
	go handlers.SendUpdate()
	// We will be sending realtime cache content update.

	r.GET("/get/:key", controllers.Getcache)
	r.POST("/set", controllers.Setcache)

	// WebSocket endpoint
	r.GET("/ws", func(c *gin.Context) {
		handleConnections(c)
	})

	go handleMessages() // Not required as of now

	r.Run(":8000") // Server will start running at port 8000
}
