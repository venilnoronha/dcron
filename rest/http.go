package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"dcron/config"
	log "github.com/Sirupsen/logrus"
)

type RESTService struct {
	REST
	port        int
	server      *http.Server
	cronService *config.CronConfigService
}

func NewRESTService(port int, cronService *config.CronConfigService) *RESTService {
	server := &http.Server{Addr: fmt.Sprintf(":%d", port)}
	restService := RESTService{port: port, server: server, cronService: cronService}
	http.HandleFunc("/list", restService.list)
	return &restService
}

func (s *RESTService) Init() error {
	log.Info("Starting HTTP server at port ", s.port)
	if err := s.server.ListenAndServe(); err != nil {
		if err.Error() != "http: Server closed" {
			return err
		}
	}
	return nil
}

func (s *RESTService) Destroy() error {
	return s.server.Shutdown(nil)
}

func (s *RESTService) list(w http.ResponseWriter, r *http.Request) {
	conf, err := (*s.cronService).Load()
	if err != nil {
		log.WithField("err", err).Error("Failed to load cron config")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	log.WithField("conf", conf).Info("Loaded cron config")
	json.NewEncoder(w).Encode(conf)
}
