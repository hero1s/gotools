package module

import (
	"github.com/hero1s/gotools/log"
	chanrpc2 "github.com/hero1s/gotools/tools/leaf/chanrpc"
	console2 "github.com/hero1s/gotools/tools/leaf/console"
	"github.com/hero1s/gotools/tools/leaf/go"
	timer2 "github.com/hero1s/gotools/tools/leaf/timer"
	"time"
)

type Skeleton struct {
	GoLen              int
	TimerDispatcherLen int
	AsynCallLen        int
	ChanRPCServer      *chanrpc2.Server
	g                  *g.Go
	dispatcher         *timer2.Dispatcher
	client             *chanrpc2.Client
	server             *chanrpc2.Server
	commandServer      *chanrpc2.Server
}

func (s *Skeleton) Init() {
	if s.GoLen <= 0 {
		s.GoLen = 0
	}
	if s.TimerDispatcherLen <= 0 {
		s.TimerDispatcherLen = 0
	}
	if s.AsynCallLen <= 0 {
		s.AsynCallLen = 0
	}

	s.g = g.New(s.GoLen)
	s.dispatcher = timer2.NewDispatcher(s.TimerDispatcherLen)
	s.client = chanrpc2.NewClient(s.AsynCallLen)
	s.server = s.ChanRPCServer

	if s.server == nil {
		s.server = chanrpc2.NewServer(0)
	}
	s.commandServer = chanrpc2.NewServer(0)
}

func (s *Skeleton) Run(closeSig chan bool) {
	for {
		select {
		case <-closeSig:
			s.commandServer.Close()
			s.server.Close()
			for !s.g.Idle() || !s.client.Idle() {
				s.g.Close()
				s.client.Close()
			}
			return
		case ri := <-s.client.ChanAsynRet:
			s.client.Cb(ri)
		case ci := <-s.server.ChanCall:
			s.server.Exec(ci)
		case ci := <-s.commandServer.ChanCall:
			s.commandServer.Exec(ci)
		case cb := <-s.g.ChanCb:
			s.g.Cb(cb)
		case t := <-s.dispatcher.ChanTimer:
			t.Cb()
		}
	}
}

func (s *Skeleton) AfterFunc(d time.Duration, cb func()) *timer2.Timer {
	if s.TimerDispatcherLen == 0 {
		log.Critical("invalid TimerDispatcherLen")
		panic("invalid TimerDispatcherLen")
	}

	return s.dispatcher.AfterFunc(d, cb)
}

func (s *Skeleton) CronFunc(cronExpr *timer2.CronExpr, cb func()) *timer2.Cron {
	if s.TimerDispatcherLen == 0 {
		log.Critical("invalid TimerDispatcherLen")
		panic("invalid TimerDispatcherLen")
	}

	return s.dispatcher.CronFunc(cronExpr, cb)
}

func (s *Skeleton) Go(f func(), cb func()) {
	if s.GoLen == 0 {
		log.Critical("invalid GoLen")
		panic("invalid GoLen")
	}

	s.g.Go(f, cb)
}

func (s *Skeleton) NewLinearContext() *g.LinearContext {
	if s.GoLen == 0 {
		log.Critical("invalid GoLen")
		panic("invalid GoLen")
	}

	return s.g.NewLinearContext()
}

func (s *Skeleton) AsynCall(server *chanrpc2.Server, id interface{}, args ...interface{}) {
	if s.AsynCallLen == 0 {
		log.Critical("invalid AsynCallLen")
		panic("invalid AsynCallLen")
	}

	s.client.Attach(server)
	s.client.AsynCall(id, args...)
}

func (s *Skeleton) RegisterChanRPC(id interface{}, f interface{}) {
	if s.ChanRPCServer == nil {
		log.Critical("invalid ChanRPCServer")
		panic("invalid ChanRPCServer")
	}

	s.server.Register(id, f)
}

func (s *Skeleton) RegisterCommand(name string, help string, f interface{}) {
	console2.Register(name, help, f, s.commandServer)
}
