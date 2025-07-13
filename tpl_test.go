package tpl_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/callthingsoff/tpl"
)

func TestFetch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.RequestURI {
		case "/a/b/cpu":
			w.Write([]byte(`{"a": {"b": {"c": [1,2,3]}}}`))
		case "/x/y/memory":
			w.Write([]byte(`{"a": {"b": {"c": [3]}}}`))
		case "/storage":
			w.Write([]byte(`{"storage": "/storages"}`))
		case "/storages":
			w.Write([]byte(`["/storage1", "/storage2"]`))
		case "/storage1":
			w.Write([]byte(`{"s":"v1"}`))
		case "/storage2":
			w.Write([]byte(`{"s":"v2"}`))
		}
	}))
	defer server.Close()

	b, err := os.ReadFile("tpl.yaml")
	if err != nil {
		t.Fatal(err)
	}

	template, err := tpl.ParseBytes(b)
	if err != nil {
		t.Fatal(err)
	}
	got, err := tpl.NewFetcher(&tpl.Option{
		HTTPS:      false,
		IP:         server.Listener.Addr().String(),
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
	want := map[string]any{
		"cpu": []int{
			1,
			2,
			3,
		},
		"memory": []int{
			3,
		},
		"ncpu":    3,
		"nmemory": 1,
		"storage": []string{
			"/storage1",
			"/storage2",
		},
		"storages": []string{
			"v1",
			"v2",
		},
		"sumcpu":    6,
		"sumdivcpu": 2,
		"summemory": 3,
	}
	x, err := json.Marshal(got)
	if err != nil {
		t.Fatal(err)
	}
	y, err := json.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(x, y) {
		t.Fatalf("not equal:\nwant:\n%s\ngot:\n%s\n", x, y)
	}
}
