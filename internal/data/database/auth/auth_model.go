package auth

type GetUserByPasswordModel struct {
	UserName     string `validate:"required"`
	PasswordHash string `validate:"required"`
}

type GetUserByRefreshTokenModel struct {
	RefreshToken string `validate:"required"`
}

type AddRefreshTokenModel struct {
	UserId       int    `validate:"required"`
	RefreshToken string `validate:"required"`
}
