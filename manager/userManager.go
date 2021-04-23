package manager

import (
	"github.com/pkg/errors"
	"landlords/helper/encryption"
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

func NewUserManager(account, password string) (*UserManager, error) {
	manager := &UserManager{}
	manager.User = &obj.User{}
	if redis.Exists(model.ACCOUNTDATA + account) {
		ad := model.GetAccountData(account)
		ps := string(encryption.AesDeCrypt([]byte(ad.PassWord)))
		if password != ps {
			return nil, errors.New("密码错误")
		}
		model.GetUserInfo(ad.UserId, manager.User)
	} else {
		ad := &obj.AccountData{}
		manager.User.Id = uuid.GetUUID()
		ad.Account = account
		ad.UserId = manager.User.Id
		ad.PassWord = string(encryption.AcesEncrypts([]byte(password)))
		model.CreateUserInfo(account, ad.UserId, ad, manager.User)
	}
	return manager, nil
}

func (u *UserManager) SetName(name string) {
	u.User.Name = name
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
