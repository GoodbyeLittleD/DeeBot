package army_of_the_hell

import (
	"fmt"
	"time"
)

func PlayConsole() {
	game := New(2)
	game.Start()
	go func() {
		for {
			if game.WaitResponsePlayerId != -1 {
				fmt.Println("输入回应：")
				fmt.Scanf("%s\n", &game.Response)
				game.GiveResponse(game.WaitResponsePlayerId, game.Response)
			}
			time.Sleep(77 * time.Millisecond)
		}
	}()

	for {
		scores := game.GetScores()
		if scores != nil {
			fmt.Println(scores)
			break
		}

		if game.SingleMode {
			var bidValue int
			fmt.Print("输入玩家1的出价：")
			fmt.Scanf("%d\n", &bidValue)
			if err := game.GivePrice(0, bidValue); err != nil {
				fmt.Println(err)
			}

			fmt.Print("输入玩家2的出价：")
			fmt.Scanf("%d\n", &bidValue)
			if err := game.GivePrice(1, bidValue); err != nil {
				fmt.Println(err)
			}
		} else {
			var bidValue1, bidValue2 int
			fmt.Print("输入玩家1的出价：")
			fmt.Scanf("%d %d\n", &bidValue1, &bidValue2)
			if err := game.GivePrices(0, bidValue1, bidValue2); err != nil {
				fmt.Println(err)
			}
			fmt.Print("输入玩家2的出价：")
			fmt.Scanf("%d %d\n", &bidValue1, &bidValue2)
			if err := game.GivePrices(1, bidValue1, bidValue2); err != nil {
				fmt.Println(err)
			}

		}
	}
}
