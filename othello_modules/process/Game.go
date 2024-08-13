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
	g.ShowBoard()
	g.CountBoard()
	g.Status = -1
}

func (g *Game) SwapBoard() {
	tmp := g.board.PlayerBoard
	g.board.PlayerBoard = g.board.OpponentBoard
	g.board.OpponentBoard = tmp
}

func InitializeGame() Game {
	var game Game
	game.board = utils.InitializeBoard()
	game.player1 = utils.InitializePlayer(game.board.Id, true)
	game.Status = -1
	return game
}

func (g *Game) IsPass() bool {
	playerLegalBoard := MakeLegalBoard(g.board.PlayerBoard, g.board.OpponentBoard)
	opponentLegalBoard := MakeLegalBoard(g.board.OpponentBoard, g.board.PlayerBoard)
	return playerLegalBoard == 0x0000000000000000 && opponentLegalBoard != 0x0000000000000000
}

func (g *Game) IsEnd() bool {
	playerLegalBoard := MakeLegalBoard(g.board.PlayerBoard, g.board.OpponentBoard)
	opponentLegalBoard := MakeLegalBoard(g.board.OpponentBoard, g.board.PlayerBoard)
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

func ConvertPutToString(put uint64) string {
	allowedFirstString := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	allowedSecondString := []string{"1", "2", "3", "4", "5", "6", "7", "8"}
	var res string
	for i := 0; i < 64; i++ {
		if put&(1<<(63-i)) != 0 {
			res += allowedFirstString[i%8]
			res += allowedSecondString[i/8]
		}
	}
	return res
}
