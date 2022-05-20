package main

import (
	"fmt"

	"go-clean-architecture/docs"
	"go-clean-architecture/internal/api"
	"go-clean-architecture/internal/api/v1/controller/auth"
	"go-clean-architecture/internal/api/v1/controller/health"
	"go-clean-architecture/internal/api/v1/controller/sample"
	sampleController "go-clean-architecture/internal/api/v2/controller/sample"
	"go-clean-architecture/internal/util/cacher"
	"go-clean-architecture/internal/util/env"
	"go-clean-architecture/internal/util/logger"
	"go-clean-architecture/internal/util/validator"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title                       Go Clean Architecture
// @version                     1.0
// @description                 This is an template project.
// @termsOfService              http://swagger.io/terms/
// @contact.name                API Support
// @contact.url                 http://www.swagger.io/support
// @contact.email               support@swagger.io
// @license.name                Apache 2.0
// @license.url                 http://www.apache.org/licenses/LICENSE-2.0.html
// @host                        localhost:8080
// @BasePath                    /api
// @accept                      json
// @produce                     json
// @schemes                     http https
// @securityDefinitions.apikey  Bearer
// @in                          header
// @name                        Authorization
func main() {
	environment := env.New()
	environment.Init()
	loggr := logger.New(environment)
	defer loggr.Sync()
	cachr := cacher.New(environment)
	validatr := validator.New()

	// router := gin.Default()
	router := gin.New()
	router.Use(api.LoggingMiddleware(loggr))
	addRoutes(router, environment, loggr, validatr, cachr)
	addSwagger(router, environment)

	addReceivers(environment, loggr, validatr, cachr)

	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	router.Run()
}

func addRoutes(
	router *gin.Engine,
	environment env.IEnvironment,
	loggr logger.ILogger,
	validatr validator.IValidator,
	cachr cacher.ICacher,
) {
	api := router.Group("api")
	health.NewHealthController(environment, loggr, validatr, cachr, nil).RegisterRoutes(api)

	v1 := api.Group("v1")
	v2 := api.Group("v2")
	auth.NewAuthController(environment, loggr, validatr, cachr, nil).RegisterRoutes(v1)
	sample.NewSampleController(environment, loggr, validatr, cachr, nil).RegisterRoutes(v1)
	sampleController.NewSampleController(environment, loggr, validatr, cachr).RegisterRoutes(v2)
}

func addSwagger(router *gin.Engine, environment env.IEnvironment) {
	docs.SwaggerInfo.Title = fmt.Sprintf("Go Clean Architecture (%v)", environment.Get(env.AppEnvironment))
	docs.SwaggerInfo.Host = environment.Get(env.AppHost)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func addReceivers(
	environment env.IEnvironment,
	loggr logger.ILogger,
	validatr validator.IValidator,
	cachr cacher.ICacher,
) {
	//go sampleReceiver.NewSampleReceiver(environment, loggr, validatr, cachr, nil).InitReceivers(1)
}
