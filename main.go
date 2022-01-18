package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Printf("%s\n", body)
	})
	http.ListenAndServe(":80", nil)
}
