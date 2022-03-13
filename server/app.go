package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"golang.org/x/net/context"

	//"io"
	"os"
	"os/signal"
	"syscall"
)

const (
	//PingSize    = 32
	//pingTimeout = time.Second * 60

	ID = "/ipfs/exec/1.0.0"

	ServiceName = "libp2p.exec"
)

type ExecService struct {
	Host host.Host
}

func NewExecService(h host.Host) *ExecService {
	ps := &ExecService{h}
	h.SetStreamHandler(ID, ps.ExecHandler)
	return ps
}

/*func (ps *ExecService) Echo(ctx context.Context, p peer.ID) <-chan Result {
	return Echo(ctx, ps.Host, p, msg)
}*/
func execError(err error) chan string {
	ch := make(chan string, 1)
	ch <- "error"
	close(ch)
	return ch
}
func Exec(ctx context.Context, h host.Host, p peer.ID) <-chan string {
	s, err := h.NewStream(network.WithUseTransient(ctx, "exec"), p, ID)
	if err != nil {
		return execError(err)
	}
	if err := s.Scope().SetService(ServiceName); err != nil {
		//log.Debugf("error attaching stream to ping service: %s", err)
		s.Reset()
		return execError(err)
	}

	ctx, cancel := context.WithCancel(ctx)

	out := make(chan string)
	go func() {
		defer close(out)
		defer cancel()
		select {
		case cmd := <-out:
			println(cmd)
		case <-ctx.Done():
			return
		}
	}()
	return out
}
func (p *ExecService) ExecHandler(s network.Stream) {
	sreader := bufio.NewReader(s)

	//linebuff := make()
	for {
		linestr, _ := sreader.ReadString('\n')
		println(linestr)
		/*sreader
		_, err := io.ReadFull(s, buf)
		if err != nil {
			errCh <- err
			return
		}

		_, err = s.Write(buf)
		if err != nil {
			errCh <- err
			return
		}*/

		//timer.Reset(pingTimeout)
	}
}
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
	echo := NewExecService(node)
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
