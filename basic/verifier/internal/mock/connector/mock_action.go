package mock_connectorclient

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"regexp"
	"sync"

	"github.com/nanachi-sh/susubot-code/basic/verifier/pkg/protos/connector"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
)

type readStream struct {
	grpc.ServerStreamingClient[connector.ReadResponse]
	update    context.Context
	readBlock sync.RWMutex
	waitBlock sync.RWMutex
}

var (
	readstream *readStream

	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
)

type myret struct{}

func init() {
	readstream = &readStream{
		update: context.Background(),
	}
	readstream.readBlock.Lock()
}

func DefaultMock() *MockConnector {
	ctrl := gomock.NewController(nil)
	mock := NewMockConnector(ctrl)

	mock.EXPECT().Write(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(_ any, in *connector.WriteRequest, _ ...any) (*connector.BasicResponse, error) {
		if in == nil {
			return nil, errors.New("")
		}
		gfl_action, _ := regexp.Match(`{"action":"get_friend_list","echo":".*"}`, in.Buf)
		sfm_action, _ := regexp.Match(`{"action":"send_private_msg","params":{"user_id":"[0-9].+","message":\[.+\]},"echo":".*"}`, in.Buf)
		m := make(map[string]any)
		if err := json.Unmarshal(in.Buf, &m); err != nil {
			return nil, err
		}
		echo := ""
		if e, ok := m["echo"].(string); ok {
			echo = e
		} else {
			return nil, errors.New("no found echo")
		}
		switch {
		default:
			return nil, errors.New("no match")
		case gfl_action:
			m := make(map[string]any)
			if echo != "" {
				m["echo"] = echo
			}
			m["status"] = "ok"
			m["retcode"] = 0
			m["message"] = ""
			m["wording"] = ""
			m["data"] = []struct {
				UserId   int64  `json:"user_id"`
				NickName string `json:"nickname"`
				Remark   string `json:"remark"`
				Sex      string `json:"sex"`
				Level    int    `json:"level"`
			}{
				struct {
					UserId   int64  "json:\"user_id\""
					NickName string "json:\"nickname\""
					Remark   string "json:\"remark\""
					Sex      string "json:\"sex\""
					Level    int    "json:\"level\""
				}{
					UserId:   100000,
					NickName: "",
					Remark:   "",
					Sex:      "",
					Level:    0,
				},
			}
			ret, err := json.Marshal(m)
			if err != nil {
				return nil, err
			}
			readstream.emit(&connector.ReadResponse{
				Body: &connector.ReadResponse_Buf{Buf: ret},
			})
			return &connector.BasicResponse{}, nil
		case sfm_action:
			logger.Println(string(in.Buf))
			m := make(map[string]any)
			if echo != "" {
				m["echo"] = echo
			}
			m["status"] = "ok"
			m["retcode"] = 0
			m["message"] = ""
			m["wording"] = ""
			m["data"] = struct {
				MessageId int64 `json:"message_id"`
			}{
				0,
			}
			ret, err := json.Marshal(m)
			if err != nil {
				return nil, err
			}
			readstream.emit(&connector.ReadResponse{
				Body: &connector.ReadResponse_Buf{Buf: ret},
			})
			return &connector.BasicResponse{}, nil
		}
	}).AnyTimes()
	mock.EXPECT().Read(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(_ any, _ *connector.Empty, _ ...any) (connector.Connector_ReadClient, error) {
		return readstream, nil
	}).AnyTimes()
	return mock
}

func (rs *readStream) emit(r *connector.ReadResponse) {
	rs.update = context.WithValue(rs.update, myret{}, r)
	rs.readBlock.Unlock()
}

func (rs *readStream) reset() {
	rs.update = context.WithValue(rs.update, myret{}, nil)
}

func (rs *readStream) Recv() (*connector.ReadResponse, error) {
	// 等待所有会话读取完毕
	rs.waitBlock.RLock()
	rs.waitBlock.RUnlock()
	// 等待新的流信息
	rs.readBlock.RLock()
	// 保证顺序，该defer为最后执行
	defer func() {
		if rs.readBlock.TryLock() { //最后一个会话，顺带锁上读队列
			// 切换等待队列为入出
			rs.waitBlock.Unlock()
			// 重置ctx内容
			rs.reset()
		}
	}()
	defer rs.readBlock.RUnlock()
	//切换等待队列为仅入
	rs.waitBlock.TryLock()
	if ret, ok := rs.update.Value(myret{}).(*connector.ReadResponse); ok {
		return ret, nil
	} else {
		return nil, errors.New("undefined error")
	}
}
