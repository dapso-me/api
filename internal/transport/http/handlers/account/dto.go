package account_handler

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerReq struct {
	Email string `json:"email"`
}

type confirmRegisterReq struct {
	Email    string `json:"email"`
	Code     string `json:"code"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type updateReq struct {
	Name string `json:"name"`
}

type updagePasswordReq struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}
