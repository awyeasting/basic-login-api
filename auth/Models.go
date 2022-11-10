package auth

type NewUser struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type LoginInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserInfo struct {
	FirstName *string `json:"first_name" sql:"first_name"`
	LastName  *string `json:"last_name" sql:"last_name"`
	Username  *string `json:"username" sql:"username"`
	Email     *string `json:"email" sql:"email"`
}
