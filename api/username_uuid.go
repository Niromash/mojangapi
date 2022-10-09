package api

import (
	"io"
	"net/http"
)

func UsernameUuid(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Query().Get("username")) == 0 {
		w.WriteHeader(400)
		w.Write([]byte("You need to provide a username"))
		return
	}

	// req to https://api.mojang.com/users/profiles/minecraft/username
	resp, err := http.Get("https://api.mojang.com/users/profiles/minecraft/" + r.URL.Query().Get("username"))
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	defer resp.Body.Close()

	// copy headers from the resp to the resp
	for k, v := range resp.Header {
		w.Header()[k] = v
	}
	io.Copy(w, resp.Body)
}
