package enmu

type CardType int

const (
	_CardTypeStatus CardType = iota
	SINGLE          CardType = 1  //单张
	DOUBLE          CardType = 2  //对子
	THREE           CardType = 3  //三不带
	THREE_AND_ONE   CardType = 4  //三带一
	BOMB            CardType = 5  //炸弹
	DOUBLE_ALONE    CardType = 6  //连对
	SINGLE_ALONE    CardType = 7  //顺子
	KING_BOMB       CardType = 8  //王炸
	PLANE           CardType = 9  //飞机
	PLANE_SINGLE    CardType = 10 //飞机带单
	PLANE_DOUBLE    CardType = 11 //飞机带双
	THREE_AND_TWO   CardType = 12 //三带二
	ERROR_TYPE      CardType = 13 //非法类型
)

const (
	ServerHost = "127.0.0.1"
	ServerPort = "8888"
	RedisPort  = "6379"
)
