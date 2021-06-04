package config

import (
	"fmt"
	"os"

	"github.com/bsync-tech/mlog"
	"github.com/spf13/viper"
)

type Config struct {
	Path string
}

type Option func(cfg *Config)

var (
	// 模块变量
	inCfg = &Config{
		Path: "./conf/config.yaml",
	}
	// 配置文件变量
	C *viper.Viper
)

func New(opts ...Option) *Config {
	for _, o := range opts {
		o(inCfg)
	}
	return inCfg
}

func WithPath(path string) Option {
	return func(cfg *Config) {
		cfg.Path = path
	}
}

func Initialize() {
	if inCfg.Path == "" {
		panic(fmt.Errorf("cfg configure file invalid %s\n", inCfg.Path))
	}
	if _, err := os.Stat(inCfg.Path); err != nil {
		panic(fmt.Errorf("cfg configure file invalid %s err %s\n", inCfg.Path, err))
	}
	C = viper.New()
	C.SetConfigFile(inCfg.Path)
	C.SetConfigType("yaml")
	mlog.Debug(inCfg.Path)
	err := C.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error cfg configure file: %s \n", err))
	}
	mlog.Debug("Module config Initialize success")
}

func UnInitialize() {
	mlog.Debug("Module config UnInitialize success")
}
