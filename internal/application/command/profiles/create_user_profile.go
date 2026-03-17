package profiles

type CreateUserProfileCommand struct {
	UserID      int
	FullName    string
	Age         int
	WeightKg    float32
	Goal        string
	Description string
}
