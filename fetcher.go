package tpl

import (
	"errors"
	"fmt"
	"sync"

	"github.com/callthingsoff/gjson"
)

// Option determines how to send http requests.
type Option struct {
	HTTPS      bool   `json:"https"`
	IP         string `json:"ip"`
	User       string `json:"user"`
	Password   string `json:"password"`
	TimeoutSec int    `json:"timeout"`
}

// Fetcher holds context while fetching resources.
type Fetcher struct {
	opt *Option

	sendFunc SendFunc

	cacheB *sync.Map
	cacheR *sync.Map
}

// NewFetcher creates a Fetcher by opt Option, if sendFunc
// is nil, default send function will be used.
func NewFetcher(opt *Option, sendFunc SendFunc) *Fetcher {
	if opt == nil {
		return nil
	}
	if sendFunc == nil {
		sendFunc = send
	}
	return &Fetcher{
		opt:      opt,
		sendFunc: sendFunc,
		cacheB:   new(sync.Map),
		cacheR:   new(sync.Map),
	}
}

// Fetch fetches all the items from tpl.
func (f *Fetcher) Fetch(tpl *Template) (any, error) {
	if tpl == nil {
		return nil, errors.New("nil Template")
	}

	m := map[string]any{}
	for _, t := range tpl.Template {
		b, err := tryCacheOrSend(t.URL, f.opt, f.cacheB, f.sendFunc)
		if err != nil {
			return nil, err
		}

		r := gjson.ParseBytes(b)
		for _, x := range t.Group {
			v := r.Get(x.JSONPath, f)
			if !v.Exists() {
				return nil, fmt.Errorf("%q not found", x.JSONPath)
			}
			m[x.ID] = v.Value()
		}
	}
	return m, nil
}
