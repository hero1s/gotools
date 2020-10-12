package cluster

import (
	conf2 "github.com/hero1s/gotools/tools/leaf/conf"
	network2 "github.com/hero1s/gotools/tools/leaf/network"
	"math"
	"time"
)

var (
	server  *network2.TCPServer
	clients []*network2.TCPClient
)

func Init() {
	if conf2.ListenAddr != "" {
		server = new(network2.TCPServer)
		server.Addr = conf2.ListenAddr
		server.MaxConnNum = int(math.MaxInt32)
		server.PendingWriteNum = conf2.PendingWriteNum
		server.LenMsgLen = 4
		server.MaxMsgLen = math.MaxUint32
		server.NewAgent = newAgent

		server.Start()
	}

	for _, addr := range conf2.ConnAddrs {
		client := new(network2.TCPClient)
		client.Addr = addr
		client.ConnNum = 1
		client.ConnectInterval = 3 * time.Second
		client.PendingWriteNum = conf2.PendingWriteNum
		client.LenMsgLen = 4
		client.MaxMsgLen = math.MaxUint32
		client.NewAgent = newAgent

		client.Start()
		clients = append(clients, client)
	}
}

func Destroy() {
	if server != nil {
		server.Close()
	}

	for _, client := range clients {
		client.Close()
	}
}

type Agent struct {
	conn *network2.TCPConn
}

func newAgent(conn *network2.TCPConn) network2.Agent {
	a := new(Agent)
	a.conn = conn
	return a
}

func (a *Agent) Run() {}

func (a *Agent) OnClose() {}
