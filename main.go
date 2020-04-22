package main

import (
	"fmt"
	"gophers/deck"
	"strings"
)

//State vvdfd
type State int8

//fdfgr
const (
	StatePlayerTurn State = iota
	StateDealerTurn
	StateHandOver
)

//Deal fdfgr
func Deal(gs GameStat) GameStat {
	res := clone(gs)
	var card deck.Card
	res.Dealer = make(Hand, 0, 5)
	res.Player = make(Hand, 0, 5)
	for i := 0; i < 2; i++ {
		card, res.Deck = Draw(res.Deck)
		res.Player = append(res.Player, card)
		card, res.Deck = Draw(res.Deck)
		res.Dealer = append(res.Dealer, card)
	}
	res.State = StatePlayerTurn
	return res
}

//Shuffle fdf
func Shuffle(gs GameStat) GameStat {
	res := clone(gs)
	res.Deck = deck.New(deck.Deck(3), deck.Shuffle)
	return res
}

//Hit fdf
func Hit(gs GameStat) GameStat {
	res := clone(gs)
	handp := res.CurrentPlayer()
	var card deck.Card
	card, res.Deck = Draw(res.Deck)
	*handp = append(*handp, card)
	if handp.Score() > 21 {
		return Stand(res)
	}

	return res
}

//Stand frgr
func Stand(gs GameStat) GameStat {
	gs = clone(gs)
	gs.State++
	return gs
}

func main() {

	var gs GameStat
	gs = Shuffle(gs)
	gs = Deal(gs)
	var input string
	for gs.State == StatePlayerTurn {
		fmt.Println("player:", gs.Player)
		fmt.Println("dealer:", gs.Dealer.DealerString())
		fmt.Println("what will you do? (h)it or (s)tand")
		fmt.Scanf("%s\n", &input)
		switch input {
		case "h":
			gs = Hit(gs)
		case "s":
			gs = Stand(gs)
		default:
			fmt.Println("invalid arg", input)
		}
	}

	for gs.State == StateDealerTurn {
		if gs.Dealer.Score() < 16 || (gs.Dealer.Score() == 17 && gs.Dealer.MinScore() != 17) {
			gs = Hit(gs)
		} else {
			gs = Stand(gs)
		}
	}
	EndHand(gs)

}

//EndHand fgrfg
func EndHand(gs GameStat) GameStat {
	res := clone(gs)
	fmt.Println("== GAME OVER ==")

	dScore, pScore := res.Dealer.Score(), res.Player.Score()

	switch {
	case pScore > 21:
		fmt.Println("you are busted")
	case dScore > 21:
		fmt.Println("dealer busted!")
	case pScore > dScore:
		fmt.Println("you won!")
	case pScore < dScore:
		fmt.Println("you lost!")
	case dScore == pScore:
		fmt.Println("draw")
	}

	fmt.Println("player:", res.Player, "score:", pScore)
	fmt.Println("dealer:", res.Dealer, "score:", dScore)

	return res
}

//DealerString fdf
func (h Hand) DealerString() string {
	return h[0].String() + ", " + "**HIDDEN**"
}

//Draw ffgf
func Draw(cs []deck.Card) (deck.Card, []deck.Card) {
	return cs[0], cs[1:]
}

//Hand df
type Hand []deck.Card

func (h Hand) String() string {
	d := make([]string, len(h))
	for i := range h {
		d[i] = h[i].String()
	}
	return strings.Join(d, ", ")
}

//MinScore vfvf
func (h Hand) MinScore() int {
	res := 0
	for _, card := range h {
		res += min(int(card.Rank), 10)
	}
	return res
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

//Score dfvd
func (h Hand) Score() int {
	min := h.MinScore()
	if min > 11 {
		return min
	}
	for _, card := range h {
		if card.Rank == deck.Ace {
			return min + 10
		}
	}
	return min
}

func clone(gs GameStat) GameStat {
	res := GameStat{
		Deck:   make([]deck.Card, len(gs.Deck)),
		State:  gs.State,
		Dealer: make(Hand, len(gs.Dealer)),
		Player: make(Hand, len(gs.Player)),
	}
	copy(res.Deck, gs.Deck)
	copy(res.Player, gs.Player)
	copy(res.Dealer, gs.Dealer)

	return res
}

//CurrentPlayer dfd
func (gs *GameStat) CurrentPlayer() *Hand {
	switch gs.State {
	case StateDealerTurn:
		return &gs.Dealer
	case StatePlayerTurn:
		return &gs.Player
	default:
		panic("something went wrong")
	}
}

//GameStat fvfv
type GameStat struct {
	Deck   []deck.Card
	State  State
	Player Hand
	Dealer Hand
}
