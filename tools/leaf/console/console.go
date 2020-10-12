package console

import (
	"bufio"
	conf2 "github.com/hero1s/gotools/tools/leaf/conf"
	network2 "github.com/hero1s/gotools/tools/leaf/network"
	"math"
	"strconv"
	"strings"
)

var server *network2.TCPServer

func Init() {
	if conf2.ConsolePort == 0 {
		return
	}

	server = new(network2.TCPServer)
	server.Addr = "localhost:" + strconv.Itoa(conf2.ConsolePort)
	server.MaxConnNum = int(math.MaxInt32)
	server.PendingWriteNum = 100
	server.NewAgent = newAgent

	server.Start()
}

func Destroy() {
	if server != nil {
		server.Close()
	}
}

type Agent struct {
	conn   *network2.TCPConn
	reader *bufio.Reader
}

func newAgent(conn *network2.TCPConn) network2.Agent {
	a := new(Agent)
	a.conn = conn
	a.reader = bufio.NewReader(conn)
	return a
}

func (a *Agent) Run() {
	for {
		if conf2.ConsolePrompt != "" {
			a.conn.Write([]byte(conf2.ConsolePrompt))
		}

		line, err := a.reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSuffix(line[:len(line)-1], "\r")

		args := strings.Fields(line)
		if len(args) == 0 {
			continue
		}
		if args[0] == "quit" {
			break
		}
		var c Command
		for _, _c := range commands {
			if name() == args[0] {
				c = _c
				break
			}
		}
		if c == nil {
			a.conn.Write([]byte("command not found, try `help` for help\r\n"))
			continue
		}
		output := run(args[1:])
		if output != "" {
			a.conn.Write([]byte(output + "\r\n"))
		}
	}
}

func (a *Agent) OnClose() {}
