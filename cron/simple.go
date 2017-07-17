package cron

import (
	"os/exec"
	"time"

	"dcron/config"
	log "github.com/Sirupsen/logrus"
	cr "gopkg.in/robfig/cron.v2"
)

type SimpleCronService struct {
	CronService
	configService *config.CronConfigService
	cron          *cr.Cron
}

func NewSimpleCronService(configService *config.CronConfigService) *SimpleCronService {
	return &SimpleCronService{configService: configService, cron: nil}
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
	log.Info("Stopping cron")
	s.cron.Stop()
	log.Info("Cron stopped")
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

	log.Info("Setting up cron")
	s.cron = cr.New()
	jobs, _ := MakeJobsFromString(conf.Config)
	for _, job := range *jobs {
		s.cron.AddFunc(job.Expression, func() {
			log.WithField("job", *job).Info("Executing job")
			out, err := exec.Command("sh", "-c", job.Command).CombinedOutput()
			if err != nil {
				log.WithFields(log.Fields{"job": *job, "err": err, "out": string(out)}).Error("Failed to execute job")
				return
			}
			log.WithFields(log.Fields{"job": *job, "out": string(out)}).Info("Finished executing job")
		})
		log.WithField("job", *job).Info("Job was set up")
	}
	s.cron.Start()
	log.Info("Cron setup complete")
}
