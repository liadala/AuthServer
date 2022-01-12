package config

import (
	"encoding/json"
	"log"
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
	} `json:"twitch"`
	Web struct {
		Listen string `json:"ip"`
		Port   int    `json:"port"`
	} `json:"web"`
}

var Config configStruct

func Load(path string) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(file, &Config)
	if err != nil {
		log.Fatal(err)
	}
}
