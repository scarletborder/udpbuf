package sender

import (
	"fmt"
	"math"
	"net"
	"time"
	"udpbuf/constant"
	"udpbuf/message"
	"udpbuf/pb/general"

	"google.golang.org/protobuf/proto"
)

type UdpSender struct{}

func (s *UdpSender) SendMessage(conn *net.UDPConn, addr *net.UDPAddr, message *message.GeneralMessage) error {
	data, err := proto.Marshal((*general.GeneralMessage)(message))
	if err != nil {
		return fmt.Errorf("failed to serialize GeneralMessage: %w", err)
	}

	messageID := uint32(time.Now().UnixNano() / 1e6) // 使用时间戳作为消息ID
	err = s.sendLargeMessage(conn, addr, messageID, data)
	if err != nil {
		return err
	}

	return nil
}

func (s *UdpSender) sendLargeMessage(conn *net.UDPConn, addr *net.UDPAddr, messageID uint32, data []byte) error {
	totalFragments := uint32(math.Ceil(float64(len(data)) / float64(constant.MTU)))

	for i := uint32(0); i < totalFragments; i++ {
		start := i * constant.MTU
		end := start + constant.MTU
		if end > uint32(len(data)) {
			end = uint32(len(data))
		}
		fragmentData := data[start:end]

		fragment := &general.FragmentMessage{
			MessageId:      messageID,
			FragmentId:     i,
			TotalFragments: totalFragments,
			FragmentData:   fragmentData,
		}

		serializedFragment, err := proto.Marshal(fragment)
		if err != nil {
			return fmt.Errorf("failed to serialize fragment: %w", err)
		}
		_, err = conn.WriteToUDP(serializedFragment, addr)
		if err != nil {
			return fmt.Errorf("failed to send fragment: %w", err)
		}
	}
	return nil
}
