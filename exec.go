package p001

import (
	"bufio"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"golang.org/x/net/context"
	"os/exec"
	"strings"
)

const (
	ID          = "/ipfs/exec/1.0.0"
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

func execError(err error) chan string {
	ch := make(chan string, 1)
	ch <- "error"
	close(ch)
	return ch
}
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
		case cmdstr := <-cmd:
			doexec(s, cmdstr)
			//println(cmd)
		case <-ctx.Done():
			return
		}
	}()
	return out
}
func doexec(s network.Stream, cmd string) (int, error) {
	sreader := bufio.NewReader(s)
	_, err := s.Write([]byte(cmd))
	_, err = sreader.ReadString('\n')
	if err != nil {
		return 0, err
	}
	return 0, nil
}
func (p *ExecService) ExecHandler(s network.Stream) {
	sreader := bufio.NewReader(s)

	//linebuff := make()
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
