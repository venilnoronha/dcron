package main

import (
	"fmt"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
)

// server is the HTTP server.
var server *http.Server

func initHTTPServer(port int) {
	http.HandleFunc("/jobs/list", listJobs)

	log.Info("Starting HTTP server on port ", port)
	server = &http.Server{Addr: fmt.Sprintf(":%d", port)}
	if err := server.ListenAndServe(); err != nil {
		log.WithField("err", err).Error("Failed to start HTTP server")
		os.Exit(1)
	}
}

func destroyHTTPServer() {
	log.Info("Shutting down HTTP server")
	if err := server.Shutdown(nil); err != nil {
		log.WithField("err", err).Error("Failed to shutdown HTTP server")
	}
	log.Info("Completed HTTP server shutdown")
}

func listJobs(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "job listing")
}
