package http

import (
	"crypto/tls"
	"fmt"
	"math"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "net/http/pprof"

	"github.com/bsync-tech/mlog"
	"github.com/bsync-tech/scaffold/config"
	resty "github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

var (
	restyClient *resty.Client
)

func Init() {
	httpClient := http.Client{}
	httpClient.Timeout = config.C.GetDuration("http_client.read_timeout")
	httpClient.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          config.C.GetInt("http_client.max_idle_connections"),
		MaxIdleConnsPerHost:   config.C.GetInt("http_client.max_idle_connections_perhost"),
		IdleConnTimeout:       config.C.GetDuration("http_client.max_idle_connection_timeout"),
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 5 * time.Second,
	}
	restyClient = resty.NewWithClient(&httpClient)
	restyClient.SetHeader("User-Agent", "FClient")
	restyClient.SetContentLength(true)
	restyClient.SetDebug(config.C.GetBool("http_client.debug"))

	if config.C.GetBool("http_client.https") == true {
		// 跳过校验
		mlog.Debug(fmt.Sprintf("%d global insecure skip verify", os.Getpid()))
		restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}

	restyClient.SetRetryCount(math.MaxInt32)
	restyClient.SetRetryWaitTime(3 * time.Second)
	restyClient.SetRetryMaxWaitTime(60 * time.Second)

	restyClient.AddRetryCondition(
		func(r *resty.Response, err error) bool {
			return r.StatusCode() != http.StatusOK
		},
	)
	restyClient.OnError(func(req *resty.Request, err error) {
		if v, ok := err.(*resty.ResponseError); ok {
			mlog.Error("request url err", zap.String("url", req.URL), zap.String("response", v.Response.String()), zap.Any("response err", v.Err), zap.Any("err", err))
		}
		mlog.Error("request url err", zap.String("url", req.URL), zap.Any("err", err))
	})

	restyClient.SetLogger(LFR)

	mlog.Info("http client init success.")
}

// 内部Resty日志
var LFR LogForResty

type LogForResty struct{}

func (l LogForResty) Errorf(format string, v ...interface{}) {
	mlog.Error(fmt.Sprintf(format, v...))
}
func (l LogForResty) Warnf(format string, v ...interface{}) {
	mlog.Warn(fmt.Sprintf(format, v...))
}
func (l LogForResty) Debugf(format string, v ...interface{}) {
	mlog.Debug(fmt.Sprintf(format, v...))
}

func GetHttpClient() *http.Client {
	return restyClient.GetClient()
}

func Get(url string, params map[string]string) (*resty.Response, error) {
	resp, err := restyClient.R().SetQueryParams(params).SetHeader("Content-Type", "application/json").Get(url)
	return resp, err
}

func Post(url string, params map[string]string, data interface{}) (*resty.Response, error) {
	resp, err := restyClient.R().EnableTrace().SetQueryParams(params).SetBody(data).SetHeader("Content-Type", "application/json").Post(url)
	return resp, err
}

func PostNotParseBody(url string, params map[string]string, data interface{}) (*resty.Response, error) {
	resp, err := restyClient.SetDoNotParseResponse(true).R().EnableTrace().SetQueryParams(params).SetBody(data).SetHeader("Content-Type", "application/json").Post(url)
	return resp, err
}

func MonitorMsg(murl, loc, msg string) error {
	body := map[string]string{
		"timestamp": strconv.FormatInt(time.Now().Unix(), 10),
		"subsys":    "test",
		"module":    "-",
		"level":     "error",
		"msg":       msg,
		"account":   "invalid",
	}
	url := murl + loc
	resp, err := Post(url, nil, body)
	if err != nil {
		return fmt.Errorf("log alarm %s msg %s with err %s response %v", url, msg, err, resp)
	}
	return nil
}
