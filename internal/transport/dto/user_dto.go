package transport

type RegisterUserRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
	UserType string `json:"user_type"`
}

type RegisterUserResponse struct {
	Token string `json:"token"`
}

type LoginUserRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	Token string `json:"token"`
}