package service

import (
	"GOIM/helper"
	"GOIM/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func ChatList(c *gin.Context) {
	roomIDHex := c.Query("room_id")
	if roomIDHex == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "房间号不能为空",
		})
		return
	}
	// 将字符串roomIDHex转换为ObjectID
	roomID, err := bson.ObjectIDFromHex(roomIDHex)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "fail to parse room object id",
		})
		return
	}
	// 判断用户是否属于该房间
	uc := c.MustGet("user_claims").(*helper.UserClaims)
	_, err = models.GetUserRoomByUserIDRoomID(uc.ID, roomID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "fail to get the correct user room",
		})
		return
	}

	pageIndex, err := strconv.ParseInt(c.Query("page_index"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "fail to parse page index as int",
		})
		return
	}
	pageSize, err := strconv.ParseInt(c.Query("page_size"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "fail to parse page size as int",
		})
		return
	}
	skip := (pageIndex - 1) * pageSize
	// 在数据库中查找聊天记录
	msgData, err := models.GetMessageListByRoomID(roomID, &pageSize, &skip)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "系统异常: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "数据加载成功",
		"data": msgData,
	})
}
