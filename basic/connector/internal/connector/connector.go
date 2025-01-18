package connector

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/netip"
	"regexp"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nanachi-sh/susubot-code/basic/connector/internal/types"
	connector_pb "github.com/nanachi-sh/susubot-code/basic/connector/pkg/protos/connector"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	conn       *websocket.Conn
	now        context.Context
	now_cancel context.CancelFunc

	readBlock sync.RWMutex
	readWait  sync.RWMutex
	writeLock sync.Mutex

	closeLock sync.Mutex
	closed    chan struct{} //未连接与Close后都视为close

	timewait_Close = time.Second * 5
)

type responseMerge struct{}

type Request struct {
	logger logx.Logger
}

func NewRequest(l logx.Logger) *Request {
	return &Request{
		logger: l,
	}
}

func reset() {
	conn = nil
	now, now_cancel = context.WithCancel(context.Background())
	readBlock = sync.RWMutex{}
	readWait = sync.RWMutex{}
	writeLock = sync.Mutex{}
	closeLock = sync.Mutex{}
	closed = make(chan struct{})
	close(closed)
}

func init() {
	reset()
}

func (r *Request) Connect(in *connector_pb.ConnectRequest) (*connector_pb.ConnectResponse, error) {
	if in.Addr == "" || in.Port == 0 {
		return &connector_pb.ConnectResponse{}, status.Error(codes.InvalidArgument, "")
	}
	if in.Port <= 0 || in.Port > 65535 {
		return &connector_pb.ConnectResponse{
			Body: &connector_pb.ConnectResponse_Err{
				Err: connector_pb.Errors_PortError,
			},
		}, nil
	}
	token := ""
	if in.Token != nil {
		token = *in.Token
	}
	resp, serr := connect(r.logger, types.ConnectRequest{
		Addr:  in.Addr,
		Port:  int(in.Port),
		Token: token,
	})
	if serr != nil {
		return &connector_pb.ConnectResponse{
			Body: &connector_pb.ConnectResponse_Err{
				Err: *serr,
			},
		}, nil
	}
	return &connector_pb.ConnectResponse{
		Body: &connector_pb.ConnectResponse_Buf{Buf: resp},
	}, nil
}

func (r *Request) Read(in *connector_pb.Empty, stream connector_pb.Connector_ReadServer) error {
	for {
		buf, serr := readFromContext(r.logger)
		if serr != nil {
			stream.Send(&connector_pb.ReadResponse{
				Body: &connector_pb.ReadResponse_Err{
					Err: *serr,
				},
			})
			return nil
		}
		if err := stream.Send(&connector_pb.ReadResponse{
			Body: &connector_pb.ReadResponse_Buf{
				Buf: buf,
			},
		}); err != nil {
			return err
		}
	}
}

func (r *Request) Write(in *connector_pb.WriteRequest) (*connector_pb.BasicResponse, error) {
	writeLock.Lock()
	defer writeLock.Unlock()
	if len(in.Buf) == 0 {
		return &connector_pb.BasicResponse{}, status.Error(codes.InvalidArgument, "")
	}
	if serr := writeToServer(r.logger, in.Buf); serr != nil {
		return &connector_pb.BasicResponse{
			Err: serr,
		}, nil
	}
	return &connector_pb.BasicResponse{}, nil
}

func (r *Request) Close(in *connector_pb.Empty) (*connector_pb.BasicResponse, error) {
	return &connector_pb.BasicResponse{
		Err: connectClose(r.logger),
	}, nil
}

func isClose() bool {
	select {
	case <-closed:
		return true
	default:
		return false
	}
}

func writeToServer(logger logx.Logger, buf []byte) *connector_pb.Errors {
	if isClose() {
		return connector_pb.Errors_Closed.Enum()
	}
	if err := conn.WriteMessage(websocket.TextMessage, buf); err != nil {
		logger.Error(err)
		return connector_pb.Errors_Undefined.Enum()
	}
	return nil
}

