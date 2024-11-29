package udpbuf

import (
	"net"
	"udpbuf/message"
	"udpbuf/receiver"
	"udpbuf/sender"
)

type Sender interface {
	// 同步发送信息
	SendMessage(conn *net.UDPConn, addr *net.UDPAddr, message *message.GeneralMessage) error
}

type Receiver interface {
	// 异步接受和处理信息
	ReceiveMessage(conn *net.UDPConn, handleMessage func(*message.GeneralMessage)) (receiver.ReceiverCancelFunc, chan error)
}

func NewSender() Sender {
	return &sender.UdpSender{}
}

func NewReceiver() Receiver {
	return &receiver.UdpReceiver{}
}
