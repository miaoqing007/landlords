package operatecard

import (
	"landlords/enmu"
	"landlords/helper/conv"
	"landlords/helper/util"
	"sort"
)

//比对该玩家的牌是否打过上一位玩家
func ComparisonTwoPlayersCards(wasteCards, newCards []string) bool {
	wasteArrayCards := GetCardsValue(wasteCards)
	if len(wasteCards) != len(wasteArrayCards) {
		return false
	}
	newArrayCards := GetCardsValue(newCards)
	if len(newCards) != len(newArrayCards) {
		return false
	}
	wasteType := JudgeCardsType(wasteArrayCards)
	if wasteType == enmu.ERROR_TYPE {
		return false
	}
	if wasteType == enmu.KING_BOMB {
		return false
	}
	newType := JudgeCardsType(newArrayCards)
	if newType == enmu.ERROR_TYPE {
		return false
	}
	if len(newArrayCards) == len(wasteArrayCards) {
		if newType == wasteType {
			return comparisonTwoPalyerCardsSize(newArrayCards, wasteArrayCards, newType)
		} else if newType == enmu.KING_BOMB || newType == enmu.BOMB {
			return true
		}
	} else if newType == enmu.KING_BOMB || newType == enmu.BOMB {
		return true
	}
	return false
}

//判断牌类型
func JudgeCardsType(arrayCards []int) enmu.CardType {
	switch len(arrayCards) {
	case 1:
		return enmu.SINGLE
	case 2:
		if judgeArrayIfAllIsSameValue(arrayCards) {
			return enmu.DOUBLE
		} else if judgeArrayIfIsKing_Bomb(arrayCards) {
			return enmu.KING_BOMB
		}
	case 3:
		if judgeArrayIfAllIsSameValue(arrayCards) {
			return enmu.THREE
		}
	case 4:
		if judgeArrayIfAllIsSameValue(arrayCards) {
			return enmu.BOMB
		} else if judgeArrayIfIsDoubleThree_And_One(arrayCards) {
			return enmu.THREE_AND_ONE
		}
	case 5:
		if judgeArrayIfIsSingle_Alone(arrayCards) {
			return enmu.SINGLE_ALONE
		} else if judgeThree_And_Two(arrayCards) {
			return enmu.THREE_AND_TWO
		}
	case 6:
		if judgeArrayIfIsSingle_Alone(arrayCards) {
			return enmu.SINGLE_ALONE
		} else if judgeArrayIfIsDouble_Alone(arrayCards) {
			return enmu.DOUBLE_ALONE
		} else if judgePlane(arrayCards) {
			return enmu.PLANE
		}
	case 7, 9, 11, 13:
		if judgeArrayIfIsSingle_Alone(arrayCards) {
			return enmu.SINGLE_ALONE
		}
	default:
		if judgeArrayIfIsSingle_Alone(arrayCards) {
			return enmu.SINGLE_ALONE
		} else if judgeArrayIfIsDouble_Alone(arrayCards) {
			return enmu.DOUBLE_ALONE
		} else if judgePlane(arrayCards) {
			return enmu.PLANE
		} else if judgePlane_Single(arrayCards) {
			return enmu.PLANE_SINGLE
		} else if judgePlane_Double(arrayCards) {
			return enmu.PLANE_DOUBLE
		}
	}
	return enmu.ERROR_TYPE
}

//获取牌面值
func GetCardsValue(cards []string) []int {
	intArrayCards := make([]int, 0)
	for _, card := range cards {
		intArrayCards = append(intArrayCards, conv.ParseInt(util.SubString(card, 1, len(card))))
	}
	return intArrayCards
}

//判断是否是对子
func judgeArrayIfAllIsSameValue(array []int) bool {
	if len(array) == 0 {
		return false
	}
	temp := array[0]
	for _, v := range array {
		if temp != v {
			return false
		}
	}
	return true
}

//判断是否是三带一
func judgeArrayIfIsDoubleThree_And_One(array []int) bool {
	num1 := 0
	num2 := 0
	temp := array[0]
	for _, v := range array {
		if temp == v {
			num1++
		} else {
			num2++
		}
	}
	if num1 == 1 && num2 == 3 || num1 == 3 && num2 == 1 {
		return true
	}
	return false
}

//判断是否是连对
func judgeArrayIfIsDouble_Alone(array []int) bool {
	if len(array) < 6 || len(array)%2 != 0 {
		return false
	}
	sort.Ints(array)
	for i := 0; i < len(array); i++ {
		if array[i] != array[i+1] {
			return false
		}
		i++
	}
	return true
}

//判断是否是王炸
func judgeArrayIfIsKing_Bomb(array []int) bool {
	if len(array) != 2 {
		return false
	}
	if array[0] == 88 && array[1] == 99 || array[0] == 99 && array[1] == 88 {
		return true
	}
	return false
}

//判断是否是顺子
func judgeArrayIfIsSingle_Alone(array []int) bool {
	sort.Ints(array)
	for i := 0; i < len(array)-1; i++ {
		if array[i]+1 != array[i+1] {
			return false
		}
	}
	return true
}

//判断是否飞机
func judgePlane(cards []int) bool {
	if len(cards) < 6 || len(cards)%3 != 0 {
		return false
	}
	sort.Ints(cards)
	for i := 0; i < len(cards)-3; {
		if cards[i] != cards[i+1] || cards[i+1] != cards[i+2] || cards[i] != cards[i+2] {
			return false
		}
		i += 3
	}
	return true
}

//判断是否是飞机带单
func judgePlane_Single(cards []int) bool {
	if len(cards) < 8 || len(cards)%4 != 0 {
		return false
	}
	m := array2map(cards)
	return judgeMap(m, len(cards), 1)
}

//判断是否是飞机带双
func judgePlane_Double(cards []int) bool {
	if len(cards) < 10 || len(cards)%5 != 0 {
		return false
	}
	m := array2map(cards)
	return judgeMap(m, len(cards), 2)
}

//判断是否是三带二
func judgeThree_And_Two(cards []int) bool {
	if len(cards) != 5 {
		return false
	}
	m := array2map(cards)
	return judgeMap(m, len(cards), 2)
}

//比对牌大小
func comparisonTwoPalyerCardsSize(newCards, wasteCards []int, Type enmu.CardType) bool {
	switch Type {
	case enmu.THREE_AND_ONE, enmu.THREE_AND_TWO:
		sort.Ints(newCards)
		sort.Ints(wasteCards)
		return newCards[2] > wasteCards[2]
	case enmu.PLANE_SINGLE, enmu.PLANE_DOUBLE:
		nc := 0
		wc := 0
		nm := array2map(newCards)
		wm := array2map(wasteCards)
		for k, v := range nm {
			if v == 3 {
				nc += k
			}
		}
		for k, v := range wm {
			if v == 3 {
				wc += k
			}
		}
		return nc > wc
	default:
		nc := 0
		wc := 0
		for _, v := range newCards {
			nc += v
		}
		for _, v := range wasteCards {
			wc += v
		}
		return nc > wc
	}
}

func array2map(cards []int) map[int]int {
	m := make(map[int]int)
	for _, v := range cards {
		if m[v] != 0 {
			m[v]++
		} else {
			m[v] = 1
		}
	}
	return m
}

func judgeMap(m map[int]int, cardsLen, oneOrTwo int) bool {
	cardNum := 0
	th := 0
	tw := 0
	for _, v := range m {
		if v == 3 {
			cardNum += 3
			th++
		} else if v == oneOrTwo {
			cardNum += oneOrTwo
			tw++
		}
	}
	if cardNum != cardsLen || th != tw {
		return false
	}
	return true
}