func connect(logger logx.Logger, req types.ConnectRequest) ([]byte, *connector_pb.Errors) {
	if !isClose() {
		return nil, connector_pb.Errors_Connected.Enum()
	}
	logger.Debug("开始连接机器人核心")
	dialer := &websocket.Dialer{}
	var addr netip.Addr
	ok := false
	for {
		a := ""
		if ip := net.ParseIP(req.Addr); ip != nil { //为IP
			a = ip.String()
		} else if ok, err := regexp.MatchString(`^[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+$`, req.Addr); ok { //为域名
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()
			ips, err := net.DefaultResolver.LookupIP(ctx, "ip", req.Addr)
			if err != nil {
				logger.Error(err)
				break
			}
			if len(ips) == 0 {
				logger.Error("无IP返回")
				break
			}
			a = ips[0].String()
		} else { //若无错误，为未知
			if err != nil {
				logger.Error(err)
				break
			} else {
				logger.Error("地址格式有误")
				break
			}
		}
		netipa, err := netip.ParseAddr(a)
		if err != nil {
			logger.Error(err)
			return nil, connector_pb.Errors_AddrError.Enum()
		}
		addr = netipa
		ok = true
		break
	}
	if !ok {
		logger.Debug("解析机器人核心地址失败")
		return nil, connector_pb.Errors_AddrError.Enum()
	}
	logger.Debug("解析机器人核心地址成功")
	tcpaddr := net.TCPAddrFromAddrPort(netip.AddrPortFrom(addr, uint16(req.Port)))
	headers := make(http.Header)
	if req.Token != "" {
		headers.Add("Authorization", "Bearer "+req.Token)
	}
	c, _, err := dialer.DialContext(context.Background(), fmt.Sprintf("ws://%v", tcpaddr.String()), headers)
	if err != nil {
		logger.Error(err)
		return nil, connector_pb.Errors_DialError.Enum()
	}
	conn = c
	closed = make(chan struct{})
	if err := readAndwrite(logger); err != nil {
		go connectClose(logger)
		logger.Error(err)
		return nil, connector_pb.Errors_DialError.Enum()
	}
	last := readLast()
	readFromContext_Reset()
	go func() {
		if serr := readToEnd(logger); serr != nil {
			writeToContext(merge{
				createTime: time.Now(),
				err:        serr,
			})
		}
	}()
	return last.buf, nil
}

// 连接后调用，连接结束或发生错误自行退出
func readToEnd(logger logx.Logger) *connector_pb.Errors {
	for {
		if isClose() {
			return connector_pb.Errors_Closed.Enum()
		}
		if serr := readAndwrite(logger); serr != nil {
			go connectClose(logger)
			return serr
		}
	}
}

func readFromServer(logger logx.Logger) ([]byte, *connector_pb.Errors) {
	if isClose() {
		return nil, connector_pb.Errors_Closed.Enum()
	}
	_, r, err := conn.NextReader()
	if err != nil {
		logger.Error(err)
		return nil, connector_pb.Errors_Undefined.Enum()
	}
	buf, err := io.ReadAll(r)
	if err != nil {
		logger.Error(err)
		return nil, connector_pb.Errors_Undefined.Enum()
	}
	return buf, nil
}

func writeToContext(m merge) {
	now = context.WithValue(now, responseMerge{}, &m)
	now_cancel()
}

func readAndwrite(logger logx.Logger) *connector_pb.Errors {
	buf, serr := readFromServer(logger)
	if serr != nil {
		return serr
	}
	//等待读取返回结束
	readWait.RLock()
	readWait.RUnlock()
	writeToContext(merge{
		buf:        buf,
		createTime: time.Now(),
	})
	return nil
}

func readFromContext_Reset() {
	now = context.WithoutCancel(now)
	now, now_cancel = context.WithCancel(now)
}

// 幂等
func connectClose(logger logx.Logger) *connector_pb.Errors {
	closeLock.Lock()
	defer closeLock.Unlock()
	if isClose() {
		return connector_pb.Errors_Closed.Enum()
	}
	ch := make(chan struct{})
	ctx, cancel := context.WithTimeout(context.Background(), timewait_Close)
	defer cancel()
	go func() {
		defer close(ch)
		if err := conn.Close(); err != nil {
			logger.Error(err)
		}
	}()
	select {
	case <-ch:
	case <-ctx.Done():
	}
	close(closed)
	conn = nil
	return nil
}

func readFromContext(logger logx.Logger) ([]byte, *connector_pb.Errors) {
	if isClose() {
		return nil, connector_pb.Errors_Closed.Enum()
	}
	//若等待队列关闭，则加入并阻塞
	readWait.RLock()
	readWait.RUnlock()
	//进入阻塞队列
	readBlock.RLock()
	//检查阻塞队列是否为空
	defer func() {
		if isClose() {
			return
		}
		if readBlock.TryLock() { //阻塞队列为空
			//打开等待队列
			readWait.Unlock()
			//重置
			readFromContext_Reset()
			//等待Wait队列空
			readWait.Lock()
			readWait.Unlock()
			//打开阻塞队列
			readBlock.Unlock()
		}
	}()
	defer readBlock.RUnlock()
	select {
	case <-closed:
		return nil, connector_pb.Errors_Closed.Enum()
	case <-now.Done():
		//第一个会话负责关闭等待队列
		readWait.TryLock()
		if last := readLast(); last == nil {
			logger.Error("读取上一个Buffer为nil")
			return nil, connector_pb.Errors_Undefined.Enum()
		} else {
			return last.buf, nil
		}
	}
}

type merge struct {
	buf        []byte
	createTime time.Time
	err        *connector_pb.Errors
}

func readLast() *merge {
	if x, ok := now.Value(responseMerge{}).(*merge); ok {
		return x
	} else {
		return nil
	}
}
