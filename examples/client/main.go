package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"zeroim/common/libnet"
	"zeroim/imrpc/imrpcclient"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
)

var (
	token   string
	toToken string
)

func main() {
	logx.DisableStat()

	conn, err := net.Dial("tcp", "127.0.0.1:9898")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to server.")

	fmt.Println("请输入聊天对象Token:")
	chatWithToken()

	fmt.Println("请输入你的Token:")

	protocol := libnet.NewIMProtocol()
	codec := protocol.NewCodec(conn)

	go readServerResponse(codec)

	err = login(codec)
	if err != nil {
		panic(err)
	}

	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')

		msgReq := &imrpcclient.PostMsg{
			Token:   token,
			ToToken: toToken,
			Msg:     text,
		}
		msgData, err := proto.Marshal(msgReq)
		if err != nil {
			panic(err)
		}
		msg := libnet.Message{
			Body: msgData,
		}
		err = codec.Send(msg)
		if err != nil {
			fmt.Printf("send error: %v\n", err)
		}
	}
}

func login(codec libnet.Codec) error {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	token = text
	loginReq := &imrpcclient.LoginRequest{
		Token:         token,
		Authorization: "Bearer token",
	}
	loginData, err := proto.Marshal(loginReq)
	if err != nil {
		panic(err)
	}
	msg := libnet.Message{Body: loginData}
	return codec.Send(msg)
}

func chatWithToken() {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	toToken = text
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
