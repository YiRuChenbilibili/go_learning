# gin响应 #
## 数据格式响应 ##
```
r.GET("/moreJSON", func(c *gin.Context) {
	// 以下方式都会输出 :   {"user": "Lena", "Message": "hey", "Number": 123}
  //结构体响应
	var msg struct {
		Name    string `json:"user" xml:"user"`
		Message string
		Number  int
	}
	msg.Name = "Lena"
	msg.Message = "hey"
	msg.Number = 123
  c.JSON(http.StatusOK, msg)
  //	c.XML(http.StatusOK, msg)
	//  c.YAML(http.StatusOK, msg)
  
  //JSON/XML/YAML响应
	c.JSON(http.StatusOK, gin.H{"user": "Lena", "Message": "hey", "Number": 123})
	c.XML(http.StatusOK, gin.H{"user": "Lena", "Message": "hey", "Number": 123})
	c.YAML(http.StatusOK, gin.H{"user": "Lena", "Message": "hey", "Number": 123})
  
  //protobuf格式
  reps := []int64{int64(1), int64(2)}
  // 定义数据
  label := "label"
  // 传protobuf格式数据
  data := &protoexample.Test{
       Label: &label,
       Reps:  reps,
  }
  c.ProtoBuf(200, data)
})
```
## 模板响应 ##
gin支持加载HTML模板, 然后根据模板参数进行配置并返回相应的数据，本质上就是字符串替换  
LoadHTMLGlob()方法可以加载模板文件  
```
package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    //加载模板文件
    r.LoadHTMLGlob("index/*")
    r.GET("/index", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.html", gin.H{"title": "我是测试", "ce": "123456"})
    })
    r.Run()
}
```
```
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>{{.title}}</title>
</head>
    <body>
        Hello,chenyiru!
    </body>
</html>
```
## 文件响应 ##
```
//获取当前文件的相对路径
router.Static("/assets", "./assets")
//
router.StaticFS("/more_static", http.Dir("my_file_system"))
//获取相对路径下的文件
router.StaticFile("/favicon.ico", "./resources/favicon.ico")
```
## 重定向 ##
```
package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    r.GET("/index", func(c *gin.Context) {
        ////支持内部和外部的重定向,定向到资源位置
        c.Redirect(http.StatusMovedPermanently, "http://www.5lmh.com")
    })
    r.Run()
}
```
