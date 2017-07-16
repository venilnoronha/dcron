package config

// CronConfig represents cron configuration.
type CronConfig struct {
	Config  string
	Version int64
}

// CronConfigService abstracts logic to load and save cron configuration.
type CronConfigService interface {
	// Load loads the cron configuration.
	Load() (*CronConfig, error)

	// Save saves the new cron configuration.
	Save(*CronConfig) error
}
