// Copyright D-Platform Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"fmt"
	"math/rand"
	"sort"

	"github.com/D-PlatformOperatingSystem/plugin/plugin/dapp/pokerbull/types"
)

// PokerCardNum   ，4 * 13
var PokerCardNum = 52

// ColorOffset
var ColorOffset uint32 = 8

// ColorBitMask    bit
var ColorBitMask = 0xFF

// CardNumPerColor
var CardNumPerColor = 13

// CardNumPerGame
var CardNumPerGame = 5

const (
	// PokerbullResultX1     1
	PokerbullResultX1 = 1
	// PokerbullResultX2     2
	PokerbullResultX2 = 2
	// PokerbullResultX3     3
	PokerbullResultX3 = 3
	// PokerbullResultX4     4
	PokerbullResultX4 = 4
	// PokerbullResultX5     5
	PokerbullResultX5 = 5
	// PokerbullLeverageMax
	PokerbullLeverageMax = PokerbullResultX1
)

// NewPoker
func NewPoker() *types.PBPoker {
	poker := new(types.PBPoker)
	poker.Cards = make([]int32, PokerCardNum)
	poker.Pointer = int32(PokerCardNum - 1)

	for i := 0; i < PokerCardNum; i++ {
		color := i / CardNumPerColor
		num := i%CardNumPerColor + 1
		poker.Cards[i] = int32(color<<ColorOffset + num)
	}
	return poker
}

// Shuffle
func Shuffle(poker *types.PBPoker, rng int64) {
	rndn := rand.New(rand.NewSource(rng))

	for i := 0; i < PokerCardNum; i++ {
		idx := rndn.Intn(PokerCardNum - 1)
		tmpV := poker.Cards[idx]
		poker.Cards[idx] = poker.Cards[PokerCardNum-i-1]
		poker.Cards[PokerCardNum-i-1] = tmpV
	}
	poker.Pointer = int32(PokerCardNum - 1)
}

// Deal
func Deal(poker *types.PBPoker, rng int64) []int32 {
	if poker.Pointer < int32(CardNumPerGame) {
		logger.Error(fmt.Sprintf("Wait to be shuffled: deal cards [%d], left [%d]", CardNumPerGame, poker.Pointer+1))
		Shuffle(poker, rng+int64(poker.Cards[0]))
	}

	rndn := rand.New(rand.NewSource(rng))
	res := make([]int32, CardNumPerGame)
	for i := 0; i < CardNumPerGame; i++ {
		idx := rndn.Intn(int(poker.Pointer))
		tmpV := poker.Cards[poker.Pointer]
		res[i] = poker.Cards[idx]
		poker.Cards[idx] = tmpV
		poker.Cards[poker.Pointer] = res[i]
		poker.Pointer--
	}

	return res
}

// Result
func Result(cards []int32) int32 {
	temp := 0
	r := -1 //

	pk := newcolorCard(cards)

	//    10
	cardsC := make([]int, len(cards))
	for i := 0; i < len(pk); i++ {
		if pk[i].num > 10 {
			cardsC[i] = 10
		} else {
			cardsC[i] = pk[i].num
		}
	}

	//
	result := make([]int, 10)
	var offset = 0
	for x := 0; x < 3; x++ {
		for y := x + 1; y < 4; y++ {
			for z := y + 1; z < 5; z++ {
				if (cardsC[x]+cardsC[y]+cardsC[z])%10 == 0 {
					for j := 0; j < len(cardsC); j++ {
						if j != x && j != y && j != z {
							temp += cardsC[j]
						}
					}

					if temp%10 == 0 {
						r = 10 //   ，
					} else {
						r = temp % 10 //   ，
					}
					result[offset] = r
					offset++
				}
			}
		}
	}

	//
	if r == -1 {
		return -1
	}

	return int32(result[0])
}

// Leverage
func Leverage(hand *types.PBHand) int32 {
	result := hand.Result

	//    [1, 6]
	if result < 7 {
		return PokerbullResultX1
	}

	//    [7, 9]
	if result >= 7 && result < 10 {
		return PokerbullResultX2
	}

	var flowers = 0
	if result == 10 {
		for _, card := range hand.Cards {
			if (int(card) & ColorBitMask) > 10 {
				flowers++
			}
		}

		//
		if flowers < 4 {
			return PokerbullResultX3
		}

		//
		if flowers == 4 {
			return PokerbullResultX4
		}

		//
		if flowers == 5 {
			return PokerbullResultX5
		}
	}

	return PokerbullResultX1
}

type pokerCard struct {
	num   int
	color int
}

type colorCardSlice []*pokerCard

func (p colorCardSlice) Len() int {
	return len(p)
}
func (p colorCardSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p colorCardSlice) Less(i, j int) bool {
	if i >= p.Len() || j >= p.Len() {
		logger.Error("length error. slice length:", p.Len(), " compare lenth: ", i, " ", j)
	}

	if p[i] == nil || p[j] == nil {
		logger.Error("nil pointer at ", i, " ", j)
	}
	return p[i].num < p[j].num
}

func newcolorCard(a []int32) colorCardSlice {
	var cardS []*pokerCard
	for i := 0; i < len(a); i++ {
		num := int(a[i]) & ColorBitMask
		color := int(a[i]) >> ColorOffset
		cardS = append(cardS, &pokerCard{num, color})
	}

	return cardS
}

// CompareResult
func CompareResult(i, j *types.PBHand) bool {
	if i.Result < j.Result {
		return true
	}

	if i.Result == j.Result {
		return Compare(i.Cards, j.Cards)
	}

	return false
}

// Compare
func Compare(a []int32, b []int32) bool {
	cardA := newcolorCard(a)
	cardB := newcolorCard(b)

	if !sort.IsSorted(cardA) {
		sort.Sort(cardA)
	}
	if !sort.IsSorted(cardB) {
		sort.Sort(cardB)
	}

	maxA := cardA[len(a)-1]
	maxB := cardB[len(b)-1]
	if maxA.num != maxB.num {
		return maxA.num < maxB.num
	}

	return maxA.color < maxB.color
}
