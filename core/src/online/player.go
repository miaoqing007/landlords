package main

import (
	"github.com/pkg/errors"
	"landlords/command/router"
)

type Player struct {
	router *router.Router
	User   *UserManager
	innerChan chan interface{}
}

func (player *Player) InitBase(account, password string) error {
	player.RegisterPlayerMsg()
	userManger, err := NewUserManager(account, password)
	if err != nil {
		return errors.New("")
	}
	player.User = userManger
	return nil
}
