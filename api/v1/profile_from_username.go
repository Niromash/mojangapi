package v1

import (
	"encoding/json"
	"io"
	"net/http"
)

type payload struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

func ProfileFromUsername(w http.ResponseWriter, r *http.Request) {
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

	p := payload{}
	if err = json.NewDecoder(resp.Body).Decode(&p); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	// req to https://sessionserver.mojang.com/session/minecraft/profile/uuid
	resp, err = http.Get("https://sessionserver.mojang.com/session/minecraft/profile/" + p.Id)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	defer resp.Body.Close()

	io.Copy(w, resp.Body)
}
