package process

import (
	"fmt"
	"main/utils"
)

type Game struct {
	board   utils.Board
	player1 utils.Player
	Status  int8 // -1のときはゲーム外、0のときは先手、1のときは後手がプレイ
}

func (g *Game) StartGame() {
	g.Status = 0
}

func (g *Game) EndGame() {
	g.Status = -1
	g.ShowBoard()
	g.CountBoard()
}

func ConvertToBit(i int, j int) uint64 {
	mask := uint64(0x8000000000000000)
	mask >>= j
	mask >>= i * 8
	return mask
}

func (g *Game) CanPut(put uint64) bool {
	legalBoard := g.MakeLegalBoard()
	return (put&legalBoard == put)
}

func (g *Game) Reverse(put uint64) {
	g.board.PlacingFrame(put)
	g.Status = 1 - g.Status
	g.board.Total += 1
}

func (g *Game) SwapBoard() {
	tmp := g.board.PlayerBoard
	g.board.PlayerBoard = g.board.OpponentBoard
	g.board.OpponentBoard = tmp
}

func (g *Game) MakeLegalBoard() uint64 {
	horizontalWatchBoard := g.board.OpponentBoard & 0x7e7e7e7e7e7e7e7e // 左右
	verticalWatchBoard := g.board.OpponentBoard & 0x00FFFFFFFFFFFF00   // 上下
	allSideWatchBoard := g.board.OpponentBoard & 0x007e7e7e7e7e7e00    // 全辺

	blankBoard := ^(g.board.PlayerBoard | g.board.OpponentBoard) // 空きマスのみ
	var tmp uint64                                               // 隣に相手の色があるか一時保存
	var legalBoard uint64                                        // 返り値

	// 8方向チェック
	// 左
	tmp = horizontalWatchBoard & (g.board.PlayerBoard << 1)
	tmp |= horizontalWatchBoard & (tmp << 1)
	tmp |= horizontalWatchBoard & (tmp << 1)
	tmp |= horizontalWatchBoard & (tmp << 1)
	tmp |= horizontalWatchBoard & (tmp << 1)
	tmp |= horizontalWatchBoard & (tmp << 1)
	legalBoard = blankBoard & (tmp << 1)

	// 右
	tmp = horizontalWatchBoard & (g.board.PlayerBoard >> 1)
	tmp |= horizontalWatchBoard & (tmp >> 1)
	tmp |= horizontalWatchBoard & (tmp >> 1)
	tmp |= horizontalWatchBoard & (tmp >> 1)
	tmp |= horizontalWatchBoard & (tmp >> 1)
	tmp |= horizontalWatchBoard & (tmp >> 1)
	legalBoard |= blankBoard & (tmp >> 1)

	// 上
	tmp = verticalWatchBoard & (g.board.PlayerBoard << 8)
	tmp |= verticalWatchBoard & (tmp << 8)
	tmp |= verticalWatchBoard & (tmp << 8)
	tmp |= verticalWatchBoard & (tmp << 8)
	tmp |= verticalWatchBoard & (tmp << 8)
	tmp |= verticalWatchBoard & (tmp << 8)
	legalBoard |= blankBoard & (tmp << 8)

	// 下
	tmp = verticalWatchBoard & (g.board.PlayerBoard >> 8)
	tmp |= verticalWatchBoard & (tmp >> 8)
	tmp |= verticalWatchBoard & (tmp >> 8)
	tmp |= verticalWatchBoard & (tmp >> 8)
	tmp |= verticalWatchBoard & (tmp >> 8)
	tmp |= verticalWatchBoard & (tmp >> 8)
	legalBoard |= blankBoard & (tmp >> 8)

	// 右斜め上
	tmp = allSideWatchBoard & (g.board.PlayerBoard << 7)
	tmp |= allSideWatchBoard & (tmp << 7)
	tmp |= allSideWatchBoard & (tmp << 7)
	tmp |= allSideWatchBoard & (tmp << 7)
	tmp |= allSideWatchBoard & (tmp << 7)
	tmp |= allSideWatchBoard & (tmp << 7)
	legalBoard |= blankBoard & (tmp << 7)

	// 左斜め上
	tmp = allSideWatchBoard & (g.board.PlayerBoard << 9)
	tmp |= allSideWatchBoard & (tmp << 9)
	tmp |= allSideWatchBoard & (tmp << 9)
	tmp |= allSideWatchBoard & (tmp << 9)
	tmp |= allSideWatchBoard & (tmp << 9)
	tmp |= allSideWatchBoard & (tmp << 9)
	legalBoard |= blankBoard & (tmp << 9)

	// 右斜め下
	tmp = allSideWatchBoard & (g.board.PlayerBoard >> 9)
	tmp |= allSideWatchBoard & (tmp >> 9)
	tmp |= allSideWatchBoard & (tmp >> 9)
	tmp |= allSideWatchBoard & (tmp >> 9)
	tmp |= allSideWatchBoard & (tmp >> 9)
	tmp |= allSideWatchBoard & (tmp >> 9)
	legalBoard |= blankBoard & (tmp >> 9)

	// 左斜め下
	tmp = allSideWatchBoard & (g.board.PlayerBoard >> 7)
	tmp |= allSideWatchBoard & (tmp >> 7)
	tmp |= allSideWatchBoard & (tmp >> 7)
	tmp |= allSideWatchBoard & (tmp >> 7)
	tmp |= allSideWatchBoard & (tmp >> 7)
	tmp |= allSideWatchBoard & (tmp >> 7)
	legalBoard |= blankBoard & (tmp >> 7)

	return legalBoard
}

