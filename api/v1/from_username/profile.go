package from_username

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

func ProfileFromUsername(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Query().Get("username")) == 0 {
		w.WriteHeader(400)
		w.Write([]byte("You need to provide a username"))
		return
	}

	usernameToUuidPayload := GetUuidFromUsername(r.URL.Query().Get("username"), w)
	if usernameToUuidPayload == nil { // if the uuid is empty, the error has already been handled
		return
	}

	profileResponse := GetProfileFromUuid(usernameToUuidPayload.Id, w)

	if err := json.NewEncoder(w).Encode(profileResponse); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
}

func GetUuidFromUsername(username string, w http.ResponseWriter) *UsernameToUuidPayload {
	// req to https://api.mojang.com/users/profiles/minecraft/username
	resp, err := http.Get("https://api.mojang.com/users/profiles/minecraft/" + username)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode == 429 {
		w.WriteHeader(429)
		w.Write([]byte("Too many requests"))
		return nil
	}
	if resp.StatusCode != 200 {
		w.WriteHeader(404)
		w.Write([]byte("User not found"))
		return nil
	}

	usernameToUuidPayload := UsernameToUuidPayload{}
	if err = json.NewDecoder(resp.Body).Decode(&usernameToUuidPayload); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return nil
	}
	return &usernameToUuidPayload
}

func GetProfileFromUuid(uuid string, w http.ResponseWriter) *ProfilePayload {
	// req to https://sessionserver.mojang.com/session/minecraft/profile/uuid
	resp, err := http.Get(fmt.Sprintf("https://sessionserver.mojang.com/session/minecraft/profile/%s?unsigned=false", uuid))
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode == 429 {
		w.WriteHeader(429)
		w.Write([]byte("Too many requests"))
		return nil
	}
	if resp.StatusCode != 200 {
		w.WriteHeader(404)
		w.Write([]byte("User not found"))
		return nil
	}

	var profilePayload ProfilePayload
	if err = json.NewDecoder(resp.Body).Decode(&profilePayload); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return nil
	}

	if len(profilePayload.Properties) > 0 {
		decodedProperty, err := base64.StdEncoding.DecodeString(profilePayload.Properties[0].Value)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return nil
		}

		skinPropertyValueDecoded := SkinPropertyValueDecoded{}
		if err = json.Unmarshal(decodedProperty, &skinPropertyValueDecoded); err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return nil
		}

		profilePayload.SkinUrl = skinPropertyValueDecoded.Textures.Skin.Url
	}
	return &profilePayload
}

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
	Textures struct {
		Skin struct {
			Url string `json:"url"`
		} `json:"SKIN"`
	} `json:"textures"`
}

type ProfileResponse struct {
	ProfilePayload
}
