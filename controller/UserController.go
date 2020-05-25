package controller

import (
	"ginLearn/common"
	"ginLearn/dto"
	"ginLearn/model"
	"ginLearn/response"
	"ginLearn/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func Register(c *gin.Context) {
	DB := common.GetDB()

	name := c.PostForm("name")
	telephone := c.PostForm("telephone")
	password := c.PostForm("password")

	if len(telephone) != 11 {
		response.Response(c, http.StatusUnprocessableEntity, 442, nil, "手机号错误,手机号必须为11位!")
		return
	}

	if len(password) < 6 {
		response.Response(c, http.StatusUnprocessableEntity, 442, nil, "密码至少为6位!")
		return
	}

	if len(name) == 0 {
		name = util.RandomString(10)
	}

	if isTelephoneExist(DB, telephone) {
		response.Response(c, http.StatusUnprocessableEntity, 442, nil, "用户已存在!")
		return
	}

	//
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "加密错误!")
		return
	}

	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hashedPassword),
	}
	DB.Create(&newUser)

	log.Println(name, telephone, password)

	response.Success(c, nil, "注册成功!")
}

func Login(c *gin.Context) {
	DB := common.GetDB()

	telephone := c.PostForm("telephone")
	password := c.PostForm("password")
	if len(telephone) != 11 {
		response.Response(c, http.StatusUnprocessableEntity, 442, nil, "手机号错误,手机号必须为11位!")
		return
	}

	if len(password) < 6 {
		response.Response(c, http.StatusUnprocessableEntity, 442, nil, "密码至少为6位!")
		return
	}

	var user model.User
	DB.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 {
		response.Response(c, http.StatusUnprocessableEntity, 442, nil, "用户不存在!")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, "密码错误!")
		return
	}

	token, err := common.ReleaseToken(user)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "系统异常!")
		log.Printf("token generate error: %v", err)
		return
	}

	response.Success(c, gin.H{"token": token}, "登陆成功!")
}

func Info(c *gin.Context) {
	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{"user": dto.ToUserDto(user.(model.User))},
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
