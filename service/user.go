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
	u, err := models.GetUserByAccount(account)
	if err != nil {
		fmt.Println("查询用户失败" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户名或密码错误",
		})
		return
	}
	if !helper.CheckPassword(password, u.Password) {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "密码错误",
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

	u.SetLastLoginNow()
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
	password := c.PostForm("password")
	nickname := c.PostForm("nickname")
	email := c.PostForm("email")
	if account == "" || password == "" ||
		nickname == "" || email == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "注册信息不完整",
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
	u := &models.User{
		ID:        bson.NewObjectID(),
		Account:   account,
		Password:  savePwd,
		Nickname:  nickname,
		Email:     email,
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
	uc := u.(*helper.UserClaims)
	user, err := models.GetUserByID(uc.ID)
	if err != nil {
		log.Printf("[DB ERROR] %v\n", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "系统错误，获取用户信息失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取用户信息成功",
		"data": user,
	})
}

func SendCode(c *gin.Context) {
	email := c.PostForm("email")
	if email == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "邮箱不能为空",
		})
		return
	}
	cnt, err := models.GetUserCountByEmail(email)
	if err != nil {
		log.Printf("[DB ERROR]: %v\n", err)
		return
	}
	if cnt > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "当前邮箱已被注册",
		})
		return
	}
	code := helper.GetCode()
	err = helper.SendCode(email, code)
	if err != nil {
		log.Printf("[ERROR] 发送验证码失败: %v\n", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "发送验证码失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "验证码发送成功",
	})
}

type UserParamResult struct {
	Account  string `bson:"account" json:"account"`
	Nickname string `bson:"nickname" json:"nickname"`
	IsFriend bool   `bson:"is_friend" json:"is_friend"`
	Gender   bool   `bson:"gender" json:"gender"` // false 女 true 男
	Email    string `bson:"email" json:"email"`
	Avatar   string `bson:"avatar" json:"avatar"`
	Status   int    `bson:"status" json:"status"` // 0 离线 1 在线
}

func UserParam(c *gin.Context) {
	account := c.Param("account")
	if account == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "account 参数不正确",
		})
		return
	}
	u, err := models.GetUserByAccount(account)
	if err != nil {
		log.Println("[DB ERROR]: " + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "查询用户数据异常",
		})
		return
	}
	data := UserParamResult{
		Account:  u.Account,
		Nickname: u.Nickname,
		Gender:   u.Gender,
		Email:    u.Email,
		Avatar:   u.Avatar,
		Status:   u.Status,
		IsFriend: false,
	}
	uc := c.MustGet("user_claims").(*helper.UserClaims)
	flag, err := models.JudgeUserIsFriend(u.ID, uc.ID)
	if err != nil {
		log.Println("[DB ERROR]: " + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "查询用户数据异常",
		})
	}
	if flag {
		data.IsFriend = true
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "数据加载成功",
		"data": data,
	})
}

func UserAdd(c *gin.Context) {
	account := c.PostForm("account")
	if account == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不正确",
		})
		return
	}
	u, err := models.GetUserByAccount(account)
	if err != nil {
		log.Printf("[DB ERROR]: %v\n", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "数据查询异常",
		})
		return
	}
	uc := c.MustGet("user_claims").(*helper.UserClaims)
	flag, err := models.JudgeUserIsFriend(uc.ID, u.ID)
	if err != nil {
		log.Printf("[DB ERROR]: %v\n", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "数据查询异常",
		})
		return
	}
	if flag {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "已经是好友，不可重复添加",
		})
		return
	}

	// 保存房间记录
	r := &models.Room{
		Id:        bson.NewObjectID(),
		UserId:    uc.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = models.InsertOntRoom(r)
	if err != nil {
		log.Printf("[DB ERROR]: %v\n", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "数据查询异常",
		})
		return
	}
	ur := &models.UserRoom{
		UserId:    uc.ID,
		RoomId:    r.Id,
		RoomType:  1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = models.InsertOneUserRoom(ur)
	if err != nil {
		log.Printf("[DB ERROR]: %v\n", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "数据查询异常",
		})
		return
	}

	ur = &models.UserRoom{
		UserId:    u.ID,
		RoomId:    r.Id,
		RoomType:  1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = models.InsertOneUserRoom(ur)
	if err != nil {
		log.Printf("[DB ERROR]: %v\n", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "数据查询异常",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "添加成功",
	})
}

func UserDelete(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不正确",
		})
		return
	}

	// 获取房间id

}
