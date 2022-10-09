package v1

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

type UsernameToUuidPayload struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type ProfilePayload struct {
	Id         string       `json:"id"`
	Name       string       `json:"name"`
	Properties []Properties `json:"properties"`
	SkinUrl    string       `json:"skinUrl,omitempty"`
}

type Properties struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	Signature string `json:"signature"`
}

type SkinPropertyValueDecoded struct {
	Url string `json:"url"`
}

type ProfileResponse struct {
	UsernameToUuidPayload
	ProfilePayload
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

	usernameToUuidPayload := UsernameToUuidPayload{}
	if err = json.NewDecoder(resp.Body).Decode(&usernameToUuidPayload); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	// req to https://sessionserver.mojang.com/session/minecraft/profile/uuid
	resp, err = http.Get(fmt.Sprintf("https://sessionserver.mojang.com/session/minecraft/profile/%s?unsigned=false", usernameToUuidPayload.Id))
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	defer resp.Body.Close()
	var profilePayload ProfilePayload
	if err = json.NewDecoder(resp.Body).Decode(&profilePayload); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	if len(profilePayload.Properties) > 0 {
		decodedProperty, err := base64.StdEncoding.DecodeString(profilePayload.Properties[0].Value)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		skinPropertyValueDecoded := SkinPropertyValueDecoded{}
		if err = json.Unmarshal(decodedProperty, &skinPropertyValueDecoded); err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		profilePayload.SkinUrl = skinPropertyValueDecoded.Url
	}

	if err = json.NewEncoder(w).Encode(ProfileResponse{
		usernameToUuidPayload,
		profilePayload,
	}); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
}
