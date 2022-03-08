package httpclient

import (
	"time"

	"github.com/rohanraj7316/logger"
)

var (
	REQUEST_TIMEOUT = "2s"
)

type Config struct {
	// Timeout gives you timeout for request
	// Default: 30s
	Timeout time.Duration

	// bool flag which help us in configuring proxy
	// Default: false
	UseProxy bool

	// url need to do the proxy
	// Default: nil
	ProxyURL string

	// LogReqResEnable helps in logging request & responses.
	// Default true
	LogReqResEnable bool

	// LogReqResBodyEnable helps in logging request and responses body
	// Default true
	LogReqResBodyEnable bool
}

var ConfigDefault = Config{
	UseProxy:            false,
	LogReqResEnable:     true,
	LogReqResBodyEnable: true,
}

func configDefault(config ...Config) Config {
	timeout, err := time.ParseDuration(REQUEST_TIMEOUT)
	if err != nil {
		logger.Error(err.Error())
		ConfigDefault.Timeout = 2 * time.Second
	} else {
		ConfigDefault.Timeout = timeout
	}

	if len(config) < 1 {
		return ConfigDefault
	}

	cfg := config[0]

	return cfg
}
