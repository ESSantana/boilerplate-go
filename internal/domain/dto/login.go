package dto

type LoginRequest struct {
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

func (r LoginRequest) IsValid() bool {
	return r.Email != "" && r.PasswordHash != ""
}

type LoginResponse struct {
	Token string `json:"token"`
}

type RecoverPasswordRequest struct {
	Email string `json:"email"`
}

func (r RecoverPasswordRequest) IsValid() bool {
	return r.Email != ""
}

type RecoverPasswordResponse struct {
	Message string `json:"message"`
}
