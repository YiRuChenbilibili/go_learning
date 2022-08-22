# Channel #
Goroutine 和 Channel 是 Go 语言并发编程的两大基石。Goroutine 用于执行并发任务，Channel 用于 goroutine 之间的同步、通信。Go提倡使用通信的方法代替共享内存，当一个Goroutine需要和其他Goroutine资源共享时，Channel就会在他们之间架起一座桥梁，并提供确保安全同步的机制。channel本质上其实还是一个队列，遵循FIFO原则。具体规则如下：

先从 Channel 读取数据的 Goroutine 会先接收到数据；
先向 Channel 发送数据的 Goroutine 会得到先发送数据的权利；

**创建Channel**     
```
通道实例 := make(chan 数据类型)

ch1 := make(chan int)                 // 创建一个整型类型的通道
ch2 := make(chan interface{})         // 创建一个空接口类型的通道, 可以存放任意格式
```  

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
Go语言中有缓冲的通道（buffered channel）是一种在被接收前能存储一个或者多个值的通道。这种类型的通道并不强制要求 goroutine 之间必须同时完成发送和接收。通道会阻塞发送和接收动作的条件也会不同。只有在通道中没有要接收的值时，接收动作才会阻塞。只有在通道没有可用缓冲区容纳被发送的值时，发送动作才会阻塞。

有缓冲通道的定义方式如下：
```
通道实例 := make(chan 通道类型, 缓冲大小)
```
通道类型：和无缓冲通道用法一致，影响通道发送和接收的数据类型。
缓冲大小：决定通道最多可以保存的元素数量。
通道实例：被创建出的通道实例。
例子:
```
package main

import (
    "sync"
    "time"
)

func main() {
    c := make(chan string, 2)

    var wg sync.WaitGroup
    wg.Add(2)

    go func() {
        defer wg.Done()

        c <- `Golang梦工厂`
        c <- `asong`
    }()

    go func() {
        defer wg.Done()

        time.Sleep(time.Second * 1)
        println(`公众号: `+ <-c)
        println(`作者: `+ <-c)
    }()

    wg.Wait()
}
```
## 消息队列编码实现 ##
**准备**      
定义接口，列出需要实现的方法：
```
type Broker interface {
	publish(topic string, msg interface{}) error
	subscribe(topic string) (<-chan interface{}, error) //返回对应的通道
	unsubscribe(topic string, sub <-chan interface{}) error
	close()
	broadcast(msg interface{}, subscribers []chan interface{})
	setConditions(capacity int)
}
```
publish：进行消息的推送，有两个参数即topic、msg，分别是订阅的主题、要传递的消息      
subscribe：消息的订阅，传入订阅的主题，即可完成订阅，并返回对应的channel通道用来接收数据       
unsubscribe：取消订阅，传入订阅的主题和对应的通道       
close：这个的作用就是很明显了，就是用来关闭消息队列的      
broadCast：这个属于内部方法，作用是进行广播，对推送的消息进行广播，保证每一个订阅者都可以收到      
setConditions：这里是用来设置条件，条件就是消息队列的容量，这样我们就可以控制消息队列的大小了        

封装成客户端可以直接调用的方法：
```
package mq


type Client struct {
	bro *BrokerImpl
}

func NewClient() *Client {
	return &Client{
		bro: NewBroker(),
	}
}

func (c *Client)SetConditions(capacity int)  {
	c.bro.setConditions(capacity)
}

func (c *Client)Publish(topic string, msg interface{}) error{
	return c.bro.publish(topic,msg)
}

func (c *Client)Subscribe(topic string) (<-chan interface{}, error){
	return c.bro.subscribe(topic)
}

func (c *Client)Unsubscribe(topic string, sub <-chan interface{}) error {
	return c.bro.unsubscribe(topic,sub)
}

func (c *Client)Close()  {
	 c.bro.close()
}

func (c *Client)GetPayLoad(sub <-chan interface{})  interface{}{
	for val:= range sub{
		if val != nil{
			return val
		}
	}
	return nil
}
```

**消息队列的结构**
```
type BrokerImpl struct {
	exit chan bool
	capacity int

	topics map[string][]chan interface{} // key： topic  value ： queue
	sync.RWMutex // 同步锁
}
```
exit：也是一个通道，这个用来做关闭消息队列用的     
capacity：即用来设置消息队列的容量      
topics：这里使用一个map结构，key即是topic，其值则是一个切片，chan类型，这里这么做的原因是我们一个topic可以有多个订阅者，所以一个订阅者对应着一个通道      
sync.RWMutex：读写锁，这里是为了防止并发情况下，数据的推送出现错误，所以采用加锁的方式进行保证      

