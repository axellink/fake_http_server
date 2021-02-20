package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type ReceivedRequest struct {
	Date    string            `json:Date`
	Method  string            `json:Method`
	Url     string            `json:URL`
	Proto   string            `json:proto`
	Headers map[string]string `json:Headers`
	Params  map[string]string `json:Params`
	Body    string            `json:Body`
}

func JSONEncode(r *http.Request) string {
	var rr ReceivedRequest

	// Get simple values
	rr.Date = time.Now().String()
	rr.Method = r.Method
	rr.Url = r.URL.EscapedPath()
	rr.Proto = r.Proto

	// Get headers
	rr.Headers = make(map[string]string)
	for key, value := range r.Header {
		v := ""
		for _, i := range value {
			v = v + i
		}
		rr.Headers[key] = v
	}

	// Get Params
	rr.Params = make(map[string]string)
	for key, value := range r.URL.Query() {
		v := ""
		for _, i := range value {
			v = v + i
		}
		rr.Params[key] = v
	}

	// Get body
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	rr.Body = string(bytes)

	// Try to convert in json
	jsonrr, err := json.Marshal(rr)
	if err != nil {
		return "{\"Error\":\"AH !\"}"
	}

	return string(jsonrr)
}
