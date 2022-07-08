# gin框架 #
gin是一个轻量级的 WEB 框架，支持 RestFull 风格 API，支持 GET，POST，PUT，PATCH，DELETE，OPTIONS 等 http 方法，支持文件上传，分组路由，Multipart/Urlencoded FORM，以及支持 JsonP，参数处理等等功能。  

## 简单路由使用 ##
基本路由 gin 框架中采用的路由库是 httprouter。
**gin路由**：   
**api 参数**：通过Context的Param方法来获取。     
```
package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	// 创建带有默认中间件的路由:
	router := gin.Default()
  //创建不带中间件的路由：
	//router := gin.New()
	//定义路径，使用冒号:代替变量(name,age为变量)
	router.GET("/user/:name/:age", func(context *gin.Context) {
		//获取变量值
		name := context.Param("name")
		age := context.Param("age")
		//截取
		age = strings.Trim(age, "/")
		message := name + " is " + age
		//返回值
		context.String(http.StatusOK, "hello %s", message)
	})
   // 3.监听端口，默认在8080
   // Run("里面不指定端口号默认为8080")
	router.Run()
  //router.Run(":8000")z指定端口
}
```
访问http://127.0.0.1:8080/user/name1/age1：
![image](https://user-images.githubusercontent.com/24589721/177819868-0d9a0b7a-bc2a-4b60-a50b-413057c6f52a.png)
**URL 参数**：通过 DefaultQuery 或 Query 方法获取  
DefaultQuery()若参数存在，返回默认值，Query()若不存在，返回空串。(API ? name=zs)
```
// url 为 http://localhost:8080/welcome?name=ningskyer时
// 输出 Hello ningskyer
// url 为 http://localhost:8080/welcome时
// 输出 Hello Guest
router.GET("/welcome", func(c *gin.Context) {
	name := c.DefaultQuery("name", "Guest") //可设置默认值
	// 是 c.Request.URL.Query().Get("lastname") 的简写
	//name := c.Query("name") //无默认值
	c.String(http.StatusOK, fmt.Sprintf("hello %s", name))
})
```
带参数  
![image](https://user-images.githubusercontent.com/24589721/177899945-cadc669b-d9ad-4b9e-8be9-ada35ba2cde0.png)  
不带参数
![image](https://user-images.githubusercontent.com/24589721/177900016-f791b92f-0dd1-4185-aeff-611514665cf5.png)  
**表单参数**：通过PostForm()方法获取，该方法默认解析的是x-www-form-urlencoded或from-data格式的参数
```
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Document</title>
</head>
<body>
    <form action="http://localhost:8080/form" method="post" action="application/x-www-form-urlencoded">
        用户名：<input type="text" name="username" placeholder="请输入你的用户名">  <br>
        密&nbsp;&nbsp;&nbsp;码：<input type="password" name="userpassword" placeholder="请输入你的密码">  <br>
        <input type="submit" value="提交">
    </form>
</body>
</html>
```
```
func main() {
    r := gin.Default()
    r.POST("/form", func(c *gin.Context) {
    	//使用DefaultPostForm为参数设置一个默认值，当前端没有传参时直接默认值赋值给相应的参数
        types := c.DefaultPostForm("type", "post")
	//解析x-www-form-urlencoded参数
        username := c.PostForm("username")
        password := c.PostForm("userpassword")
        // c.String(http.StatusOK, fmt.Sprintf("username:%s,password:%s,type:%s", username, password, types))
        c.String(http.StatusOK, fmt.Sprintf("username:%s,password:%s,type:%s", username, password, types))
    })
    r.Run()
}
```
![image](https://user-images.githubusercontent.com/24589721/177903157-c1443c8a-7b46-49b5-88b8-82ca2fa1a3e7.png)
![image](https://user-images.githubusercontent.com/24589721/177903193-4e7b6a7a-b41c-4873-affc-778825ce1a85.png)  
上传文件(from-data)


