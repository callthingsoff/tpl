package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/callthingsoff/gjson"
)

type Option struct {
	IP         string `json:"ip"`
	User       string `json:"user"`
	Password   string `json:"password"`
	TimeoutSec int    `json:"timeout"`
}

type Fetcher struct {
	opt    *Option
	cacheB *sync.Map
	cacheR *sync.Map
}

func NewFetcher(opt *Option) *Fetcher {
	if opt == nil {
		return nil
	}
	if opt.TimeoutSec <= 0 {
		opt.TimeoutSec = 3
	}
	return &Fetcher{
		opt:    opt,
		cacheB: new(sync.Map),
		cacheR: new(sync.Map),
	}
}

func (f *Fetcher) Fetch(tpl *Template) (any, error) {
	if tpl == nil {
		return nil, errors.New("nil Template")
	}

	m := map[string]any{}
	for _, t := range tpl.Template {
		b, err := tryCacheOrSend(t.URL, f.opt, f.cacheB)
		if err != nil {
			return nil, err
		}

		r := gjson.ParseBytes(b)
		for _, x := range t.Val {
			v := r.Get(x.JSONPath, f.opt, f.cacheB, f.cacheR)
			if !v.Exists() {
				return nil, fmt.Errorf("%s does not exit", x.JSONPath)
			}
			m[x.ID] = v.Value()
		}
	}
	return m, nil
}

func tryCacheOrSend(url string, opt *Option, cacheB *sync.Map) ([]byte, error) {
	url = "http://" + opt.IP + url
	v, ok := cacheB.Load(url)
	if ok {
		return v.([]byte), nil
	}
	b, err := send(url, opt)
	if err != nil {
		return nil, err
	}
	cacheB.LoadOrStore(url, b)
	return b, err
}

var hc = &http.Client{Transport: &http.Transport{
	TLSClientConfig: &tls.Config{
		InsecureSkipVerify: true,
	},
}}

func send(url string, opt *Option) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(opt.TimeoutSec)*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	rsp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	return io.ReadAll(rsp.Body)
}
