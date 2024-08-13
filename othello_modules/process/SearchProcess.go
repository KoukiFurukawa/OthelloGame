package process

import (
	"main/utils"
)

const inf int = 1000000

var CELL_WEIGHT = [64]int{
	30, -12, 0, -1, -1, 0, -12, 30,
	-12, -15, -3, -3, -3, -3, -15, -12,
	0, -3, 0, -1, -1, 0, -3, 0,
	-1, -3, -1, -1, -1, -1, -3, -1,
	-1, -3, -1, -1, -1, -1, -3, -1,
	0, -3, 0, -1, -1, 0, -3, 0,
	-12, -15, -3, -3, -3, -3, -15, -12,
	30, -12, 0, -1, -1, 0, -12, 30,
}

var CELL_WEIGHT_FINAL = [64]int{
	2, 1, 1, 1, 1, 1, 1, 2,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	2, 1, 1, 1, 1, 1, 1, 2,
}

func Evaluation(playerBoard uint64, opponentBoard uint64) (int, int) {
	var playerEvaluation, opponentEvaluation int
	for i := 0; i < 64; i++ {
		if playerBoard&(1<<(63-i)) != 0 {
			playerEvaluation += CELL_WEIGHT[i]
		} else if opponentBoard&(1<<(63-i)) != 0 {
			opponentEvaluation += CELL_WEIGHT[i]
		}
	}
	return playerEvaluation, opponentEvaluation
}

func EvaluationFinal(playerBoard uint64, opponentBoard uint64) (int, int) {
	var playerEvaluation, opponentEvaluation int
	for i := 0; i < 64; i++ {
		if playerBoard&(1<<(63-i)) != 0 {
			playerEvaluation += CELL_WEIGHT_FINAL[i]
		} else if opponentBoard&(1<<(63-i)) != 0 {
			opponentEvaluation += CELL_WEIGHT_FINAL[i]
		}
	}
	return playerEvaluation, opponentEvaluation
}

func SwapBoard(playerBoard uint64, opponentBoard uint64) (uint64, uint64) {
	return opponentBoard, playerBoard
}

func Put(playerBoard uint64, opponentBoard uint64, put uint64) (uint64, uint64) {
	var rev uint64
	for k := 0; k < 8; k++ {
		var rev_ uint64
		mask := utils.Transfer(put, k)
		for (mask != 0) && ((mask & opponentBoard) != 0) {
			rev_ |= mask
			mask = utils.Transfer(mask, k)
		}
		if (mask & playerBoard) != 0 {
			rev |= rev_
		}
	}
	//反転する
	playerBoard ^= put | rev
	opponentBoard ^= rev
	return playerBoard, opponentBoard
}

func (g *Game) Search(depth uint8) uint64 {

	var NegaMax func(playerBoard uint64, opponentBoard uint64, depth uint8, passed bool) int
	NegaMax = func(playerBoard uint64, opponentBoard uint64, depth uint8, passed bool) int {

		// 葉に辿り着いたら評価 -----------------------------------------------------
		if depth == 0 {
			playerScore, _ := Evaluation(playerBoard, opponentBoard)
			return playerScore
		}

		maxScore := -inf

		// 葉ノードでなければ子ノードに対して再帰 -------------------------------------
		legalBoard := MakeLegalBoard(playerBoard, opponentBoard)
		mask := uint64(0x8000000000000000)
		var newPlayerBoard, newOpponentBoard, place uint64
		for i := 0; i < 64; i++ {
			place = mask >> i
			if legalBoard&place == place {
				newPlayerBoard, newOpponentBoard = Put(playerBoard, opponentBoard, place)
				maxScore = max(maxScore, -NegaMax(newPlayerBoard, newOpponentBoard, depth-1, false))
			}
		}

		// パスの処理 手番を交代して同じ深さで再帰する ----------------------------------------------------
		if maxScore == -inf {
			// 2回連続パスなら評価関数を実行
			if passed {
				playerScore, _ := Evaluation(playerBoard, opponentBoard)
				return playerScore
			}
			playerBoard, opponentBoard = SwapBoard(playerBoard, opponentBoard)
			return -NegaMax(playerBoard, opponentBoard, depth, true)
		}
		return maxScore
	}

	// 処理 ------------------------------------------------------------------------------
	var res uint64
	maxScore := -inf
	legalBoard := MakeLegalBoard(g.board.PlayerBoard, g.board.OpponentBoard)
	mask := uint64(0x8000000000000000)
	var newPlayerBoard, newOpponentBoard, place uint64
	for i := 0; i < 64; i++ {
		place = mask >> i
		if legalBoard&place == place {
			newPlayerBoard, newOpponentBoard = Put(g.board.PlayerBoard, g.board.OpponentBoard, place)
			score := -NegaMax(newPlayerBoard, newOpponentBoard, depth-1, false)
			if maxScore < score {
				maxScore = score
				res = place
			}
		}
	}
	return res
}

