package log

import (
	"io/ioutil"

	"github.com/bsync-tech/mlog"
	"gopkg.in/yaml.v2"
)

type Log struct {
	Path string
}

type Option func(log *Log)

var (
	// 模块变量
	inLog = &Log{
		Path: "./conf/log.yaml",
	}
)

func New(opts ...Option) *Log {
	for _, o := range opts {
		o(inLog)
	}
	return inLog
}

func WithPath(path string) Option {
	return func(m *Log) {
		m.Path = path
	}
}

func Initialize() {
	f, err := ioutil.ReadFile(inLog.Path)
	if err != nil {
		panic(err)
	}

	var c mlog.LogConfig
	err = yaml.Unmarshal(f, &c)
	if err != nil {
		panic(err)
	}
	mlog.Init(&c)
	mlog.Debug("Module log Initialize success")
}

func UnInitialize() {
	mlog.Sync()
	mlog.Debug("Module log UnInitialize success")
}
