package cron

import (
	"time"

	"dcron/config"
	log "github.com/Sirupsen/logrus"
)

type SimpleCronService struct {
	CronService
	configService *config.CronConfigService
}

func NewSimpleCronService(configService *config.CronConfigService) *SimpleCronService {
	return &SimpleCronService{configService: configService}
}

func (s *SimpleCronService) Init() {
	s.load()
	ch := (*s.configService).Watch()
	for {
		<-ch
		log.Info("ConfigService emitted a config update event")
		s.reset()
		s.load()
	}
}

func (s *SimpleCronService) reset() {
}

func (s *SimpleCronService) load() {
	var conf *config.CronConfig
	for {
		var err error
		conf, err = (*s.configService).Load()
		if err == nil {
			log.WithField("conf", conf).Info("Successfully loaded cron config")
			break
		}
		log.WithField("err", err).Error("Failed to load cron config, retrying in 2 secs")
		time.Sleep(2 * time.Second)
	}

	jobs, _ := MakeJobsFromString(conf.Config)
	log.Info(jobs)
}
