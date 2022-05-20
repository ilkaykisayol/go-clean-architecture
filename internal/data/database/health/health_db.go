package health

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	"go-clean-architecture/internal/util/env"
)

type IHealthDb interface {
	Ping() error
}

type HealthDb struct {
	environment      env.IEnvironment
	connectionString string
	driverName       string
	timeout          time.Duration
}

// Returns a new HealthDb.
func NewHealthDb(environment env.IEnvironment) IHealthDb {
	db := HealthDb{
		environment:      environment,
		driverName:       "postgres",
		connectionString: environment.Get(env.PostgresqlConnectionString),
		timeout:          time.Second * 5,
	}

	return &db
}

func (d *HealthDb) Ping() error {

	connection, err := sql.Open(d.driverName, d.connectionString)
	if err != nil {
		return err
	}
	defer connection.Close()

	pingErr := connection.Ping()
	if pingErr != nil {
		return pingErr
	}

	return nil
}
