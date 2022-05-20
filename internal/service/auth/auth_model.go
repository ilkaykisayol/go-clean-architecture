package auth

type LoginServiceModel struct {
	UserName string `validate:"required"`
	Password string `validate:"required"`
}

type GetAccessTokenServiceModel struct {
	RefreshToken string `validate:"required"`
}

type GetProgrammaticAccessTokenServiceModel struct {
	UserName   string `validate:"required"`
	Password   string `validate:"required"`
	ExpiryDays int32  `validate:"required"`
}
