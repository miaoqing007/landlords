package manager

type UserManager struct {
	UserId   string
	UserName string
}

func NewUserManager(userId string) (*UserManager, error) {
	return &UserManager{}, nil
}
