package manager

import (
	"app/obj"
)

type UserManager struct {
	*obj.User
}

func NewUserManager(id string) (*UserManager, error) {
	manager := &UserManager{}
	manager.User = &obj.User{}
	manager.User.Id = id
	return manager, nil
}
