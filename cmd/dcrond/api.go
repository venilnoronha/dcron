package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
)

// server is the HTTP server.
var server *http.Server

func initHTTPServer(port int) {
	server = &http.Server{Addr: fmt.Sprintf(":%d", port)}
	http.HandleFunc("/list", list)

	log.Info("Starting HTTP server on port ", port)
	if err := server.ListenAndServe(); err != nil {
		if err.Error() != "http: Server closed" {
			log.WithField("err", err).Error("Failed to start HTTP server")
			os.Exit(1)
		}
	}
}

func destroyHTTPServer() {
	log.Info("Shutting down HTTP server")
	if err := server.Shutdown(nil); err != nil {
		log.WithField("err", err).Error("Failed to shutdown HTTP server")
	}
	log.Info("Completed HTTP server shutdown")
}

func list(w http.ResponseWriter, r *http.Request) {
	conf, err := cronConfigService.Load()
	if err != nil {
		log.WithField("err", err).Error("Failed to load cron config")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	log.WithField("conf", conf).Info("Loaded cron config")
	json.NewEncoder(w).Encode(conf)
}
