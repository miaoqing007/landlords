package manager

import (
	"app/helper/uuid"
	"app/model"
	"app/obj"
	"app/redis"
)

type UserManager struct {
	*obj.User
}

func NewUserManager(platformid string) (*UserManager, error) {
	manager := &UserManager{}
	manager.User = &obj.User{}
	if redis.Exists(model.PLATFORMID + platformid) {
		model.GetUserInfo(platformid, manager.User)
	} else {
		manager.User.Id = uuid.GetUUID()
		manager.User.PlatformId = platformid
		pd := &obj.PlatformData{}
		pd.Id = platformid
		pd.UserId = manager.User.Id
		model.CreateUserInfo(platformid, pd.UserId, pd, manager.User)
	}
	return manager, nil
}

func (u *UserManager) SetNameId(name string) {
	u.Name = name
	redis.HSet(model.NAMEIDKEY, name, u.Id)
	model.UpdateUserInfo(u.Id, u.User)
}
