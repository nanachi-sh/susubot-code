package connector

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/netip"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nanachi-sh/susubot-code/connector/LLOneBot/protos/connector"
)

type Connector struct {
	addr       net.Addr
	conn       *websocket.Conn
	now        context.Context
	now_cancel context.CancelFunc
	readWG     sync.WaitGroup
	closeLock  sync.Mutex
	closed     chan struct{}
	reting     int
}

type responseBuf struct{}

func New() *Connector {
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan struct{})
	close(ch)
	c := &Connector{
		now:        ctx,
		now_cancel: cancel,
		closed:     ch,
	}
	// go c.test()
	return c
}

func (c *Connector) Connect(req *connector.ConnectRequest) ([]byte, error) {
	select {
	case <-c.closed:
	default:
		return nil, errors.New("已连接服务器")
	}
	dialer := &websocket.Dialer{}
	addr := req.Addr
	port := req.Port
	if port <= 0 || port > 65535 {
		return nil, errors.New("服务器端口范围不正确")
	}
	if addr == "" {
		return nil, errors.New("服务器地址为空")
	}
	if ip := net.ParseIP(addr); ip != nil { //为IP
		addr = ip.String()
	} else if ok, err := regexp.MatchString(`^[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+$`, addr); ok { //为域名
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		ips, err := net.DefaultResolver.LookupIP(ctx, "ip", addr)
		if err != nil {
			return nil, err
		}
		if len(ips) == 0 {
			return nil, errors.New("域名无IP返回")
		}
		addr = ips[0].String()
	} else { //若无错误，为未知
		if err != nil {
			return nil, err
		} else {
			return nil, errors.New("服务器地址设置有误")
		}
	}
	netipaddr, err := netip.ParseAddr(addr)
	if err != nil {
		return nil, err
	}
	c.addr = net.TCPAddrFromAddrPort(netip.AddrPortFrom(netipaddr, uint16(port)))
	headers := make(http.Header)
	if req.Token != nil {
		headers.Add("Authorization", "Bearer "+*req.Token)
	}
	conn, _, err := dialer.DialContext(context.Background(), fmt.Sprintf("ws://%v", c.addr.String()), headers)
	if err != nil {
		return nil, err
	}
	c.conn = conn
	c.closed = make(chan struct{})
	if err := c.readAndwrite(); err != nil {
		go c.close()
		return nil, err
	}
	buf := c.readLast()
	c.readReset()
	go c.readToEnd()
	return buf, nil
}

// 连接后调用，连接结束或发生错误自行退出
func (c *Connector) readToEnd() {
	for {
		select {
		case <-c.closed:
			return
		default:
		}
		if err := c.readAndwrite(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			if err := c.close(); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			return
		}
	}
}

func (c *Connector) read() ([]byte, error) {
	select {
	case <-c.closed:
		return nil, errors.New("服务器已断开或未连接")
	default:
	}
	_, r, err := c.conn.NextReader()
	if err != nil {
		return nil, err
	}
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (c *Connector) write(buf []byte) {
	c.now = context.WithValue(c.now, responseBuf{}, buf)
	c.now_cancel()
}

func (c *Connector) readAndwrite() error {
	buf, err := c.read()
	if err != nil {
		return err
	}
	//确保读取返回已结束
	select {
	case <-c.now.Done(): //若正在返回则等待
		c.readWG.Wait()
	default:
	}
	c.write(buf)
	return nil
}

func (c *Connector) readReset() {
	c.now = context.WithoutCancel(c.now)
	c.now, c.now_cancel = context.WithCancel(c.now)
}

// func (c *Connector) test() {
// 	for {
// 		time.Sleep(time.Second * 2)
// 		select {
// 		case <-c.now.Done(): //若正在返回则等待
// 			c.readLock.Lock()
// 			c.readLock.Unlock()
// 		default:
// 		}
// 		c.write([]byte{byte(rand.IntN(256)), byte(rand.IntN(256)), byte(rand.IntN(256)), byte(rand.IntN(256)), byte(rand.IntN(256)), byte(rand.IntN(256)), byte(rand.IntN(256)), byte(rand.IntN(256)), byte(rand.IntN(256))})
// 	}
// }

// 幂等
func (c *Connector) close() error {
	c.closeLock.Lock()
	defer c.closeLock.Unlock()
	select {
	case <-c.closed:
		return errors.New("Closed")
	default:
	}
	defer func() {
		close(c.closed)
		c.conn = nil
		c.addr = nil
	}()
	if err := c.conn.Close(); err != nil {
		return err
	}
	return nil
}

func (c *Connector) Read(a int64) ([]byte, error) {
	//检查是否在返回过程中
	fmt.Printf("%v %v: Waiting\n", a, time.Now().Format("2006-01-02 15:04:05.000"))
	c.readWG.Wait()
	fmt.Printf("%v %v: WaitDone\n", a, time.Now().Format("2006-01-02 15:04:05.000"))
	//
	fmt.Printf("%v %v: WaitAdd\n", a, time.Now().Format("2006-01-02 15:04:05.000"))
	c.readWG.Add(1)
	defer c.readWG.Done()
	//检查是否为最后一个
	defer func() {
		if c.reting == 0 {
			c.readReset()
		}
	}()
	c.reting++
	fmt.Printf("%v %v: Block\n", a, time.Now().Format("2006-01-02 15:04:05.000"))
	select {
	case <-c.closed:
		fmt.Printf("%v %v: closed\n", a, time.Now().Format("2006-01-02 15:04:05.000"))
		return nil, errors.New("连接已断开或未连接")
	case <-c.now.Done():
		defer func() { c.reting-- }()
		fmt.Printf("%v %v: response return\n", a, time.Now().Format("2006-01-02 15:04:05.000"))
		if buf := c.readLast(); buf == nil {
			return nil, errors.New("异常错误")
		} else {
			return buf, nil
		}
	}
}

func (c *Connector) ReadLast() []byte {
	return c.readLast()
}

func (c *Connector) readLast() []byte {
	switch x := c.now.Value(responseBuf{}).(type) {
	case []byte:
		return x
	default:
		return nil
	}
}

func (c *Connector) Write(buf []byte) error {
	select {
	case <-c.closed:
		return errors.New("连接已断开或未连接")
	default:
	}
	if err := c.conn.WriteMessage(websocket.TextMessage, buf); err != nil {
		return err
	}
	return nil
}

func (c *Connector) Close() error {
	return c.close()
}
