package main

import (
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 建立一个p2p节点，默认配置
	node, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/2000"),
		libp2p.Ping(false))
	if err != nil {
		panic(err)
	}

	// 打印节点的所有地址

	fmt.Println("Listen addresses:", node.Addrs())
	for _, v := range node.Addrs() {
		fmt.Printf("%s/p2p/%s\n", v, node.ID())
	}

	//node.Addrs()[0].String()
	ps := ping.NewPingService(node)
	fmt.Println(ps)
	//fmt.Println(node.ID())
	fmt.Println("Listen addresses:", node.Network().ListenAddresses())
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	fmt.Println("Received signal, shutting down...")

	// 关闭节点，然后退出
	if err := node.Close(); err != nil {
		panic(err)
	}
}
