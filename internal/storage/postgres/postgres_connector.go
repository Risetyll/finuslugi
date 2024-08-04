package postgres

import "os"

type PostgresConnector struct{}

func (PostgresConnector) GetProvider() string {
	return os.Getenv("DB_HOSTNAME")
}

func (PostgresConnector) GetConnect() string {
	return os.Getenv("DATABASE_URL")
}
