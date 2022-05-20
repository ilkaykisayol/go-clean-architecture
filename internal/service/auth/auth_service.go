package auth

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"go-clean-architecture/internal/data/database/auth"
	"go-clean-architecture/internal/util/customerror"
	"strconv"
	"time"

	"go-clean-architecture/internal/util"
	"go-clean-architecture/internal/util/cacher"
	"go-clean-architecture/internal/util/env"
	"go-clean-architecture/internal/util/logger"
	"go-clean-architecture/internal/util/validator"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type IAuthService interface {
	Login(ch chan *LoginServiceResponse, model *LoginServiceModel)
	GetAccessToken(ch chan *LoginServiceResponse, model *GetAccessTokenServiceModel)
	GetProgrammaticAccessToken(
		ch chan *GetProgrammaticAccessTokenServiceResponse,
		model *GetProgrammaticAccessTokenServiceModel,
	)
}

type AuthService struct {
	environment env.IEnvironment
	loggr       logger.ILogger
	validatr    validator.IValidator
	cachr       cacher.ICacher
	authDb      auth.IAuthDb
}

// NewAuthService
// Returns a new AuthService.
func NewAuthService(
	environment env.IEnvironment,
	loggr logger.ILogger,
	validatr validator.IValidator,
	cachr cacher.ICacher,
	authDb auth.IAuthDb,
) IAuthService {
	service := AuthService{
		environment: environment,
		loggr:       loggr,
		validatr:    validatr,
		cachr:       cachr,
	}

	if authDb != nil {
		service.authDb = authDb
	} else {
		service.authDb = auth.NewAuthDb(environment, loggr, validatr, cachr)
	}

	return &service
}

// Login
// Logs in user with given username and password.
// Returns an error if user is not found.
func (s *AuthService) Login(ch chan *LoginServiceResponse, model *LoginServiceModel) {
	modelErr := s.validatr.ValidateStruct(model)
	if modelErr != nil {
		ch <- &LoginServiceResponse{Error: customerror.New(modelErr, customerror.LogLevelInfo)}
		return
	}

	passwordHashBytes := sha256.Sum256([]byte(model.Password))
	passwordHash := fmt.Sprintf("%x", passwordHashBytes[:]) // TODO: check sprintf usage

	chGetUserByPasswordResponse := make(chan *auth.GetUserByPasswordResponse)
	defer close(chGetUserByPasswordResponse)

	go s.authDb.GetUserByPassword(
		chGetUserByPasswordResponse, &auth.GetUserByPasswordModel{
			UserName:     model.UserName,
			PasswordHash: passwordHash,
		},
	)

	getUserByPasswordResponse := <-chGetUserByPasswordResponse
	if getUserByPasswordResponse.Error != nil {
		ch <- &LoginServiceResponse{Error: getUserByPasswordResponse.Error}
		return
	}

	if !getUserByPasswordResponse.IsActive || getUserByPasswordResponse.IsProgrammatic {
		ch <- &LoginServiceResponse{Error: errors.New("user is not found")}
		return
	}

	chAddRefreshTokenResponse := make(chan *auth.AddRefreshTokenResponse)
	defer close(chAddRefreshTokenResponse)

	refreshToken := uuid.NewString()
	go s.authDb.AddRefreshToken(
		chAddRefreshTokenResponse, &auth.AddRefreshTokenModel{
			UserId:       int(getUserByPasswordResponse.Id),
			RefreshToken: refreshToken,
		},
	)

	addRefreshTokenResponse := <-chAddRefreshTokenResponse
	if addRefreshTokenResponse.Error != nil {
		ch <- &LoginServiceResponse{Error: addRefreshTokenResponse.Error}
		return
	}

	tokenString, err := s.generateJwt(getUserByPasswordResponse.Id, getUserByPasswordResponse.UserName, getUserByPasswordResponse.Email, 0)

	if err != nil {
		ch <- &LoginServiceResponse{Error: err}
		return
	}

	ch <- &LoginServiceResponse{
		JwtToken:     tokenString,
		RefreshToken: refreshToken,
	}
}

