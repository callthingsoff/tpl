package tpl_test

import (
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/callthingsoff/tpl"
)

func TestRun(t *testing.T) {
	b, err := os.ReadFile("tpl.yaml")
	if err != nil {
		t.Fatal(err)
	}

	template, err := tpl.ParseBytes(b)
	if err != nil {
		t.Fatal(err)
	}
	val, err := tpl.NewFetcher(&tpl.Option{
		IP:         "localhost",
		User:       "x",
		Password:   "y",
		TimeoutSec: 3,
	}, func(url string, opt *tpl.Option) ([]byte, error) {
		rsp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer rsp.Body.Close()
		return io.ReadAll(rsp.Body)
	}).Fetch(template)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%v", val)
}
