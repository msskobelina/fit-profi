package entity

type LoginFailed struct{}

func (LoginFailed) GetName() string {
	return "fit_profi_login_failed_total"
}

func (LoginFailed) GetDescription() string {
	return "total number of failed login"
}

func (LoginFailed) GetLabels() []string {
	return []string{"reason"}
}
