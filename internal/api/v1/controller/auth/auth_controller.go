package auth

import (
	"fmt"
	"go-clean-architecture/internal/service/auth"
	"net/http"

	"go-clean-architecture/internal/api"
	"go-clean-architecture/internal/util/cacher"
	"go-clean-architecture/internal/util/env"
	"go-clean-architecture/internal/util/helper"
	"go-clean-architecture/internal/util/logger"
	"go-clean-architecture/internal/util/validator"

	"github.com/gin-gonic/gin"
)

type IAuthController interface {
	RegisterRoutes(routerGroup *gin.RouterGroup)
	Login(context *gin.Context)
	GetAccessToken(context *gin.Context)
	GetProgrammaticAccessToken(context *gin.Context)
}

type AuthController struct {
	path        string
	environment env.IEnvironment
	loggr       logger.ILogger
	validatr    validator.IValidator
	cachr       cacher.ICacher
	authService auth.IAuthService
}

// NewAuthController
// Returns a new AuthController.
func NewAuthController(environment env.IEnvironment, loggr logger.ILogger, validatr validator.IValidator, cachr cacher.ICacher, authService auth.IAuthService) IAuthController {
	controller := AuthController{
		path:        "auth",
		environment: environment,
		loggr:       loggr,
		validatr:    validatr,
		cachr:       cachr,
	}

	if authService != nil {
		controller.authService = authService
	} else {
		controller.authService = auth.NewAuthService(environment, loggr, validatr, cachr, nil)
	}

	return &controller
}

// RegisterRoutes
// Registers routes to gin.
func (c *AuthController) RegisterRoutes(routerGroup *gin.RouterGroup) {
	routes := routerGroup.Group(c.path)
	routes.POST("login", c.Login)
	routes.POST("access-token", c.GetAccessToken)
	routes.POST("access-token/programmatic", c.GetProgrammaticAccessToken)
	routes.Use(api.AuthenticationMiddleware(c.environment))
	routes.GET("", c.Get)
}

// Login
// @basePath     /api
// @router       /v1/auth/login [post]
// @tags         Auth
// @summary      Log in and get JWT token via username and password.
// @description  Generates JWT token if user can log in successfully.
// @security     Bearer
// @accept       json
// @produce      json
// @Param        X-Culture   header    string  true  "Request culture"   default(tr-TR)
// @Param        X-Timezone  header    string  true  "Request timezone"  default(Europe/Istanbul)
// @success      200         {object}  api.ApiResponse
// @failure      400         {object}  api.ApiResponse
// @failure      500         {object}  api.ApiResponse
// @Param        model       body      LoginModel  true  "Username and password"
func (c *AuthController) Login(context *gin.Context) {
	var model LoginModel
	err := context.ShouldBindJSON(&model)
	if err != nil {
		context.Error(err)
		return
	}

	ch := make(chan *auth.LoginServiceResponse)
	defer close(ch)

	go c.authService.Login(ch, &auth.LoginServiceModel{
		UserName: model.UserName,
		Password: model.Password,
	})

	serviceResponse := <-ch

	if serviceResponse.Error != nil {
		context.Error(serviceResponse.Error)
		return
	}

	context.JSON(http.StatusOK, api.Ok(serviceResponse))
}

// GetAccessToken
// @basePath     /api
// @router       /v1/auth/access-token [post]
// @tags         Auth
// @summary      Gets access token for users.
// @description  Generates JWT token if user can log in successfully.
// @security     Bearer
// @accept       json
// @produce      json
// @Param        X-Culture   header    string  true  "Request culture"   default(tr-TR)
// @Param        X-Timezone  header    string  true  "Request timezone"  default(Europe/Istanbul)
// @success      200         {object}  api.ApiResponse
// @failure      400         {object}  api.ApiResponse
// @failure      500         {object}  api.ApiResponse
// @Param        model       body      GetAccessTokenModel  true  "RefreshToken"
func (c *AuthController) GetAccessToken(context *gin.Context) {
	var model GetAccessTokenModel
	err := context.ShouldBindJSON(&model)
	if err != nil {
		context.Error(err)
		return
	}

	ch := make(chan *auth.LoginServiceResponse)
	defer close(ch)

	go c.authService.GetAccessToken(ch, &auth.GetAccessTokenServiceModel{
		RefreshToken: model.RefreshToken,
	})

	serviceResponse := <-ch

	if serviceResponse.Error != nil {
		context.Error(serviceResponse.Error)
		return
	}

	context.JSON(http.StatusOK, api.Ok(serviceResponse))
}

// GetProgrammaticAccessToken
// @basePath     /api
// @router       /v1/auth/access-token/programmatic [post]
// @tags         Auth
// @summary      Gets access token with expiry days for programmatic users.
// @description  Generates JWT token if user can log in successfully.
// @security     Bearer
// @accept       json
// @produce      json
// @Param        X-Culture   header    string  true  "Request culture"   default(tr-TR)
// @Param        X-Timezone  header    string  true  "Request timezone"  default(Europe/Istanbul)
// @success      200         {object}  api.ApiResponse
// @failure      400         {object}  api.ApiResponse
// @failure      500         {object}  api.ApiResponse
// @Param        model       body      GetProgrammaticAccessTokenModel  true  "Username, password and expiryDays"
func (c *AuthController) GetProgrammaticAccessToken(context *gin.Context) {
	var model GetProgrammaticAccessTokenModel
	err := context.ShouldBindJSON(&model)
	if err != nil {
		context.Error(err)
		return
	}

	ch := make(chan *auth.GetProgrammaticAccessTokenServiceResponse)
	defer close(ch)

	go c.authService.GetProgrammaticAccessToken(ch, &auth.GetProgrammaticAccessTokenServiceModel{
		UserName:   model.UserName,
		Password:   model.Password,
		ExpiryDays: model.ExpiryDays,
	})

	serviceResponse := <-ch

	if serviceResponse.Error != nil {
		context.Error(serviceResponse.Error)
		return
	}

	if err != nil {
		context.Error(err)
		return
	}

	context.JSON(http.StatusOK, api.Ok(serviceResponse))
}

// Get
// @basePath     /api
// @router       /v1/auth [get]
// @tags         Auth
// @summary      Gets user info.
// @description  Gets user info.
// @security     Bearer
// @accept       json
// @produce      json
// @Param        X-Culture   header    string  true  "Request culture"   default(tr-TR)
// @Param        X-Timezone  header    string  true  "Request timezone"  default(Europe/Istanbul)
// @success      200         {object}  api.ApiResponse
// @failure      400         {object}  api.ApiResponse
// @failure      401         {object}  api.ApiResponse
// @failure      500         {object}  api.ApiResponse
func (c *AuthController) Get(context *gin.Context) {
	id := helper.GetUserId(context)
	username := helper.GetUserName(context)
	email := helper.GetUserEmail(context)

	message := fmt.Sprintf("id:%d user:%s email:%s", id, username, email)

	context.JSON(http.StatusOK, api.Ok(message))
}
