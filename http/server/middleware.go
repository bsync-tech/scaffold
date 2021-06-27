package server

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/bsync-tech/scaffold/config"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

// IP 白名单拦截
func IPAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		isMatched := false
		for _, host := range config.C.GetStringSlice("http.allow_ip") {
			if c.ClientIP() == host {
				isMatched = true
			}
		}
		if !isMatched {
			c.JSON(http.StatusInternalServerError, errors.New(fmt.Sprintf("%v, not in iplist", c.ClientIP())))
			c.Abort()
			return
		}
		c.Next()
	}
}

// 捕获所有panic，并且返回错误信息
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				//先做一下日志记录
				// mlog.Debug(string(debug.Stack()))
				if config.C.GetString("mode") != "release" {
					c.JSON(http.StatusInternalServerError, errors.New("内部错误"))
					return
				} else {
					c.JSON(http.StatusInternalServerError, errors.New(fmt.Sprint(err)))
					return
				}
			}
		}()
		c.Next()
	}
}

// 请求 id
func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.Request.Header.Get("X-Request-Id")
		if requestID == "" {
			uuid4 := uuid.NewV4()
			requestID = uuid4.String()
		}

		c.Set("RequestId", requestID)

		c.Writer.Header().Set("X-Request-Id", requestID)
		c.Next()
	}
}

func GinLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)

		logger.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.String("request-id", c.GetString("RequestId")),
			zap.Duration("cost", cost),
		)
	}
}
