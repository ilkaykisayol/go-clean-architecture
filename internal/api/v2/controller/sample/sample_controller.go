package sample

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go-clean-architecture/internal/api"
	"go-clean-architecture/internal/util/cacher"
	"go-clean-architecture/internal/util/env"
	"go-clean-architecture/internal/util/logger"
	"go-clean-architecture/internal/util/validator"
)

type ISampleController interface {
	RegisterRoutes(routerGroup *gin.RouterGroup)
	Get(context *gin.Context)
}

type SampleController struct {
	path        string
	environment env.IEnvironment
	loggr       logger.ILogger
	validatr    validator.IValidator
	cachr       cacher.ICacher
}

// NewSampleController
// Returns a new SampleController.
func NewSampleController(environment env.IEnvironment, loggr logger.ILogger, validatr validator.IValidator, cachr cacher.ICacher) ISampleController {
	controller := SampleController{
		path:        "sample",
		environment: environment,
		loggr:       loggr,
		validatr:    validatr,
		cachr:       cachr,
	}
	return &controller
}

// RegisterRoutes
// Registers routes to gin.
func (c *SampleController) RegisterRoutes(routerGroup *gin.RouterGroup) {
	routes := routerGroup.Group(c.path)
	routes.Use(api.AuthenticationMiddleware(c.environment))
	routes.GET("", c.Get)
}

// Get
// @basePath     /api
// @router       /v2/sample [get]
// @tags         Sample
// @summary      Gets a sample response.
// @description  Gets a sample response.
// @security     Bearer
// @accept       json
// @produce      json
// @Param        X-Culture   header    string  true  "Request culture"   default(tr-TR)
// @Param        X-Timezone  header    string  true  "Request timezone"  default(Europe/Istanbul)
// @success      200         {object}  api.ApiResponse
// @failure      400         {object}  api.ApiResponse
// @failure      401         {object}  api.ApiResponse
// @failure      500         {object}  api.ApiResponse
// @Param        pageSize    query     int  false  "Size of the page."
// @Param        pageIndex   query     int  false  "Index of the page."
func (c *SampleController) Get(context *gin.Context) {
	pageSize := context.Query("pageSize")
	pageIndex := context.Query("pageIndex")

	if pageSize == "" {
		pageSize = "Empty"
	}
	if pageIndex == "" {
		pageIndex = "Empty"
	}

	message := fmt.Sprintf("Sample v2 with PageSize:%v and PageIndex:%v", pageSize, pageIndex)

	context.JSON(http.StatusOK, api.Ok(message))
}
