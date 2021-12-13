package config

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	configLoaded   bool
	dialTimeout    = 5 * time.Second
	contextTimeout = 5 * time.Second
	reloadDelay    = time.Second * 10
)

// Cnf ...
// Let's start with some sensible defaults
var Cnf = &Config{
	Database: DatabaseConfig{
		Type:         "postgres",
		Host:         "postgres",
		Port:         5432,
		User:         "go_oauth2_server",
		Password:     "",
		DatabaseName: "go_oauth2_server",
		MaxIdleConns: 5,
		MaxOpenConns: 5,
	},
	Oauth: OauthConfig{
		AccessTokenLifetime:  3600,    // 1 hour
		RefreshTokenLifetime: 1209600, // 14 days
		AuthCodeLifetime:     3600,    // 1 hour
	},
	IsDevelopment: true,
}

// NewConfig loads configuration from etcd and returns *Config struct
// It also starts a goroutine in the background to keep config up-to-date
func NewConfig(mustLoadOnce bool, keepReloading bool, backendType string) *Config {
	if configLoaded {
		return Cnf
	}

	var backend Backend

	switch backendType {
	case "etcd":
		backend = new(etcdBackend)
	case "consul":
		backend = new(consulBackend)
	default:
		log.Fatal("%s is not a valid backend", backendType)
		os.Exit(1)
	}

	backend.InitConfigBackend()

	// If the config must be loaded once successfully
	if mustLoadOnce && !configLoaded {
		// Read from remote config the first time
		newCnf, err := backend.LoadConfig()

		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		// Refresh the config
		backend.RefreshConfig(newCnf)

		// Set configLoaded to true
		configLoaded = true
		log.Info("Successfully loaded config for the first time")
	}

	if keepReloading {
		// Open a goroutine to watch remote changes forever
		go func() {
			for {
				// Delay after each request
				<-time.After(reloadDelay)

				// Attempt to reload the config
				newCnf, err := backend.LoadConfig()
				if err != nil {
					log.Error(err)
					continue
				}

				// Refresh the config
				backend.RefreshConfig(newCnf)

				// Set configLoaded to true
				configLoaded = true
				//log.Info("Successfully reloaded config")
			}
		}()
	}

	return Cnf
}
