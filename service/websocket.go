package service

import (
	"GOIM/helper"
	"GOIM/models"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var upgrader = websocket.Upgrader{}

// 连接写入锁
var connLocks = make(map[bson.ObjectID]*sync.Mutex)
var connLocksMutex sync.Mutex

// room id -> user id -> conn
var roomConnections = make(map[bson.ObjectID]map[bson.ObjectID]*websocket.Conn)
var roomConnLock sync.RWMutex

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

	uc := c.MustGet("user_claims").(*helper.UserClaims)

	// 获取用户所在的房间列表
	userRooms, err := models.GetUserRoomByUserID(uc.ID)
	if err != nil {
		log.Printf("Failed to get user rooms : %v\n", err)
		conn.Close()
		return
	}

	// 加入所有房间
	for _, userRoom := range userRooms {
		joinRoom(uc.ID, userRoom.RoomId, conn)
	}

	// 断开时清理
	defer func() {
		for _, userRoom := range userRooms {
			leaveRoom(uc.ID, userRoom.RoomId)
		}

		// 清理连接锁
		connLocksMutex.Lock()
		if _, exists := connLocks[uc.ID]; exists {
			delete(connLocks, uc.ID)
			log.Printf("Remove user %v \n", uc.ID)
		}
		connLocksMutex.Unlock()
		conn.Close()
	}()

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
			continue
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
			continue
		}

		// 广播消息到房间
		broadcastToRoom(ms.RoomId, []byte(ms.Message), uc.ID)
	}

}

// 用户连接时加入房间
func joinRoom(userID, roomID bson.ObjectID, conn *websocket.Conn) {
	roomConnLock.Lock()
	defer roomConnLock.Unlock()

	if roomConnections[roomID] == nil {
		roomConnections[roomID] = make(map[bson.ObjectID]*websocket.Conn)

	}
	roomConnections[roomID][userID] = conn
}

// 用户断开时离开房间
func leaveRoom(userID, roomID bson.ObjectID) {
	roomConnLock.Lock()
	defer roomConnLock.Unlock()

	if roomConns, exists := roomConnections[roomID]; exists {
		delete(roomConns, userID)
		if len(roomConns) == 0 {
			delete(roomConnections, roomID)
		}
	}
}

// 向房间广播信息
func broadcastToRoom(roomID bson.ObjectID, message []byte, senderId bson.ObjectID) {
	roomConnLock.Lock()
	defer roomConnLock.Unlock()

	if roomConns, exists := roomConnections[roomID]; exists {
		for userId, conn := range roomConns {
			if senderId == userId {
				continue
			}
			connLocksMutex.Lock()
			lock, exists := connLocks[userId]
			if !exists {
				lock = &sync.Mutex{}
				connLocks[userId] = lock
			}
			connLocksMutex.Unlock()

			lock.Lock()
			err := conn.WriteMessage(websocket.TextMessage, message)
			lock.Unlock()
			if err != nil {
				log.Printf("Failed to send message to user: %v - %v", userId, err)
			}
		}
	}
}
