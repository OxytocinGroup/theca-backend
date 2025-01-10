package pkg

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

type LoginResponse struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	Username string `json:"username"`
}
