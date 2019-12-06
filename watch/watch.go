package watch

import (
	"github.com/howeyc/fsnotify"
)

type Log interface {
	Error(format string, v ...interface{})
	Info(format string, v ...interface{})
}

/*
监控目录下的文件,可监控多信目录, 目前只提供了监控文件的改变，没有监控文件创建删除移动等
callback用于文件有改变时执行的回调函数
*/

func WatchPath(callback func() error, log Log, paths ...string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				log.Info("文件(%v)有更新", ev.Name)
				err := callback()
				if err != nil {
					log.Error("重新加载配置文件(%v)失败,error:%v", ev.Name, err.Error())
				}
			case err := <-watcher.Error:
				if err != nil {
					log.Error("监控文件有错误发生,error:%v", err.Error())
				}

			}
		}
	}()

	for _, path := range paths {
		err := watcher.AddWatch(path, fsnotify.FSN_MODIFY)
		if err != nil {
			return err
		}
	}

	// don't call watcher.Close(), because it should run forever until program dead
	return nil
}
