# gin框架 #
gin是一个轻量级的 WEB 框架，支持 RestFull 风格 API，支持 GET，POST，PUT，PATCH，DELETE，OPTIONS 等 http 方法，支持文件上传，分组路由，Multipart/Urlencoded FORM，以及支持 JsonP，参数处理等等功能。  

## gin路由 ##
基本路由 gin 框架中采用的路由库是 httprouter。 
### api 参数 ###
通过Context的Param方法来获取。     
```
package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	// 创建带有默认中间件的路由:
	// 默认使用了2个中间件Logger(), Recovery()
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
### URL 参数 ###
通过 DefaultQuery 或 Query 方法获取  
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
### 上传文件(from-data) ###  
multipart/form-data格式用于文件上传  
gin文件上传与原生的net/http方法类似，不同在于*gin把原生的request封装到c.Request中*
$\color{#FF0000}{单个文件:}$ 
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
    <form action="http://localhost:8080/upload" method="post" enctype="multipart/form-data">
          上传文件:<input type="file" name="file" >
          <input type="submit" value="提交">
    </form>
</body>
</html>
```
```
	r := gin.Default()
	//限制表单上传大小 8MB，默认为32MB
	r.MaxMultipartMemory = 8 << 20 
	r.POST("/upload", func(c *gin.Context) {
		//获取上传文件
		file, err := c.FormFile("file")
		if err != nil {
			c.String(500, "上传出错！")
		}
		//保存到服务器指定位置
		pre := "D:/golang/upload/"
		c.SaveUploadedFile(file, pre+file.Filename)
		c.String(http.StatusOK, fmt.Sprintf("%s 上传成功！", file.Filename))
	})
	r.Run()
```

![image](https://user-images.githubusercontent.com/24589721/177907891-a4b88fab-878a-40ed-ae4e-850dfa51db36.png)  
![image](https://user-images.githubusercontent.com/24589721/177907920-a0d03843-a0f3-4a7c-bb6b-2f655e26b4a0.png)  
$\color{#FF0000}{多个文件:}$ 
```
	r.POST("/upload", func(c *gin.Context) {
		//获取复合型表单
		form, err := c.MultipartForm()
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get err %s", err.Error()))
		}
		//获取所有文件
		files := form.File["files"]
		pre := "D:/golang/upload/"
		//遍历所有文件并保存到服务器指定位置
		for _, file := range files {
			if err := c.SaveUploadedFile(file, pre+file.Filename); err != nil {
				c.String(400, fmt.Sprintf("文件上传失败！"))
			}
		}
		c.String(http.StatusOK, fmt.Sprintf("%d 个文件上传成功！", len(files)))
	})
```
![image](https://user-images.githubusercontent.com/24589721/177910308-31020b87-5e79-4098-99ec-c0a3530d016f.png)
![image](https://user-images.githubusercontent.com/24589721/177910327-04c1a727-dbf5-43ea-9660-18cbc864ddb3.png)
### routes group ###
通过routes group管理相同的url： 
```
func main() {
   r := gin.Default()
   // 路由组1 ，处理GET请求
   v1 := r.Group("/v1")
   // {} 是书写规范
   {
      v1.GET("/login", login)
      v1.GET("submit", submit)
   }
   //路由组2，处理post
   v2 := r.Group("/v2")
   {
      v2.POST("/login", login)
      v2.POST("/submit", submit)
   }
   r.Run(":8000")
}

func login(c *gin.Context) {
    //add
}

func submit(c *gin.Context) {
    //add
}
```
 ### 路由拆分与注册 ###
当项目的规模增大后就不太适合继续在项目的main.go文件中去实现路由注册相关逻辑了，可以把路由部分的代码都拆分出来，形成一个或多个单独的文件或包。
**路由拆分到不同的APP**
当项目规模实在太大时，就会更倾向于把业务拆分的更详细一些，例如把不同的业务代码拆分成不同的APP。  
在项目目录下单独定义一个app目录，用来存放我们不同业务线的代码文件，这样就很容易进行横向扩展。大致目录结构如下：
```
gin_demo
├── app
│   ├── blog
│   │   ├── handler.go
│   │   └── router.go
│   └── shop
│       ├── handler.go
│       └── router.go
├── go.mod
├── go.sum
├── main.go
└── routers
    └── routers.go
```
其中app/blog/router.go用来定义post相关路由信息:
```
func Routers(e *gin.Engine) {
    e.GET("/post", postHandler)
    e.GET("/comment", commentHandler)
}
```
app/shop/router.go用来定义shop相关路由信息：
```
func Routers(e *gin.Engine) {
    e.GET("/goods", goodsHandler)
    e.GET("/checkout", checkoutHandler)
}
```
routers/routers.go中根据需要定义Include函数用来注册子app中定义的路由，Init函数用来进行路由的初始化操作：
```
//Option类型是一个函数类型
type Option func(*gin.Engine)

var options = []Option{}

// 注册app的路由配置, Option类型切片
func Include(opts ...Option) {
    options = append(options, opts...)
}

// 初始化
func Init() *gin.Engine {
    //初始化gin.Engine结构体
    r := gin.New()
    for _, opt := range options {
    	//用app/blog/router.go和app/shop/router.go中定义的Routers(e *gin.Engine)进行初始化
        opt(r)
    }
    return r
}
```
main.go中按如下方式先注册子app中的路由，然后再进行路由的初始化：
```
func main() {
    // 加载多个APP的路由配置
    routers.Include(shop.Routers, blog.Routers)
    // 初始化路由
    r := routers.Init()
    if err := r.Run(); err != nil {
        fmt.Println("startup service failed, err:%v\n", err)
    }
}
```
