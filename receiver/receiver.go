package receiver

import (
	"bytes"
	"fmt"
	"net"
	"sync"
	"udpbuf/message"
	"udpbuf/pb/general"

	"google.golang.org/protobuf/proto"
)

type UdpReceiver struct {
	messageCache map[uint32]*fragmentBuffer
	cacheLock    sync.Mutex
}

type ReceiverCancelFunc func()

type fragmentBuffer struct {
	fragments      map[uint32][]byte
	totalFragments uint32
}

func (r *UdpReceiver) ReceiveMessage(conn *net.UDPConn, handleMessage func(*message.GeneralMessage)) (ReceiverCancelFunc, chan error) {
	stopChan := make(chan struct{})
	errorChan := make(chan error, 100)

	go func() {
		for {
			select {
			case <-stopChan:
				// 收到关闭信号，退出 Goroutine
				// fmt.Println("Receiver stopped.")
				return
			default:
				data := make([]byte, 2048)
				n, _, err := conn.ReadFromUDP(data)
				if err != nil {
					// 记录错误到错误通道
					select {
					case errorChan <- err:
					default:
						// 如果错误通道已满，不阻塞
					}
					continue
				}

				fragment := &general.FragmentMessage{}
				err = proto.Unmarshal(data[:n], fragment)
				if err != nil {
					select {
					case errorChan <- fmt.Errorf("failed to parse FragmentMessage: %w", err):
					default:
					}
					continue
				}

				r.cacheLock.Lock()
				messageData, complete := r.processFragment(fragment)
				r.cacheLock.Unlock()

				if complete {
					message := &message.GeneralMessage{}
					err := proto.Unmarshal(messageData, (*general.GeneralMessage)(message))
					if err != nil {
						select {
						case errorChan <- fmt.Errorf("failed to parse GeneralMessage: %w", err):
						default:
						}
						continue
					}

					// 调用业务处理函数
					handleMessage(message)
				}
			}
		}
	}()

	// 返回取消函数和错误通道
	return func() {
		stopChan <- struct{}{}
		close(stopChan)
		close(errorChan)
	}, errorChan
}

func (r *UdpReceiver) processFragment(fragment *general.FragmentMessage) ([]byte, bool) {
	if r.messageCache == nil {
		r.messageCache = make(map[uint32]*fragmentBuffer)
	}

	if _, exists := r.messageCache[fragment.MessageId]; !exists {
		r.messageCache[fragment.MessageId] = &fragmentBuffer{
			fragments:      make(map[uint32][]byte),
			totalFragments: fragment.TotalFragments,
		}
	}

	buffer := r.messageCache[fragment.MessageId]
	buffer.fragments[fragment.FragmentId] = fragment.FragmentData

	if uint32(len(buffer.fragments)) == buffer.totalFragments {
		var completeMessage bytes.Buffer
		for i := uint32(0); i < buffer.totalFragments; i++ {
			completeMessage.Write(buffer.fragments[i])
		}
		delete(r.messageCache, fragment.MessageId)
		return completeMessage.Bytes(), true
	}

	return nil, false
}
