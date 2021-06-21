package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func startWebserver() {

	http.HandleFunc("/auth", func(rw http.ResponseWriter, r *http.Request) {
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
	})

	http.HandleFunc("/api/refresh", func(rw http.ResponseWriter, r *http.Request) {
		var (
		//err error
		)
		req := strings.ReplaceAll(r.URL.Path, "/api/refresh", "")
		/*
			url := fmt.Sprint("https://id.twitch.tv/oauth2/token",
				"?client_id="+config.ClientId,
				"&client_secret="+config.ClientSecret,
				"&grant_type=refresh_token",
				"&refresh_token="+req,
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
		*/
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		//rw.Write(b)
		rw.Write([]byte(req))
	})

	http.Handle("/", http.FileServer(http.Dir("./web")))
	log.Fatal(http.ListenAndServe(":80", nil))
}
