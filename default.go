package tpl

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"sync"
	"time"

	"k8s.io/klog"
)

const defaultTimeoutSec = 10

// SendFunc is function to send http request, and get response.
type SendFunc func(url string, opt *Option) ([]byte, error)

func tryCacheOrSend(url string, opt *Option, cache *sync.Map, sendFunc SendFunc) ([]byte, error) {
	var prefix string
	if opt.HTTPS {
		prefix = "https://"
	} else {
		prefix = "http://"
	}
	url = prefix + opt.IP + url
	v, ok := cache.Load(url)
	if ok {
		return v.([]byte), nil
	}
	b, err := sendFunc(url, opt)
	if err != nil {
		return nil, err
	}
	cache.LoadOrStore(url, b)
	return b, err
}

var hc = &http.Client{Transport: &http.Transport{
	TLSClientConfig: &tls.Config{
		InsecureSkipVerify: true,
	},
}}

func determineTimeout(sec int) time.Duration {
	if sec <= 0 {
		sec = defaultTimeoutSec
	}
	return time.Duration(sec) * time.Second
}

func sendOrLog(url string, opt *Option) ([]byte, error) {
	b, err := send(url, opt)
	if err != nil {
		klog.Errorf("send failed: %v", err)
		return nil, err
	}
	return b, nil
}

func send(url string, opt *Option) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), determineTimeout(opt.TimeoutSec))
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(opt.User, opt.Password)

	rsp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	return io.ReadAll(rsp.Body)
}
