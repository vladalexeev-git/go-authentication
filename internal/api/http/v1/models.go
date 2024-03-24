package v1

type errorResponse struct {
	Error string `json:"error"`
}

type accountCreateRequest struct {
	Email    string `json:"email" binding:"required,email,lte=255"`
	Username string `json:"username" binding:"required,alphanum,gte=4,lte=16"`
	Password string `json:"password" binding:"required,gte=8,lte=64"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type tokenRequest struct {
	Password string `json:"password" binding:"required"`
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
}
