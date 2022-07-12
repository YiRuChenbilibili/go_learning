# gorm之关联 #
## Belongs To ##
belongs to 会与另一个模型建立了一对一的连接。 这种模型的每一个实例都**属于**另一个模型的一个实例。  
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
**引用**   
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
  Code string //对应的引用(不再使用ID)
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
       //放入Company: company实体
	user := User{Name: "chenyiru", Company: company}
	db.Create(&user)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
}
```
![image](https://user-images.githubusercontent.com/24589721/178443062-ab183f75-6301-49a4-82c0-640ba73b1967.png)
**外键约束**可以通过OnUpdate, OnDelete配置标签来增加关联关系的级联操作，如下面的例子，通过GORM可以完成用户和公司的级联更新和级联删除操作：
```
type User struct {
  gorm.Model
  Name      string
  CompanyID int
  //级联更新和级联删除
  Company   Company `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Company struct {
  ID   int
  Name string
}
```
## Has One ##
has one 与另一个模型建立一对一的关联，但它和一对一关系有些许不同。 这种关联表明一个模型的每个实例都包含或拥有另一个模型的一个实例。  
例如，应用包含 user 和 credit card 模型，且每个 user 只能有一张 credit card。
```
// User 有一张 CreditCard，UserID 是外键
type User struct {
  gorm.Model
  CreditCard CreditCard
}

type CreditCard struct {
  gorm.Model
  Number string
  UserID uint //外键，将 user 的 ID 保存到自己的 UserID 字段（user添加CreditCard时）
}
```
**重写外键**  
对于 has one 关系，同样必须存在外键字段。拥有者将把属于它的模型的主键保存到这个字段。这个字段的名称通常由 has one 模型的类型加上其 主键 生成，对于上面的例子，它是 UserID。
为 user 添加 credit card 时，它会将 user 的 ID 保存到自己的 UserID 字段。如果你想要使用另一个字段来保存该关系，你同样可以使用标签 foreignKey 来更改它，例如：
```
type User struct {
  gorm.Model
  CreditCard CreditCard `gorm:"foreignKey:UserName"`
  // 使用 UserName 作为外键
}

type CreditCard struct {
  gorm.Model
  Number   string
  UserName string //外键
}
```
**引用**   
默认情况下，拥有者实体会将 has one 对应模型的主键保存为外键，可以使用标签 references 来更改它它，用另一个字段来保存，例如下个这个使用 Name 来保存的例子：
```
type User struct {
  gorm.Model
  //`gorm:"index"`创建索引
  Name       string     `gorm:"index"` //引用
  CreditCard CreditCard `gorm:"foreignkey:UserName;references:name"`
}

type CreditCard struct {
  gorm.Model
  Number   string
  UserName string //外键
}
```
**多态关联**   
GORM 为 has one 和 has many 提供了多态关联支持，它会将拥有者实体的表名、主键值都保存到多态类型的字段中。
```
type Cat struct {
  ID    int
  Name  string
  Toy   Toy `gorm:"polymorphic:Owner;"`
}

type Dog struct {
  ID   int
  Name string
  Toy  Toy `gorm:"polymorphic:Owner;"`
}

type Toy struct {
  ID        int
  Name      string
  OwnerID   int
  OwnerType string
}

db.Create(&Dog{Name: "dog1", Toy: Toy{Name: "toy1"}})
// INSERT INTO `dogs` (`name`) VALUES ("dog1")
// INSERT INTO `toys` (`name`,`owner_id`,`owner_type`) VALUES ("toy1","1","dogs")
```
**related**
```
var card CreditCard
db.Model(&user).Related(&card, "CreditCard")
//// SELECT * FROM credit_cards WHERE user_id = 123; // 123 是用户表的主键
// CreditCard  是用户表的字段名，这意味着获取用户的信用卡关系并写入变量 card。
// 像上面的例子，如果字段名和变量类型名一样，它就可以省略， 像：
db.Model(&user).Related(&card)
```
## Has Many ##
has many 与另一个模型建立了一对多的连接。 不同于 has one，拥有者可以有零或多个关联模型。例如，应用包含 user 和 credit card 模型，且每个 user 可以有多张 credit card。
```
// User 有多张 CreditCard，UserID 是外键
type User struct {
  gorm.Model
  CreditCards []CreditCard
}

type CreditCard struct {
  gorm.Model
  Number string
  UserID uint
}
```
**重写外键**   
```
type User struct {
  gorm.Model
  CreditCards []CreditCard `gorm:"foreignKey:UserRefer"`
}

type CreditCard struct {
  gorm.Model
  Number    string
  UserRefer uint
}
```
**重写引用**   
```
type User struct {
  gorm.Model
  MemberNumber string
  CreditCards  []CreditCard `gorm:"foreignKey:UserNumber;references:MemberNumber"`
}

type CreditCard struct {
  gorm.Model
  Number     string
  UserNumber string
}
```
**多态关联**  GORM 为 has one 和 has many 提供了多态关联支持，它会将拥有者实体的表名、主键都保存到多态类型的字段中。可以使用标签 polymorphicValue 来更改多态类型的值
```
type Dog struct {
  ID   int
  Name string
  Toys []Toy `gorm:"polymorphic:Owner;polymorphicValue:master"`
}

type Toy struct {
  ID        int
  Name      string
  OwnerID   int
  OwnerType string
}

db.Create(&Dog{Name: "dog1", Toy: []Toy{{Name: "toy1"}, {Name: "toy2"}}})
// INSERT INTO `dogs` (`name`) VALUES ("dog1")
// INSERT INTO `toys` (`name`,`owner_id`,`owner_type`) VALUES ("toy1","1","master"), ("toy2","1","master")
```
