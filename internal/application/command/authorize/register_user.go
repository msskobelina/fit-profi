package authorize

type RegisterUserCommand struct {
	FullName string
	Email    string
	Password string
}

type RegisterUserResult struct {
	Token    string
	UserID   int
	FullName string
	Email    string
}
