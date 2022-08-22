# Channel #
Goroutine 和 Channel 是 Go 语言并发编程的两大基石。Goroutine 用于执行并发任务，Channel 用于 goroutine 之间的同步、通信。Go提倡使用通信的方法代替共享内存，当一个Goroutine需要和其他Goroutine资源共享时，Channel就会在他们之间架起一座桥梁，并提供确保安全同步的机制。channel本质上其实还是一个队列，遵循FIFO原则。具体规则如下：

先从 Channel 读取数据的 Goroutine 会先接收到数据；
先向 Channel 发送数据的 Goroutine 会得到先发送数据的权利；

**创建Channel**     
```通道实例 := make(chan 数据类型)```  

数据类型：通道内传输的元素类型。
通道实例：通过make创建的通道句柄。

## 无缓冲通道的使用 ##
Go语言中无缓冲的通道（unbuffered channel）是指在接收前没有能力保存任何值的通道。这种类型的通道要求发送 goroutine 和接收 goroutine 同时准备好，才能完成发送和接收操作。

无缓冲通道的定义方式如下：
```
通道实例 := make(chan 通道类型)
```
通道类型：和无缓冲通道用法一致，影响通道发送和接收的数据类型。
缓冲大小：0
通道实例：被创建出的通道实例。
例子:
```
package main

import (
    "sync"
    "time"
)

func main() {
    c := make(chan string)

    var wg sync.WaitGroup
    wg.Add(2)

    go func() {
        defer wg.Done()
        c <- `Golang梦工厂`
    }()

    go func() {
        defer wg.Done()

        time.Sleep(time.Second * 1)
        println(`Message: `+ <-c)
    }()

    wg.Wait()
}
```
## 带缓冲的通道的使用 ##
