package main

import (
	"ginLearn/controller"
	"ginLearn/middleware"
	"github.com/gin-gonic/gin"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.Use(middleware.CORSMiddleware(), middleware.RecoveryMiddleware())
	r.POST("api/auth/register", controller.Register)
	r.POST("api/auth/login", controller.Login)
	r.GET("api/auth/info", middleware.AuthMiddleware(), controller.Info)

	categoryRoutes := r.Group("/categories")
	categoryController := controller.NewCategoryController()
	categoryRoutes.POST("", categoryController.Create)
	categoryRoutes.PUT("/:id", categoryController.Update)
	categoryRoutes.GET("/:id", categoryController.Show)
	categoryRoutes.DELETE("/:id", categoryController.Delete)

	postsRoutes := r.Group("/posts")
	postsRoutes.Use(middleware.AuthMiddleware())
	postsController := controller.NewPostsController()
	postsRoutes.POST("", postsController.Create)
	postsRoutes.PUT("/:id", postsController.Update)
	postsRoutes.GET("/:id", postsController.Show)
	postsRoutes.DELETE("/:id", postsController.Delete)
	postsRoutes.POST("page/list", postsController.PageList)

	return r
}
