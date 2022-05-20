package health

import (
	"go-clean-architecture/internal/service/health"
	"net/http"

	"go-clean-architecture/internal/api"
	"go-clean-architecture/internal/util/cacher"
	"go-clean-architecture/internal/util/env"
	"go-clean-architecture/internal/util/logger"
	"go-clean-architecture/internal/util/validator"

	"github.com/gin-gonic/gin"
)

type IHealthController interface {
	RegisterRoutes(routerGroup *gin.RouterGroup)
	Ping(context *gin.Context)
	Service(context *gin.Context)
}

type HealthController struct {
	environment   env.IEnvironment
	loggr         logger.ILogger
	validatr      validator.IValidator
	cachr         cacher.ICacher
	healthService health.IHealthService
}

// NewHealthController
// Returns a new HealthController.
func NewHealthController(environment env.IEnvironment,
	loggr logger.ILogger,
	validatr validator.IValidator,
	cachr cacher.ICacher,
	healthService health.IHealthService) IHealthController {

	controller := HealthController{
		environment:   environment,
		loggr:         loggr,
		validatr:      validatr,
		cachr:         cachr,
		healthService: healthService,
	}

	if healthService != nil {
		controller.healthService = healthService
	} else {
		controller.healthService = health.NewHealthService(environment, loggr, validatr, cachr, nil)
	}
	return &controller
}

func (c *HealthController) RegisterRoutes(routerGroup *gin.RouterGroup) {
	routerGroup.GET("ping", c.Ping)
	routerGroup.GET("service", c.Service)
}

// Ping
// @basePath     /api
// @router       /ping [get]
// @tags         Health
// @summary      Send a ping request.
// @description  Send a ping request.
// @accept       json
// @produce      json
// @success      200  {object}  api.ApiResponse
// @failure      400  {object}  api.ApiResponse
// @failure      500  {object}  api.ApiResponse
func (c *HealthController) Ping(context *gin.Context) {
	context.JSON(http.StatusOK, api.Ok("Ping OK"))
}

// Service
// @basePath     /api
// @router       /service [get]
// @tags         Health
// @summary      Send a service check request.
// @description  Send a service check request.
// @accept       json
// @produce      json
// @success      200  {object}  api.ApiResponse
// @failure      400  {object}  api.ApiResponse
// @failure      500  {object}  api.ApiResponse
func (c *HealthController) Service(context *gin.Context) {

	healthCheckCh := make(chan *health.HealthCheckServiceResponse)
	defer close(healthCheckCh)

	go c.healthService.HealthCheck(healthCheckCh)

	res := <-healthCheckCh

	context.JSON(http.StatusOK, api.Ok(res.HealthMessage))
}
