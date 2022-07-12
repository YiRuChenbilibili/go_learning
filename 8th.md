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
## 查询记录 ##
gorm查询数据本质上就是提供一组函数，帮我们快速拼接sql语句，尽量减少编写sql语句的工作量。  
gorm查询结果我们一般都是**保存到结构体(struct)变量**，所以在执行查询操作之前需要根据自己想要查询的数据定义结构体类型。
```
//定义接收查询结果的结构体变量
user := User{}
//user := []User{} 当获取结果为数组时的定义(结构体数组)
// 获取第一条记录，按主键排序
db.First(&user)
//// SELECT * FROM users ORDER BY id LIMIT 1;

// 获取一条记录，不指定排序
db.Take(&user)
//// SELECT * FROM users LIMIT 1;

// 获取最后一条记录，按主键排序
db.Last(&user)
//// SELECT * FROM users ORDER BY id DESC LIMIT 1;

// 获取所有的记录，Find函数返回的是一个数组，所以定义一个数组用来接收结果user := []User{}
db.Find(&users)
//// SELECT * FROM users;

// 通过主键进行查询 (仅适用于主键是数字类型)
db.First(&user, 10)
//// SELECT * FROM users WHERE id = 10;

fmt.Println(user)
```
![image](https://user-images.githubusercontent.com/24589721/178425847-c10a59db-70e6-416f-8ce2-6f10f7900608.png)

**查询语句**   
where：
```
// 获取第一条匹配的记录
db.Where("name = ?", "jinzhu").First(&user)
//// SELECT * FROM users WHERE name = 'jinzhu' limit 1;
```
Not:
```
// 不包含
db.Not("name", []string{"jinzhu", "jinzhu 2"}).Find(&users)
//// SELECT * FROM users WHERE name NOT IN ("jinzhu", "jinzhu 2");
```
Or:
```
// Struct
db.Where("name = 'jinzhu'").Or(User{Name: "jinzhu 2"}).Find(&users)
//// SELECT * FROM users WHERE name = 'jinzhu' OR name = 'jinzhu 2';
```
gorm可以使用struct或map实现查询   
**当通过struct进行查询的时候，GORM 将会查询这些字段的非零值， 意味着你的字段包含 0， ''， false 或者其他 零值, 将不会出现在查询语句中， 例如:**
```
db.Where(&User{Name: "jinzhu", Age: 0}).Find(&users)
//// SELECT * FROM users WHERE name = "jinzhu";
```
可以考虑适用指针类型或者 scanner/valuer 来避免这种情况:
```
// 使用指针类型
type User struct {
  gorm.Model
  Name string
  Age  *int
}

// 使用 scanner/valuer
type User struct {
  gorm.Model
  Name string
  Age  sql.NullInt64
}
```
**行内条件查询**   
需要注意的是，当使用链式调用传入行内条件查询时，这些查询不会被传参给后续的中间方法。  
```
// 通过主键进行查询 (仅适用于主键是数字类型)
db.First(&user, 23)
//// SELECT * FROM users WHERE id = 23 LIMIT 1;
// 非数字类型的主键查询
db.First(&user, "id = ?", "string_primary_key")
//// SELECT * FROM users WHERE id = 'string_primary_key' LIMIT 1;
```
**额外的查询选项**
```
// 为查询 SQL 添加额外的选项
db.Set("gorm:query_option", "FOR UPDATE").First(&user, 10)
//// SELECT * FROM users WHERE id = 10 FOR UPDATE;
```
**Row()和Scan()
```
rows, err := db.Table("orders").Select("date(created_at) as date, sum(amount) as total").Group("date(created_at)").Rows()
//数据库连接需要关闭
defer rows.Close()
for rows.Next() {
    err = rows.Scan(user)
} 

type Result struct {
    Date  time.Time
    Total int64
}
//使用Scan()将信息快速绑定到结构体
db.Table("orders").Select("date(created_at) as date, sum(amount) as total").Group("date(created_at)").Having("sum(amount) > ?", 100).Scan(&results)
```
## 记录更新 ##
**更新所有数据**
Save 方法在执行 SQL 更新操作时将包含所有字段，即使这些字段没有被修改。
```
//更新表中第一条数据
db.First(&user)
user.Name = "jinzhu 2"
user.Age = 100
db.Save(&user)
```
**更新已更改的字段**只想更新已经修改了的字段，可以使用 Update，Updates 方法。
```
// 如果单个属性被更改了，更新它
db.Model(&user).Update("name", "hello")
//// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE id=111;

// 使用组合条件更新单个属性
db.Model(&user).Where("active = ?", true).Update("name", "hello")
//// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE id=111 AND active=true;

// 使用 `map` 更新多个属性，只会更新那些被更改了的字段
db.Model(&user).Updates(map[string]interface{}{"name": "hello", "age": 18, "actived": false})
//// UPDATE users SET name='hello', age=18, actived=false, updated_at='2013-11-17 21:34:10' WHERE id=111;

// 使用 `struct` 更新多个属性，只会更新那些被修改了的和非空的字段
db.Model(&user).Updates(User{Name: "hello", Age: 18})
//// UPDATE users SET name='hello', age=18, updated_at = '2013-11-17 21:34:10' WHERE id = 111;

// 警告： 当使用结构体更新的时候, GORM 只会更新那些非空的字段
// 例如下面的更新，没有东西会被更新，因为像 "", 0, false 是这些字段类型的空值
db.Model(&user).Updates(User{Name: "", Age: 0, Actived: false})
```
**更新选中的字段**只想更新或者忽略某些字段，可以使用 Select，Omit方法。
```
db.Model(&user).Select("name").Updates(map[string]interface{}{"name": "hello", "age": 18, "actived": false})
//// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE id=111;

db.Model(&user).Omit("name").Updates(map[string]interface{}{"name": "hello", "age": 18, "actived": false})
//// UPDATE users SET age=18, actived=false, updated_at='2013-11-17 21:34:10' WHERE id=111;
```
**更新列Hook方法**
上面的更新操作更新时会执行模型的 BeforeUpdate 和 AfterUpdate 方法，来更新 UpdatedAt 时间戳，并且保存他的 关联。如果你不想执行这些操作，可以使用 UpdateColumn，UpdateColumns 方法。
```
// Update single attribute, similar with `Update`
db.Model(&user).UpdateColumn("name", "hello")
//// UPDATE users SET name='hello' WHERE id = 111;

// Update multiple attributes, similar with `Updates`
db.Model(&user).UpdateColumns(User{Name: "hello", Age: 18})
//// UPDATE users SET name='hello', age=18 WHERE id = 111;
```
**批量更新**不会执行hook函数：
```
db.Table("users").Where("id IN (?)", []int{10, 11}).Updates(map[string]interface{}{"name": "hello", "age": 18})
//// UPDATE users SET name='hello', age=18 WHERE id IN (10, 11);

// 使用结构体更新将只适用于非零值，或者使用 map[string]interface{}
db.Model(User{}).Updates(User{Name: "hello", Age: 18})
//// UPDATE users SET name='hello', age=18;

// 使用 `RowsAffected` 获取更新影响的记录数
db.Model(User{}).Updates(User{Name: "hello", Age: 18}).RowsAffected
```
**额外的更新选项**
```
// 在更新 SQL 语句中添加额外的 SQL 选项
db.Model(&user).Set("gorm:update_option", "OPTION (OPTIMIZE FOR UNKNOWN)").Update("name", "hello")
//// UPDATE users SET name='hello', updated_at = '2013-11-17 21:34:10' WHERE id=111 OPTION (OPTIMIZE FOR UNKNOWN);
```
