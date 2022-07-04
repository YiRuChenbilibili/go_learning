# 函数、方法和接口 #
 Go 程序函数启动顺序的示意图：  
 ![image](https://user-images.githubusercontent.com/24589721/177110756-912d27c9-923f-4f4e-99dd-02b0bdfc2531.png)
 
 ## 函数 ##
 函数有具名和匿名之分：
 ```
 // 具名函数
func Add(a, b int) int {
	return a+b
}

// 匿名函数
var Add = func(a, b int) int {
	return a+b
}
```
Go 语言中的函数可以有多个参数和多个返回值,*函数还支持可变数量的参数，可变数量的参数必须是最后出现的参数，可变数量的参数其实是一个切片类型的参数*:
```
// 多个参数和多个返回值
func Swap(a, b int) (int, int) {
	return b, a
}

// 可变数量的参数
// more 对应 []int 切片类型
func Sum(a int, more ...int) int {
	for _, v := range more {
		a += v
	}
	return a
}
```
函数的返回值可命名：
```
func Find(m map[int]int, key int) (value int, ok bool) {
	value, ok = m[key]
	return
}
```
可以通过 defer 语句在 return 语句之后修改返回值:
```
func Inc() (v int) {
	defer func(){ v++ } ()
	return 42
}
//输出42+1,43
```
*Go语言中defer语句会将起后面跟随的语句进行延迟处理，在defer归属的函数即将返回时，将延迟处理的语句按照defer的逆序进行执行，也就是说先被defer的语句最后执行，最后被defer的语句，最先被执行。*  
**任何可以通过函数参数修改调用参数的情形，都是因为函数参数中显式或隐式传入了指针参数。**

## 方法 ##
C++ 语言中方法对应一个类对象的成员函数，是关联到具体对象上的虚表中的。*Go 语言的方法是关联到类型的，这样可以在编译阶段完成方法的静态绑定。*
go中构建某一类型独有方法：
```
文件对象结构体
type File struct {
	fd int
}
// 关闭文件
func (f *File) Close() error {
	// ...
}

// 读文件数据
func (f *File) Read(offset int64, data []byte) int {
	// ...
}
```
这时，CloseFile 和 ReadFile 函数是 File 类型独有的方法（而不是 File 对象方法）。它们也不再占用包级空间中的名字资源，同时File 类型也已经明确了它们操作对象。  
*我们可以给任何自定义类型添加一个或多个方法。每种类型对应的方法必须和类型的定义在同一个包中，因此是无法给 int 这类内置类型添加方法的（因为方法的定义和类型的定义不在一个包中）。对于给定的类型，每个方法的名字必须是唯一的，同时方法和函数一样也**不支持重载。***  
方法是由函数演变而来，只是将函数的第一个对象参数移动到了函数名前面了而已。我们依然可以按照原始的过程式思维来使用方法, 通过叫**方法表达式**的特性可以将方法还原为普通类型的函数：
```
// 不依赖具体的文件对象
// func CloseFile(f *File) error
var CloseFile = (*File).Close //还原

// 不依赖具体的文件对象
// func ReadFile(f *File, offset int64, data []byte) int
var ReadFile = (*File).Read  //还原

// 文件处理
f, _ := OpenFile("foo.dat")
ReadFile(f, 0, data)
CloseFile(f)
```
```
