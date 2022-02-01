package httpclient

import "time"

type Options struct {
	// gives you timeout for request
	Timeout time.Duration

	// bool flag which help us in configuring proxy
	UseProxy bool

	// url need to do the proxy
	ProxyURL string
}
