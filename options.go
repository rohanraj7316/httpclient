package httpclient

import (
	"time"

	"github.com/rohanraj7316/logger"
)

type Options struct {
	// gives you timeout for request
	Timeout time.Duration

	// bool flag which help us in configuring proxy
	UseProxy bool

	// url need to do the proxy
	ProxyURL string

	// default false. true when you need request+response logging
	LogReqResEnable bool

	// default false. true when you need request+response body logging
	LogReqResBodyEnable bool

	LoggerOptions *logger.Options
}
