package enmu

type CardType int

const (
	_CardTypeStatus CardType = iota
	SINGLE          CardType = 1 //单张
	DOUBLE          CardType = 2 //对子
	THREE           CardType = 3 //三不带
	THREE_AND_ONE   CardType = 4 //三带一
	BOMB            CardType = 5 //炸弹
	DOUBLE_ALONE    CardType = 6 //连对
	SINGLE_ALONE    CardType = 7 //顺子
	KING_BOMB       CardType = 8 //王炸
	ERROR_TYPE      CardType = 9 //非法类型
)
