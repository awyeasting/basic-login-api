package google

type AuthError interface {
	Error() string
	Type() string
}

type InvalidGoogleLogin struct {
}

func (e InvalidGoogleLogin) Error() string {
	return "Invalid Google Login"
}

func (e InvalidGoogleLogin) Type() string {
	return "invalid_google_login_error"
}
