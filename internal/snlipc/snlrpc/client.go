package snlrpc

import (
	"fmt"
	"net/rpc"

	"github.com/sirupsen/logrus"
)

const localhost = "127.0.0.1"

func Send(channelName string, message string, port int) {
	logrus.Debug("Connecting to rpc server...")
	client, err := rpc.DialHTTP("tcp", fmt.Sprintf("%s:%d", localhost, port))
	if err != nil {
		logrus.Fatal("rpc connection dialing:", err)
	}

	args := &Args{channelName, message}
	var reply string
	err = client.Call("MsgBus.Receive", args, &reply)
	if err != nil {
		logrus.Fatal("push error:", err)
	}
}
