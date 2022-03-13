package main

import (
	"encoding/hex"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"p001"

	//"io"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 建立一个p2p节点，默认配置
	priv, _ := hex.DecodeString("4e1518672e45fb2746ec5a217330ed24d815d44537da647e973c06d0b0069053")
	pk, err := crypto.UnmarshalSecp256k1PrivateKey(priv)
	if err != nil {
		fmt.Println(err.Error())
	}
	node, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/2000"),
		libp2p.Ping(false),
		libp2p.Identity(pk),
	)
	if err != nil {
		panic(err)
	}

	// 打印节点的所有地址

	fmt.Println("Listen addresses:", node.Addrs())
	for _, v := range node.Addrs() {
		fmt.Printf("%s/p2p/%s\n", v, node.ID())
	}
	//node.SetStreamHandler()

	//node.Addrs()[0].String()
	ps := ping.NewPingService(node)

	fmt.Println(ps)
	echo := p001.NewExecService(node)
	fmt.Println(echo)
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
