package cron

import (
	"bufio"
	"errors"
	"strings"

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
	Expression string
	Command    string
}

func NewJobFromString(job string) (*CronJob, error) {
	tokens := strings.SplitN(job, " ", 6)
	cronExpr := strings.Join(tokens[0:len(tokens)-1], " ")
	_, err := cr.Parse(cronExpr)
	if err != nil {
		return nil, errors.New("Failed to parse cron entry " + job)
	}
	return &CronJob{Expression: cronExpr, Command: tokens[5]}, nil
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
