package snlrpc

import (
	"net"
	"net/http"
	"net/rpc"

	"github.com/sirupsen/logrus"
)

// MsgBus TODO
type MsgBus int

// Args TODO
type Args struct {
	Channel string
	Value   string
}

// Receive receives a message from the CLI call and sets it as a variable
func (t *MsgBus) Receive(args *Args, reply *string) error {
	logrus.Infof("RECEIVED: %+v", *args)
	*reply = "OKiDoKi"
	return nil
}

// Setup sets up the communcation IPC communcation via tcp or unix socket
func Setup() {
	bus := new(MsgBus)
	rpc.Register(bus)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		logrus.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}
