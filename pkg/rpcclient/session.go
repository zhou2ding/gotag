package rpcclient

import (
	"encoding/json"
	"gotag/pkg/internal/rpctool"
	"net"
	"time"
)

func createSession() *session {
	return &session{}
}

type session struct {
	conn net.Conn
}

func (s *session) open(addr string, timeout int) error {
	//建立连接
	conn, err := net.DialTimeout("tcp", addr, time.Millisecond*time.Duration(timeout))
	if err != nil {
		return err
	}

	s.conn = conn

	return nil
}

func (s *session) close() {
	if s.conn != nil {
		_ = s.conn.Close()
	}
}

//传入的JsonRequest的字段中Domain可选，Key和Value必填，Id字段的值忽略
func (s *session) rpcCall(jsonReq *jsonRequest, binaryReq []byte) error {
	val, err := json.Marshal(jsonReq)
	if err != nil {
		return err
	}

	pkts := rpctool.GetRPCMsgSplit().SplitPacket(val, binaryReq)

	for _, pkt := range pkts {
		_, err = s.writeToConn(pkt)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *session) rpcRecv() (*jsonResponse, []byte, error) {
	builder := rpctool.CreateRPCMsgBuilder(s.conn)
	payload, err := builder.DoRead()
	if err != nil {
		return nil, nil, err
	}

	var jsonRes jsonResponse
	err = json.Unmarshal(payload.JsonRequest, &jsonRes)
	if err != nil {
		return nil, nil, err
	}

	return &jsonRes, payload.BinaryRequest, nil
}

func (s *session) writeToConn(data []byte) (int, error) {
	dataLen := len(data)
	nWrites := 0
	for {
		n, err := s.conn.Write(data[nWrites:])
		if err != nil {
			return n, err
		}

		nWrites += n
		if nWrites >= dataLen {
			break
		}
	}

	return nWrites, nil
}
