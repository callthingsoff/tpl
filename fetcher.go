package tpl

import (
	"errors"
	"fmt"
	"sync"

	"github.com/callthingsoff/gjson"
)

// Option holds options to request http/https.
type Option struct {
	HTTPS      bool   `json:"https"`
	IP         string `json:"ip"`
	User       string `json:"user"`
	Password   string `json:"password"`
	TimeoutSec int    `json:"timeoutSec"`
}

// Fetcher holds contexts while fetching resources.
type Fetcher struct {
	opt *Option

	sendFunc SendFunc

	cacheB *sync.Map
	cacheR *sync.Map
}

// NewFetcher creates a Fetcher by opt, if sendFunc
// is nil, sendOrLog function will be used.
func NewFetcher(opt *Option, sendFunc SendFunc) *Fetcher {
	if opt == nil {
		return nil
	}
	if sendFunc == nil {
		sendFunc = sendOrLog
	}
	return &Fetcher{
		opt:      opt,
		sendFunc: sendFunc,
		cacheB:   new(sync.Map),
		cacheR:   new(sync.Map),
	}
}

// Fetch fetches all the resources from tpl, and reports an error if failed.
func (f *Fetcher) Fetch(tpl *Template) (any, error) {
	if tpl == nil {
		return nil, errors.New("nil Template")
	}

	m := map[string]any{}
	for _, it := range tpl.Items {
		b, err := tryCacheOrSend(it.URL, f.opt, f.cacheB, f.sendFunc)
		if err != nil {
			return nil, err
		}

		r := gjson.ParseBytes(b)
		for _, x := range it.Group {
			v := r.Get(x.JSONPath, f)
			if !v.Exists() {
				return nil, fmt.Errorf("%q: %q, not found in %q", x.ID, x.Name, x.JSONPath)
			}
			m[x.ID] = v.Value()
		}
	}
	return m, nil
}
