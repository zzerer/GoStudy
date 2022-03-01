package main

import (
	"os"

	"github.com/sirupsen/logrus"
)

type MyHook struct {
	msg string
}

func (h *MyHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *MyHook) Fire(entry *logrus.Entry) error {  // error级别才触发
	if entry.Level == logrus.ErrorLevel {
		entry.Data["error_msg"] = h.msg
	}
	return nil
}

func main() {
	//最简单的用法
	logrus.Info("Something noteworthy happened!")
	//INFO[0000] Something noteworthy happened!

	//设置全局日志级别
	logrus.SetLevel(logrus.InfoLevel) //默认级别
	logrus.SetLevel(logrus.DebugLevel)

	//不同的log level
	logrus.Trace("Something very low level.")
	logrus.Debug("Useful debugging information.")
	logrus.Info("Something noteworthy happened!")
	logrus.Warn("You should probably take a look at this.")
	logrus.Error("Something failed but I'm not quitting.")
	// 执行完后会执行 os.Exit(1)
	//logrus.Fatal("Bye.")
	//执行完后会执行 panic()
	//logrus.Panic("I'm bailing.")

	//设置日志格式TEXT/JSON
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.Info("Something noteworthy happened!")
	//INFO[0000] Something noteworthy happened!
	logrus.SetFormatter(&logrus.TextFormatter{ //纯文本的格式，"k=v"
		DisableColors: true,
		FullTimestamp: true,
	})
	logrus.Info("Something noteworthy happened!")
	//time="2022-03-01T15:13:21+08:00" level=info msg="Something noteworthy happened!"

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.Info("Something noteworthy happened!")
	//{"level":"info","msg":"Something noteworthy happened!","time":"2022-03-01T14:45:54+08:00"}

	//带上方法名
	logrus.SetReportCaller(true)
	logrus.Info("Something noteworthy happened!")
	//{"file":"/Users/zengyuzhao/gopath/pkg/mod/github.com/sirupsen/logrus@v1.4.2/logger.go:192","func":"github.com/sirupsen/logrus.(*Logger).Log","level":"info","msg":"Something noteworthy happened!","time":"2022-03-01T15:17:58+08:00"}
	logrus.SetReportCaller(false)

	//logger 实例
	var logger = logrus.New()
	logger.Info("logger log")
	//INFO[0000] logger log

	//添加自定义字段
	logrus.WithFields(logrus.Fields{
		"animal": "walrus",
		"size":   10,
	}).Info("A walrus appears")
	//{"animal":"walrus","level":"info","msg":"A walrus appears","size":10,"time":"2022-03-01T15:39:52+08:00"}

	mylogger := logrus.WithFields(logrus.Fields{ //带默认字段的logger
		"common": "this is a common field",
		"other":  "I also should be logged always",
	})
	mylogger.Info("mylogger log")
	//{"common":"this is a common field","level":"info","msg":"mylogger log","other":"I also should be logged always","time":"2022-03-01T15:55:03+08:00"}

	//重定向输出
	var logger_new = logrus.New()
	//logger_new.Out = os.Stdout

	// 可以写入到文件
	file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logger_new.Out = file
		logger_new.Info("success to log to file")
	} else {
		logger_new.Info("Failed to log to file, using default stderr")
	}

	//hook
	logger_withhook := logrus.New()
	my_hook := &MyHook{msg: "this is error message"}
	logger_withhook.AddHook(my_hook)
	logger_withhook.Info("should not call hook")
	logger_withhook.Error("should call hook")
	//INFO[0000] should not call hook                         
	//ERRO[0000] should call hook                              error_msg="this is error message"
}
