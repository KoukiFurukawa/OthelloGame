package main

import (
	"fmt"
	"main/process"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	var put uint64

	game := process.InitializeGame()
	game.StartGame()

	for game.Status != -1 {

		// 手番表示 ----------------------------------------
		hand := "⚫️"
		if game.Status == 1 {
			hand = "⚪️"
		}
		fmt.Println(hand, "の手番です。")
		game.ShowBoard()

		// 終局判定 --------------------------------------
		if game.IsEnd() {
			game.EndGame()
			break
		}

		// パス判定 ---------------------------------------
		if game.IsPass() {
			fmt.Println("置ける場所がありません。")
			game.SwapBoard()
			game.Status = 1 - game.Status
			continue
		}

		if game.Status == 1 {
			if game.IsFinalStage() {
				put = game.LastSearch(15)
			} else {
				put = game.Search(8)
			}
			fmt.Println(process.ConvertPutToString(put))

		} else if game.Status == 0 {

			fmt.Println("配置可能マス")
			game.CanPutList()

			// 入力受付 ---------------------------------------
			j, i, err := process.ScanUserInput()
			if err != nil {
				fmt.Println(err)
				continue
			}

			// 置く ------------------------------------------
			put = process.ConvertToBit(i, j)
			if !game.CanPut(put) {
				fmt.Println("can't put")
				continue
			}
		}
		game.Reverse(put)
		game.SwapBoard()
		game.ShowBoard()
	}

	// r.Run() // 0.0.0.0:8080 でサーバーを立てます。
}
