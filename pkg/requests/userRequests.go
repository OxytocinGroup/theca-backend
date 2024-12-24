package requests

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=8"`
}

type EmailVerifyRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ChangePasswordRequest struct {
	Password string `json:"password" binding:"required"`
}

type RequestPasswordReset struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPassword struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required"`
}
