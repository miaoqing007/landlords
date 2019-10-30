package manager

import "app/obj"

type UserManager struct {
	*obj.User
}

func NewUserManager(userId string) (*UserManager, error) {
	return &UserManager{}, nil
}
