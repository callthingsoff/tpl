package tpl

import (
	"sync"

	"github.com/callthingsoff/gjson"
)

type Extra struct {
	Opt *Option

	CacheB *sync.Map
	CacheR *sync.Map
}

func init() {
	parseExtra := func(extra ...any) *Extra {
		if len(extra) != 1 {
			return nil
		}
		e, _ := extra[0].(*Extra)
		return e
	}

	gjson.AddModifier("store", func(json, arg string, extra ...any) string {
		e := parseExtra(extra...)
		if e.CacheR == nil {
			return ""
		}
		e.CacheR.LoadOrStore(arg, json)
		return json
	})

	gjson.AddModifier("load", func(json, arg string, extra ...any) string {
		e := parseExtra(extra...)
		if e.CacheR == nil {
			return ""
		}
		v, ok := e.CacheR.Load(arg)
		if !ok {
			return ""
		}
		return v.(string)
	})

	gjson.AddModifier("url", func(json, arg string, extra ...any) string {
		e := parseExtra(extra...)
		if e.CacheR == nil || e.CacheB == nil {
			return ""
		}
		r := gjson.Parse(json)

		b, err := tryCacheOrSend(r.String(), e.Opt, e.CacheB)
		if err != nil {
			return ""
		}
		return string(b)
	})
}
