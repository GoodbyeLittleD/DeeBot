package army_of_the_hell

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestRandomPlay(t *testing.T) {
	game := New(2)
	game.Start()
	game.SetName(0, "玩家1")
	game.SetName(1, "玩家2")
	go func() {
		for {
			if game.GetScores() != nil {
				break
			}

			if game.WaitResponsePlayerId != -1 {
				response := "是"
				if rand.Intn(2) == 0 {
					response = "否"
				}
				game.GiveResponse(game.WaitResponsePlayerId, response)
			}

			time.Sleep(77 * time.Millisecond)
		}
	}()

	for i := 0; i < 10000; i++ {
		scores := game.GetScores()
		if scores != nil {
			fmt.Println(scores)
			break
		}

		var bidValue int
		bidValue = rand.Intn(128)
		if err := game.GivePrice(0, bidValue); err != nil {
			bidValue = rand.Intn(11)
			if err := game.GivePrice(0, bidValue); err != nil {
				fmt.Println(err)
			}
		}
		time.Sleep(1 * time.Microsecond)

		if game.GetScores() != nil {
			break
		}

		bidValue = rand.Intn(128)
		if err := game.GivePrice(1, bidValue); err != nil {
			bidValue = rand.Intn(11)
			if err := game.GivePrice(1, bidValue); err != nil {
				fmt.Println(err)
			}
		}
		time.Sleep(1 * time.Microsecond)

		if game.GetScores() != nil {
			break
		}
	}
}

func TestGame(t *testing.T) {
	game := New(2)
	game.Start()
	go func() {
		if game.WaitResponsePlayerId != -1 {
			fmt.Println("输入回应：")
			fmt.Scanf("%s\n", &game.Response)
			game.GiveResponse(game.WaitResponsePlayerId, game.Response)
		}
	}()

	for {
		scores := game.GetScores()
		if scores != nil {
			fmt.Println(scores)
			break
		}

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
	}
}