func ScanUserInput() (int, int, error) {
	var userInput string
	fmt.Printf("置く場所を入力 : ")
	fmt.Scan(&userInput)
	i, j, err := isValid(userInput)
	if err != nil {
		return -1, -1, fmt.Errorf("invalid input")
	}
	return i, j, err
}

func isValid(input string) (int, int, error) {
	if len(input) != 2 {
		return -1, -1, fmt.Errorf("allowed input lengths is 2")
	}
	allowedFirstString := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	allowedSecondString := []string{"1", "2", "3", "4", "5", "6", "7", "8"}
	if Contains(allowedFirstString, string(input[0])) && Contains(allowedSecondString, string(input[1])) {
		return Find(allowedFirstString, string(input[0])), Find(allowedSecondString, string(input[1])), nil
	}
	return -1, -1, fmt.Errorf("invalid input")
}

func Contains(slice []string, key string) bool {
	for _, s := range slice {
		if s == key {
			return true
		}
	}
	return false
}

func Find(slice []string, key string) int {
	for i, s := range slice {
		if s == key {
			return i
		}
	}
	return -1
}

func InitializeGame() Game {
	var game Game
	game.board = utils.InitializeBoard()
	game.player1 = utils.InitializePlayer(game.board.Id, true)
	game.Status = -1
	return game
}

func (g *Game) IsPass() bool {
	playerLegalBoard := g.MakeLegalBoard()
	g.SwapBoard()
	opponentLegalBoard := g.MakeLegalBoard()
	g.SwapBoard()
	return playerLegalBoard == 0x0000000000000000 && opponentLegalBoard != 0x0000000000000000
}

func (g *Game) IsEnd() bool {
	playerLegalBoard := g.MakeLegalBoard()
	g.SwapBoard()
	opponentLegalBoard := g.MakeLegalBoard()
	g.SwapBoard()
	return playerLegalBoard == 0x0000000000000000 && opponentLegalBoard == 0x0000000000000000
}

func (g *Game) ShowBoard() {
	board := make([]string, 64)
	var playerBoard uint64
	var opponentBoard uint64
	if g.Status == 0 {
		playerBoard = g.board.PlayerBoard
		opponentBoard = g.board.OpponentBoard
	} else {
		playerBoard = g.board.OpponentBoard
		opponentBoard = g.board.PlayerBoard
	}
	for i := 0; i < 64; i++ {
		if playerBoard&(1<<(63-i)) != 0 {
			board[i] = "⚫️"
		} else if opponentBoard&(1<<(63-i)) != 0 {
			board[i] = "⚪️"
		} else {
			board[i] = "ー"
		}
	}

	fmt.Println(" A B C D E F G H")
	for i := 0; i < 64; i++ {
		if i%8 == 0 {
			fmt.Print(i/8 + 1)
		}
		fmt.Printf(board[i])
		if (i+1)%8 == 0 {
			fmt.Printf("\n")
		}
	}
	fmt.Println()
}

func (g *Game) CountBoard() {
	var playerBoard uint64
	var opponentBoard uint64
	var playerCount uint8
	var opponentCount uint8
	if g.Status == 0 {
		playerBoard = g.board.PlayerBoard
		opponentBoard = g.board.OpponentBoard
	} else {
		playerBoard = g.board.OpponentBoard
		opponentBoard = g.board.PlayerBoard
	}

	for i := 0; i < 64; i++ {
		if playerBoard&(1<<(63-i)) != 0 {
			playerCount++
		} else if opponentBoard&(1<<(63-i)) != 0 {
			opponentCount++
		}
	}

	fmt.Println()
	fmt.Println("⚫️ : ", playerCount)
	fmt.Println("⚪️ : ", opponentCount)
	if playerCount == opponentCount {
		fmt.Println("引き分け")
	} else if playerCount > opponentCount {
		fmt.Println("⚫️の勝ち")
	} else {
		fmt.Println("⚪️の勝ち")
	}
}
