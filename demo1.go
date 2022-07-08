//GO 函数式选项模式(Functional Options Pattern) 初始化结构体
/*
Option模式的优缺点
优点:
    1. 支持传递多个参数，并且在参数个数、类型发生变化时保持兼容性
    2. 任意顺序传递参数
    3. 支持默认值
    4. 方便拓展
缺点:
    1. 增加许多function，成本增大
    2. 参数不太复杂时，尽量少用
*/
package main

import "fmt"

type Client struct {
	Id        int64
	AppKey    string
	AppSecret string
}

type Option func(*Client) //  go函数的参数都是值传递 因此想要修改Client(默认值) 必须传递指针

func WithAppKey(appKey string) Option {
	return func(client *Client) {
		client.AppKey = appKey
	}
}

func WithAppSecret(appSecret string) Option {
	return func(client *Client) {
		client.AppSecret = appSecret
	}
}

//
//  NewClient
//  @Description 将一个函数的参数设置为可选的功能
//  @param id 固定参数，也可以将所有都放进可选参数 opts 中
//  @param opts
//  @return Client 返回 *Client 和 Client 都可以
//
func NewClient(id int64, opts ...Option) Client {
	o := Client{
		Id:        id,
		AppKey:    "key_123456",
		AppSecret: "secret_123456",
	}

	for _, opt := range opts {
		opt(&o) //  go函数的参数都是值传递 因此想要修改Client(默认值) 必须传递指针
	}

	return o
}

func main() {
	//  使用默认值
	fmt.Println(NewClient(1)) //  {1 key_123456 secret_123456}
	//  使用传入的值
	fmt.Println(NewClient(2, WithAppKey("change_key_222"))) //  {2 change_key_222 secret_123456}
	//  不按照顺序传入
	fmt.Println(NewClient(3, WithAppSecret("change_secret_333"))) //  {3 key_123456 change_secret_333}
	fmt.Println(NewClient(4, WithAppSecret("change_secret_444"), WithAppKey("change_key_444"))) //  {4 change_key_444 change_secret_444}
}

