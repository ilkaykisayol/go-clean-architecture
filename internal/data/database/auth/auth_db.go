package auth

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

type IAuthDb interface {
	GetUserByPassword(ch chan *GetUserByPasswordResponse, model *GetUserByPasswordModel)
	GetUserByRefreshToken(ch chan *GetUserByRefreshTokenResponse, model *GetUserByRefreshTokenModel)
	AddRefreshToken(ch chan *AddRefreshTokenResponse, model *AddRefreshTokenModel)
}

type AuthDb struct {
	loggr            logger.ILogger
	validatr         validator.IValidator
	cachr            cacher.ICacher
	environment      env.IEnvironment
	connectionString string
	driverName       string
	timeout          time.Duration
}

// NewAuthDb
// Returns a new AuthDb.
func NewAuthDb(environment env.IEnvironment, loggr logger.ILogger, validatr validator.IValidator, cachr cacher.ICacher) IAuthDb {
	db := AuthDb{
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

// GetUserByPassword
// Gets user response from postgresql by password.
func (d *AuthDb) GetUserByPassword(ch chan *GetUserByPasswordResponse, model *GetUserByPasswordModel) {
	modelErr := d.validatr.ValidateStruct(model)
	if modelErr != nil {
		ch <- &GetUserByPasswordResponse{Error: modelErr}
		return
	}

	connection, err := sql.Open(d.driverName, d.connectionString)
	if err != nil {
		ch <- &GetUserByPasswordResponse{Error: err}
		return
	}
	defer connection.Close()

	ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
	defer cancel()

	query := `select id, username, email, is_active, is_programmatic from users where username = $1 and password_hash = $2`

	var user GetUserByPasswordResponse
	var email sql.NullString
	dbErr := connection.QueryRowContext(ctx, query, model.UserName, model.PasswordHash).Scan(&user.Id, &user.UserName, &email, &user.IsActive, &user.IsProgrammatic)

	if dbErr != nil || dbErr == sql.ErrNoRows {
		ch <- &GetUserByPasswordResponse{Error: dbErr}
		return
	}

	if email.Valid {
		user.Email = email.String
	}

	ch <- &user
}

// GetUserByPassword
// Gets user response from postgresql by password.
func (d *AuthDb) GetUserByRefreshToken(ch chan *GetUserByRefreshTokenResponse, model *GetUserByRefreshTokenModel) {
	modelErr := d.validatr.ValidateStruct(model)
	if modelErr != nil {
		ch <- &GetUserByRefreshTokenResponse{Error: modelErr}
		return
	}

	connection, err := sql.Open(d.driverName, d.connectionString)
	if err != nil {
		ch <- &GetUserByRefreshTokenResponse{Error: err}
		return
	}
	defer connection.Close()

	ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
	defer cancel()

	query := `
	select u.id, username, email, is_active, is_programmatic
	from users as u
	inner join users_refresh_tokens as ur on u.id = ur.user_id
	where ur.refresh_token = $1
	and ur.expiry_date > now()`

	var user GetUserByRefreshTokenResponse
	var email sql.NullString
	dbErr := connection.QueryRowContext(ctx, query, model.RefreshToken).Scan(&user.Id, &user.UserName, &email, &user.IsActive, &user.IsProgrammatic)

	if dbErr != nil || dbErr == sql.ErrNoRows {
		ch <- &GetUserByRefreshTokenResponse{Error: dbErr}
		return
	}

	if email.Valid {
		user.Email = email.String
	}

	ch <- &user
}

// AddRefreshToken
// Adds refresh token to db
func (d *AuthDb) AddRefreshToken(ch chan *AddRefreshTokenResponse, model *AddRefreshTokenModel) {
	modelErr := d.validatr.ValidateStruct(model)
	if modelErr != nil {
		ch <- &AddRefreshTokenResponse{Error: modelErr}
		return
	}

	connection, err := sql.Open(d.driverName, d.connectionString)
	if err != nil {
		ch <- &AddRefreshTokenResponse{Error: err}
		return
	}
	defer connection.Close()

	ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
	defer cancel()

	query := `insert into users_refresh_tokens (user_id, refresh_token, expiry_date) values ($1, $2, current_timestamp + interval '30' day)`

	result, err := connection.ExecContext(ctx, query, model.UserId, model.RefreshToken)

	if err != nil {
		ch <- &AddRefreshTokenResponse{Error: err}
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		ch <- &AddRefreshTokenResponse{Error: err}
		return
	}
	if rows != 1 {
		ch <- &AddRefreshTokenResponse{Error: err}
		return
	}

	ch <- &AddRefreshTokenResponse{}
}
