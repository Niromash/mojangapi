package api

import (
	"io"
	"net/http"
)

func Profile(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Query().Get("uuid")) == 0 {
		w.WriteHeader(400)
		w.Write([]byte("You need to provide a uuid"))
		return
	}
	// req to https://sessionserver.mojang.com/session/minecraft/profile/uuid
	resp, err := http.Get("https://sessionserver.mojang.com/session/minecraft/profile/" + r.URL.Query().Get("uuid"))
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
