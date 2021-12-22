package server

import (
	"net/http"
	"os"
	"time"

	"github.com/bsync-tech/mlog"
	"github.com/bsync-tech/scaffold/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gopkg.in/tylerb/graceful.v1"
)

type Server struct {
	Port string
}

type Option func(server *Server)

var (
	httpHandler *graceful.Server
)

func init() {
	empty, _ := os.Open("/dev/null")
	gin.DefaultWriter = empty
}

func HttpServerRun() {
	gin.SetMode(config.C.GetString("mode"))

	r := InitRouter()
	httpHandler = &graceful.Server{
		Server: &http.Server{
			Addr:         config.C.GetString("http_server.listen"),
			Handler:      r,
			ReadTimeout:  time.Duration(config.C.GetInt("http_server.read_timeout")) * time.Second,
			WriteTimeout: time.Duration(config.C.GetInt("http_server.write_timeout")) * time.Second,
		},
	}
	mlog.Info("HttpServerRun", zap.Any("listen", config.C.GetString("http_server.listen")))
	if err := httpHandler.ListenAndServe(); err != nil {
		mlog.Error("HttpServerRun", zap.Any("addr", config.C.GetString("http_server.addr")), zap.Any("err", err))
	}
	mlog.Debug("HttpServerRun server exit")
}

func InitRouter() *gin.Engine {
	router := gin.New()
	router.Use(RecoveryMiddleware(), RequestId(), GinLogger(mlog.GetLogger()))
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	return router
}

func ServerStop() {

	mlog.Debug("ServerStop stopped")
}
