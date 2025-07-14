package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"zeroim/common/libnet"

	"github.com/zeromicro/go-zero/core/logx"
)

func main() {
	logx.DisableStat()

	conn, err := net.Dial("tcp", "127.0.0.1:9898")
	if err != nil {
		fmt.Printf("Error connection:%v\n", err)
		return
	}
	fmt.Println("Connected to server.")

	protocol := libnet.NewIMProtocol()
	coder := protocol.NewCodec(conn)

	go readServerResponse(coder)

	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		if err := coder.Send(libnet.Message{
			Body: []byte(text),
		}); err != nil {
			fmt.Printf("send message err:%v\n", err)
		}
	}
}
func readServerResponse(codec libnet.Codec) {
	for {
		msg, err := codec.Receive()
		if err != nil {
			logx.Errorf("Error reading from server:%v", err)
			break
		}
		fmt.Println("Server response: " + string(msg.Body))
		fmt.Println("请输入:")
	}
}
