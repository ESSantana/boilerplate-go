package dto

type LoginRequest struct {
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type RecoverPasswordRequest struct {
	Email string `json:"email"`
}

type RecoverPasswordResponse struct {
	Message string `json:"message"`
}
