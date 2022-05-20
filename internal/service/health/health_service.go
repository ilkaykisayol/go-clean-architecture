package health

import (
	"fmt"
	"strings"

	"go-clean-architecture/internal/data/database/health"
	"go-clean-architecture/internal/util/cacher"
	"go-clean-architecture/internal/util/env"
	"go-clean-architecture/internal/util/logger"
	"go-clean-architecture/internal/util/validator"
)

type IHealthService interface {
	HealthCheck(ch chan *HealthCheckServiceResponse)
}

type HealthService struct {
	environment env.IEnvironment
	loggr       logger.ILogger
	validatr    validator.IValidator
	cachr       cacher.ICacher
	healthDb    health.IHealthDb
}

// NewHealthService returns new HealthService
func NewHealthService(
	environment env.IEnvironment,
	loggr logger.ILogger,
	validatr validator.IValidator,
	cachr cacher.ICacher,
	healthDb health.IHealthDb,
) IHealthService {

	service := HealthService{
		environment: environment,
		loggr:       loggr,
		validatr:    validatr,
		cachr:       cachr,
		healthDb:    healthDb,
	}

	if healthDb != nil {
		service.healthDb = healthDb
	} else {
		service.healthDb = health.NewHealthDb(environment)
	}

	return &service
}

// UpsertSupplier calls the proxy to add or update supplier
func (s *HealthService) HealthCheck(ch chan *HealthCheckServiceResponse) {

	var healthResponses []string

	hasError := false
	err := s.cachr.Ping()

	if err != nil {
		healthResponses = append(healthResponses, fmt.Sprintf("Redis connection is unhealthy. Error: %s.", err))
		hasError = true
	} else {
		healthResponses = append(healthResponses, "Redis connection is healthy.")
	}

	err = s.healthDb.Ping()

	if err != nil {
		healthResponses = append(healthResponses, fmt.Sprintf("Db connection is unhealthy. Error: %s.", err))
		hasError = true
	} else {
		healthResponses = append(healthResponses, "Db connection is healthy.")
	}

	healthMessage := strings.Join(healthResponses, " ")
	if hasError {
		s.loggr.Error(healthMessage)
	} else {
		s.loggr.Info(healthMessage)
	}

	ch <- &HealthCheckServiceResponse{healthMessage}
}
