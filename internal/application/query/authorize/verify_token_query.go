package authorize

type VerifyTokenQuery struct {
	Token string
}

type VerifyTokenResult struct {
	UserID int
	Role   string
}
