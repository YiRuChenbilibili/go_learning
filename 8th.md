# CRUD接口 #
增删改查
## 增加记录 ##
```
package main

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string
	Age      int
	Birthday time.Time
}

func main() {
	dsn := "root:chen0309@tcp(localhost:3306)/gorm?charset=utf8&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//GORM 的 AutoMigrate() 方法用于自动迁移 ORM 的 Schemas。所谓 “迁移” 就是刷新数据库中的表格定义，使其保持最新（只增不减）。
	//AutoMigrate 会创建（新的）表、缺少的外键、约束、列和索引，并且会更改现有列的类型（如果其大小、精度、是否为空可更改的话）。但不会删除未使用的列，以保护现存的数据。
	//自动创建users表
	db.AutoMigrate(&User{})
	//记录
	user := User{Name: "Jinzhu", Age: 18, Birthday: time.Now()}
	//创建记录
	db.Create(&user)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
}
```
![image](https://user-images.githubusercontent.com/24589721/178227966-004c2470-f4fa-423e-8606-632825f216ea.png)

