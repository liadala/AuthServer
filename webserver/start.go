package webserver

import (
	"TwitchTokGen/config"
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/nicklaw5/helix/v2"
)

//go:embed www/html/*
var webfiles embed.FS

//go:embed www/static/*
var fsstatic embed.FS

var (
	router *mux.Router = mux.NewRouter()
	twitch *helix.Client
)

func Start(options *helix.Options) {
	var err error
	twitch, err = helix.NewClient(options)
	if err != nil {
		log.Fatal(err)
	}

	// Prepare Filesystem
	static, err := fs.Sub(fsstatic, "www/static")
	if err != nil {
		log.Fatal(err)
	}

	templ, err := template.ParseFS(webfiles, "www/html/*")
	if err != nil {
		log.Fatal(err)
	}

	// Routing
	router.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.FS(static))))

	router.HandleFunc("/api/create", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		keys, ok := r.URL.Query()["code"]
		if !ok || len(keys[0]) < 1 {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("{\"status\":\"Bad Request\",\"stauscode\":400}"))
			return
		}
		response, err := twitch.RequestUserAccessToken(keys[0])
		if err != nil {
			log.Println(err)
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("{\"status\":\"Bad Request\",\"stauscode\":400}"))
			return
		}
		json.NewEncoder(rw).Encode(response.Data)
	})

	router.HandleFunc("/api/refresh/{key}", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		var vars = mux.Vars(r)
		if vars["key"] == "" {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("{\"status\":\"Bad Request\",\"stauscode\":400}"))
			return
		}

		response, err := twitch.RefreshUserAccessToken(vars["key"])
		if err != nil {
			log.Println(err)
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("{\"status\":\"Bad Request\",\"stauscode\":400}"))
			return
		}
		json.NewEncoder(rw).Encode(response.Data)
	})

	router.HandleFunc("/api/revoke/{key}", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		var vars = mux.Vars(r)
		if vars["key"] == "" {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("{\"status\":\"Bad Request\",\"stauscode\":400}"))
			return
		}

		response, err := twitch.RevokeUserAccessToken(vars["key"])
		if err != nil {
			log.Println(err)
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("{\"status\":\"Bad Request\",\"stauscode\":400}"))
			return
		}
		json.NewEncoder(rw).Encode(struct {
			Status  int    `json:"status"`
			Message string `json:"message"`
		}{
			Status:  response.ResponseCommon.StatusCode,
			Message: response.ErrorMessage,
		})
	})

	router.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		response, err := parseTemplate(templ, "index.html", struct {
			ClientId    string
			RedirectURL string
			Scopes      []string
		}{
			config.Config.Twitch.Api.ClientId,
			config.Config.Twitch.Api.RedirectURL,
			config.Config.Twitch.Api.AviableScopeList,
		})
		if err != nil {
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
		fmt.Fprint(rw, response)
	})

	srv := &http.Server{
		Addr:         ":80",
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func parseTemplate(templ *template.Template, name string, data interface{}) (string, error) {
	buf := new(bytes.Buffer)
	err := templ.ExecuteTemplate(buf, name, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
