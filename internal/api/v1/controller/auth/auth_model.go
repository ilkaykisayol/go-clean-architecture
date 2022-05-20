package auth

type LoginModel struct {
	UserName string `json:"UserName"`
	Password string `json:"Password"`
}

type GetAccessTokenModel struct {
	RefreshToken string `json:"RefreshToken"`
}

type GetProgrammaticAccessTokenModel struct {
	UserName   string `json:"UserName"`
	Password   string `json:"Password"`
	ExpiryDays int32  `json:"ExpiryDays"`
}
