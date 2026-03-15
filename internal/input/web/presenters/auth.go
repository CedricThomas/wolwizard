package presenters

type AuthResponse struct {
	Token string `json:"token"`
}

type ErrorResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

func AuthSuccess(token string) AuthResponse {
	return AuthResponse{
		Token: token,
	}
}

func AuthError(err error) ErrorResponse {
	return ErrorResponse{
		Status: "error",
		Error:  err.Error(),
	}
}

func AuthVerifySuccess() AuthResponse {
	return AuthResponse{
		Token: "valid",
	}
}
