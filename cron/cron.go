package cron

// CronService encapsulates logic for monitoring and executing cron jobs.
type CronService interface {
	// Init initializes the cron service.
	Init()

	// Destroy destroys the cron service.
	Destroy()
}
