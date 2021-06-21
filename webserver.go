package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func startWebserver() {
	router := mux.NewRouter()
	fs := http.FileServer(http.Dir("./web"))

	router.HandleFunc("/api/refresh/{key}", refresh)
	router.HandleFunc("api/revoke/{key}", revoke)
	router.HandleFunc("/auth", auth)
	router.PathPrefix("/").Handler(fs)

	//http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":80", router))
}

func revoke(rw http.ResponseWriter, r *http.Request) {
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
		"?client_id="+config.ClientId,
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

func refresh(rw http.ResponseWriter, r *http.Request) {
	var (
		err  error
		vars = mux.Vars(r)
	)
	if vars["key"] == "" {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("{\"status\":\"Bad Request\",\"stauscode\":400}"))
	}

	url := fmt.Sprint("https://id.twitch.tv/oauth2/token",
		"?client_id="+config.ClientId,
		"&client_secret="+config.ClientSecret,
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

func auth(rw http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["code"]
	if !ok || len(keys[0]) < 1 {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("{\"status\":\"Bad Request\",\"stauscode\":400}"))
	}
	url := fmt.Sprint("https://id.twitch.tv/oauth2/token",
		"?client_id="+config.ClientId,
		"&client_secret="+config.ClientSecret,
		"&code="+keys[0],
		"&grant_type=authorization_code",
		"&redirect_uri="+config.RedirectURL,
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
