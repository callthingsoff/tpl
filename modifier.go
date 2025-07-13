package tpl

import (
	"sync"

	"github.com/callthingsoff/gjson"
)

type extraType struct {
	opt *Option

	cacheB *sync.Map
	cacheR *sync.Map
}

func init() {
	parseExtra := func(extra ...any) *extraType {
		if len(extra) != 1 {
			return nil
		}
		e, _ := extra[0].(*extraType)
		return e
	}

	gjson.AddModifier("store", func(json, arg string, extra ...any) string {
		e := parseExtra(extra...)
		if e.cacheR == nil {
			return ""
		}
		e.cacheR.LoadOrStore(arg, json)
		return json
	})

	gjson.AddModifier("load", func(json, arg string, extra ...any) string {
		e := parseExtra(extra...)
		if e.cacheR == nil {
			return ""
		}
		v, ok := e.cacheR.Load(arg)
		if !ok {
			return ""
		}
		return v.(string)
	})

	gjson.AddModifier("url", func(json, arg string, extra ...any) string {
		e := parseExtra(extra...)
		if e.cacheR == nil || e.cacheB == nil {
			return ""
		}
		r := gjson.Parse(json)

		b, err := tryCacheOrSend(r.String(), e.opt, e.cacheB)
		if err != nil {
			return ""
		}
		return string(b)
	})
}
