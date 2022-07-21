package main

import (
	"net/http"
	"webtest/controller"

	"github.com/gin-gonic/gin"
)

func main() {
	r := setupRouter()
	_ = r.Run(":8000")
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})

	userRepo := controller.New()
	r.LoadHTMLGlob("index/*")

	//创建新内容
	r.GET("/Create", func(c *gin.Context) {
		c.HTML(http.StatusOK, "Create.html", gin.H{"title": "CreateUser"})
	})
	r.POST("/CreateUser", userRepo.CreateUser)

	//获取数据库中所有信
	r.GET("/AllUsers", userRepo.GetAllUsers)
	r.GET("/UserById", func(c *gin.Context) {
		c.HTML(http.StatusOK, "UserId.html", gin.H{"title": "UserId"})
	})

	//根据ID查找
	r.POST("/GetUserById", userRepo.GetUserById)

	//根据Name查找
	r.GET("/UserByName", func(c *gin.Context) {
		c.HTML(http.StatusOK, "UserName.html", gin.H{"title": "UserName"})
	})
	r.POST("/GetUserByName", userRepo.GetUserByName)

	//更新
	r.GET("/Update", func(c *gin.Context) {
		c.HTML(http.StatusOK, "Update.html", gin.H{"title": "Update"})
	})
	r.POST("/UpdateUser", userRepo.UpdateUser)

	//删除
	r.GET("/Delete", func(c *gin.Context) {
		c.HTML(http.StatusOK, "Delete.html", gin.H{"title": "Delete"})
	})
	r.POST("/DeleteUsers", userRepo.DeleteUsers)

	return r
}
