package service

import (
	"GOIM/models"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func ExportUserExcel(c *gin.Context) {

	userList, err := models.GetUserList()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "获取所有用户信息失败" + err.Error(),
		})
		return
	}

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Println("excel file format fail: ", err)
			return
		}
	}()

	streamWriter, err := f.NewStreamWriter("Sheet1")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "创建Excel流失败: " + err.Error(),
		})
		return
	}
	headers := []any{"Account", "Nickname", "Gender"}
	if err := streamWriter.SetRow("A1", headers); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "写入表头失败: " + err.Error(),
		})
		return
	}
	for i, v := range userList {
		row := []any{v.Account, v.Nickname, v.Gender}
		cell, _ := excelize.CoordinatesToCellName(1, i+2)
		if err := streamWriter.SetRow(cell, row); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  fmt.Sprintf("写入第%d行失败：%s\n", i+2, err.Error()),
			})
			return
		}
	}
	if err := streamWriter.Flush(); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "刷新excel流失败: " + err.Error(),
		})
		return
	}
	buf, err := f.WriteToBuffer()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "fail to generate excel: " + err.Error(),
		})
		return
	}
	fileName := "userlist_" + time.Now().Format("20060102_150405") + ".xlsx"

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	encodedFileName := url.QueryEscape(fileName)
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"; filename*=UTF-8''%s`, encodedFileName, encodedFileName))

	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}
