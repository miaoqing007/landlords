package operatecard

import (
	"app/enmu"
	"app/helper/conv"
	"app/helper/util"
	"sort"
)

func ComparisonTwoPlayersCards(wasteCards, newCards []string) bool {
	wasteArrayCards := getCardsValue(wasteCards)
	if len(wasteCards) != len(wasteArrayCards) {
		return false
	}
	newArrayCards := getCardsValue(newCards)
	if len(newCards) != len(newArrayCards) {
		return false
	}
	wasteType := judgeCardsType(wasteArrayCards)
	if wasteType == enmu.ERROR_TYPE {
		return false
	}
	newType := judgeCardsType(newArrayCards)
	if newType == enmu.ERROR_TYPE {
		return false
	}
	if wasteType == newType {
		judgeTwoPlayerCardsInSameType(newType, newArrayCards, wasteArrayCards)
	} else {

	}
	return true
}

func judgeTwoPlayerCardsInSameType(Type enmu.CardType, newArrayCards, wasteArrayCards []int) bool {
	switch Type {
	case enmu.SINGLE:
		if newArrayCards[0] <= wasteArrayCards[0] {
			return false
		}
	case enmu.DOUBLE:
	case enmu.THREE:
	}
	return true
}

func judgeCardsType(arrayCards []int) enmu.CardType {
	switch len(arrayCards) {
	case 1:
		return enmu.SINGLE
	case 2:
		if judgeArrayIfAllIsSameValue(arrayCards) {
			return enmu.DOUBLE
		} else {
			if judgeArrayIfIsKing_Bomb(arrayCards) {
				return enmu.KING_BOMB
			}
		}
	case 3:
		if judgeArrayIfAllIsSameValue(arrayCards) {
			return enmu.THREE
		}
	case 4:
		if judgeArrayIfAllIsSameValue(arrayCards) {
			return enmu.BOMB
		} else {
			if judgeArrayIfIsDoubleThree_And_One(arrayCards) {
				return enmu.THREE_AND_ONE
			} else {
				if judgeArrayIfIsDouble_Alone(arrayCards) {
					return enmu.DOUBLE_ALONE
				}
			}
		}
	case 5:
		if judgeArrayIfIsSingle_Alone(arrayCards) {
			return enmu.SINGLE_ALONE
		}
	}
	return enmu.ERROR_TYPE
}

func getCardsValue(cards []string) []int {
	intArrayCards := make([]int, 0)
	for _, card := range cards {
		intArrayCards = append(intArrayCards, conv.ParseInt(util.SubString(card, 1, len(card))))
	}
	return intArrayCards
}

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

func judgeArrayIfIsDouble_Alone(array []int) bool {
	if len(array)%2 != 0 {
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

func judgeArrayIfIsKing_Bomb(array []int) bool {
	if len(array) != 2 {
		return false
	}
	if array[0] == 88 && array[1] == 99 || array[0] == 99 && array[1] == 88 {
		return true
	}
	return false
}

func judgeArrayIfIsSingle_Alone(array []int) bool {
	sort.Ints(array)
	for i := 0; i < len(array)-1; i++ {
		if array[i]+1 != array[i+1] {
			return false
		}
	}
	return true
}
