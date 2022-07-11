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
