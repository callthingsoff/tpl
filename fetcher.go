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

	SendFunc SendFunc `json:"-"`
}

type Fetcher struct {
	opt *Option

	cacheB *sync.Map
	cacheR *sync.Map
}

func NewFetcher(opt *Option) *Fetcher {
	if opt == nil {
		return nil
	}
	if opt.SendFunc == nil {
		opt.SendFunc = send
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
		for _, x := range t.Group {
			v := r.Get(x.JSONPath, &extraType{opt: f.opt, cacheB: f.cacheB, cacheR: f.cacheR})
			if !v.Exists() {
				return nil, fmt.Errorf("%q not found", x.JSONPath)
			}
			m[x.ID] = v.Value()
		}
	}
	return m, nil
}
