package gateway

import (
	"fmt"
	"log"
	"net"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"

	"github.com/qingw1230/plato/common/config"
	"golang.org/x/sys/unix"
)

var ep *ePool

// tcpNum 允许接入的最大 tcp 连接数
var tcpNum int32

// ePool 管理运行的 ep 模型
type ePool struct {
	eChan  chan *connection // 存储获取到的新连接
	tables sync.Map         // 连接文件描述符与连接的映射
	eSize  int              // 创建的 ep 模型的数量
	done   chan struct{}    // 用于终止网关

	ln *net.TCPListener                 // 监听套接字
	fn func(c *connection, ep *epoller) // 事件发生时的回调
}

// initEpoll 初始化网关
func initEpoll(ln *net.TCPListener, f func(c *connection, ep *epoller)) {

	setLimit()
	ep = newEPool(ln, f)
	ep.createAcceptProcess()
	ep.startEpoll()
}

func newEPool(ln *net.TCPListener, cb func(c *connection, ep *epoller)) *ePool {
	return &ePool{
		eChan:  make(chan *connection, config.GetGatewayEpollerChanNum()),
		tables: sync.Map{},
		eSize:  config.GetGatewayEpollerNum(),
		done:   make(chan struct{}),
		ln:     ln,
		fn:     cb,
	}
}

// createAcceptProcess 创建与 CPU 核数相同的协程用于获取连接
func (e *ePool) createAcceptProcess() {
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				conn, e := e.ln.AcceptTCP()
				if e != nil {
					fmt.Errorf("accept err: %v\n", e)
				}
				// 限流熔断，超过限制后拒绝新连接
				if !checkTCP() {
					conn.Close()
					continue
				}
				setTCPConfig(conn)

				c := connection{
					conn: conn,
					fd:   socketFD(conn),
				}
				ep.addTask(&c)
			}
		}()
	}
}

func (e *ePool) startEpoll() {
	for i := 0; i < e.eSize; i++ {
		go e.startEProc()
	}
}

// startEProc 各个 ep 开始处理事件
func (e *ePool) startEProc() {
	ep, err := newEpoller()
	if err != nil {
		panic(err)
	}

	// 从 eChan 中取走连接，并开始监听
	go func() {
		for {
			select {
			case <-e.done:
				return
			case conn := <-e.eChan:
				addTCPNum()
				// TODO(qingw1230): 连续多次发起大量连接请求，第一次可以快速处理，后续处理速度慢
				if err := ep.add(conn); err != nil {
					fmt.Printf("failed to add connection %v\n", err)
					conn.Close()
					continue
				}
				fmt.Printf("EpollerPoll new connection[%v] tcpSize:%d\n", conn.RemoteAddr(), tcpNum)
			}
		}
	}() // go func() {

	for {
		select {
		case <-e.done:
			return
		default:
			connections, err := ep.wait(200)

			if err != nil && err != syscall.EINTR {
				fmt.Printf("failed to epoll wait %v\n", err)
				continue
			}
			for _, conn := range connections {
				if conn == nil {
					break
				}
				e.fn(conn, ep)
			}
		}
	} // for {
}

// addTask 将获取到的新连接添加到 chan 中
func (e *ePool) addTask(c *connection) {
	e.eChan <- c
}

// epoller 底层 ep 模型控制器
type epoller struct {
	fd int // ep 模型的文件描述符
}

// newEpoller 创建一个 ep 模型
func newEpoller() (*epoller, error) {
	fd, err := unix.EpollCreate1(0)
	if err != nil {
		return nil, err
	}
	return &epoller{
		fd: fd,
	}, nil
}

// add 向 ep 添加新连接
func (e *epoller) add(conn *connection) error {
	fd := conn.fd
	ev := &unix.EpollEvent{
		Events: unix.EPOLLIN | unix.EPOLLHUP,
		Fd:     int32(fd),
	}
	err := unix.EpollCtl(e.fd, syscall.EPOLL_CTL_ADD, fd, ev)
	if err != nil {
		return err
	}
	ep.tables.Store(fd, conn)
	return nil
}

// remove 从 ep 中移除指定连接
func (e *epoller) remove(c *connection) error {
	subTCPNum()
	fd := c.fd
	err := unix.EpollCtl(e.fd, syscall.EPOLL_CTL_DEL, fd, nil)
	if err != nil {
		return err
	}
	ep.tables.Delete(fd)
	return nil
}

// wait 等待 ep 中的事件发生，返回有事件发生的连接列表
func (e *epoller) wait(msec int) ([]*connection, error) {
	events := make([]unix.EpollEvent, config.GetGatewayEpollWaitQueueSize())
	n, err := unix.EpollWait(e.fd, events, msec)
	if err != nil {
		return nil, err
	}

	var connections []*connection
	for i := 0; i < n; i++ {
		if conn, ok := ep.tables.Load(int(events[i].Fd)); ok {
			connections = append(connections, conn.(*connection))
		}
	}
	return connections, nil
}

func socketFD(conn *net.TCPConn) int {
	tcpConn := reflect.Indirect(reflect.ValueOf(*conn)).FieldByName("conn")
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")
	return int(pfdVal.FieldByName("Sysfd").Int())
}

// setLimit 将 go 进程可以打开的文件数量设为最大值
func setLimit() {
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
	rLimit.Cur = rLimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
	log.Printf("set cur limit: %d", rLimit.Cur)
}

func addTCPNum() {
	atomic.AddInt32(&tcpNum, 1)
}

func getTCPNum() int32 {
	return atomic.LoadInt32(&tcpNum)
}

func subTCPNum() {
	atomic.AddInt32(&tcpNum, -1)
}

// checkTCP 检查是否还能创建 TCP 连接
func checkTCP() bool {
	num := getTCPNum()
	maxTCPNum := config.GetGatewayMaxTCPNum()
	return num <= maxTCPNum
}

// setTCPConfig 为 TCP 连接设置保活
func setTCPConfig(c *net.TCPConn) {
	c.SetKeepAlive(true)
}
