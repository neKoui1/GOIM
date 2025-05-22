package service

import (
	"GOIM/helper"
	"GOIM/models"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)


func Login(c *gin.Context) {

	account := c.PostForm("account")
	password := c.PostForm("password")
	if account == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户名或密码不能为空",
		})
		return
	}
	u,err := models.GetUserByAccount(account)
	if err != nil {
		fmt.Println("查询用户失败" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户名或密码错误",
		})
		return
	}
	if !helper.CheckPassword(password, u.Password) {
		c.JSON(http.StatusOK, gin.H {
			"code": -1,
			"msg":"密码错误",
		})
		return 
	}

	token, err := helper.GenerateToken(u.ID, u.Email)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "生成token失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "登录成功",
		"data": gin.H{
			"token": token,
		},
	})
}

func Register(c *gin.Context) {
	account := c.PostForm("account")
	password:=c.PostForm("password")
	nickname:=c.PostForm("nickname")
	email := c.PostForm("email")
	if account == "" || password == "" || 
	nickname == "" || email == "" {
		c.JSON(http.StatusOK,gin.H{
			"code":-1,
			"msg":"注册信息不完整",
		})
		return
	}
	cnt, err := models.GetUserCountByAccount(account)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "系统错误，获取用户数量失败",
		})
		return
	}
	if cnt > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "账号已注册",
		})
		return
	}


	savePwd, err := helper.HashPassword(password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "系统错误，密码加密失败",
		})
		return
	}
	u:= &models.User{
		ID: bson.NewObjectID(),
		Account: account,
		Password: savePwd,
		Nickname: nickname,
		Email: email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		LastLogin: time.Now(),
	}
	err = models.InsertOntUser(u)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "系统错误，注册失败",
		})
		return
	}

	token, err := helper.GenerateToken(u.ID, u.Email)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "系统错误，生成token失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "注册成功",
		"data": gin.H{
			"token": token,
		},
	})
}


func GetUserInfo(c *gin.Context) {
	u, _ := c.Get("user_claims")
	uc :=u.(*helper.UserClaims)
	user, err := models.GetUserByID(uc.ID)
	if err != nil {
		log.Printf("[DB ERROR] %v\n", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg": "系统错误，获取用户信息失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg": "获取用户信息成功",
		"data": user,
	})
}

func SendCode(c *gin.Context) {
	email := c.PostForm("email")
	if email == "" {
		c.JSON(http.StatusOK, gin.H {
			"code":-1,
			"msg":"邮箱不能为空",
		})
		return
	}
	cnt, err := models.GetUserCountByEmail(email)
	if err!= nil {
		log.Printf("[DB ERROR]: %v\n", err)
		return
	}
	if cnt > 0 {
		c.JSON(http.StatusOK,gin.H {
			"code":-1,
			"msg":"当前邮箱已被注册",
		})
		return
	}
	err = helper.SendCode(email, "123456")
	if err != nil {
		log.Printf("[ERROR] 发送验证码失败: %v\n", err)
		c.JSON(http.StatusOK, gin.H{
			"code":-1,
			"msg":"发送验证码失败",
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"code":200,
		"msg":"验证码发送成功",
	})
}
