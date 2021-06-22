package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

func tokenRevoke(rw http.ResponseWriter, r *http.Request) {
	var (
		err  error
		vars = mux.Vars(r)
	)
	if vars["key"] == "" {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("{\"status\":\"Bad Request\",\"stauscode\":400}"))
	}

	url := fmt.Sprint("https://id.twitch.tv/oauth2/revoke",
		"?client_id="+config.Twitch.Api.ClientId,
		"&token="+vars["key"],
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
