package authorize

type ResetPasswordCommand struct {
	Token    string
	Password string
}
