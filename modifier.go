package tpl

import (
	"github.com/callthingsoff/gjson"
)

func init() {
	parseExtra := func(extra ...any) *Fetcher {
		if len(extra) != 1 {
			return nil
		}
		e, _ := extra[0].(*Fetcher)
		return e
	}

	gjson.AddModifier("store", func(json, arg string, extra ...any) string {
		f := parseExtra(extra...)
		if f.cacheR == nil {
			return ""
		}
		f.cacheR.LoadOrStore(arg, json)
		return json
	})

	gjson.AddModifier("load", func(json, arg string, extra ...any) string {
		f := parseExtra(extra...)
		if f.cacheR == nil {
			return ""
		}
		v, ok := f.cacheR.Load(arg)
		if !ok {
			return ""
		}
		return v.(string)
	})

	gjson.AddModifier("url", func(json, arg string, extra ...any) string {
		f := parseExtra(extra...)
		if f.cacheR == nil || f.cacheB == nil {
			return ""
		}
		r := gjson.Parse(json)

		b, err := tryCacheOrSend(r.String(), f.opt, f.cacheB, f.sendFunc)
		if err != nil {
			return ""
		}
		return string(b)
	})
}
