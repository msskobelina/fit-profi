package authorize

type LoginUserCommand struct {
	Email    string
	Password string
}

type LoginUserResult struct {
	Token    string
	UserID   int
	FullName string
	Email    string
}
