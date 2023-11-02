package config

import "time"

var defaultConfig = map[string]interface{}{
	"auth.refresh_subject":                   RefreshTokenSubject,
	"auth.access_subject":                    AccessTokenSubject,
	"auth.refresh_expiration_time":           RefreshTokenExpirationDuration,
	"auth.access_expiration_time":            AccessTokenExpirationDuration,
	"application.gracefull_shutdown_timeout": time.Second * 5,
}