// GetAccessToken
// Logs in user with given refresh token
// Returns an error if user is not found.
func (s *AuthService) GetAccessToken(ch chan *LoginServiceResponse, model *GetAccessTokenServiceModel) {
	modelErr := s.validatr.ValidateStruct(model)
	if modelErr != nil {
		ch <- &LoginServiceResponse{Error: modelErr}
		return
	}

	chGetUserByRefreshTokenResponse := make(chan *auth.GetUserByRefreshTokenResponse)
	defer close(chGetUserByRefreshTokenResponse)

	go s.authDb.GetUserByRefreshToken(
		chGetUserByRefreshTokenResponse, &auth.GetUserByRefreshTokenModel{
			RefreshToken: model.RefreshToken,
		},
	)

	getUserByRefreshTokenResponse := <-chGetUserByRefreshTokenResponse
	if getUserByRefreshTokenResponse.Error != nil {
		ch <- &LoginServiceResponse{Error: getUserByRefreshTokenResponse.Error}
		return
	}

	if !getUserByRefreshTokenResponse.IsActive || getUserByRefreshTokenResponse.IsProgrammatic {
		ch <- &LoginServiceResponse{Error: errors.New("user is not found")}
		return
	}

	chAddRefreshTokenResponse := make(chan *auth.AddRefreshTokenResponse)
	defer close(chAddRefreshTokenResponse)

	refreshToken := uuid.NewString()
	go s.authDb.AddRefreshToken(
		chAddRefreshTokenResponse, &auth.AddRefreshTokenModel{
			UserId:       int(getUserByRefreshTokenResponse.Id),
			RefreshToken: refreshToken,
		},
	)

	addRefreshTokenResponse := <-chAddRefreshTokenResponse
	if addRefreshTokenResponse.Error != nil {
		ch <- &LoginServiceResponse{Error: addRefreshTokenResponse.Error}
		return
	}

	tokenString, err := s.generateJwt(getUserByRefreshTokenResponse.Id, getUserByRefreshTokenResponse.UserName, getUserByRefreshTokenResponse.Email, 0)

	if err != nil {
		ch <- &LoginServiceResponse{Error: err}
		return
	}

	ch <- &LoginServiceResponse{
		JwtToken:     tokenString,
		RefreshToken: refreshToken,
	}
}

// GetProgrammaticAccessToken
// Logs in user with given username and password.
// Returns an error if user is not found.
func (s *AuthService) GetProgrammaticAccessToken(
	ch chan *GetProgrammaticAccessTokenServiceResponse,
	model *GetProgrammaticAccessTokenServiceModel,
) {
	modelErr := s.validatr.ValidateStruct(model)
	if modelErr != nil {
		ch <- &GetProgrammaticAccessTokenServiceResponse{Error: modelErr}
		return
	}

	passwordHashBytes := sha256.Sum256([]byte(model.Password))
	passwordHash := fmt.Sprintf("%x", passwordHashBytes[:]) // TODO: check sprintf usage

	chGetUserByPasswordResponse := make(chan *auth.GetUserByPasswordResponse)
	defer close(chGetUserByPasswordResponse)

	go s.authDb.GetUserByPassword(
		chGetUserByPasswordResponse, &auth.GetUserByPasswordModel{
			UserName:     model.UserName,
			PasswordHash: passwordHash,
		},
	)

	dataResponse := <-chGetUserByPasswordResponse
	if dataResponse.Error != nil {
		ch <- &GetProgrammaticAccessTokenServiceResponse{Error: dataResponse.Error}
		return
	}

	if !dataResponse.IsActive || !dataResponse.IsProgrammatic {
		ch <- &GetProgrammaticAccessTokenServiceResponse{Error: errors.New("programmatic user is not found")}
		return
	}

	tokenString, err := s.generateJwt(dataResponse.Id, dataResponse.UserName, dataResponse.Email, model.ExpiryDays)

	if err != nil {
		ch <- &GetProgrammaticAccessTokenServiceResponse{Error: err}
		return
	}

	ch <- &GetProgrammaticAccessTokenServiceResponse{
		JwtToken: tokenString,
	}
}

func (s *AuthService) generateJwt(id int64, username string, email string, expiryDays int32) (string, error) {
	expiresAt := jwt.NewNumericDate(time.Now().Add(time.Minute * 30))
	if expiryDays != 0 {
		expiresAt = jwt.NewNumericDate(time.Now().AddDate(0, 0, int(expiryDays)))
	}

	// TODO: if we implement permission feature, also fetch permissions from db and attach them to the jwt claims
	claims := util.CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Delivery Hero",
			Subject:   strconv.FormatInt(id, 10),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: expiresAt,
		},
		Username: username,
		Email:    email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.environment.Get(env.JwtSecret)))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
