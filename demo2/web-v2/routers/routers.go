package routers

import (
	"net/http"
	"web-v2/controller"

	"github.com/gin-gonic/gin"
)

func SetupRouters() {
	r := gin.Default()

	r.GET("ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})
	userRepo := controller.New()
	r.LoadHTMLGlob("index/**/*")
	r.GET("/All", func(c *gin.Context) {
		c.HTML(http.StatusOK, "All.html", gin.H{"title": "All"})
	})

	//用户信息操作路由组
	UserGroup := r.Group("/user")
	UserGroup.Use(controller.LoginMiddleware())
	{
		//创建新内容
		UserGroup.GET("/Create", func(c *gin.Context) {
			c.HTML(http.StatusOK, "Create.html", gin.H{"title": "CreateUser"})
		})
		//根据ID查找
		UserGroup.GET("/UserById", func(c *gin.Context) {
			c.HTML(http.StatusOK, "UserId.html", gin.H{"title": "UserId"})
		})
		//根据Name查找
		UserGroup.GET("/UserByName", func(c *gin.Context) {
			c.HTML(http.StatusOK, "UserName.html", gin.H{"title": "UserName"})
		})
		//更新
		UserGroup.GET("/Update", func(c *gin.Context) {
			c.HTML(http.StatusOK, "Update.html", gin.H{"title": "Update"})
		})
		//删除
		UserGroup.GET("/Delete", func(c *gin.Context) {
			c.HTML(http.StatusOK, "Delete.html", gin.H{"title": "Delete"})
		})
		//获取数据库中所有信息
		UserGroup.GET("/AllUsers", userRepo.GetAllUsers)

		UserGroup.POST("/CreateUser", userRepo.CreateUser)
		UserGroup.POST("/GetUserById", userRepo.GetUserById)
		UserGroup.POST("/GetUserByName", userRepo.GetUserByName)
		UserGroup.POST("/UpdateUser", userRepo.UpdateUser)
		UserGroup.POST("/DeleteUsers", userRepo.DeleteUsers)

	}
	//观点信息操作路由组
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{"title": "login"})
	})
	r.POST("/loginstate", userRepo.UserLogin)
	ViewGroup := r.Group("/view")
	ViewGroup.Use(controller.UserLoginMiddleWare())
	{
		ViewGroup.POST("/login", userRepo.UserLogin)
		ViewGroup.GET("/Create", func(c *gin.Context) {
			c.HTML(http.StatusOK, "CreateView.html", gin.H{"title": "CreateView"})
		})
		ViewGroup.GET("/Delete", func(c *gin.Context) {
			c.HTML(http.StatusOK, "DeleteView.html", gin.H{"title": "DeleteView"})
		})
		ViewGroup.POST("/CreateView", userRepo.CreateViews)
		ViewGroup.GET("/GetViewsbyUserId", userRepo.GetViewsbyUserId)
		ViewGroup.POST("DeleteView", userRepo.DeleteViewsbyId)

	}

	r.Run(":8000")

}
