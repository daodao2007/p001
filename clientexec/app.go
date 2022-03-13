package main

import (
	"flag"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/multiformats/go-multiaddr"
	"golang.org/x/net/context"
	"p001"
)

func main() {
	dest := flag.String("d", "", "Destination multiaddr string")
	flag.Parse()
	node, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/0"),
		libp2p.Ping(false),
	)
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
	node.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)
	if err != nil {
		panic(err)
	}
	//ps := ping.NewPingService(node)
	ps := p001.NewExecService(node)
	cmdch := make(chan string, 1)
	ch := ps.ExecStart(context.Background(), info.ID, cmdch)
	cmdch <- "go\n"
	result := <-ch
	fmt.Println("recv", result)
	//fmt.Println("sending 5 ping messages to", addr)
	ch := ps.Ping(context.Background(), info.ID)
	/*for i := 0; i < 5; i++ {
		res := <-ch
		fmt.Println("got ping response!", "RTT:", res.RTT)
	}*/

	// shut the node down
	if err := node.Close(); err != nil {
		panic(err)
	}
}
