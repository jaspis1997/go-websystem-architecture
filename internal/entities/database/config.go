package database

import (
	"fmt"
	"net/url"
)

type Config interface {
	DSN() string
}

type PostgresConfig struct {
	User     string
	Password string
	Host     string
	Port     uint16
	DB       string
	SSLMode  bool
}

func (config PostgresConfig) DSN() string {
	const postgres = "postgres"
	const defaultPort = 5432
	userinfo := func() *url.Userinfo {
		user := func() string {
			if config.User == "" {
				return postgres
			}
			return config.User
		}
		if config.Password == "" {
			return url.User(user())
		}
		return url.UserPassword(user(), config.Password)
	}
	host := func() string {
		if config.Host == "" {
			return postgres
		}
		return config.Host
	}
	db := func() string {
		if config.DB == "" {
			return postgres
		}
		return config.DB
	}
	query := func() string {
		query := make(url.Values)
		if !config.SSLMode {
			query.Add("sslmode", "disable")
		} else {
			query.Add("sslmode", "require")
		}
		return query.Encode()
	}
	port := config.Port
	if port == 0 {
		port = defaultPort
	}
	url := url.URL{
		Scheme:   postgres,
		Host:     fmt.Sprintf("%s:%d", host(), port),
		User:     userinfo(),
		Path:     db(),
		RawQuery: query(),
	}
	return url.String()
}

type SQLiteConfig struct {
	Filename string
}

func (config SQLiteConfig) DSN() string {
	return config.Filename
}
