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
在方法表达式中，因为得到的 ReadFile 和 CloseFile 函数参数中含有 File 这个特有的类型参数，这使得 File 相关的方法无法**和其它不是 File 类型但是有着相同 Read 和 Close 方法的对象无缝适配**。此时，可以通过结合闭包特性（类似lambda表达式）来消除方法表达式中第一个参数类型的差异（*先通过打开对象获得对象类型，再将方法绑定到对象,从而可以处理不同类型的对象（打开不同类型的对象）*）：
```
// 先打开文件对象
f, _ := OpenFile("foo.dat")

// 绑定到了 f 对象
// func Close() error
var Close = func() error {
	return (*File).Close(f) 
	//方法表达式:传入结构体进行实际调用，需要具体的变量f
}

// 绑定到了 f 对象
// func Read(offset int64, data []byte) int
var Read = func(offset int64, data []byte) int {
	return (*File).Read(f, offset, data)
}

// 文件处理
Read(0, data)
Close()
```
**方法表达式:传入结构体进行实际调用**  
方法值（ method value ）其实是一个带有闭包的函数变量，其底层实现原理和带有闭包的匿名函数类似， 接收值被隐式地绑定到方法值（ method value ）的闭包环境中。后续调用不需要再显式地传递接收者。  
也可以用方法值特性可以简化实现（直接通过打开对象，再将方法绑定到对象）：
```
// 先打开文件对象，
f, _ := OpenFile("foo.dat")

// 方法值: 绑定到了 f 对象
// func Close() error
var Close = f.Close //为值传递方式

// 方法值: 绑定到了 f 对象
// func Read(offset int64, data []byte) int
var Read = f.Read

// 文件处理
Read(0, data)
Close()
```
**方法值：直接将方法声明赋值给新变量**  
Go语言**不支持传统面向对象中的继承特性**。Go 语言中，通过在结构体内置匿名的成员来实现继承：
```
import "image/color"

type Point struct{ X, Y float64 }

type ColoredPoint struct {
	Point
	Color color.RGBA
}
```
通过嵌入匿名的成员，我们*不仅可以继承匿名成员的内部成员，而且可以继承匿名成员类型所对应的方法。*我们一般会将 Point 看作基类，把 ColoredPoint 看作是它的继承类或子类。不过这种方式继承的方法并不能实现 C++ 中虚函数的多态特性。*所有继承来的方法的接收者参数依然是那个匿名成员本身，而不是当前的变量。*
```
var cp ColoredPoint
cp.X = 1
fmt.Println(cp.Point.X) // "1"
cp.Point.Y = 2
fmt.Println(cp.Y)       // "2"
```
**在传统的面向对象语言（eg.C++ 或 Java）的继承中，子类的方法是在运行时动态绑定到对象的，因此基类实现的某些方法看到的 this 可能不是基类类型对应的对象，这个特性会导致基类方法运行的不确定性。而在 Go 语言通过嵌入匿名的成员来“继承”的基类方法，this 就是实现该方法的类型的对象，Go 语言中方法是编译时静态绑定的。如果需要虚函数的多态特性，我们需要借助 Go 语言接口来实现。**
## 接口 ##
Go 语言中的面向对象，如果一个对象只要看起来像是某种接口类型的实现，那么它就可以作为该接口类型使用。*Go 语言的接口类型是延迟绑定，可以实现类似虚函数的多态功能。*  
go语言的接口，是一种新的类型定义，它把所有的具有共性的方法定义在一起，**任何其他类型只要实现了这些方法就是实现了这个接口。**  
**实现接口必须实现接口中的所有方法**  
fmt.Fprintf 函数的签名为:```func Fprintf(w io.Writer, format string, args ...interface{}) (int, error)```
其中 io.Writer 用于输出的接口，error 是内置的错误接口，它们的定义如下：
```
type io.Writer interface {
	Write(p []byte) (n int, err error)
}

type error interface {
	Error() string
}
```
只要自己定制的类型实现了io.Writer中的方法，就实现了io.Writer，此时就可以作为fmt.Fprintf 函数的输出
```
//定义UpperWriter结构体
type UpperWriter struct {
	io.Writer  //结构体中匿名成员,接口类型
}
//UpperWriter结构体实现io.Writer接口方法
func (p *UpperWriter) Write(data []byte) (n int, err error) {
	return p.Writer.Write(bytes.ToUpper(data))  //匿名成员所对应的方法
}
//可以输出UpperWriter
func main() {
	fmt.Fprintln(&UpperWriter{os.Stdout}, "hello, world")
}
```
有时候对象和接口之间太灵活了，导致需要人为地限制这种无意之间的适配。常见的做法是定义一个含特殊方法来区分接口。比如 runtime 包中的 Error 接口就定义了一个特有的 RuntimeError 方法，用于避免其它类型无意中适配了该接口：
```
type runtime.Error interface {
	error

	// RuntimeError is a no-op function but
	// serves to distinguish types that are run time
	// errors from ordinary errors: a type is a
	// run time error if it has a RuntimeError method.
	RuntimeError()
}
```
再严格一点的做法是给接口定义一个私有方法。只有满足了这个私有方法的对象才可能满足这个接口，而私有方法的名字是包含包的绝对路径名的，因此只能在包内部实现这个私有方法才能满足这个接口。  
这种防护措施也不是绝对的。**通过*在结构体中嵌入匿名类型成员*，可以继承匿名类型的方法。其实这个被嵌入的匿名成员不一定是普通类型，也可以是*接口类型*。我们可以通过嵌入匿名的 testing.TB 接口来伪造私有的 private 方法，因为接口方法是延迟绑定，编译时 private 方法是否真的存在并不重要。**
