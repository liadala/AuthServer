package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

func tokenRefresh(rw http.ResponseWriter, r *http.Request) {
	var (
		err  error
		vars = mux.Vars(r)
	)
	if vars["key"] == "" {
		rw.Header().Set("Content-Type", "application/json")
		//rw.Header().Set("RateLimit-Limit", config.HttpRequest.RateLimit)
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("{\"status\":\"Bad Request\",\"stauscode\":400}"))
	}

	url := fmt.Sprint("https://id.twitch.tv/oauth2/token",
		"?client_id="+config.Twitch.Api.ClientId,
		"&client_secret="+config.Twitch.Api.ClientSecret,
		"&grant_type=refresh_token",
		"&refresh_token="+vars["key"],
	)
	request, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("{\"status\":\"Bad Request\",\"stauscode\":400}"))
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("{\"status\":\"Bad Request\",\"stauscode\":400}"))
	}
	defer response.Body.Close()

	b, err := io.ReadAll(response.Body)
	if err != nil {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("{\"status\":\"Bad Request\",\"stauscode\":400}"))
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(b)
}
