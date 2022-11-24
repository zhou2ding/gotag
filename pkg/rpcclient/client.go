package rpcclient

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"gotag/pkg/idgen"
	"gotag/pkg/internal/rpctool"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Client struct {
	session    *session
	sendCh     chan *rpcReqPkt
	reqMap     map[int]chan *rpcResPkt
	subMap     map[int]chan []byte
	notifyCh   chan interface{}
	stopCh     chan struct{}
	notifyOnce sync.Once
	logger     *zap.Logger
	errCode    ErrorCode
	errInfo    error

	commonProcessCh    chan interface{}
	subscribeProcessCh chan interface{}

	flowWatched int32

	bindInfo interface{}
}

func CreateRPCClient(logger *zap.Logger) *Client {
	return &Client{
		sendCh:             make(chan *rpcReqPkt, 10),
		reqMap:             make(map[int]chan *rpcResPkt),
		stopCh:             make(chan struct{}),
		subMap:             make(map[int]chan []byte),
		commonProcessCh:    make(chan interface{}, 20),
		subscribeProcessCh: make(chan interface{}, 20),
		logger:             logger,
	}
}

func (c *Client) Open(addr string, timeout int) error {
	s := createSession()
	err := s.open(addr, timeout)
	if err != nil {
		return err
	}

	c.session = s

	go c.send()      //此协程将在发送出错、调用Close()或保活失败的情况下退出
	go c.receive()   //此协程将在接收出错，或调用Close()导致session.close()被调用时退出
	go c.doCore()    //此协程将在读写错误、调用Close()或保活失败的情况下退出
	go c.keepalive() //此协程将在读写错误、调用Close()或保活失败的情况下退出

	return nil
}

func (c *Client) Close() {
	if c.session != nil {
		c.session.close()
	}
	c.reportException(ErrorNormal, nil)
}

func (c *Client) Login(username string, password string, timeout int) error {
	//first login
	loginVal := loginValue{
		User: username,
	}
	val, err := json.Marshal(&loginVal)
	if err != nil {
		return err
	}

	reply, err := c.Call("common", "System.login", val, nil, timeout)
	if err != nil {
		return err
	}

	if reply.RetVal != 0 {
		return errors.New(reply.ErrMsg)
	}

	err = json.Unmarshal(reply.JsonResp, &loginVal)
	if err != nil {
		return err
	}

	//second login
	c.logger.Info("login info", zap.String("username", username), zap.String("pwd", password), zap.String("random", loginVal.Random))
	encPwd := rpctool.GetEncryptHelperInstance().Encrypt(username, password, loginVal.Random)
	loginVal.Password = encPwd
	val, err = json.Marshal(loginVal)
	if err != nil {
		return err
	}

	reply, err = c.Call("common", "System.login", val, nil, timeout)
	if err != nil {
		return err
	}

	if reply.RetVal != 0 {
		return errors.New(reply.ErrMsg)
	}

	return nil
}

func (c *Client) NotifyFail(notifyCh chan interface{}, bind interface{}) {
	c.notifyCh = notifyCh
	c.bindInfo = bind
}

func (c *Client) Subscribe(domain string, key string, timeout int) (*rpcSubReply, error) {
	reply, err := c.Call(domain, key, nil, nil, timeout)
	if err != nil {
		return nil, err
	}
	if reply.RetVal != 0 {
		r := &rpcSubReply{
			RetVal: reply.RetVal,
			ErrMsg: reply.ErrMsg,
		}
		return r, nil
	}

	var subValue rpcSubValue
	err = json.Unmarshal(reply.JsonResp, &subValue)
	if err != nil {
		return nil, err
	}

	ch := make(chan []byte, 10)
	c.subscribeProcessCh <- &subscribeProcessing{
		notifyId: subValue.NotifyId,
		notifyCh: ch,
		mode:     OperationOn,
	}

	return &rpcSubReply{
		RetVal: 0,
		SubId:  subValue.NotifyId,
		MsgCh:  ch,
	}, nil
}

func (c *Client) Unsubscribe(domain string, key string, subId int, timeout int) error {
	val := &rpcSubValue{
		NotifyId: subId,
	}
	jVal, _ := json.Marshal(val)

	reply, err := c.Call(domain, key, jVal, nil, 2000)
	if err != nil {
		return err
	}

	if reply.RetVal != 0 {
		return errors.New(reply.ErrMsg)
	}

	c.subscribeProcessCh <- &subscribeProcessing{
		notifyId: subId,
		mode:     OperationOff,
	}

	return nil
}

