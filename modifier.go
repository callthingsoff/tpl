package tpl

import (
	"github.com/callthingsoff/gjson"
)

func init() {
	parseExtra := func(extra ...any) *Option {
		if len(extra) != 1 {
			return nil
		}
		e, _ := extra[0].(*Option)
		return e
	}

	gjson.AddModifier("store", func(json, arg string, extra ...any) string {
		opt := parseExtra(extra...)
		if opt.CacheR == nil {
			return ""
		}
		opt.CacheR.LoadOrStore(arg, json)
		return json
	})

	gjson.AddModifier("load", func(json, arg string, extra ...any) string {
		opt := parseExtra(extra...)
		if opt.CacheR == nil {
			return ""
		}
		v, ok := opt.CacheR.Load(arg)
		if !ok {
			return ""
		}
		return v.(string)
	})

	gjson.AddModifier("url", func(json, arg string, extra ...any) string {
		opt := parseExtra(extra...)
		if opt.CacheR == nil {
			return ""
		}
		r := gjson.Parse(json)

		b, err := opt.TryCacheOrSend(r.String(), opt)
		if err != nil {
			return ""
		}
		return string(b)
	})
}
