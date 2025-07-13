package tpl_test

import (
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
	}).Fetch(template)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%v", val)
}
