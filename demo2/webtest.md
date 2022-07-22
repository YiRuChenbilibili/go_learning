## web-v2 ##  
运行程序`go run main.go`   
进入`http://localhost:8000/All`页面   
![image](https://user-images.githubusercontent.com/24589721/180251701-35246951-7b34-4ab7-99b5-ab3c076d846d.png)    
管理员部分可编辑或查看用户信息，用户部分可编辑或查看观点信息。

## 优化总结 ##
1、对于web服务，在第二周的用户信息创建及编辑的基础上，增加了与用户相关的观点信息编辑。观点信息关联于用户信息（一个用户可以有多个观点信息，一条观点信息只能属于一个用户），能够同步跟随用户信息更新和删除。同时给用户信息增加了密码项，设置密码默认值为123456  
2、给用户信息编辑部分增加了简单的admin登录，使用了基础的gin.BasicAuth()中间件方法。    
```
//构建中间件用于admin登录
func LoginMiddleware() gin.HandlerFunc {
	// 设置登陆用户和密码
	accounts := gin.Accounts{"admin": "admin@123456", "amy": "amy@123456"}
	return gin.BasicAuth(accounts)
}
```
3、给观点信息编辑部分加入了使用cookie的登录模块，并使用中间件传递用户ID信息。   
4、观点信息只能由登录用户进行编辑，一个用户只能查看，创建，删除自己的观点信息。
5、使用路由群组重新整理了routers.go中的post，get等方法，分为对用户信息进行操作的路由群组和对观点信息进行操作的路由群组。      
6、将所有功能汇总到localhost:8000/All界面，可以直接点击按钮转跳各个界面。
