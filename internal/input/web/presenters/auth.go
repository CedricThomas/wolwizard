package presenters

type AuthResponse struct {
	Token string `json:"token"`
}

type ErrorResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

type AuthVerifyResponse struct {
	Valid    bool   `json:"status"`
	Username string `json:"username"`
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

func AuthVerifySuccess(username string) AuthVerifyResponse {
	return AuthVerifyResponse{
		Valid:    true,
		Username: username,
	}
}
