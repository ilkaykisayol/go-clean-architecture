package sample

import (
	"fmt"
	"go-clean-architecture/internal/service/sample"
	"net/http"
	"strconv"

	"go-clean-architecture/internal/api"
	"go-clean-architecture/internal/util/cacher"
	"go-clean-architecture/internal/util/env"
	"go-clean-architecture/internal/util/logger"
	"go-clean-architecture/internal/util/validator"

	"github.com/gin-gonic/gin"
)

type ISampleController interface {
	RegisterRoutes(routerGroup *gin.RouterGroup)
	Get(context *gin.Context)
	Add(context *gin.Context)
	Update(context *gin.Context)
	GetProxy(context *gin.Context)
	GetDatabase(context *gin.Context)
	GetCache(context *gin.Context)
	PublishPubSubMessage(context *gin.Context)
	PostSampleXml(context *gin.Context)
}

type SampleController struct {
	path          string
	environment   env.IEnvironment
	loggr         logger.ILogger
	validatr      validator.IValidator
	cachr         cacher.ICacher
	sampleService sample.ISampleService
}

// NewSampleController
// Returns a new SampleController.
func NewSampleController(environment env.IEnvironment, loggr logger.ILogger, validatr validator.IValidator, cachr cacher.ICacher, sampleService sample.ISampleService) ISampleController {
	controller := SampleController{
		path:        "sample",
		environment: environment,
		loggr:       loggr,
		validatr:    validatr,
		cachr:       cachr,
	}

	if sampleService != nil {
		controller.sampleService = sampleService
	} else {
		controller.sampleService = sample.NewSampleService(environment, loggr, validatr, cachr, nil, nil, nil, nil)
	}
	return &controller
}

// RegisterRoutes
// Registers routes to gin.
func (c *SampleController) RegisterRoutes(routerGroup *gin.RouterGroup) {
	routes := routerGroup.Group(c.path)
	routes.Use(api.AuthenticationMiddleware(c.environment))
	routes.GET("", c.Get)
	routes.POST("", c.Add)
	routes.PUT(":id", c.Update)
	routes.GET("proxy", c.GetProxy)
	routes.GET("database", c.GetDatabase)
	routes.GET("cache", c.GetCache)
	routes.POST("pub-sub", c.PublishPubSubMessage)
	routes.POST("xml", c.PostSampleXml)
}

// Get
// @basePath     /api
// @router       /v1/sample [get]
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

	message := fmt.Sprintf("Sample v1 with PageSize:%v and PageIndex:%v", pageSize, pageIndex)

	context.JSON(http.StatusOK, api.Ok(message))
}

// Add
// @basePath     /api
// @router       /v1/sample [post]
// @tags         Sample
// @summary      Adds new sample.
// @description  Adds new sample.
// @security     Bearer
// @accept       json
// @produce      json
// @Param        X-Culture   header    string  true  "Request culture"   default(tr-TR)
// @Param        X-Timezone  header    string  true  "Request timezone"  default(Europe/Istanbul)
// @success      200         {object}  api.ApiResponse
// @failure      400         {object}  api.ApiResponse
// @failure      401         {object}  api.ApiResponse
// @failure      500         {object}  api.ApiResponse
// @Param        Model       body      AddSampleModel  true  "Request model"
func (c *SampleController) Add(context *gin.Context) {
	var model AddSampleModel
	err := context.ShouldBindJSON(&model)
	if err != nil {
		context.Error(err)
		return
	}

	context.JSON(http.StatusOK, api.Ok("Sample is created successfully."))
}

// Update
// @basePath     /api
// @router       /v1/sample/{id} [put]
// @tags         Sample
// @summary      Updates given sample.
// @description  Updates given sample.
// @security     Bearer
// @accept       json
// @produce      json
// @Param        X-Culture   header    string  true  "Request culture"   default(tr-TR)
// @Param        X-Timezone  header    string  true  "Request timezone"  default(Europe/Istanbul)
// @success      200         {object}  api.ApiResponse
// @failure      400         {object}  api.ApiResponse
// @failure      401         {object}  api.ApiResponse
// @failure      500         {object}  api.ApiResponse
// @Param        id          path      int                true  "Sample Id"
// @Param        Model       body      UpdateSampleModel  true  "Request model"
func (c *SampleController) Update(context *gin.Context) {
	id := context.Param("id")
	sampleId, err := strconv.Atoi(id)
	if err != nil {
		context.Error(err)
		return
	}

	var model UpdateSampleModel
	err = context.ShouldBindJSON(&model)
	if err != nil {
		context.Error(err)
		return
	}

	ch := make(chan *sample.UpdateSampleServiceResponse)
	defer close(ch)
	go c.sampleService.UpdateSample(ch, &sample.UpdateSampleServiceModel{
		SampleId:     sampleId,
		SampleStatus: model.SampleStatus,
		ModifiedBy:   model.ModifiedBy,
	})

	serviceResponse := <-ch
	if serviceResponse.Error != nil {
		context.Error(serviceResponse.Error)
		return
	}

	context.JSON(http.StatusOK, api.Ok(serviceResponse))
}

// GetProxy
// @basePath     /api
// @router       /v1/sample/proxy [get]
// @tags         Sample
// @summary      Gets a sample response via proxy.
// @description  Gets a sample response via proxy.
// @security     Bearer
// @accept       json
// @produce      json
// @Param        X-Culture   header    string  true  "Request culture"   default(tr-TR)
// @Param        X-Timezone  header    string  true  "Request timezone"  default(Europe/Istanbul)
// @success      200         {object}  api.ApiResponse
// @failure      400         {object}  api.ApiResponse
// @failure      401         {object}  api.ApiResponse
// @failure      500         {object}  api.ApiResponse
func (c *SampleController) GetProxy(context *gin.Context) {
	ch := make(chan *sample.GetSampleServiceResponse)
	defer close(ch)
	go c.sampleService.GetGoogle(ch, &sample.GetSampleServiceModel{
		Id:         8,
		SampleName: "Trying some proxy requests..",
	})

	serviceResponse := <-ch
	if serviceResponse.Error != nil {
		context.Error(serviceResponse.Error)
		return
	}

	context.JSON(http.StatusOK, api.Ok(serviceResponse))
}

