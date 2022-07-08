# gin数据解析和绑定 #
## Json 数据解析和绑定 ##
```
// 定义接收json数据的结构体
type Login struct {
   // binding:"required"修饰的字段，若接收为空值，则报错，是必须字段
   User    string `form:"username" json:"user" uri:"user" xml:"user" binding:"required"`
   Pssword string `form:"password" json:"password" uri:"password" xml:"password" binding:"required"`
}

func main() {
   // 1.创建路由
   // 默认使用了2个中间件Logger(), Recovery()
   r := gin.Default()
   // JSON绑定
   r.POST("loginJSON", func(c *gin.Context) {
      // 声明接收的变量
      var json Login
      // 将request的body中的数据，自动按照json格式解析到结构体
      if err := c.ShouldBindJSON(&json); err != nil {
         // 返回错误信息
         // gin.H封装了生成json数据的工具
         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
         return
      }
      // 判断用户名密码是否正确,返回状态
      if json.User != "root" || json.Pssword != "admin" {
         c.JSON(http.StatusBadRequest, gin.H{"status": "304"})
         return
      }
      c.JSON(http.StatusOK, gin.H{"status": "200"})
   })
   r.Run(":8000")
}
```
## 表单数据解析和绑定 ##
```
        //结构体
        var form Login
        // Bind()默认解析并绑定form格式
        // 根据请求头中content-type自动推断
        if err := c.Bind(&form); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
```
##  URI数据解析和绑定 ##
```
        //结构体
        var login Login
        // ShouldBindUri()默认解析并绑定url格式
        // 根据请求头中content-type自动推断
        if err := c.ShouldBindUri(&login); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
