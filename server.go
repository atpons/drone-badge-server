package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/BurntSushi/toml"

	"github.com/atpons/drone-badge-server/badge"

	db "github.com/atpons/drone-badge-server/database"
	"github.com/atpons/drone-badge-server/drone"
	"github.com/gorilla/mux"
)

type Request struct {
	BuildNum      int    `json:"build_number"`
	RepoNamespace string `json:"repo_namespace"`
	RepoName      string `json:"repo_name"`
}

var (
	RequestError     = fmt.Errorf("request parse Error")
	DroneClientError = fmt.Errorf("drone client error")
)

type Config struct {
	Drone DroneConfig
}

type DroneConfig struct {
	Token string
	Host  string
}

func main() {
	var config Config
	_, err := toml.DecodeFile("config.toml", &config)
	if err != nil {
		panic(err)
	}

	c, err := drone.NewDrone(config.Drone.Host, config.Drone.Token)
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/generate", GenerateHandler(c))
	r.HandleFunc("/badge/{repo_id:[0-9]+}/{stage_number:[0-9]+}", ReturnBadge(c))

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func GenerateHandler(client *drone.Drone) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			RenderError(w, r, RequestError)
			return
		}
		if err := json.Unmarshal(body, &req); err != nil {
			RenderError(w, r, RequestError)
			return
		}

		err = client.SetStage(req.RepoNamespace, req.RepoName, req.BuildNum)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			RenderError(w, r, DroneClientError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		fmt.Println(req, err)
	}
}

func ReturnBadge(client *drone.Drone) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		repoId, err := strconv.Atoi(vars["repo_id"])
		if err != nil {
			RenderError(w, r, RequestError)
			return
		}
		stageNum, err := strconv.Atoi(vars["stage_number"])
		if err != nil {
			RenderError(w, r, RequestError)
			return
		}
		status, err := db.GetValue(client.Db, repoId, stageNum)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			RenderError(w, r, DroneClientError)
			return
		}
		badge.ReturnImage(w, r, BadgeRoute, repoId, stageNum, status)
	}
}

func RenderError(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, err.Error())
}
