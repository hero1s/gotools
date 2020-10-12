package leaf

import (
	"github.com/hero1s/gotools/log"
	cluster2 "github.com/hero1s/gotools/tools/leaf/cluster"
	console2 "github.com/hero1s/gotools/tools/leaf/console"
	module2 "github.com/hero1s/gotools/tools/leaf/module"
	"os"
	"os/signal"
)
//独立进程启动
func Run(mods ...module2.Module) {
	log.Info("Leaf starting up")
	// module
	for i := 0; i < len(mods); i++ {
		module2.Register(mods[i])
	}
	module2.Init()

	// cluster
	cluster2.Init()

	// console
	console2.Init()

	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	log.Info("Leaf closing down (signal: %v)", sig)
	console2.Destroy()
	cluster2.Destroy()
	module2.Destroy()
}
//内部模块启动
func RunInside(end chan bool, mods ...module2.Module) {
	log.Info("Leaf starting up")

	// module
	for i := 0; i < len(mods); i++ {
		module2.Register(mods[i])
	}
	module2.Init()

	// cluster
	cluster2.Init()

	// console
	console2.Init()

	// close
	sig := <-end
	log.Info("Leaf closing down (signal: %v)", sig)
	console2.Destroy()
	cluster2.Destroy()
	module2.Destroy()
}