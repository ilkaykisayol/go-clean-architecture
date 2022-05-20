package api

import (
	"errors"
	"go-clean-architecture/internal/util/customerror"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go-clean-architecture/internal/util"
	"go-clean-architecture/internal/util/env"
	"go-clean-architecture/internal/util/helper"
	"go-clean-architecture/internal/util/logger"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

// AuthenticationMiddleware
// Checks JWT token if it's valid or not.
func AuthenticationMiddleware(environment env.IEnvironment) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		tokenString := strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(authHeader, "Bearer", ""), "bearer", ""))

		token, err := jwt.ParseWithClaims(
			tokenString, &util.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(environment.Get(env.JwtSecret)), nil
			},
		)

		if err != nil || !token.Valid {
			c.Header("Content-Type", "application/json; charset=utf-8")
			c.AbortWithError(http.StatusUnauthorized, errors.New("JWT token is invalid"))
			return
		}

		claims, ok := token.Claims.(*util.CustomClaims)
		if !ok {
			c.Header("Content-Type", "application/json; charset=utf-8")
			c.AbortWithError(http.StatusUnauthorized, errors.New("JWT token is invalid"))
			return
		}

		c.Set(helper.UserId, claims.Subject)
		c.Set(helper.UserName, claims.Username)
		c.Set(helper.UserEmail, claims.Email)

		/* Later, if we implement permissions feature;
		We can get claims and check if token has required claims to access any given endpoint. */

		c.Next()
	}
}

// LoggingMiddleware
// Logs HTTP requests with a predefined structure.
func LoggingMiddleware(loggr logger.ILogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		host := c.Request.Host
		route := c.FullPath()
		remoteAddr := c.Request.RemoteAddr
		clientIp := c.ClientIP()
		protocol := c.Request.Proto
		method := c.Request.Method
		uri := c.Request.RequestURI
		queryString := c.Request.URL.RawQuery
		elapsedMilliseconds := time.Since(start).Milliseconds()
		statusCode := c.Writer.Status()
		hasError := len(c.Errors.Errors()) > 0

		var customerrors []*customerror.Error
		for _, err := range c.Errors {
			var customerr *customerror.Error
			if errors.As(err.Err, &customerr) {
				customerrors = append(customerrors, customerr)
			} else {
				customerrors = append(customerrors, &customerror.Error{
					Loglevel: customerror.LogLevelError,
					Err:      err.Err,
				})
			}
		}

		var errInfo []error
		var errWarn []error
		var errError []error

		for _, customerr := range customerrors {
			switch customerr.Loglevel {
			case customerror.LogLevelInfo:
				errInfo = append(errInfo, customerr.Err)
			case customerror.LogLevelWarn:
				errWarn = append(errWarn, customerr.Err)
			case customerror.LogLevelError:
				errError = append(errError, customerr.Err)
			default:
				errError = append(errError, customerr.Err)
			}
		}

		pingRoute := "/api/ping"
		if route != pingRoute {
			logMessage := protocol + " " + method + " " + uri + " responded " + strconv.Itoa(statusCode) + " in " + strconv.Itoa(int(elapsedMilliseconds)) + " ms"
			if hasError {
				if len(errError) > 0 {
					loggr.Error(
						logMessage,
						zap.String("host", host),
						zap.String("route", route),
						zap.String("protocol", protocol),
						zap.String("uri", uri),
						zap.String("method", method),
						zap.String("remoteAddr", remoteAddr),
						zap.String("clientIp", clientIp),
						zap.String("queryString", queryString),
						zap.Int("statusCode", statusCode),
						zap.Int64("elapsedMilliseconds", elapsedMilliseconds),
						zap.Errors("errors", errError),
					)
				}

				if len(errWarn) > 0 {
					loggr.Warn(
						logMessage,
						zap.String("host", host),
						zap.String("route", route),
						zap.String("protocol", protocol),
						zap.String("uri", uri),
						zap.String("method", method),
						zap.String("remoteAddr", remoteAddr),
						zap.String("clientIp", clientIp),
						zap.String("queryString", queryString),
						zap.Int("statusCode", statusCode),
						zap.Int64("elapsedMilliseconds", elapsedMilliseconds),
						zap.Errors("errors", errWarn),
					)
				}

				if len(errInfo) > 0 {
					loggr.Info(
						logMessage,
						zap.String("host", host),
						zap.String("route", route),
						zap.String("protocol", protocol),
						zap.String("uri", uri),
						zap.String("method", method),
						zap.String("remoteAddr", remoteAddr),
						zap.String("clientIp", clientIp),
						zap.String("queryString", queryString),
						zap.Int("statusCode", statusCode),
						zap.Int64("elapsedMilliseconds", elapsedMilliseconds),
						zap.Errors("errors", errInfo),
					)
				}
			} else {
				loggr.Info(
					logMessage,
					zap.String("host", host),
					zap.String("route", route),
					zap.String("protocol", protocol),
					zap.String("uri", uri),
					zap.String("method", method),
					zap.String("remoteAddr", remoteAddr),
					zap.String("clientIp", clientIp),
					zap.String("queryString", queryString),
					zap.Int("statusCode", statusCode),
					zap.Int64("elapsedMilliseconds", elapsedMilliseconds),
					zap.Errors("errors", errError),
				)
			}
		}

		if hasError {
			c.AbortWithStatusJSON(statusCode, Error(c.Errors.Last().Error()))
		}
	}
}
