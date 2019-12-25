package leaf

import (
	"github.com/hero1s/gotools/leaf/cluster"
	"github.com/hero1s/gotools/leaf/console"
	"github.com/hero1s/gotools/leaf/module"
	"github.com/hero1s/gotools/log"
	"os"
	"os/signal"
)
//独立进程启动
func Run(mods ...module.Module) {
	log.Info("Leaf starting up")
	// module
	for i := 0; i < len(mods); i++ {
		module.Register(mods[i])
	}
	module.Init()

	// cluster
	cluster.Init()

	// console
	console.Init()

	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	log.Info("Leaf closing down (signal: %v)", sig)
	console.Destroy()
	cluster.Destroy()
	module.Destroy()
}
//内部模块启动
func RunInside(end chan bool, mods ...module.Module) {
	log.Debug("Leaf starting up")

	// module
	for i := 0; i < len(mods); i++ {
		module.Register(mods[i])
	}
	module.Init()

	// cluster
	cluster.Init()

	// console
	console.Init()

	// close
	sig := <-end
	log.Debug("Leaf closing down (signal: %v)", sig)
	console.Destroy()
	cluster.Destroy()
	module.Destroy()
}