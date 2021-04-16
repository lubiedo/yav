package models

import (
	"fmt"
	"net/http"
	"strings"
)

type Headers map[string]string

func (h *Headers) Set(s string) error {
	parsed := strings.SplitN(s, ":", 2)
	if len(parsed) != 2 {
		return fmt.Errorf("invalid HTTP header")
	}

	if *h == nil {
		*h = make(map[string]string)
	}
	(*h)[parsed[0]] = strings.Trim(parsed[1], " ")
	return nil
}

func (h *Headers) String() string {
	var ret string
	for k, v := range *h {
		ret = k + ": " + v + "\n"
	}
	return ret
}

func (h *Headers) AddToHttp(header *http.Header) {
	for k, v := range *h {
		header.Add(k, v)
	}
}
