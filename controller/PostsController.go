package controller

import (
	"ginLearn/common"
	"ginLearn/model"
	"ginLearn/response"
	"ginLearn/vo"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strconv"
)

type IPostsController interface {
	RestController
	PageList(ctx *gin.Context)
}

type PostsController struct {
	DB *gorm.DB
}

func (p PostsController) Create(ctx *gin.Context) {
	var requestPost vo.CreatePostRequst

	if err := ctx.ShouldBind(&requestPost); err != nil {
		response.Fail(ctx, nil, "数据验证错误!")
		return
	}

	user, _ := ctx.Get("user")

	post := model.Post{
		UserId:     user.(model.User).ID,
		CategoryId: requestPost.CategoryId,
		Title:      requestPost.Title,
		HeadImg:    requestPost.HeadImg,
		Content:    requestPost.Content,
	}

	if err := p.DB.Create(&post).Error; err != nil {
		panic(err)
		return
	}

	response.Success(ctx, nil, "创建成功!")
}

func (p PostsController) Update(ctx *gin.Context) {
	var requestPost vo.CreatePostRequst

	if err := ctx.ShouldBind(&requestPost); err != nil {
		response.Fail(ctx, nil, "数据验证错误!")
		return
	}

	postId := ctx.Params.ByName("id")

	var post model.Post
	if p.DB.Where("id = ?", postId).First(&post).RecordNotFound() {
		response.Fail(ctx, nil, "文章不存在!")
		return
	}

	user, _ := ctx.Get("user")
	userId := user.(model.User).ID
	if userId != post.UserId {
		response.Fail(ctx, nil, "文章不属于您, 请勿非法操作!")
		return
	}

	if err := p.DB.Model(&post).Update(requestPost).Error; err != nil {
		response.Fail(ctx, nil, "更新失败!")
		return
	}

	response.Success(ctx, gin.H{"post": post}, "更新成功!")
}

func (p PostsController) Show(ctx *gin.Context) {
	postId := ctx.Params.ByName("id")

	var post model.Post
	if p.DB.Preload("Category").Where("id = ?", postId).First(&post).RecordNotFound() {
		response.Fail(ctx, nil, "文章不存在!")
		return
	}

	response.Success(ctx, gin.H{"post": post}, "成功!")
}

func (p PostsController) Delete(ctx *gin.Context) {
	postId := ctx.Params.ByName("id")

	var post model.Post
	if p.DB.Where("id = ?", postId).First(&post).RecordNotFound() {
		response.Fail(ctx, nil, "文章不存在!")
		return
	}

	user, _ := ctx.Get("user")
	userId := user.(model.User).ID
	if userId != post.UserId {
		response.Fail(ctx, nil, "文章不属于您, 请勿非法操作!")
		return
	}

	p.DB.Delete(&post)

	response.Success(ctx, gin.H{"post": post}, "删除成功!")
}

func (p PostsController) PageList(ctx *gin.Context) {
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	var posts []model.Post
	p.DB.Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&posts)

	var total int
	p.DB.Model(model.Post{}).Count(&total)

	response.Success(ctx, gin.H{"data": posts, "total": total}, "成功!")
}

func NewPostsController() IPostsController {
	db := common.GetDB()
	db.AutoMigrate(model.Post{})

	return PostsController{DB: db}
}
