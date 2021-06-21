package main

import (
	"encoding/json"
	"os"
)

type configStruct struct {
	ClientId     string
	ClientSecret string
	RedirectURL  string
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