func (c *Client) Call(domain string, key string, jsonReq []byte, binaryReq []byte, timeout int) (*RpcReply, error) {
	jsReq := &jsonRequest{
		Domain: domain,
		Key:    key,
		Value:  json.RawMessage(jsonReq),
		Id:     uint32(idgen.GetIdGenerator().GetId()),
	}
	now := time.Now()
	defer func() {
		cost := time.Since(now)
		c.logger.Info("Call", zap.String("key", key), zap.Int64("cost", int64(cost/time.Millisecond)))
	}()

	replyCh := make(chan *rpcResPkt, 1)
	select {
	case c.sendCh <- &rpcReqPkt{jsonReq: jsReq, binaryReq: binaryReq, replyCh: replyCh}:
	case <-c.stopCh:
		return nil, errors.New("client has been closed, or errors occur")
	}

	select {
	case <-time.After(time.Millisecond * time.Duration(timeout)):
		return nil, errors.New("timeout")
	case <-c.stopCh:
		return nil, errors.New("client has been closed, or errors occur")
	case resPkt, ok := <-replyCh:
		if !ok {
			return nil, errors.New("client error")
		}
		if resPkt.jsonRes.RetVal != 0 {
			reply := RpcReply{
				RetVal: resPkt.jsonRes.RetVal,
				ErrMsg: resPkt.jsonRes.ErrMsg,
			}
			return &reply, nil
		}

		reply := RpcReply{
			RetVal:     0,
			JsonResp:   resPkt.jsonRes.Value,
			BinaryResp: resPkt.binaryReq,
		}

		return &reply, nil
	}
}

func (c *Client) send() {
	for {
		select {
		case jsReq := <-c.sendCh:
			c.commonProcessCh <- &requestProcessing{
				reqId:   int(jsReq.jsonReq.Id),
				replyCh: jsReq.replyCh,
			}

			err := c.session.rpcCall(jsReq.jsonReq, jsReq.binaryReq)
			if err != nil {
				c.reportException(ErrConnSendException, err)
			} else {
				atomic.StoreInt32(&c.flowWatched, 1)
			}
		case <-c.stopCh:
			return
		}
	}
}

func (c *Client) receive() {
	for {
		jsonRes, binaryRes, err := c.session.rpcRecv()
		if err != nil {
			c.reportException(ErrConnRecvException, err)
			break
		} else {
			atomic.StoreInt32(&c.flowWatched, 1)
		}
		if jsonRes.Domain == "common" {
			pkt := &rpcResPkt{
				jsonRes:   jsonRes,
				binaryReq: binaryRes,
			}
			c.commonProcessCh <- &replyProcessing{
				reqId:  int(jsonRes.Id),
				resPkt: pkt,
			}
		} else {
			c.subscribeProcessCh <- &notificationProcessing{
				notifyId:  int(jsonRes.Id),
				notifyMsg: jsonRes.Value,
			}
		}
	}
}

func (c *Client) doCore() {
	for {
		select {
		case iface := <-c.commonProcessCh:
			switch val := iface.(type) {
			case *requestProcessing:
				c.reqMap[val.reqId] = val.replyCh
			case *replyProcessing:
				if ch, ok := c.reqMap[val.reqId]; ok {
					ch <- val.resPkt
					delete(c.reqMap, val.reqId)
				}
			}
		case iface := <-c.subscribeProcessCh:
			switch val := iface.(type) {
			case *subscribeProcessing:
				if val.mode == OperationOn {
					c.subMap[val.notifyId] = val.notifyCh
				} else {
					close(c.subMap[val.notifyId])
					delete(c.subMap, val.notifyId)
				}
			case *notificationProcessing:
				if ch, ok := c.subMap[val.notifyId]; ok {
					ch <- val.notifyMsg
				}
			}
		case <-c.stopCh:
			if len(c.reqMap) != 0 {
				for _, ch := range c.reqMap {
					close(ch)
				}
				c.reqMap = make(map[int]chan *rpcResPkt)
			}
			if len(c.subMap) != 0 {
				for _, mch := range c.subMap {
					close(mch)
				}
				c.subMap = make(map[int]chan []byte)
			}
			return
		}
	}
}

func (c *Client) keepalive() {
	ticker := time.NewTicker(time.Duration(SessionKeepaliveTime) * time.Second)
	defer ticker.Stop()
	count := 0
	for {
		select {
		case <-ticker.C:
			if atomic.SwapInt32(&c.flowWatched, 0) > 0 { //刚过去的保活周期内有数据交互，不需要再发送 保活包
				continue
			}

			_, err := c.Call("common", "System.heartbeat", nil, nil, 3000)
			if err == nil {
				count = 0
			} else {
				c.logger.Warn("keepalive failed", zap.Error(err))
				count++
				if count >= 3 {
					c.logger.Warn("keepalive timeout")
					c.reportException(ErrorKeepaliveFail, errors.New("keepalive fail"))
				}
			}
		case <-c.stopCh:
			return
		}
	}
}

func (c *Client) reportException(ec ErrorCode, err error) {
	c.notifyOnce.Do(func() {
		if c.notifyCh != nil {
			c.notifyCh <- c.bindInfo
		}
		close(c.stopCh)
		c.errCode = ec
		c.errInfo = err
	})
}

func (c *Client) GetErrorInfo() (ErrorCode, error) {
	return c.errCode, c.errInfo
}

func GeneratePwd(serial, random string) string {
	m5 := md5.New()
	m5.Write([]byte(serial))
	m5.Write([]byte("GoTag@2022"))
	m5.Write([]byte(random))
	cipherStr := hex.EncodeToString(m5.Sum(nil))
	return cipherStr[0:16]
}
