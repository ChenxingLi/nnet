package node

import (
	"errors"

	"github.com/gogo/protobuf/proto"
	"github.com/nknorg/nnet/message"
	"github.com/nknorg/nnet/protobuf"
)

const (
	// msg length is encoded by 32 bit int
	msgLenBytes = 4
)

// RemoteMessage is the received msg from remote node
type RemoteMessage struct {
	RemoteNode *RemoteNode
	Msg        *protobuf.Message
}

// NewRemoteMessage creates a RemoteMessage with remote node rn and msg
func NewRemoteMessage(rn *RemoteNode, msg *protobuf.Message) (*RemoteMessage, error) {
	remoteMsg := &RemoteMessage{
		RemoteNode: rn,
		Msg:        msg,
	}
	return remoteMsg, nil
}

// NewPingMessage creates a PING message for heartbeat
func NewPingMessage() (*protobuf.Message, error) {
	id, err := message.GenID()
	if err != nil {
		return nil, err
	}

	msgBody := &protobuf.Ping{}

	buf, err := proto.Marshal(msgBody)
	if err != nil {
		return nil, err
	}

	msg := &protobuf.Message{
		MessageType: protobuf.PING,
		RoutingType: protobuf.DIRECT,
		MessageId:   id,
		Message:     buf,
	}

	return msg, nil
}

// NewPingReply creates a PING reply for heartbeat
func NewPingReply(replyToID []byte) (*protobuf.Message, error) {
	msgBody := &protobuf.Pong{}

	buf, err := proto.Marshal(msgBody)
	if err != nil {
		return nil, err
	}

	msg := &protobuf.Message{
		MessageType: protobuf.PING,
		RoutingType: protobuf.DIRECT,
		ReplyToId:   replyToID,
		Message:     buf,
	}
	return msg, nil
}

// NewGetNodeMessage creates a GET_NODE message to get node info
func NewGetNodeMessage() (*protobuf.Message, error) {
	id, err := message.GenID()
	if err != nil {
		return nil, err
	}

	msgBody := &protobuf.GetNode{}

	buf, err := proto.Marshal(msgBody)
	if err != nil {
		return nil, err
	}

	msg := &protobuf.Message{
		MessageType: protobuf.GET_NODE,
		RoutingType: protobuf.DIRECT,
		MessageId:   id,
		Message:     buf,
	}

	return msg, nil
}

// NewGetNodeReply creates a GET_NODE reply to send node info
func NewGetNodeReply(replyToID []byte, n *protobuf.Node) (*protobuf.Message, error) {
	msgBody := &protobuf.GetNodeReply{
		Node: n,
	}

	buf, err := proto.Marshal(msgBody)
	if err != nil {
		return nil, err
	}

	msg := &protobuf.Message{
		MessageType: protobuf.GET_NODE,
		RoutingType: protobuf.DIRECT,
		ReplyToId:   replyToID,
		Message:     buf,
	}

	return msg, nil
}

// handleRemoteMessage handles a remote message and returns error
func (ln *LocalNode) handleRemoteMessage(remoteMsg *RemoteMessage) error {
	switch remoteMsg.Msg.MessageType {
	case protobuf.PING:
		replyMsg, err := NewPingReply(remoteMsg.Msg.MessageId)
		if err != nil {
			return err
		}

		err = remoteMsg.RemoteNode.SendMessageAsync(replyMsg)
		if err != nil {
			return err
		}

	case protobuf.GET_NODE:
		replyMsg, err := NewGetNodeReply(remoteMsg.Msg.MessageId, remoteMsg.RemoteNode.LocalNode.Node.Node)
		if err != nil {
			return err
		}

		err = remoteMsg.RemoteNode.SendMessageAsync(replyMsg)
		if err != nil {
			return err
		}

	default:
		return errors.New("Unknown message type")
	}

	return nil
}
