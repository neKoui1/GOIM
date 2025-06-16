package service

import (
	"GOIM/helper"
	"GOIM/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var upgrader = websocket.Upgrader{}
var wc = make(map[bson.ObjectID]*websocket.Conn)

type MessageStruct struct {
	Message string        `json:"message"`
	RoomId  bson.ObjectID `json:"room_id"`
}

func WebSocketMessage(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "系统异常" + err.Error(),
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

		// 判断用户是否属于消息体的房间
		_, err := models.GetUserRoomByUserIDRoomID(uc.ID, ms.RoomId)
		if err != nil {
			log.Printf("User id: %v, Room id: %v NOT EXISTS, err: %v\n",
				uc.ID, ms.RoomId, err)
			return
		}

		// 保存消息
		msg := &models.Message{
			UserId: uc.ID,
			RoomId: ms.RoomId,
			Data:   ms.Message,
		}
		err = models.InsertOntMessage(msg)
		if err != nil {
			log.Printf("[DB ERROR]: %v\n", err)
			return
		}

		// 获取在特定房间的在线用户
		userRooms, err := models.GetUserRoomByRoomID(ms.RoomId)
		if err != nil {
			log.Printf("[DB ERROR]: %v\n", err)
			return
		}
		for _, room := range userRooms {
			if cc, ok := wc[room.UserId]; ok {
				err = cc.WriteMessage(websocket.TextMessage, []byte(ms.Message))
				if err != nil {
					log.Printf("Write Message Error: %v\n", err)
					return
				}
			}
		}

	}

}
