package from_username

import (
	"net/http"
)

func UsernameUuid(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Query().Get("username")) == 0 {
		w.WriteHeader(400)
		w.Write([]byte("You need to provide a username"))
		return
	}

	username := GetUuidFromUsername(r.URL.Query().Get("username"), w)
	if username == "" { // if the uuid is empty, the error has already been handled
		return
	}

	w.Write([]byte(username))
}
