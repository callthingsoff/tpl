package tpl

import (
	"errors"
	"fmt"
	"sync"

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

	gjson.AddModifier("sum", func(json, arg string, extra ...any) string {
		r := gjson.Parse(json)
		if !r.IsArray() {
			return ""
		}

		var s int64
		for _, x := range r.Array() {
			s += x.Int()
		}
		return fmt.Sprintf("%d", s)
	})

	gjson.AddModifier("div", func(json, arg string, extra ...any) string {
		d := gjson.Parse(arg).Int()
		if d == 0 {
			return ""
		}

		r := gjson.Parse(json)
		v := r.Int() / d
		return fmt.Sprintf("%d", v)
	})

	gjson.AddModifier("url", func(json, arg string, extra ...any) string {
		r := gjson.Parse(json)
		if !r.Exists() {
			return ""
		}

		f := parseExtra(extra...)
		if f.cacheB == nil {
			return ""
		}

		b, err := tryCacheOrSend(r.String(), f.opt, f.cacheB, f.sendFunc)
		if err != nil {
			return ""
		}
		return string(b)
	})

	gjson.AddModifier("urls", func(json, arg string, extra ...any) string {
		r := gjson.Parse(json)
		if !r.IsArray() {
			return ""
		}

		f := parseExtra(extra...)
		if f.cacheB == nil {
			return ""
		}

		if arg == "async" {
			return asyncSend(f, r.Array())
		}
		return syncSend(f, r.Array())
	})

}

func syncSend(f *Fetcher, arr []gjson.Result) string {
	bs := []byte("[")
	for _, x := range arr {
		b, err := tryCacheOrSend(x.String(), f.opt, f.cacheB, f.sendFunc)
		if err != nil {
			return ""
		}
		bs = append(bs, b...)
		bs = append(bs, ',')
	}
	bs = append(bs, ']')
	return string(bs)
}

func asyncSend(f *Fetcher, arr []gjson.Result) string {
	vs := make([][]byte, len(arr))
	errs := make([]error, len(arr))

	var wg sync.WaitGroup
	for i, x := range arr {
		wg.Add(1)
		go func() {
			defer wg.Done()

			b, err := tryCacheOrSend(x.String(), f.opt, f.cacheB, f.sendFunc)
			vs[i] = b
			errs[i] = err
		}()
	}

	wg.Wait()

	if err := errors.Join(errs...); err != nil {
		return ""
	}

	bs := []byte("[")
	for _, b := range vs {
		bs = append(bs, b...)
		bs = append(bs, ',')
	}
	bs = append(bs, ']')
	return string(bs)
}
