package service

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func ExportUserExcel(c *gin.Context) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Println("excel file format fail: ", err)
			return
		}
	}()
}
