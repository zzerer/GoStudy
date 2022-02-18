package main

import (
	"errors"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

type SeverityHook struct{}

func (h SeverityHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if level != zerolog.NoLevel {
		e.Str("hook", level.String())
	}
}

func main() {

	//最简单的用法
	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Print("hello world") //log.Print()方法默认级别DEBUG
	// {"level":"debug","time":1645175280,"message":"hello world"}

	//添加自定义字段
	log.Debug().
		Str("Scale", "833 cents").
		Float64("Interval", 833.09).
		Msg("Fibonacci is everywhere")
	// {"level":"debug","Scale":"833 cents","Interval":833.09,"time":1645175280,"message":"Fibonacci is everywhere"}
	log.Debug().
		Str("Name", "Tom").
		Send() //同Msg("")
	// {"level":"debug","Name":"Tom","time":1645175280}

	//log level
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Info().Msg("hello world")
	log.Debug().Msg("hello world") // 这句不会打印出来
	// {"level":"info","time":1645175280,"message":"hello world"}

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Info().Msg("hello world")
	log.Debug().Msg("hello world")
	// {"level":"info","time":1645175280,"message":"hello world"}
	// {"level":"debug","time":1645175280,"message":"hello world"}

	//error log
	err := errors.New("seems we have an error here")
	log.Error().Err(err).Msg("")
	// {"level":"error","error":"seems we have an error here","time":1645175280}

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack //支持打印错误栈，必须设置
	log.Error().Stack().Err(err).Msg("")

	//logger实例
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	logger.Info().Str("foo", "bar").Msg("hello world")
	// {"level":"info","foo":"bar","time":1645175280,"message":"hello world"}

	sublogger := log.With().Str("component", "foo").Logger() //带预定义字段的logger
	sublogger.Info().Msg("hello world")
	// {"level":"info","component":"foo","time":1645175766,"message":"hello world"}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}) //覆盖原输出格式，可以定义颜色，分隔符等
	log.Info().Str("foo", "bar").Msg("Hello world")
	// 9:18AM INF Hello world foo=bar
	log.Logger = log.Output(os.Stderr) //还原

	//自定义标签
	zerolog.TimestampFieldName = "t"
	zerolog.LevelFieldName = "l"
	zerolog.MessageFieldName = "m"
	log.Info().Msg("hello world")
	// {"l":"info","t":1645176304,"m":"hello world"}

	//带上调用文件和行数
	log.Logger = log.With().Caller().Logger()
	log.Info().Msg("hello world")
	// {"l":"info","t":1645176439,"caller":"/Users/zengyuzhao/gopath/src/github/GoStudy/zerolog/zerolog.go:74","m":"hello world"}

	//Hook
	hooked := log.Hook(SeverityHook{})
	hooked.Warn().Msg("")
	// {"l":"warn","t":1645176744,"caller":"/Users/zengyuzhao/gopath/src/github/GoStudy/zerolog/zerolog.go:87","hook":"warn"}

	//同时多种途径输出
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
	multi := zerolog.MultiLevelWriter(consoleWriter, os.Stdout)
	logger = zerolog.New(multi).With().Timestamp().Logger()
	logger.Info().Msg("Hello World!")
	// 9:37AM INF Hello World!
	// {"l":"info","t":1645177020,"m":"Hello World!"}
}
