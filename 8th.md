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
**默认值**   
可以通过标签定义字段的默认值(*会在mysql中进行默认值设置*），例如：
```
type Animal struct {
    ID   int64
    //设置默认值为galeone
    Name string `gorm:"default:'galeone'"`
    Age  int64
}
```
然后 SQL 会排除那些没有值或者有零值的字段，在记录插入数据库之后，gorm将从数据库中加载这些字段的值。 
```
var animal = Animal{Age: 99, Name: ""}
db.Create(&animal)
// INSERT INTO animals("age") values('99');
// SELECT name from animals WHERE ID=111; // 返回的主键是 111
// animal.Name => 'galeone' //名字为默认值的名字
```
**注意:所有包含零值的字段，像 0，''，false 或者其他的 零值 不会被保存到数据库中，但会使用这个字段的默认值(即会导致无法赋零值，只赋默认值)。应该考虑使用指针类型或者其他的值来避免这种情况:**
```
// Use pointer value
type User struct {
  gorm.Model
  Name string
  Age  *int `gorm:"default:18"`
}

// Use scanner/valuer
type User struct {
  gorm.Model
  Name string
  Age  sql.NullInt64 `gorm:"default:18"`
}
```
**在Hook中设置字段值**  
如果想在 BeforeCreate 函数中更新字段的值，应该使用 scope.SetColumn，例如：
```
func (user *User) BeforeCreate(scope *gorm.Scope) error {
  scope.SetColumn("ID", uuid.New())
  return nil
}
```
**创建额外选项**
```
// 为插入 SQL 语句添加额外选项
db.Set("gorm:insert_option", "ON CONFLICT").Create(&product)
//on conflict 唯一键
// INSERT INTO products (name, code) VALUES ("name", "code") ON CONFLICT;
```
## 删除记录 ##
