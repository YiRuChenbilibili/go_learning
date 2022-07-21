package controller

import (
	"errors"
	"net/http"
	"webtest/data"
	"webtest/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func New() *Repo {
	db := database.Connect2DB()
	db.AutoMigrate(&data.User{})
	return &Repo{db: db}
}

//创造users
func (repo *Repo) CreateUser(c *gin.Context) {
	var user data.User

	//将request的body中的数据，自动按照form格式解析到结构体
	if err := c.Bind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := data.CreateUser(repo.db, &user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

//获取users
func (repo *Repo) GetAllUsers(c *gin.Context) {
	var users []data.User
	if err := data.GetAllUsers(repo.db, &users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, users)
}

//通过名字查询用户信息
func (repo *Repo) GetUserById(c *gin.Context) {
	var user data.User
	id := c.PostForm("id")
	if err := data.GetUserById(repo.db, &user, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"ID is not exist! error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

//通过名字查询用户信息
func (repo *Repo) GetUserByName(c *gin.Context) {
	var users []data.User
	name := c.PostForm("name")
	if err := data.GetUserByName(repo.db, &users, name); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"Name is not exist! error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

//更新数据，取出要更新的数据并进行更新
func (repo *Repo) UpdateUser(c *gin.Context) {
	var user data.User
	//取出
	id := c.PostForm("id")
	if err := data.GetUserById(repo.db, &user, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"ID is not exist! error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//更新
	c.Bind(&user)
	if err := data.UpdateUser(repo.db, &user); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, user)
}

//删除用户信息
func (repo *Repo) DeleteUsers(c *gin.Context) {
	var user data.User
	id := c.PostForm("id")
	if err := data.DeleteUsers(repo.db, &user, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Delete successfully!"})
}
