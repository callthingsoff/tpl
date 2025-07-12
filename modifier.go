package main

import (
	"sync"

	"github.com/callthingsoff/gjson"
)

var light gjson.Result

func init() {
	light = gjson.Parse(`{}`)

	parseExtra := func(extra ...any) (*Option, *sync.Map, *sync.Map) {
		if len(extra) != 3 {
			return nil, nil, nil
		}
		opt, ok := extra[0].(*Option)
		if !ok {
			return nil, nil, nil
		}
		cacheB, ok := extra[1].(*sync.Map)
		if !ok {
			return nil, nil, nil
		}
		cacheR, ok := extra[2].(*sync.Map)
		if !ok {
			return nil, nil, nil
		}
		return opt, cacheB, cacheR
	}

	gjson.AddModifier("store", func(json, arg string, extra ...any) string {
		_, _, cacheR := parseExtra(extra...)
		if cacheR == nil {
			return ""
		}
		cacheR.LoadOrStore(arg, json)
		return json
	})

	gjson.AddModifier("load", func(json, arg string, extra ...any) string {
		_, _, cacheR := parseExtra(extra...)
		if cacheR == nil {
			return ""
		}
		v, ok := cacheR.Load(arg)
		if !ok {
			return ""
		}
		return v.(string)
	})
}
