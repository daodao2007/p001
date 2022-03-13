# libp2p入门,开始区块链开发
libp2p 库是由ipfs组织开发的一个p2p开发包，现在不少区块链项目都构建在libp2p的基础之上。

libp2p是一个跨语言的p2p 基础实现。
支持的语言有，go,js,rust,python几种，目前go,js,rust有官方的入门教程。

我选择用go 语言来开始这个入门教程。

## 环境
首先我使用的go 版本是 1.18，版本不同对于安装libp2p依赖有一些差异。

## 建立项目


```
mkdir p001
cd p001
go mod init p001
go get github.com/libp2p/go-libp2p@v0.18.0-rc5
```
在安装依赖的时候如果使用golang 1.17 可以不指定版本
## 开始一个节点
```
package main

import (
	"fmt"
	"github.com/libp2p/go-libp2p"
)

func main() {
	// 建立一个p2p节点，默认配置
	node, err := libp2p.New()
	if err != nil {
		panic(err)
	}

	// 打印节点的所有地址
	fmt.Println("Listen addresses:", node.Addrs())

	// 关闭节点，然后退出
	if err := node.Close(); err != nil {
		panic(err)
	}
}
```
输出如下：
```
➜ go run .\app.go
Listen addresses: [/ip4/192.168.2.22/tcp/60026 /ip4/127.0.0.1/tcp/60026 /ip6/::1/tcp/60027]
```
使用节点默认配置，上面监听的都是本地地址，端口也是随机的。
也可以使用自定义配置，如绑定一个公网ip和端口。
## 绑定ip、端口
做如下修改
```
node, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/2000"),)
```
### 地址定义

数据传输建立在节点连接的基础上，在可以拨号远程节点并建立连接之前，需要知道远程节点的监听地址。因为每种传输协议都有自己的地址格式，所以Libp2p使用一种称为“multiaddr”的编码方案来统一不同的协议的地址格式。

TCP/IP传输协议“multiaddr”的描述如下:

/ip4/127.0.0.1/tcp/9999
UDP传输协议“multiaddr”的描述如下:

/ip4/127.0.0.1/udp/9998
用这种描述方式来代替127.0.0.1:9999的好处是什么呢？

“multiaddr”能更明确的描述使用的协议，如127.0.0.1属于IPv4协议，9999属于TCP协议，9998属于UDP协议。

以上为“multiaddr”描述的节点监听地址，当拨号一个节点时也是使用“multiaddr”，但需要加上远程节点的ID，例如：

/ip4/192.168.100.100/tcp/9999/p2p/QmcEPrat8ShnCph8WjkREzt5CPXF2RwhYxYBALDcLC1iV6
这样就知道对方使用IP4，地址:192.168.100.100，TCP协议，端口：9999，是一个P2P节点，节点ID：QmcEPrat8ShnCph8WjkREzt5CPXF2RwhYxYBALDcLC1iV6。

上面是复制的别人文章，大家这么使用就好了。
## 提供一个ping 服务

在libp2p的包里面，默认已经支持了ping服务，我们只需要导入ping的包就可以了
这么使用：

```
ps := ping.NewPingService(node)
```
这样就可以在程序里增加ping服务
打印节点的multiaddr 
```
	for _, v := range node.Addrs() {
		fmt.Printf("%s/p2p/%s\n", v, node.ID())
	}
```
输出如下：
```
Listen addresses: [/ip4/127.0.0.1/tcp/2000]
/ip4/127.0.0.1/tcp/2000/p2p/QmVeG1AMxy5D5MDMZSyLN4VAKMEFjPVvma3Wj1oVuA4dKy
```

## 一个使用ping服务的客户端

几个步骤

1. 初始化节点 
2. 把ping服务节点的地址添加到Peerstore 
3. 添加自己的ping服务 
4. 通过ping服务发起ping,服务节点回复pong

前面三步的代码如下：
```
	node, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/0"),
		libp2p.Ping(false),
	) //初始化节点 
	if err != nil {
		panic(err)
	}
	addr, err := multiaddr.NewMultiaddr(*dest) //"/ip4/127.0.0.1/tcp/2000/p2p/QmVeG1AMxy5D5MDMZSyLN4VAKMEFjPVvma3Wj1oVuA4dKy")
	if err != nil {
		panic(err)
	}
	info, err := peer.AddrInfoFromP2pAddr(addr)
	if err != nil {
		panic(err)
	}
	node.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL) //把ping服务节点的地址添加到Peerstore 
	if err != nil {
		panic(err)
	}
```
添加ping服务和发起ping:
```
    ps := ping.NewPingService(node)
	fmt.Println("sending 5 ping messages to", addr)
	ch := ps.Ping(context.Background(), info.ID) //发起ping，id是提供pong服务的node id
	for i := 0; i < 5; i++ {
		res := <-ch
		fmt.Println("got ping response!", "RTT:", res.RTT)
	}

```

运行结果如下：
```
go run .\client\app.go -d /ip4/127.0.0.1/tcp/2000/p2p/QmVeG1AMxy5D5MDMZSyLN4VAKMEFjPVvma3Wj1oVuA4dKy
sending 5 ping messages to /ip4/127.0.0.1/tcp/2000/p2p/QmVeG1AMxy5D5MDMZSyLN4VAKMEFjPVvma3Wj1oVuA4dKy
got ping response! RTT: 0s
got ping response! RTT: 0s
got ping response! RTT: 1.0011ms
got ping response! RTT: 0s
got ping response! RTT: 0s
```

我这个例子和官方的还是有些不一样，
官方的例子，其实比较传统，比如他有connect 代码，其实在这个例子中Ping的时候指定了，node id,他是会自动去连接的，不需要显式的去发起connect

代码提交到了github，
[golang use libp2p](https://github.com/daodao2007/p001) 详细的可以看代码

参考文档：
[docs.libp2p.io](https://docs.libp2p.io/tutorials/getting-started/go/)
[官方 golang examples](https://github.com/libp2p/go-libp2p/tree/master/examples)
