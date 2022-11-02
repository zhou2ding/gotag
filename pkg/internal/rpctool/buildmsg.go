package rpctool

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/pkg/errors"
)

const (
	PacketMaxLength     = 4 * 1024
	PacketHeaderLength  = 32
	PacketTailLength    = 4
	PacketPayloadLength = PacketMaxLength - PacketHeaderLength - PacketTailLength
)

type msgHeader struct {
	Prefix          [4]byte
	Version         uint16
	Seq             uint16
	PktLength       uint32
	SubPktLength    uint32
	SubPktCount     uint32
	SubPktIndex     uint32
	PktJsonLength   uint32
	PktBinaryLength uint32
}

type pktTraceInfo struct {
	pktLength     uint32
	jsonLength    uint32
	binLength     uint32
	subPktCount   uint32
	nextSubPktIdx uint32
	dataBuf       []byte
	dataPos       uint32
}

type RpcMsgPayload struct {
	JsonRequest   []byte
	BinaryRequest []byte
}

type RpcMsgBuilder struct {
	reader  io.Reader
	InfoMap map[int]*pktTraceInfo
}

func CreateRPCMsgBuilder(reader io.Reader) *RpcMsgBuilder {
	return &RpcMsgBuilder{
		reader:  reader,
		InfoMap: make(map[int]*pktTraceInfo),
	}
}

func (builder *RpcMsgBuilder) DoRead() (*RpcMsgPayload, error) {
	for {
		//先读取RPC消息头
		headerLen := binary.Size(msgHeader{})
		headerbuf := make([]byte, headerLen)
		_, err := builder.readToBuf(builder.reader, headerbuf)
		if err != nil {
			return nil, err
		}

		header, err := builder.resolveMsgHeader(headerbuf)
		if err != nil {
			return nil, err
		}

		//判断RPC消息是否分包
		if header.SubPktCount == 1 { //不分包
			bodyBuf := make([]byte, header.PktLength+PacketTailLength)
			_, err = builder.readToBuf(builder.reader, bodyBuf)
			if err != nil {
				return nil, err
			}

			jsonData := bodyBuf[0:header.PktJsonLength]
			binData := bodyBuf[header.PktJsonLength : header.PktJsonLength+header.PktBinaryLength]

			return &RpcMsgPayload{JsonRequest: jsonData, BinaryRequest: binData}, nil
		} else { //分包
			var traceInfo *pktTraceInfo
			if _, ok := builder.InfoMap[int(header.Seq)]; ok { //已有记录
				traceInfo = builder.InfoMap[int(header.Seq)]
				traceInfo.jsonLength += header.PktJsonLength
				traceInfo.binLength += header.PktBinaryLength
			} else { //没有记录则新增一条记录
				traceInfo = &pktTraceInfo{
					pktLength:     header.PktLength,
					jsonLength:    header.PktJsonLength,
					binLength:     header.PktBinaryLength,
					subPktCount:   header.SubPktCount,
					nextSubPktIdx: 0,
					dataBuf:       make([]byte, header.PktLength+PacketTailLength),
					dataPos:       0,
				}
				builder.InfoMap[int(header.Seq)] = traceInfo
			}

			if header.SubPktIndex != traceInfo.nextSubPktIdx {
				return nil, errors.New("Packet lost or packet order error!!")
			}

			if header.SubPktLength > traceInfo.pktLength-traceInfo.dataPos {
				return nil, errors.New("Too Many Data")
			}

			//继续接收该分包
			if header.SubPktLength > traceInfo.pktLength-traceInfo.dataPos {
				//分包携带数据量超过整包剩余数据量（头部描述的整包数据量-已接收包数据量）
				return nil, errors.New("Too Many Data")
			}
			nRead, err := builder.readToBuf(builder.reader, traceInfo.dataBuf[traceInfo.dataPos:traceInfo.dataPos+header.SubPktLength+PacketTailLength])
			if err != nil {
				return nil, err
			}

			traceInfo.nextSubPktIdx += 1
			//去除四个字节的消息尾
			traceInfo.dataPos += uint32(nRead) - PacketTailLength

			if traceInfo.nextSubPktIdx < traceInfo.subPktCount {
				continue
			} else { //接收完成
				if traceInfo.dataPos != traceInfo.pktLength {
					return nil, errors.New("Packet Data Not Enough!!")
				}

				jsonData := traceInfo.dataBuf[0:traceInfo.jsonLength]
				binData := traceInfo.dataBuf[traceInfo.jsonLength : traceInfo.jsonLength+traceInfo.binLength]

				//完成一个，删除之记录
				delete(builder.InfoMap, int(header.Seq))

				return &RpcMsgPayload{JsonRequest: jsonData, BinaryRequest: binData}, nil
			}
		}
	}
}

func (builder *RpcMsgBuilder) resolveMsgHeader(msg []byte) (*msgHeader, error) {
	var header msgHeader
	if len(msg) < binary.Size(header) {
		return nil, nil
	}

	buf := bytes.NewBuffer(msg)
	err := binary.Read(buf, binary.BigEndian, &header)
	if err != nil {
		return nil, err
	}

	if header.Prefix[0] != 'Z' || header.Prefix[1] != 'D' || header.Prefix[2] != 'B' || header.Version != RpcVersion {
		return nil, errors.New("Invalid RPC Header")
	}

	return &header, nil
}

func (builder *RpcMsgBuilder) readToBuf(r io.Reader, buf []byte) (int, error) {
	l := len(buf)
	nRead := 0
	for {
		n, err := r.Read(buf[nRead:])
		if err != nil {
			return n, err
		}

		nRead += n
		if nRead >= l {
			break
		}
	}

	return nRead, nil
}
