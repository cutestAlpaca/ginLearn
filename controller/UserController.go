package controller

import (
	"ginLearn/common"
	"ginLearn/model"
	"ginLearn/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

func Register(c *gin.Context) {
	DB := common.GetDB()

	name := c.PostForm("name")
	telephone := c.PostForm("telephone")
	password := c.PostForm("password")

	if len(telephone) != 11 {
		c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{ //gin.H
			"code": 442,
			"msg":  "手机号错误,手机号必须为11位!",
		})
		return
	}

	if len(password) < 6 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 442,
			"msg":  "密码至少为6位!",
		})
		return
	}

	if len(name) == 0 {
		name = util.RandomString(10)
	}

	if isTelephoneExist(DB, telephone) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 442,
			"msg":  "用户已存在!",
		})
		return
	}

	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  password,
	}
	DB.Create(&newUser)

	log.Println(name, telephone, password)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "注册成功!",
	})
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)

	if user.ID != 0 {
		return true
	}

	return false
}
