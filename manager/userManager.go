package manager

import (
	"landlords/helper/uuid"
	"landlords/model"
	"landlords/obj"
	"landlords/redis"
)

type UserManager struct {
	*obj.User
	roomId    string
	piecewise int //分段
}

func NewUserManager(platformid string) (*UserManager, error) {
	manager := &UserManager{}
	manager.User = &obj.User{}
	if redis.Exists(model.PLATFORMID + platformid) {
		model.GetUserInfo(platformid, manager.User)
	} else {
		pd := &obj.PlatformData{}
		manager.User.Id = uuid.GetUUID()
		manager.User.PlatformId = platformid
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

func (u *UserManager) SetRoomId(roomId string) {
	u.roomId = roomId
}

func (u *UserManager) GetRoomId() string {
	return u.roomId
}

func (u *UserManager) SetPiecewise(piecewise int) {
	u.piecewise = piecewise
}

func (u *UserManager) GetPiecewise() int {
	return u.piecewise
}
