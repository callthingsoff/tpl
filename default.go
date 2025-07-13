package tpl

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"sync"
	"time"
)

const defaultTimeoutSec = 10

type SendFunc func(url string, opt *Option) ([]byte, error)

func tryCacheOrSend(url string, opt *Option, cache *sync.Map) ([]byte, error) {
	url = "http://" + opt.IP + url
	v, ok := cache.Load(url)
	if ok {
		return v.([]byte), nil
	}
	b, err := opt.SendFunc(url, opt)
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
