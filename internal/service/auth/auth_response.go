package auth

type LoginServiceResponse struct {
	Error        error `json:"-"`
	JwtToken     string
	RefreshToken string
}

type GetProgrammaticAccessTokenServiceResponse struct {
	Error        error `json:"-"`
	JwtToken     string
	RefreshToken string
}
