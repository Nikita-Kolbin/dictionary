package config

import (
	"fmt"
	"net/url"
)

func (c *ListenerConfig) GetHostPort() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *PostgresConfig) GetDSN() string {
	dsn := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(c.Username, c.Password),
		Host:     c.HostPort,
		Path:     c.DBName,
		RawQuery: getPostgresConnParams(),
	}
	return dsn.String()
}

func getPostgresConnParams() string {
	params := url.Values{}
	params.Add("sslmode", "disable")
	return params.Encode()
}
