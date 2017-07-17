package cron

import (
	"bufio"
	"errors"
	"strconv"
	"strings"
	"time"
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
	Minute  string
	Hour    string
	Day     string
	Month   string
	Week    string
	Command string
}

func (j *CronJob) NextAt() time.Time {
	// TODO
	return time.Now()
}

func NewJobFromString(job string) (*CronJob, error) {
	ts := strings.SplitN(job, " ", 6)
	if !isValidToken(ts[0]) || !isValidToken(ts[1]) || !isValidToken(ts[2]) || !isValidToken(ts[3]) || !isValidToken(ts[4]) {
		return nil, errors.New("Failed to parse cron job string for line " + job)
	}
	// TODO: add more validation
	return &CronJob{Minute: ts[0], Hour: ts[1], Day: ts[2], Month: ts[3], Week: ts[4], Command: ts[5]}, nil
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

func isValidToken(token string) bool {
	_, err := strconv.Atoi(token)
	return token == "*" || err == nil
}
