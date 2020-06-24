package model

import (
	"landlords/obj"
	"landlords/redis"
)

var (
	USERKEY    = "User_Key:"
	NAMEIDKEY  = "Name_Id_Key"
	PLATFORMID = "Platform_Id_Key:"
)

func CreateUserInfo(platformid, id string, val_1, val_2 interface{}) error {
	if err := redis.HMSet(PLATFORMID+platformid, val_1, 0); err != nil {
		return err
	}
	return redis.HMSet(USERKEY+id, val_2, 0)
}

func GetUserInfo(platformId string, user *obj.User) error {
	pd := &obj.PlatformData{}
	if err := redis.HGetAll(PLATFORMID+platformId, pd); err != nil {
		return err
	}
	return redis.HGetAll(USERKEY+pd.UserId, user)
}

func UpdateUserInfo(id string, val interface{}) error {
	return redis.HMSet(USERKEY+id, val, 0)
}
