package cron

import (
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
}
