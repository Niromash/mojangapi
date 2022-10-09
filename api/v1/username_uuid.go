package v1

import (
	"fmt"
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

	w.WriteHeader(200)
	w.Header().Set("Cache-Control", fmt.Sprintf("max-age=0, s-max-age=%d", 86400))
	io.Copy(w, resp.Body)
}
