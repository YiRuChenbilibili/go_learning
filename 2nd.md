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
