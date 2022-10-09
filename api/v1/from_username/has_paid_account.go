package from_username

import (
	"net/http"
)

func HasPaidAccount(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Query().Get("username")) == 0 {
		w.WriteHeader(400)
		w.Write([]byte("You need to provide a username"))
		return
	}

	GetUuidFromUsername(r.URL.Query().Get("username"), w) // This function will handle the error
}
