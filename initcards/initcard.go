package initcards

import (
	"landlords/helper/conv"
	"github.com/golang/glog"
	"math/rand"
	"time"
)

var cards = make([]string, 54)

//初始化一副牌
func InitNewCards() {
	start := 0
	for i := 3; i <= 16; i++ {
		if i == 16 {
			//大小王
			cards[start] = "Q88"
			cards[start+1] = "K99"
		} else {
			cards[start] = "A" + conv.FormatInt(i)
			cards[start+1] = "B" + conv.FormatInt(i)
			cards[start+2] = "C" + conv.FormatInt(i)
			cards[start+3] = "D" + conv.FormatInt(i)
			start += 4
		}
	}
	glog.Info("初始化牌完成")
}

func ShuffCards() []string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := len(cards)
	for n > 0 {
		randIndex := r.Intn(n)
		cards[n-1], cards[randIndex] = cards[randIndex], cards[n-1]
		n--
	}
	return cards
}
