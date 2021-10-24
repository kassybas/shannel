package snlrpc

import (
	"net"
	"net/http"
	"net/rpc"

	"github.com/sirupsen/logrus"
)

type MsgBus int

type Args struct {
	Channel string
	Value   string
}

func (t *MsgBus) Receive(args *Args, reply *string) error {
	logrus.Infof("RECEIVED: %+v", *args)
	*reply = "OKiDoKi"
	return nil
}

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
