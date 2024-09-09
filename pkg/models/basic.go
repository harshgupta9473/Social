package models

type NewRegRequest struct {
	// UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogInResponse struct {
	// UserName  string `json:"username"`
	Email     string `json:"email"`
}


type ProfileReq struct {
	UserName  string   `json:"username"`
	Email     string   `json:"email"`
	FirsName  string   `json:"firstname"`
	LastName  string   `json:"lastname"`
	Bio       string   `json:"bio"`
	Interests []string `json:"interests"`
	Location  string   `json:"location"`
}

