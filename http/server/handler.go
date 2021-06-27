package server

import (
	"fmt"
	"time"

	"github.com/bsync-tech/mlog"
	"github.com/bsync-tech/scaffold/config"
	"github.com/gin-gonic/gin"
	"github.com/gookit/goutil/strutil"
	"go.uber.org/zap"
)

// alarm with level error_02 will be saved with timer, if the timer expires, new alarm
// will be saved again, or count will be increased

type Msg struct {
	// msg count of sha1
	count int
	// subsys
	subsys string
	// module
	module string
	// leve
	level string
	// msg
	msg string
	// ts
	ts time.Time
}

type ReqBody struct {
	// subsys
	Subsys string `json:"subsys"`
	// module
	Module string `json:"module"`
	// leve
	Level string `json:"level"`
	// msg
	Msg string `json:"msg"`
	// ts
	TS string `json:"ts"`
}

var (
	messages map[string]*Msg = make(map[string]*Msg)
	done     chan struct{}   = make(chan struct{})
)

func LogMessage(c *gin.Context) {
	var req ReqBody
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatusJSON(200, gin.H{
			"errcode": 101,
			"message": fmt.Sprintf("invalid body format: %v", err),
		})
	}

	key := req.Subsys + req.Module + req.Msg
	hash := strutil.Md5(key)
	if _, ok := messages[hash]; !ok {
		var err error
		messages[hash] = new(Msg)
		messages[hash].count = 1
		messages[hash].subsys = req.Subsys
		messages[hash].module = req.Module
		messages[hash].level = req.Level
		messages[hash].msg = req.Msg
		messages[hash].ts, err = time.Parse("2006-01-02 15:04:05 000", req.TS)
		if err != nil {
			c.AbortWithStatusJSON(200, gin.H{
				"errcode": 101,
				"message": fmt.Sprintf("invalid timestamp format: %v %s", err, req.TS),
			})
		}
		mlog.Debug(fmt.Sprintf("log message ok %v", messages[hash]))
	} else {
		messages[hash].count++
		mlog.Debug(fmt.Sprintf("log message count +1 %s %d", hash, messages[hash].count))
	}
	c.AbortWithStatusJSON(200, gin.H{
		"errcode": 0,
		"message": "ok",
	})
}

func ClearExpired() {
	for k, v := range messages {
		loc, _ := time.LoadLocation("UTC")
		now, err := time.ParseInLocation("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"), loc)
		if err != nil {
			mlog.Error("format local time err %v", zap.Error(err))
		}
		if v.ts.Add(config.C.GetDuration("alarm_expired")).Before(now) {
			delete(messages, k)
			mlog.Debug(fmt.Sprintf("msg clear %s %s", k, v.ts.Format("2006-01-02 15:04:05")))
		}
	}
}

func ClearExpiredGoRoutine() {
	for {
		select {
		case <-done:
			mlog.Debug("clear exprired goroutine done")
			return
		default:
			ClearExpired()
			time.Sleep(5 * time.Second)
		}
	}
}

func Close() {
	close(done)
}