func (g *Game) NegaAlphaSearch(depth uint8) uint64 {

	var NegaAlpha func(playerBoard uint64, opponentBoard uint64, depth uint8, passed bool, alpha int, beta int) int
	NegaAlpha = func(playerBoard uint64, opponentBoard uint64, depth uint8, passed bool, alpha int, beta int) int {

		// 葉に辿り着いたら評価 -----------------------------------------------------
		if depth == 0 {
			playerScore, _ := Evaluation(playerBoard, opponentBoard)
			return playerScore
		}

		maxScore := -inf

		// 葉ノードでなければ子ノードに対して再帰 -------------------------------------
		legalBoard := MakeLegalBoard(playerBoard, opponentBoard)
		mask := uint64(0x8000000000000000)
		var newPlayerBoard, newOpponentBoard, place uint64
		var g int
		for i := 0; i < 64; i++ {
			place = mask >> i
			if legalBoard&place == place {
				newPlayerBoard, newOpponentBoard = Put(playerBoard, opponentBoard, place)
				g = -NegaAlpha(newPlayerBoard, newOpponentBoard, depth-1, false, -beta, -alpha)
				if g >= beta {
					return g
				}
				alpha = max(alpha, g)
				maxScore = max(maxScore, g)
			}
		}

		// パスの処理 手番を交代して同じ深さで再帰する ----------------------------------------------------
		if maxScore == -inf {
			// 2回連続パスなら評価関数を実行
			if passed {
				playerScore, _ := Evaluation(playerBoard, opponentBoard)
				return playerScore
			}
			playerBoard, opponentBoard = SwapBoard(playerBoard, opponentBoard)
			return -NegaAlpha(playerBoard, opponentBoard, depth, true, -beta, -alpha)
		}
		return maxScore
	}

	// 処理 ------------------------------------------------------------------------------
	var res uint64
	alpha := -inf
	beta := inf
	legalBoard := MakeLegalBoard(g.board.PlayerBoard, g.board.OpponentBoard)
	mask := uint64(0x8000000000000000)
	var newPlayerBoard, newOpponentBoard, place uint64
	for i := 0; i < 64; i++ {
		place = mask >> i
		if legalBoard&place == place {
			newPlayerBoard, newOpponentBoard = Put(g.board.PlayerBoard, g.board.OpponentBoard, place)
			score := -NegaAlpha(newPlayerBoard, newOpponentBoard, depth-1, false, -beta, -alpha)
			if alpha < score {
				alpha = score
				res = place
			}
		}
	}
	return res
}

func (g *Game) LastSearch(depth uint8) uint64 {

	var NegaMax func(playerBoard uint64, opponentBoard uint64, depth uint8, passed bool) int
	NegaMax = func(playerBoard uint64, opponentBoard uint64, depth uint8, passed bool) int {

		// 葉に辿り着いたら評価 -----------------------------------------------------
		if depth == 0 {
			playerScore, _ := EvaluationFinal(playerBoard, opponentBoard)
			return playerScore
		}

		maxScore := -inf

		// 葉ノードでなければ子ノードに対して再帰 -------------------------------------
		legalBoard := MakeLegalBoard(playerBoard, opponentBoard)
		mask := uint64(0x8000000000000000)
		var newPlayerBoard, newOpponentBoard, place uint64
		for i := 0; i < 64; i++ {
			place = mask >> i
			if legalBoard&place == place {
				newPlayerBoard, newOpponentBoard = Put(playerBoard, opponentBoard, place)
				maxScore = max(maxScore, -NegaMax(newPlayerBoard, newOpponentBoard, depth-1, false))
			}
		}

		// パスの処理 手番を交代して同じ深さで再帰する ----------------------------------------------------
		if maxScore == -inf {
			// 2回連続パスなら評価関数を実行
			if passed {
				playerScore, _ := EvaluationFinal(playerBoard, opponentBoard)
				return playerScore
			}
			playerBoard, opponentBoard = SwapBoard(playerBoard, opponentBoard)
			return -NegaMax(playerBoard, opponentBoard, depth, true)
		}
		return maxScore
	}

	// 処理 ------------------------------------------------------------------------------
	var res uint64
	maxScore := -inf
	legalBoard := MakeLegalBoard(g.board.PlayerBoard, g.board.OpponentBoard)
	mask := uint64(0x8000000000000000)
	var newPlayerBoard, newOpponentBoard, place uint64
	for i := 0; i < 64; i++ {
		place = mask >> i
		if legalBoard&place == place {
			newPlayerBoard, newOpponentBoard = Put(g.board.PlayerBoard, g.board.OpponentBoard, place)
			score := -NegaMax(newPlayerBoard, newOpponentBoard, depth-1, false)
			if maxScore < score {
				maxScore = score
				res = place
			}
		}
	}
	return res
}
