package service

import (
	"GOIM/helper"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var upgrader = websocket.Upgrader{}
var wc = make(map[bson.ObjectID]*websocket.Conn)

type MessageStruct struct {
	Message string `json:"message"`
	RoomId bson.ObjectID `json:"room_id"`
}

func WebSocketMessage(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg":"系统异常"+err.Error(),
		})
		return
	}
	defer conn.Close()

	uc := c.MustGet("user_claims").(*helper.UserClaims)
	wc[uc.ID] = conn

	for {
		ms := new(MessageStruct)
		err = conn.ReadJSON(ms)
		if err != nil {
			log.Printf("Read Error: %v\n", err)
			return
		}

		for _, c := range wc {
			err = c.WriteMessage(websocket.TextMessage, []byte(ms.Message))
			if err != nil {
				log.Printf("Write Message Error: %v\n", err)
				return
			}
		}

	}

}