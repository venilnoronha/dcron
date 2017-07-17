package cron

import (
	"bufio"
	"errors"
	"strings"
	"time"

	cr "gopkg.in/robfig/cron.v2"
)

// CronService encapsulates logic for monitoring and executing cron jobs.
type CronService interface {
	// Init initializes the cron service.
	Init()

	// Destroy destroys the cron service.
	Destroy()
}

// CronJob represents a simple cron job.
type CronJob struct {
	Schedule cr.Schedule
	Command  string
}

func (j *CronJob) NextAt() time.Time {
	return j.Schedule.Next(time.Now())
}

func NewJobFromString(job string) (*CronJob, error) {
	tokens := strings.SplitN(job, " ", 6)
	cronExpr := strings.Join(tokens[0:len(tokens)-1], " ")
	schedule, err := cr.Parse(cronExpr)
	if err != nil {
		return nil, errors.New("Failed to parse cron entry " + job)
	}
	return &CronJob{Schedule: schedule, Command: tokens[5]}, nil
}

func MakeJobsFromString(jobs string) (*[]*CronJob, error) {
	var cronJobs []*CronJob
	scanner := bufio.NewScanner(strings.NewReader(jobs))
	for scanner.Scan() {
		line := scanner.Text()
		job, err := NewJobFromString(line)
		if err != nil {
			return nil, err
		}
		cronJobs = append(cronJobs, job)
	}
	return &cronJobs, nil
}
