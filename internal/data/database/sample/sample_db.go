package sample

import (
	"context"
	"database/sql"
	"time"

	"go-clean-architecture/internal/util/cacher"
	"go-clean-architecture/internal/util/env"
	"go-clean-architecture/internal/util/logger"
	"go-clean-architecture/internal/util/validator"

	_ "github.com/lib/pq"
)

type ISampleDb interface {
	GetSample(ch chan *GetSampleDbResponse, model *GetSampleDbModel)
}

type SampleDb struct {
	loggr            logger.ILogger
	validatr         validator.IValidator
	cachr            cacher.ICacher
	environment      env.IEnvironment
	connectionString string
	driverName       string
	timeout          time.Duration
}

// NewSampleDb
// Returns a new SampleDb.
func NewSampleDb(environment env.IEnvironment, loggr logger.ILogger, validatr validator.IValidator, cachr cacher.ICacher) ISampleDb {
	db := SampleDb{
		environment:      environment,
		loggr:            loggr,
		validatr:         validatr,
		cachr:            cachr,
		driverName:       "postgres",
		connectionString: environment.Get(env.PostgresqlConnectionString),
		timeout:          time.Second * 5,
	}

	return &db
}

// GetSample
// Gets sample response from postgresql db.
func (d *SampleDb) GetSample(ch chan *GetSampleDbResponse, model *GetSampleDbModel) {
	modelErr := d.validatr.ValidateStruct(model)
	if modelErr != nil {
		ch <- &GetSampleDbResponse{Error: modelErr}
		return
	}

	connection, err := sql.Open(d.driverName, d.connectionString)
	if err != nil {
		ch <- &GetSampleDbResponse{Error: err}
		return
	}
	defer connection.Close()

	ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
	defer cancel()
	err = connection.PingContext(ctx)
	if err != nil {
		ch <- &GetSampleDbResponse{Error: err}
		return
	}

	ch <- &GetSampleDbResponse{SampleName: "sample name here!!", Id: 1}
}
