package card

import (
	twoonone_pb "github.com/nanachi-sh/susubot-code/plugin/twoonone/pkg/protos/twoonone"
)

type Card struct {
	Number CardNumber
}

type CardNumber int

const (
	THREE CardNumber = iota
	FOUR
	FIVE
	SIX
	SEVEN
	EIGHT
	NINE
	TEN
	J
	Q
	K
	A
	TWO
	JOKER
	KING
)

func FormatInternalCard2Protobuf(c Card) *twoonone_pb.Card {
	return &twoonone_pb.Card{
		Number: twoonone_pb.Card_Number(c.Number),
	}
}

func FormatInternalCards2Protobuf(cs []Card) []*twoonone_pb.Card {
	ret := []*twoonone_pb.Card{}
	for _, v := range cs {
		ret = append(ret, FormatInternalCard2Protobuf(v))
	}
	return ret
}

func FormatProtobuf2InternalCard(c *twoonone_pb.Card) Card {
	return Card{
		Number: CardNumber(c.Number),
	}
}

func FormatProtobuf2InternalCards(cs []*twoonone_pb.Card) []Card {
	ret := []Card{}
	for _, v := range cs {
		ret = append(ret, FormatProtobuf2InternalCard(v))
	}
	return ret
}
