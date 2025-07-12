package main

import (
	"fmt"
	"log"
)

func main() {
	tpl, err := ParseTemplate("tpl.yaml")
	if err != nil {
		log.Fatal(err)
	}

	v, err := NewFetcher(&Option{
		IP:       "localhost",
		User:     "x",
		Password: "y",
	}).Fetch(tpl)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(v)
}
