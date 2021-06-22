package main

import (
	"encoding/json"
	"os"
)

type configStruct struct {
	Twitch struct {
		Api struct {
			ClientId         string   `json:"clientid"`
			ClientSecret     string   `json:"clientsecret"`
			RedirectURL      string   `json:"redirecturl"`
			AviableScopeList []string `json:"aviablescopelist"`
		} `json:"api"`
	} `json:"twitch:"`
	HttpRequest struct {
		RateLimit      int `json:"ratelimit"`
		RegenerateRate int `json:"regeneraterate"`
	} `json:"httprequest"`
}

var config configStruct

func loadConfig(path string) (configStruct, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return configStruct{}, err
	}
	var jsondata configStruct
	err = json.Unmarshal(file, &jsondata)
	if err != nil {
		return configStruct{}, err
	}
	return jsondata, nil
}
