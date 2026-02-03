package entity

type UserCreated struct{}

func (UserCreated) GetName() string {
	return "fit_profi_user_created_total"
}

func (UserCreated) GetDescription() string {
	return "total number of created users"
}

func (UserCreated) GetLabels() []string {
	return []string{"source"}
}
