package utils

import (
	"math/rand"
)

// 先手は黒(0), 後手は白(1)
type Board struct {
	Id            int64
	PlayerBoard   uint64
	OpponentBoard uint64
	Total         uint8
}

func (b *Board) BoardReset() {
	b.PlayerBoard = 0x0000000810000000
	b.OpponentBoard = 0x0000001008000000
	b.Total = 4
}

func (b *Board) BoardID() {
	boardId := rand.Int63()
	b.Id = boardId
}

func InitializeBoard() Board {
	var board Board
	board.BoardID()
	board.BoardReset()
	return board
}

func (b *Board) PlacingFrame(put uint64) {
	var rev uint64
	for k := 0; k < 8; k++ {
		var rev_ uint64
		mask := Transfer(put, k)
		for (mask != 0) && ((mask & b.OpponentBoard) != 0) {
			rev_ |= mask
			mask = Transfer(mask, k)
		}
		if (mask & b.PlayerBoard) != 0 {
			rev |= rev_
		}
	}
	//反転する
	b.PlayerBoard ^= put | rev
	b.OpponentBoard ^= rev
}

// 8方向にシフトするヘルパー関数
func Transfer(put uint64, k int) uint64 {
	switch k {
	case 0: //上
		return (put << 8) & 0xffffffffffffff00
	case 1: //右上
		return (put << 7) & 0x7f7f7f7f7f7f7f00
	case 2: //右
		return (put >> 1) & 0x7f7f7f7f7f7f7f7f
	case 3: //右下
		return (put >> 9) & 0x007f7f7f7f7f7f7f
	case 4: //下
		return (put >> 8) & 0x00ffffffffffffff
	case 5: //左下
		return (put >> 7) & 0x00fefefefefefefe
	case 6: //左
		return (put << 1) & 0xfefefefefefefefe
	case 7: //左上
		return (put << 9) & 0xfefefefefefefe00
	default:
		return 0
	}
}
