package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"web-v2/data"

	"github.com/gin-gonic/gin"
)

//通过cookie实现用户登录
func (repo *Repo) UserLogin(c *gin.Context) {
	Id := c.PostForm("userid")
	Password := c.PostForm("password")
	count := data.TryLogin(repo.db, Id, Password)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "UserID or Password is not correct!"})
		return
	} else {
		//验证成功，设置cookie
		//cookie名,cookie值,cookie有效时长, cookie 所在的目录,所在域,是否只能通过 https 访问,是否可以通过 js代码进行操作
		c.SetCookie("abc", Id, 3600, "/", "localhost", false, true)
		c.String(http.StatusOK, "Login success!")
	}

}

//构建中间件用于用户登录
func UserLoginMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		//如果正确登录
		cookie, err := c.Cookie("abc")
		if err == nil {
			c.Set("Userid", cookie)
			//继续后续函数
			c.Next()
		} else {
			//如果发生错误
			c.JSON(http.StatusUnauthorized, gin.H{"message": "请先完成登录！", "error": err.Error()})
			//不再调用后续函数
			c.Abort()
			return
		}
	}
}

//登录后可创建评论
func (repo *Repo) CreateViews(c *gin.Context) {
	var view data.View
	c.Bind(&view)
	//中间件取值
	userid, _ := c.Get("Userid")
	// 先显式转换，.(string) 把interface转换成string类型，再利用strconv.Atoi把string 转换成int
	view.UserID, _ = strconv.Atoi(userid.(string))
	if err := data.CreateViews(repo.db, &view); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": view.View})
}

//查看登录用户的所有评论
func (repo *Repo) GetViewsbyUserId(c *gin.Context) {
	var views []data.View
	userid, _ := c.Get("Userid")
	id := userid.(string)
	if err := data.GetViewsbyUserId(repo.db, &views, id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for _, view := range views {
		c.String(http.StatusOK, fmt.Sprintf("ViewID:%d  View:%s\n", view.ID, view.View))
	}
}

//删除观点,用户只能删除属于自己的观点
func (repo *Repo) DeleteViewsbyId(c *gin.Context) {
	var view data.View
	id := c.PostForm("viewid")
	if userid, _ := c.Get("Userid"); userid != nil {
		err, count := data.DeleteViewsbyId(repo.db, &view, id, userid.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if count == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "该观点不属于您或该观点不存在！无法删除！"})
			return
		}
	}
	c.JSON(http.StatusOK, "Successfully deleted!")
}
