package test

import (
	"log"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// var addr = flag.String("addr", "localhost:8080", "http service address")
var upgrader = websocket.Upgrader{}
var ws = make(map[*websocket.Conn]struct{})

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil) 
	if err != nil {
		log.Println("upgrade: ", err)
		return
	}
	defer c.Close()
	ws[c] = struct{}{}

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read: ", err)
			break
		}
		log.Printf("recv: %s\n", message)
		for conn := range ws {
			err = conn.WriteMessage(mt, message) 
			if err != nil {
				log.Println("write: ", err)
				break
			}
		}
	}
}
// func TestWebsocketServer(t *testing.T) {
// 	http.HandleFunc("/echo", echo)
// 	log.Fatal(http.ListenAndServe(*addr, nil))
// }

func TestGinWebsocketServer(t *testing.T) {
	r := gin.Default()
	r.GET("/echo", func (c *gin.Context)  {
		echo(c.Writer, c.Request)
	})
	r.Run(":8080")

}