package auth

type GetUserByPasswordResponse struct {
	Error          error `json:"-"`
	Id             int64
	UserName       string
	Email          string
	IsActive       bool
	IsProgrammatic bool
}

type GetUserByRefreshTokenResponse struct {
	Error          error `json:"-"`
	Id             int64
	UserName       string
	Email          string
	IsActive       bool
	IsProgrammatic bool
}

type AddRefreshTokenResponse struct {
	Error error `json:"-"`
}
