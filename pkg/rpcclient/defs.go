package rpcclient

import "encoding/json"

const (
	SessionKeepaliveTime = 10
)

type ErrorCode int

const (
	ErrorNormal ErrorCode = iota
	ErrorKeepaliveFail
	ErrConnSendException
	ErrConnRecvException
)

type loginValue struct {
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
	Random   string `json:"random,omitempty"`
}

type jsonRequest struct {
	Domain string          `json:"domain"`
	Key    string          `json:"key"`
	Value  json.RawMessage `json:"value"`
	Id     uint32          `json:"id"`
}

type jsonResponse struct {
	Domain string          `json:"domain"`
	Key    string          `json:"key,omitempty"`
	Value  json.RawMessage `json:"value,omitempty"`
	Id     uint32          `json:"id"`
	RetVal int32           `json:"retval"`
	ErrMsg string          `json:"errmsg,omitempty"`
}

type rpcReqPkt struct {
	jsonReq   *jsonRequest
	binaryReq []byte
	replyCh   chan *rpcResPkt
}

type rpcResPkt struct {
	jsonRes   *jsonResponse
	binaryReq []byte
}

type RpcReply struct {
	RetVal     int32
	ErrMsg     string
	JsonResp   []byte
	BinaryResp []byte
}

type rpcSubReply struct {
	RetVal int32
	ErrMsg string
	SubId  int
	MsgCh  chan []byte
}

type rpcSubValue struct {
	NotifyId int `json:"notifyId"`
}

type requestProcessing struct {
	reqId   int
	replyCh chan *rpcResPkt
}

type replyProcessing struct {
	reqId  int
	resPkt *rpcResPkt
}

type OperationMode int

const (
	OperationOff OperationMode = iota
	OperationOn
)

type subscribeProcessing struct {
	notifyId int
	notifyCh chan []byte
	mode     OperationMode
}

type notificationProcessing struct {
	notifyId  int
	notifyMsg []byte
}

type UserInfo struct {
	Name string `json:"name"`
	No   string `json:"no"`
}
