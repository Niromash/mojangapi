package from_username

import (
	"encoding/json"
	"net/http"
)

func UsernameUuid(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Query().Get("username")) == 0 {
		w.WriteHeader(400)
		w.Write([]byte("You need to provide a username"))
		return
	}

	usernameToUuidPayload := GetUuidFromUsername(r.URL.Query().Get("username"), w)
	if usernameToUuidPayload == nil { // if the uuid is empty, the error has already been handled
		return
	}

	bytes, err := json.Marshal(usernameToUuidPayload)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(bytes)
}
