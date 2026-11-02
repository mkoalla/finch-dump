package main

import (
	"finchd/lib/api"
	"log"
	"net/url"
)

func main() {
	r, err := api.Get("", &url.Values{})
	log.Println(r, err)
}