**Publish和broadcast**      
传入数据并进行广播：
```
func (b *BrokerImpl) publish(topic string, pub interface{}) error {
	select {
	case <-b.exit:
		return errors.New("broker closed")
	default:
	}

	b.RLock()
	subscribers, ok := b.topics[topic]
	b.RUnlock()
	if !ok {
		return nil
	}

	b.broadcast(pub, subscribers)
	return nil
}


func (b *BrokerImpl) broadcast(msg interface{}, subscribers []chan interface{}) {
	count := len(subscribers)
	concurrency := 1

	//当数量过多时，用for循环推送消息
	switch {
	case count > 1000:
		concurrency = 3
	case count > 100:
		concurrency = 2
	default:
		concurrency = 1
	}
	pub := func(start int) {
		for j := start; j < count; j += concurrency {
			select {
			//正常推送数据
			case subscribers[j] <- msg: 
			//超时机制，超过5毫秒就停止推送
			case <-time.After(time.Millisecond * 5):
			//结束
			case <-b.exit:
				return
			}
		}
	}
	for i := 0; i < concurrency; i++ {
		go pub(i)
	}
}
```
其中使用了select语句： **select 的代码形式和 switch 非常相似， 不过 select 的 case 里的操作语句只能是【IO 操作】 。**       
使用 select 实现 timeout 机制：
```
timeout := make (chan bool, 1)
go func() {
    time.Sleep(1e9) // sleep one second
    timeout <- true
}()
ch := make (chan int)
select {
case <- ch:
    fmt.Println("get ch")
case <- timeout:
    fmt.Println("timeout!")
}
```

**subscribe 和 unsubScribe**
```
func (b *BrokerImpl) subscribe(topic string) (<-chan interface{}, error) {
	select {
	case <-b.exit:
		return nil, errors.New("broker closed")
	default:
	}

	ch := make(chan interface{}, b.capacity) //订阅的主题创建一个channel，按订阅数量设置容量
	b.Lock()
	b.topics[topic] = append(b.topics[topic], ch) //将订阅者加入到对应的topic
	b.Unlock()
	return ch, nil //将订阅者加入到对应的topic 
}
func (b *BrokerImpl) unsubscribe(topic string, sub <-chan interface{}) error {
	select {
	case <-b.exit:
		return errors.New("broker closed")
	default:
	}

	b.RLock()
	subscribers, ok := b.topics[topic]
	b.RUnlock()

	if !ok {
		return nil
	}
	// delete subscriber删除
	var newSubs []chan interface{}
	for _, subscriber := range subscribers {
		//如果与要删除的chan相同，则丢弃
		if subscriber == sub {
			continue
		}
		//否则重新加入切片
		newSubs = append(newSubs, subscriber)
	}

	b.Lock()
	//更新topic
	b.topics[topic] = newSubs
	b.Unlock()
	return nil
}
```
**close**
```
func (b *BrokerImpl) close()  {
	select {
	//b.exit不为空且返回零值时，表示chan已经关闭了，防止关闭已经关闭的chan
	case <-b.exit:
		return
	default:
		//关闭b.exit
		close(b.exit)
		b.Lock()
		b.topics = make(map[string][]chan interface{})
		b.Unlock()
	}
	return
}
``` 
这句代码b.topics = make(map[string][]chan interface{})比较重要，这里主要是为了保证下一次使用该消息队列不发生冲突。
```
close函数是一个内建函数， 用来关闭channel，这个channel要么是双向的， 要么是只写的（chan<- Type）。       
这个方法应该只由发送者调用， 而不是接收者。 
当最后一个发送的值都被接收者从关闭的channel(下简称为c)中接收时, 
接下来所有接收的值都会非阻塞直接成功，返回channel元素的零值。 
如下的代码： 
如果c已经关闭（c中所有值都被接收）， x, ok := <- c， 读取ok将会得到false。
```

**setConditions GetPayLoad**      
```
//获取消息队列容量
func (b *BrokerImpl)setConditions(capacity int)  {
	b.capacity = capacity
}
//封装一个方法来获取订阅的消息
func (c *Client)GetPayLoad(sub <-chan interface{})  interface{}{
	for val:= range sub{
		if val != nil{
			return val
		}
	}
	return nil
}
```
**测试**
https://blog.csdn.net/qq_39397165/article/details/108686391
