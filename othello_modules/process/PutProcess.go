package process

import "fmt"

func ConvertToBit(i int, j int) uint64 {
	mask := uint64(0x8000000000000000)
	mask >>= j
	mask >>= i * 8
	return mask
}

func (g *Game) CanPut(put uint64) bool {
	legalBoard := MakeLegalBoard(g.board.PlayerBoard, g.board.OpponentBoard)
	return (put&legalBoard == put)
}

func (g *Game) Reverse(put uint64) {
	g.board.PlacingFrame(put)
	g.Status = 1 - g.Status
	g.board.Total += 1
}

func (g *Game) CanPutList() {
	var place uint64
	var placeString string
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			place = ConvertToBit(i, j)
			if g.CanPut(place) {
				placeString = ConvertPutToString(place)
				fmt.Printf("%s, ", placeString)
			}
		}
	}
	fmt.Println()
}

func MakeLegalBoard(PlayerBoard uint64, OpponentBoard uint64) uint64 {
	horizontalWatchBoard := OpponentBoard & 0x7e7e7e7e7e7e7e7e // 左右
	verticalWatchBoard := OpponentBoard & 0x00FFFFFFFFFFFF00   // 上下
	allSideWatchBoard := OpponentBoard & 0x007e7e7e7e7e7e00    // 全辺

	blankBoard := ^(PlayerBoard | OpponentBoard) // 空きマスのみ
	var tmp uint64                               // 隣に相手の色があるか一時保存
	var legalBoard uint64                        // 返り値

	// 8方向チェック
	// 左
	tmp = horizontalWatchBoard & (PlayerBoard << 1)
	tmp |= horizontalWatchBoard & (tmp << 1)
	tmp |= horizontalWatchBoard & (tmp << 1)
	tmp |= horizontalWatchBoard & (tmp << 1)
	tmp |= horizontalWatchBoard & (tmp << 1)
	tmp |= horizontalWatchBoard & (tmp << 1)
	legalBoard = blankBoard & (tmp << 1)

	// 右
	tmp = horizontalWatchBoard & (PlayerBoard >> 1)
	tmp |= horizontalWatchBoard & (tmp >> 1)
	tmp |= horizontalWatchBoard & (tmp >> 1)
	tmp |= horizontalWatchBoard & (tmp >> 1)
	tmp |= horizontalWatchBoard & (tmp >> 1)
	tmp |= horizontalWatchBoard & (tmp >> 1)
	legalBoard |= blankBoard & (tmp >> 1)

	// 上
	tmp = verticalWatchBoard & (PlayerBoard << 8)
	tmp |= verticalWatchBoard & (tmp << 8)
	tmp |= verticalWatchBoard & (tmp << 8)
	tmp |= verticalWatchBoard & (tmp << 8)
	tmp |= verticalWatchBoard & (tmp << 8)
	tmp |= verticalWatchBoard & (tmp << 8)
	legalBoard |= blankBoard & (tmp << 8)

	// 下
	tmp = verticalWatchBoard & (PlayerBoard >> 8)
	tmp |= verticalWatchBoard & (tmp >> 8)
	tmp |= verticalWatchBoard & (tmp >> 8)
	tmp |= verticalWatchBoard & (tmp >> 8)
	tmp |= verticalWatchBoard & (tmp >> 8)
	tmp |= verticalWatchBoard & (tmp >> 8)
	legalBoard |= blankBoard & (tmp >> 8)

	// 右斜め上
	tmp = allSideWatchBoard & (PlayerBoard << 7)
	tmp |= allSideWatchBoard & (tmp << 7)
	tmp |= allSideWatchBoard & (tmp << 7)
	tmp |= allSideWatchBoard & (tmp << 7)
	tmp |= allSideWatchBoard & (tmp << 7)
	tmp |= allSideWatchBoard & (tmp << 7)
	legalBoard |= blankBoard & (tmp << 7)

	// 左斜め上
	tmp = allSideWatchBoard & (PlayerBoard << 9)
	tmp |= allSideWatchBoard & (tmp << 9)
	tmp |= allSideWatchBoard & (tmp << 9)
	tmp |= allSideWatchBoard & (tmp << 9)
	tmp |= allSideWatchBoard & (tmp << 9)
	tmp |= allSideWatchBoard & (tmp << 9)
	legalBoard |= blankBoard & (tmp << 9)

	// 右斜め下
	tmp = allSideWatchBoard & (PlayerBoard >> 9)
	tmp |= allSideWatchBoard & (tmp >> 9)
	tmp |= allSideWatchBoard & (tmp >> 9)
	tmp |= allSideWatchBoard & (tmp >> 9)
	tmp |= allSideWatchBoard & (tmp >> 9)
	tmp |= allSideWatchBoard & (tmp >> 9)
	legalBoard |= blankBoard & (tmp >> 9)

	// 左斜め下
	tmp = allSideWatchBoard & (PlayerBoard >> 7)
	tmp |= allSideWatchBoard & (tmp >> 7)
	tmp |= allSideWatchBoard & (tmp >> 7)
	tmp |= allSideWatchBoard & (tmp >> 7)
	tmp |= allSideWatchBoard & (tmp >> 7)
	tmp |= allSideWatchBoard & (tmp >> 7)
	legalBoard |= blankBoard & (tmp >> 7)

	return legalBoard
}