// GetDatabase
// @basePath     /api
// @router       /v1/sample/database [get]
// @tags         Sample
// @summary      Gets a sample response from database.
// @description  Gets a sample response from database.
// @security     Bearer
// @accept       json
// @produce      json
// @Param        X-Culture   header    string  true  "Request culture"   default(tr-TR)
// @Param        X-Timezone  header    string  true  "Request timezone"  default(Europe/Istanbul)
// @success      200         {object}  api.ApiResponse
// @failure      400         {object}  api.ApiResponse
// @failure      401         {object}  api.ApiResponse
// @failure      500         {object}  api.ApiResponse
func (c *SampleController) GetDatabase(context *gin.Context) {
	ch := make(chan *sample.GetSampleServiceResponse)
	defer close(ch)
	go c.sampleService.GetDatabase(ch, &sample.GetSampleServiceModel{
		Id:         7,
		SampleName: "Trying some database requests..",
	})

	serviceResponse := <-ch
	if serviceResponse.Error != nil {
		context.Error(serviceResponse.Error)
		return
	}

	context.JSON(http.StatusOK, api.Ok(serviceResponse))
}

// GetCache
// @basePath     /api
// @router       /v1/sample/cache [get]
// @tags         Sample
// @summary      Gets a sample response from redis cache.
// @description  Gets a sample response from redis cache.
// @security     Bearer
// @accept       json
// @produce      json
// @Param        X-Culture   header    string  true  "Request culture"   default(tr-TR)
// @Param        X-Timezone  header    string  true  "Request timezone"  default(Europe/Istanbul)
// @success      200         {object}  api.ApiResponse
// @failure      400         {object}  api.ApiResponse
// @failure      401         {object}  api.ApiResponse
// @failure      500         {object}  api.ApiResponse
func (c *SampleController) GetCache(context *gin.Context) {
	ch := make(chan *sample.GetSampleServiceResponse)
	defer close(ch)
	go c.sampleService.GetCache(ch, &sample.GetSampleServiceModel{
		Id:         1,
		SampleName: "Trying some redis client requests..",
	})

	serviceResponse := <-ch
	if serviceResponse.Error != nil {
		context.Error(serviceResponse.Error)
		return
	}

	context.JSON(http.StatusOK, api.Ok(serviceResponse))
}

// PublishPubSubMessage
// @basePath     /api
// @router       /v1/sample/pub-sub [post]
// @tags         Sample
// @summary      Publishes a sample message to pub sub.
// @description  Publishes a sample message to pub sub.
// @security     Bearer
// @accept       json
// @produce      json
// @Param        X-Culture   header    string  true  "Request culture"   default(tr-TR)
// @Param        X-Timezone  header    string  true  "Request timezone"  default(Europe/Istanbul)
// @success      200         {object}  api.ApiResponse
// @failure      400         {object}  api.ApiResponse
// @failure      401         {object}  api.ApiResponse
// @failure      500         {object}  api.ApiResponse
// @Param        Model       body      PublishPubSubMessageModel  true  "Request model"
func (c *SampleController) PublishPubSubMessage(context *gin.Context) {
	var model PublishPubSubMessageModel
	modelErr := context.ShouldBindJSON(&model)
	if modelErr != nil {
		context.Error(modelErr)
		return
	}

	ch := make(chan *sample.PublishPubSubMessageServiceResponse)
	defer close(ch)
	go c.sampleService.PublishPubSubMessage(ch, &sample.PublishPubSubMessageServiceModel{
		Message: model.Message,
		Count:   model.Count,
	})

	serviceResponse := <-ch
	if serviceResponse.Error != nil {
		context.Error(serviceResponse.Error)
		return
	}

	context.JSON(http.StatusOK, api.Ok(serviceResponse))
}

// PostSampleXml
// @basePath     /api
// @router       /v1/sample/xml [post]
// @tags         Sample
// @summary      Post a sample xml via proxy.
// @description  Post a sample xml via proxy.
// @security     Bearer
// @accept       json
// @produce      json
// @Param        X-Culture   header    string  true  "Request culture"   default(tr-TR)
// @Param        X-Timezone  header    string  true  "Request timezone"  default(Europe/Istanbul)
// @success      200         {object}  api.ApiResponse
// @failure      400         {object}  api.ApiResponse
// @failure      401         {object}  api.ApiResponse
// @failure      500         {object}  api.ApiResponse
// @Param        Model       body      PostSampleXmlModel  true  "Request model"
func (c *SampleController) PostSampleXml(context *gin.Context) {
	var model PostSampleXmlModel
	modelErr := context.ShouldBindJSON(&model)
	if modelErr != nil {
		context.Error(modelErr)
		return
	}

	ch := make(chan *sample.PostSampleXmlServiceResponse)
	defer close(ch)
	go c.sampleService.PostSampleXml(ch, &sample.PostSampleXmlServiceModel{
		SampleName: model.SampleName,
		SampleType: model.SampleType,
		SampleCode: model.SampleCode,
	})

	serviceResponse := <-ch
	if serviceResponse.Error != nil {
		context.Error(serviceResponse.Error)
		return
	}

	context.JSON(http.StatusOK, api.Ok(serviceResponse))
}
