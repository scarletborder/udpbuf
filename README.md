# udpbuf

EN | [ZH_CN](#zh_cn)

## EN

### Overview

The `udpbuf` library provides a simple and efficient way to send and receive fragmented messages over UDP. It ensures large messages are split into fragments for transmission and reassembled correctly upon receipt. The library is designed for both synchronous sending and asynchronous receiving of UDP messages, with easy integration into existing projects.

---

### Features

1. **Fragmented Message Transmission**: Automatically splits large messages into smaller fragments to accommodate UDP limitations.
2. **Reliable Reassembly**: Fragments are reassembled on the receiving end to reconstruct the original message.
3. **Customizable Message Handling**: Users can define their message processing logic.
4. **Error Handling**: Built-in mechanisms for handling transmission and parsing errors.

---

### Installation

```bash
go get github.com/scarletborder/udpbuf
```

---

### Usage

#### 1. Sending Messages
The `Sender` interface provides a `SendMessage` function to send large messages.

```go
import (
    "net"
    "udpbuf"
    "udpbuf/pb/general"
)

func main() {
    conn, _ := net.DialUDP("udp", nil, &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 8080})
    defer conn.Close()

    sender := udpbuf.NewSender()

    message := &general.GeneralMessage{
        Content: "Hello, UDP!",
    }

    err := sender.SendMessage(conn, conn.RemoteAddr().(*net.UDPAddr), message)
    if err != nil {
        fmt.Println("Error sending message:", err)
    }
}
```

#### 2. Receiving Messages
The `Receiver` interface provides an asynchronous `ReceiveMessage` function with a callback to handle received messages.

```go
import (
    "fmt"
    "net"
    "udpbuf"
    "udpbuf/pb/general"
)

func main() {
    conn, _ := net.ListenUDP("udp", &net.UDPAddr{Port: 8080})
    defer conn.Close()

    receiver := udpbuf.NewReceiver()

    cancelFunc, errorChan := receiver.ReceiveMessage(conn, func(msg *general.GeneralMessage) {
        fmt.Println("Received message:", msg.Content)
    })

    // Handle errors (if any)
    go func() {
        for err := range errorChan {
            fmt.Println("Error:", err)
        }
    }()

    // Stop the receiver after some time (optional)
    // time.Sleep(time.Minute)
    // cancelFunc()
}
```

---

### API Reference

#### **Sender**
- `SendMessage(conn *net.UDPConn, addr *net.UDPAddr, message *general.GeneralMessage) error`: Sends a message synchronously.

#### **Receiver**
- `ReceiveMessage(conn *net.UDPConn, handleMessage func(*general.GeneralMessage)) (ReceiverCancelFunc, chan error)`: Receives messages asynchronously, providing a cancellation function and an error channel.


## ZH_CN

### 概述

`udpbuf` 库为 UDP 数据传输提供了一种简单高效的方式。它能够自动将大型消息分片发送，并在接收端正确组装。库支持同步消息发送和异步消息接收，可轻松集成到现有项目中。

---

### 特性

1. **分片传输**: 自动将大消息拆分为小分片以适应 UDP 限制。
2. **可靠组装**: 接收端重组分片以还原原始消息。
3. **可自定义的消息处理**: 用户可自定义消息处理逻辑。
4. **错误处理**: 内置的错误传输和解析机制。

---

### 安装

```bash
go get github.com/scarletborder/udpbuf
```

---

### 使用方法

#### 1. 发送消息
`Sender` 接口通过 `SendMessage` 方法同步发送消息。

```go
import (
    "net"
    "github.com/scarletborder/udpbuf"
    "github.com/scarletborder/udpbuf/pb/general"
)

func main() {
    conn, _ := net.DialUDP("udp", nil, &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 8080})
    defer conn.Close()

    sender := udpbuf.NewSender()

    message := &general.GeneralMessage{
        Content: "Hello, UDP!",
    }

    err := sender.SendMessage(conn, conn.RemoteAddr().(*net.UDPAddr), message)
    if err != nil {
        fmt.Println("发送消息失败:", err)
    }
}
```

#### 2. 接收消息
`Receiver` 接口通过 `ReceiveMessage` 方法异步接收消息，并提供回调处理收到的消息。

```go
import (
    "fmt"
    "net"
    "udpbuf"
    "udpbuf/pb/general"
)

func main() {
    conn, _ := net.ListenUDP("udp", &net.UDPAddr{Port: 8080})
    defer conn.Close()

    receiver := udpbuf.NewReceiver()

    cancelFunc, errorChan := receiver.ReceiveMessage(conn, func(msg *general.GeneralMessage) {
        fmt.Println("接收到消息:", msg.Content)
    })

    // 处理错误（如果有）
    go func() {
        for err := range errorChan {
            fmt.Println("错误:", err)
        }
    }()

    // 停止接收器（可选）
    // time.Sleep(time.Minute)
    // cancelFunc()
}
```

---

### API 参考

#### **Sender**
- `SendMessage(conn *net.UDPConn, addr *net.UDPAddr, message *general.GeneralMessage) error`: 同步发送消息。

#### **Receiver**
- `ReceiveMessage(conn *net.UDPConn, handleMessage func(*general.GeneralMessage)) (ReceiverCancelFunc, chan error)`: 异步接收消息，提供取消函数和错误通道。
