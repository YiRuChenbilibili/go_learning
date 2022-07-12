# gorm之关联 #
## Belongs To ##
belongs to 会与另一个模型建立了一对一的连接。 这种模型的每一个实例都“属于”另一个模型的一个实例。  
例如，应用包含 user 和 company，并且每个 user 能且只能被分配给一个 company。下面的类型就表示这种关系。 注意，在 User 对象中，有一个和 Company 一样的 CompanyID。 默认情况下， CompanyID 被隐含地用来在 User 和 Company 之间创建一个外键关系， 因此必须包含在 User 结构体中才能填充 Company 内部结构体。
```
// `User` 属于 `Company`，`CompanyID` 是外键
type User struct {
  gorm.Model
  Name      string
  //CompanyID是默认外键
  CompanyID int
  Company   Company
}

type Company struct {
//也可以直接使用gorm.Model，包含ID
  ID   int
  Name string
}
```
**重写外键**   
要定义一个 belongs to 关系，数据库的表中必须存在外键。默认情况下，外键的名字，使用拥有者的类型名称加上表的主键的字段名字。例如，定义一个User实体属于Company实体，那么外键的名字一般使用CompanyID。GORM同时提供自定义外键名字的方式：
```
type User struct {
  gorm.Model
  Name         string
  CompanyRefer int  //外键，对应的还是Company的ID
  Company      Company `gorm:"foreignKey:CompanyRefer"`
  // 使用 CompanyRefer 作为外键
}

type Company struct {
  ID   int //默认的关联外键
  Name string
}
```
**关联外键**   
对于 belongs to 关系，GORM 通常使用数据库表，主表（拥有者）的主键值作为外键参考。正如上面的例子，使用主表Company中的主键字段ID作为外键的参考值。如果在user实体中设置了company实体，那么GORM会自动把Company中的ID属性保存到User的CompanyID属性中。  
同样的，也可以使用标签 references 来更改它，例如：
```
type User struct {
  gorm.Model
  Name      string
  CompanyID string //外键
  Company   Company `gorm:"references:Code"` // 使用 Code 作为关联外键
}

type Company struct {
  ID   int
  Code string //对应的关联外键(不再使用ID)
  Name string
}
```
**使用属于**
```
db.Model(&Company).Related(&User)
//// SELECT * FROM User WHERE Company = 111; // 111 is Company's ID
```
```
package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// `User` 属于 `Company`，`CompanyID` 是外键
type User struct {
	gorm.Model
	Name string
	//CompanyID是默认外键
	CompanyID int
	Company   Company
}

type Company struct {
	//也可以直接使用gorm.Model，包含ID
	gorm.Model
	Name string
}

func main() {
	dsn := "root:chen0309@tcp(localhost:3306)/gorm?charset=utf8&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//db.AutoMigrate(&User{})
	//db.AutoMigrate(&Company{})
	company := Company{Name: "mhy"}
	db.Create(&company)
  // Company: company实体
	user := User{Name: "chenyiru", Company: company}
	db.Create(&user)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
}
```
![image](https://user-images.githubusercontent.com/24589721/178443062-ab183f75-6301-49a4-82c0-640ba73b1967.png)

