package manager

import "sync"

var _player *Players

type Players struct {
	idInMap sync.Map //map[id]player
}

func init() {
	_player = &Players{}
}

func (p *Players) getPlayer(id string) *Player {
	if v, ok := p.idInMap.Load(id); ok {
		return v.(*Player)
	}
	return nil
}

func (p *Players) addPlayer(id string, player *Player) {
	p.idInMap.Store(id, player)
}

func (p *Players) deletePlayer(id string) {
	p.idInMap.Delete(id)
}

func GetPlayer(id string) *Player {
	return _player.getPlayer(id)
}

func AddPlayer(id string, player *Player) {
	_player.addPlayer(id, player)
}

func DeletePlayer(id string) {
	_player.deletePlayer(id)
}
