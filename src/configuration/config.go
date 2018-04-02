package configuration

import (
	"time"
)

// Config holds the general application settings.
type Config struct {
	// Debug
	Debug      bool
	Copyrights bool
	Quiet      bool

	// Config
	ConfigFile string

	// Collector Settings
	Concurrency       int
	Interval          time.Duration
	Timeout           time.Duration
	AllowInsecureSSH  bool
	DisableCollection bool
	Daemonize         bool

	// Git
	GitPush bool

	// User Settings
	Workspace   string
	DefaultUser string
	DefaultPass string

	// Mail
	EmailEnabled bool
	EmailTo      string
	EmailFrom    string
	EmailSubject string

	// SMTP Plain Auth
	SMTPHost string
	SMTPPort int
	SMTPUser string
	SMTPPass string

	// Webserver
	WebserverEnabled bool
	HTTPListen       string

	// Devices
	Devices []DeviceConfig
}

// DeviceConfig holds the device specific settings.
type DeviceConfig struct {
	Name           string
	Host           string
	Type           string
	User           string
	Pass           string
	Port           int
	Disabled       bool
	CustomTimeout  time.Duration
	CommandTimeout time.Duration
	FailChan       chan bool
	FailureWarning int
}

// GetName provides a simple implementation for the Collector interface.
func (d DeviceConfig) GetName() string {
	return d.Name
}

// getConfigDefaults provides reasonable defaults for Contra!
func getConfigDefaults() *Config {
	return &Config{
		ConfigFile:   "contra.conf",
		Concurrency:  30,
		Interval:     300 * time.Second,
		Timeout:      120 * time.Second,
		Workspace:    "/workspace",
		EmailSubject: "Changes from Contra!",
		SMTPHost:     "smtphost",
		SMTPPort:     25,
		HTTPListen:   "localhost:5002",
		Daemonize:    false,
	}
}
