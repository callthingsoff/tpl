package tpl

import (
	"errors"
	"fmt"
	"sync"

	"github.com/callthingsoff/gjson"
)

type Option struct {
	IP         string `json:"ip"`
	User       string `json:"user"`
	Password   string `json:"password"`
	TimeoutSec int    `json:"timeout"`

	CacheB         *sync.Map                                     `json:"-"`
	CacheR         *sync.Map                                     `json:"-"`
	TryCacheOrSend func(url string, opt *Option) ([]byte, error) `json:"-"`
}

type Fetcher struct {
	opt *Option
}

func NewFetcher(opt *Option) *Fetcher {
	if opt == nil {
		return nil
	}
	opt.CacheB = new(sync.Map)
	opt.CacheR = new(sync.Map)
	if opt.TryCacheOrSend == nil {
		opt.TryCacheOrSend = DefaultTryCacheOrSend
	}
	return &Fetcher{
		opt: opt,
	}
}

func (f *Fetcher) Fetch(tpl *Template) (any, error) {
	if tpl == nil {
		return nil, errors.New("nil Template")
	}

	m := map[string]any{}
	for _, t := range tpl.Template {
		b, err := f.opt.TryCacheOrSend(t.URL, f.opt)
		if err != nil {
			return nil, err
		}

		r := gjson.ParseBytes(b)
		for _, x := range t.Group {
			v := r.Get(x.JSONPath, f.opt)
			if !v.Exists() {
				return nil, fmt.Errorf("%q not found", x.JSONPath)
			}
			m[x.ID] = v.Value()
		}
	}
	return m, nil
}
