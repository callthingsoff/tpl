package main

import (
	"fmt"
	"log"
	"sync"
	"testing"

	"github.com/callthingsoff/gjson"
)

func TestXxx(t *testing.T) {
	cache := new(sync.Map)

	const s = `{"a": {"b": {"c": [1,2,3]}}}`

	r := gjson.Parse(s)
	v1 := r.Get("a.b.c|@store:arr", cache).String()
	fmt.Println(v1)

	v1 = r.Get("a.b.c.#", cache).String()
	fmt.Println(v1)

	tpl, err := ParseTemplate("tpl.yaml")
	if err != nil {
		log.Fatal(err)
	}
	for _, t := range tpl.Template {
		fmt.Println(t.URL)
		fmt.Println(t.Val)
	}
}
