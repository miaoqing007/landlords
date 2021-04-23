package model

import (
	"landlords/obj"
	"landlords/redis"
)

var (
	USERKEY     = "User_Key:"
	NAMEIDKEY   = "Name_Id_Key"
	ACCOUNTDATA = "Account_Data_Key:"
)

func CreateUserInfo(account, id string, val_1, val_2 interface{}) error {
	if err := redis.HMSet(ACCOUNTDATA+account, val_1, 0); err != nil {
		return err
	}
	return redis.HMSet(USERKEY+id, val_2, 0)
}

func GetUserInfo(useId string, user *obj.User) error {
	return redis.HGetAll(USERKEY+useId, user)
}

func GetAccountData(account string) *obj.AccountData {
	pd := &obj.AccountData{}
	if err := redis.HGetAll(ACCOUNTDATA+account, pd); err != nil {
		return nil
	}
	return pd
}

func UpdateUserInfo(id string, val interface{}) error {
	if err := redis.HMSet(USERKEY+id, val, 0); err != nil {
		return err
	}
	return nil
}
