package api

import (
	"fmt"
	"io"
	"net/http"
)

func Ip(w http.ResponseWriter, r *http.Request) {
	// req get to https://ipinfo.io/json
	resp, err := http.Get("https://ipinfo.io/json")
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		return
	}
	defer resp.Body.Close()
	io.Copy(w, resp.Body)
}
