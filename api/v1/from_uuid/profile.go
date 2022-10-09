package from_uuid

import (
	"encoding/json"
	"github.com/niromash/verceltest/api/v1/from_username"
	"net/http"
)

func ProfileFromUuid(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Query().Get("uuid")) == 0 {
		w.WriteHeader(400)
		w.Write([]byte("You need to provide a uuid"))
		return
	}

	profileResponse := from_username.GetProfileFromUuid(r.URL.Query().Get("uuid"), w)

	if err := json.NewEncoder(w).Encode(profileResponse); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
}
