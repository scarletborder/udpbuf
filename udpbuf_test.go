package udpbuf_test

import (
	"net"
	"os"
	"testing"
	"udpbuf"
	"udpbuf/message"
)

func TestSendLittle(t *testing.T) {
	addr, _ := net.ResolveUDPAddr("udp", ":9000")
	conn, _ := net.ListenUDP("udp", addr)

	receiver := udpbuf.NewReceiver()
	sender := udpbuf.NewSender()

	// 接收消息处理
	end_Chan := make(chan int, 1)

	cancel, _ := receiver.ReceiveMessage(conn, func(msg *message.GeneralMessage) {
		t.Logf("Received message of type: %d\n content: %v", msg.Type, msg.Content)

	})
	defer func() {
		cancel()
		end_Chan <- 1
	}()

	// 测试发送消息
	targetAddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:9000")
	testMessage := &message.GeneralMessage{
		Type:    1,
		Content: []byte("Hello, World!"),
	}
	err := sender.SendMessage(conn, targetAddr, testMessage)
	if err != nil {
		t.Errorf("Error sending message: %s", err)
	}
	// <-end_Chan
}

func TestSendBigBlob(t *testing.T) {
	addr, _ := net.ResolveUDPAddr("udp", ":9000")
	conn, _ := net.ListenUDP("udp", addr)

	receiver := udpbuf.NewReceiver()
	sender := udpbuf.NewSender()

	// 接收消息处理
	end_Chan := make(chan int, 1)

	cancel, _ := receiver.ReceiveMessage(conn, func(msg *message.GeneralMessage) {
		t.Logf("Received message of type: %d", msg.Type)
		end_Chan <- 1
	})
	defer cancel()

	// 测试发送消息
	targetAddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:9000")
	cont, _ := os.ReadFile("./32258607_p0.png")

	testMessage := &message.GeneralMessage{
		Type:    1,
		Content: cont,
	}
	err := sender.SendMessage(conn, targetAddr, testMessage)
	if err != nil {
		t.Errorf("Error sending message: %s", err)
	}
	<-end_Chan
}
