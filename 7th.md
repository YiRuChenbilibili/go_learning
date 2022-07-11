# gorm基础 #
Object-Relationl Mapping，即对象关系映射，这里的Relationl指的是关系型数据库  
它的作用是在关系型数据库和对象之间作一个映射，这样，我们在具体的操作数据库的时候，就不需要再去和复杂的SQL语句打交道，只要像平时操作对象一样操作它就可以了。  
## 模型定义 ##
模型一般都是普通的 Golang 的结构体，也有支持指针和接口。
## 结构标签 ##
gorm结构体中支持以下标签：  
![image](https://user-images.githubusercontent.com/24589721/178170247-40f2b845-9704-44f5-91cc-bbe7614f1799.png)  
以及以下关联标签：
![image](https://user-images.githubusercontent.com/24589721/178170301-05760d08-e478-408a-8c4c-62eafe691a0a.png)
##  gorm.Model ##
gorm.Model 是一个包含一些基本字段的结构体, 包含的字段有 ID，CreatedAt， UpdatedAt， DeletedAt。  
你可以用它来嵌入到你的模型中，或者也可以用它来建立自己的模型。
```
// gorm.Model 定义
type Model struct {
  ID        uint `gorm:"primary_key"`
  CreatedAt time.Time
  UpdatedAt time.Time
  DeletedAt *time.Time
}

// 将字段 `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt` 注入到 `User` 模型中
type User struct {
  gorm.Model
  Name string
}

// 自定义 gorm.Model 模型
type User struct {
  ID   int
  Name string
}
```
 ## ID 作为主键 ##
 ```
 type User struct {
  ID   string // 字段名 `ID` 将被作为默认的主键名
}

// 设置字段 `AnimalID` 为默认主键
type Animal struct {
  AnimalID int64 `gorm:"primary_key"`
  Name     string
  Age      int64
}
```
##  表名 ##
**复数表明**  
表名是结构体名称的复数形式
```
type User struct {} // 默认的表名是 `users`

// 设置 `User` 的表名为 `profiles`
func (User) TableName() string {
  return "profiles"
}

func (u User) TableName() string {
    if u.Role == "admin" {
        return "admin_users"
    } else {
        return "users"
    }
}

// 如果设置禁用表名复数形式属性为 true，`User` 的表名将是 `user`
db.SingularTable(true)
```
**指定表名**   
```
// 用 `User` 结构体创建 `delete_users` 表
db.Table("deleted_users").CreateTable(&User{})

var deleted_users []User
db.Table("deleted_users").Find(&deleted_users)
//// SELECT * FROM deleted_users;

db.Table("deleted_users").Where("name = ?", "jinzhu").Delete()
//// DELETE FROM deleted_users WHERE name = 'jinzhu';
```
**蛇形列名**
```
type User struct {
  ID        uint      // 字段名是 `id`
  Name      string    // 字段名是 `name`
  Birthday  time.Time // 字段名是 `birthday`
  CreatedAt time.Time // 字段名是 `created_at`
}

// 重写列名
type Animal struct {
    AnimalId    int64     `gorm:"column:beast_id"`         // 设置列名为 `beast_id`
    Birthday    time.Time `gorm:"column:day_of_the_beast"` // 设置列名为 `day_of_the_beast`
    Age         int64     `gorm:"column:age_of_the_beast"` // 设置列名为 `age_of_the_beast`
}
```
## 时间戳跟踪 ##
**CreatedAt**  
对于有 CreatedAt 字段的模型，它将被设置为首次创建记录的当前时间。
```
db.Create(&user) // 将设置 `CreatedAt` 为当前时间

// 你可以使用 `Update` 方法来更改默认时间
db.Model(&user).Update("CreatedAt", time.Now())
```
**UpdatedAt**
对于有 UpdatedAt 字段的模型，它将被设置为记录更新时的当前时间。
```
db.Save(&user) // 将设置 `UpdatedAt` 为当前时间

db.Model(&user).Update("name", "jinzhu") // 将设置 `UpdatedAt` 为当前时间,更新记录时updateat自动更新
```
**DeletedAt**
对于有 DeletedAt 字段的模型，当删除它们的实例时，它们并没有被从数据库中删除，只是将 DeletedAt 字段设置为当前时间。参考 Soft Delete
