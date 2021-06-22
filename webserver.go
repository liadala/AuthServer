package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func startWebserver() {
	router := mux.NewRouter()
	static := http.FileServer(http.Dir("./static"))

	router.HandleFunc("/api/refresh/{key}", tokenRefresh)
	router.HandleFunc("api/revoke/{key}", tokenRevoke)
	router.HandleFunc("/auth", tokenCreate)

	router.PathPrefix("/static").Handler(http.StripPrefix("/static/", static))
	router.PathPrefix("/").HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("RateLimit-Limit", fmt.Sprint(rateLimit.RateLimit_Limit))
		rw.Header().Set("RateLimit-Remaining", fmt.Sprint(rateLimit.RateLimit_Remaining))
		rw.Header().Set("RateLimit-Reset", fmt.Sprint(rateLimit.RateLimit_Reset))
		if rateLimit.RateLimit_Remaining > 0 {
			tmpl := template.Must(template.ParseGlob("./web/*.html"))
			var data = struct {
				ClientId    string
				RedirectURL string
				Scopes      []string
			}{
				config.Twitch.Api.ClientId,
				config.Twitch.Api.RedirectURL,
				config.Twitch.Api.AviableScopeList,
			}
			err := tmpl.Execute(rw, data)
			if err != nil {
				panic(err)
			}
			rateLimit.RateLimit_Remaining -= 1
		} else {
			rw.WriteHeader(http.StatusTooManyRequests)
			rw.Header().Set("Content-Type", "application/json")
			errMsg := fmt.Sprintf("{\"status\":\"Bad Request\",\"stauscode\":%d}", http.StatusTooManyRequests)
			rw.Write([]byte(errMsg))
		}

	})

	//http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":80", router))
}
