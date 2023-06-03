package army_of_the_hell

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"golang.org/x/exp/slices"
)

var WINSTATE_BOSS_COUNT = 3

type Game struct {
	PlayerNum int
	Turn      int

	PlayerCredits  []int
	PlayerEntities [][]*Entity
	PlayerLeaved   []bool

	CurrentEntityPool         []*Entity
	CurrentBiddingEntity      *Entity // current bidding entity is not removed from pool
	CurrentBiddingEntityIndex int
	CurrentPlayerReady        []bool
	CurrentPlayerBid          []int
	CurrentPlayerNickname     []string

	ProcessingLock sync.Mutex

	洞穴重生的邪恶之犬_triggered bool
	海法斯特盔甲制作者_triggered bool
	黑暗长老_天黑_countdown   int8
	牙皮_triggered        bool
	火之眼_countdown       int8
	召唤者_countdown       int8
	古代无魂之卡_countdown    int8
	西希之王_playerid       int8
	优先抽取议会成员_triggered  bool
	衣卒尔_playerid        int8
	衣卒尔_player_history  []bool
	督军山克_countdown      int8
	破坏者卡兰索_triggered    bool
	诅咒的阿克姆_triggered    bool
	血腥的巴特克_triggered    bool
	不洁的凡塔_triggered     bool
	古难记录者_triggered     bool
	达克法恩_triggered      bool
	火之眼_unlocked        bool
	墨菲斯托_unlocked       bool
	三个野蛮人试炼_passed      []bool
	破坏者卡兰索_unlocked     bool

	SingleMode bool
	DoubleMode bool

	WaitResponsePlayerId int
	Response             string

	PrintFunc func(string)
}

func New(number_players int) *Game {
	game := new(Game)
	game.PlayerNum = number_players

	game.PlayerCredits = nil
	for i := 0; i < number_players; i++ {
		game.PlayerCredits = append(game.PlayerCredits, 100)
	}
	game.PlayerEntities = make([][]*Entity, number_players)
	game.PlayerLeaved = make([]bool, number_players)

	game.CurrentEntityPool = make([]*Entity, 0)
	game.CurrentBiddingEntity = nil
	game.CurrentPlayerReady = make([]bool, number_players)
	game.CurrentPlayerBid = make([]int, number_players)
	game.CurrentPlayerNickname = make([]string, number_players)
	game.WaitResponsePlayerId = -1

	game.西希之王_playerid = -1
	game.衣卒尔_playerid = -1
	game.三个野蛮人试炼_passed = make([]bool, number_players)

	game.SingleMode = true

	game.PrintFunc = func(msg string) {
		fmt.Println(msg)
	}
	return game
}

func (game *Game) Start() {
	game.CurrentEntityPool = []*Entity{
		&尸体发火,
		// &剥壳凹槽,
	}

	if game.PlayerNum > 2 {
		WINSTATE_BOSS_COUNT = 2
	} else {
		WINSTATE_BOSS_COUNT = 3
	}
	game.PrintFunc(fmt.Sprintf("游戏开始，当前游戏人数：%d，获胜所需BOSS数量：%d", game.PlayerNum, WINSTATE_BOSS_COUNT))

	game.startNewTurn()
}

func (game *Game) SetName(playerId int, name string) {
	game.CurrentPlayerNickname[playerId] = name
}

func (game *Game) GetScores() []int {
	if game_ended, winner := game.checkWinState(); !game_ended {
		return nil
	} else {
		if winner >= 0 {
			game.PrintFunc(fmt.Sprintf("游戏结束，玩家 %s 获胜！", game.CurrentPlayerNickname[winner]))
		}
		scores := make([]int, game.PlayerNum)
		for i := 0; i < game.PlayerNum; i++ {
			if i == winner {
				scores[i] = 1
			} else {
				scores[i] = 0
			}
		}
		return scores
	}
}

func (game *Game) GivePrice(playerId int, price int) error {
	game.ProcessingLock.Lock()
	defer game.ProcessingLock.Unlock()

	if game.CurrentPlayerReady[playerId] {
		return errors.New("本回合已经出价了")
	}
	if err := game.checkPriceValid(playerId, price); err != nil {
		return err
	}
	game.CurrentPlayerBid[playerId] = price
	game.CurrentPlayerReady[playerId] = true

	fmt.Println(game.CurrentPlayerReady)

	all_player_ready := true
	for _, ready := range game.CurrentPlayerReady {
		if !ready {
			all_player_ready = false
		}
	}
	if all_player_ready {
		game.endTurn()

		if win, _ := game.checkWinState(); win {
			return nil
		}

		game.startNewTurn()
	}
	return nil
}

func (game *Game) GiveResponse(playerId int, response string) {
	if game.WaitResponsePlayerId != playerId {
		return
	}
	game.Response = response
	game.WaitResponsePlayerId = -1
}

func (game *Game) AcceptTrial(playerId int) error {
	if game.三个野蛮人试炼_passed[playerId] {
		return errors.New("你已经通过了试炼")
	}
	if game.PlayerCredits[playerId] < 50 {
		return errors.New("能力点不足，你需要50点能力点")
	}
	game.PlayerCredits[playerId] -= 50
	game.三个野蛮人试炼_passed[playerId] = true
	if !game.破坏者卡兰索_unlocked {
		game.破坏者卡兰索_unlocked = true
		game.CurrentEntityPool = append(game.CurrentEntityPool, &破坏者卡兰索)
	}
	game.PrintFunc(game.CurrentPlayerNickname[playerId] + "支付了50能力点，通过了三个野蛮人试炼")
	return nil
}

func (game *Game) PlayerLeave(playerId int) bool {
	game.PlayerLeaved[playerId] = true
	game.CurrentPlayerReady[playerId] = true
	game.CurrentPlayerBid[playerId] = 0
	game.PrintFunc(game.CurrentPlayerNickname[playerId] + "离开了游戏")

	all_player_leaved := true
	for _, leaved := range game.PlayerLeaved {
		if !leaved {
			all_player_leaved = false
			break
		}
	}
	if all_player_leaved {
		game.PrintFunc("所有玩家都离开了游戏，游戏结束。")
		return true
	}

	if game.WaitResponsePlayerId == playerId {
		game.Response = "否"
		game.WaitResponsePlayerId = -1
	}
	return false
}

func (game *Game) startNewTurn() {
	game.Turn++
	var poolString string
	for _, entity := range game.CurrentEntityPool {
		poolString += entity.Name + "\n"
	}
	game.PrintFunc(fmt.Sprintf("第%d回合开始。本回合抽取池：\n%s", game.Turn, poolString))

	if len(game.CurrentEntityPool) == 0 {
		game.PrintFunc("拍卖池中已经无角色，游戏结束。")
	} else {
		game.CurrentBiddingEntityIndex = rand.Int() % len(game.CurrentEntityPool)
		game.CurrentBiddingEntity = game.CurrentEntityPool[game.CurrentBiddingEntityIndex]

		skipNormalPrintFunc := false

		// 古代无魂之卡：触发后第6回合改为拍卖【督瑞尔】
		if game.古代无魂之卡_countdown > 0 {
			game.古代无魂之卡_countdown--
			if game.古代无魂之卡_countdown == 0 {
				game.CurrentEntityPool = append(game.CurrentEntityPool, &督瑞尔)
				game.CurrentBiddingEntity = &督瑞尔
				game.CurrentBiddingEntityIndex = len(game.CurrentEntityPool) - 1
				game.PrintFunc("【古代无魂之卡】效果触发，本回合拍卖的角色为【督瑞尔】。")
			}

			// 达克法恩：跳过下次抽牌，改为抽到督军山克。
		} else if game.达克法恩_triggered {
			game.达克法恩_triggered = false
			if slices.Contains(game.CurrentEntityPool, &督军山克) {
				game.CurrentBiddingEntityIndex = slices.Index(game.CurrentEntityPool, &督军山克)
				game.CurrentBiddingEntity = &督军山克
			}

			// 督军山克：接下来第三次有效抽牌抽到矫正者怪异，已经抽到则失效。
		} else if game.督军山克_countdown > 0 {
			game.督军山克_countdown--
			if game.督军山克_countdown == 0 {
				if slices.Contains(game.CurrentEntityPool, &矫正者怪异) {
					game.PrintFunc("督军山克效果已触发。")
					game.CurrentBiddingEntityIndex = slices.Index(game.CurrentEntityPool, &矫正者怪异)
					game.CurrentBiddingEntity = &矫正者怪异
				}
			}

			// 议会成员：优先抽取
		} else if game.优先抽取议会成员_triggered {
			game.优先抽取议会成员_triggered = false
			for i := 0; i < len(game.CurrentEntityPool); i++ {
				if game.CurrentEntityPool[i].Name == "邪恶之手伊斯梅尔" ||
					game.CurrentEntityPool[i].Name == "火焰之指吉列布" ||
					game.CurrentEntityPool[i].Name == "冰拳托克" {
					game.PrintFunc("优先抽取议会成员效果已触发。")
					game.CurrentBiddingEntityIndex = i
					game.CurrentBiddingEntity = game.CurrentEntityPool[i]
					break
				}
			}

			// 破坏者卡兰索：优先抽取
		} else if game.破坏者卡兰索_triggered {
			game.破坏者卡兰索_triggered = false
			for i := 0; i < len(game.CurrentEntityPool); i++ {
				if game.CurrentEntityPool[i].Name == "诅咒的阿克姆" {
					game.PrintFunc("优先抽取诅咒的阿克姆效果已触发。")
					game.CurrentBiddingEntityIndex = i
					game.CurrentBiddingEntity = game.CurrentEntityPool[i]
					break
				}
			}

			// 诅咒的阿克姆：优先抽取
		} else if game.诅咒的阿克姆_triggered {
			game.诅咒的阿克姆_triggered = false
			for i := 0; i < len(game.CurrentEntityPool); i++ {
				if game.CurrentEntityPool[i].Name == "血腥的巴特克" {
					game.PrintFunc("优先抽取血腥的巴特克效果已触发。")
					game.CurrentBiddingEntityIndex = i
					game.CurrentBiddingEntity = game.CurrentEntityPool[i]
					break
				}
			}

			// 血腥的巴特克：优先抽取
		} else if game.血腥的巴特克_triggered {
			game.血腥的巴特克_triggered = false
			for i := 0; i < len(game.CurrentEntityPool); i++ {
				if game.CurrentEntityPool[i].Name == "不洁的凡塔" {
					game.PrintFunc("优先抽取不洁的凡塔效果已触发。")
					game.CurrentBiddingEntityIndex = i
					game.CurrentBiddingEntity = game.CurrentEntityPool[i]
					break
				}
			}

			// 不洁的凡塔：优先抽取
		} else if game.不洁的凡塔_triggered {
			game.不洁的凡塔_triggered = false
			for i := 0; i < len(game.CurrentEntityPool); i++ {
				if game.CurrentEntityPool[i].Name == "古难记录者" {
					game.PrintFunc("优先抽取古难记录者效果已触发。")
					game.CurrentBiddingEntityIndex = i
					game.CurrentBiddingEntity = game.CurrentEntityPool[i]
					break
				}
			}

			// 古难记录者：优先抽取
		} else if game.古难记录者_triggered {
			game.古难记录者_triggered = false
			for i := 0; i < len(game.CurrentEntityPool); i++ {
				if game.CurrentEntityPool[i].Name == "巴尔" {
					game.PrintFunc("优先抽取巴尔效果已触发。")
					game.CurrentBiddingEntityIndex = i
					game.CurrentBiddingEntity = game.CurrentEntityPool[i]
					break
				}
			}

			// 洞穴重生的邪恶之犬：当第一次被抽到时，跳过本轮所有竞拍
		} else if !game.洞穴重生的邪恶之犬_triggered && game.CurrentBiddingEntity.Name == "洞穴重生的邪恶之犬" {
			game.洞穴重生的邪恶之犬_triggered = true
			game.PrintFunc("洞穴重生的邪恶之犬被抽到，跳过本轮所有竞拍。")
			game.handleEndTurnCredits()
			game.startNewTurn()
			return

			// 三个野蛮人试炼：当被抽到时，改为开启试炼
		} else if game.CurrentBiddingEntity.Name == "三个野蛮人" {
			game.PrintFunc("【三个野蛮人】试炼已开启。玩家可以随时输入【接受试炼】，支付50能力点通过三个野蛮人的试炼。一方通过试炼后解锁【破坏者卡兰索】。")
			game.CurrentEntityPool = append(game.CurrentEntityPool[:game.CurrentBiddingEntityIndex], game.CurrentEntityPool[game.CurrentBiddingEntityIndex+1:]...)
			game.handleEndTurnCredits()
			game.startNewTurn()
			return

		} else {
			// 海法斯特盔甲制作者：第一次被抽到时，改为抽到ACT5-达克法恩
			if !game.海法斯特盔甲制作者_triggered && game.CurrentBiddingEntity.Name == "海法斯特盔甲制作者" {
				game.海法斯特盔甲制作者_triggered = true
				game.CurrentEntityPool = append(game.CurrentEntityPool, &达克法恩)
				game.CurrentBiddingEntity = &达克法恩
				game.CurrentBiddingEntityIndex = len(game.CurrentEntityPool) - 1
				game.PrintFunc("【海法斯特盔甲制作者】效果触发，本回合拍卖的角色改为【ACT5-达克法恩】。")
			}

			// 黑暗长老：抽到后三回合【天黑】；牙皮：永久取消【天黑】
			if game.CurrentBiddingEntity.Name == "牙皮" {
				game.牙皮_triggered = true
			}
			if game.黑暗长老_天黑_countdown > 0 && !game.牙皮_triggered {
				game.黑暗长老_天黑_countdown--
				game.PrintFunc("由于黑暗长老效果，天黑了，本回合拍卖的角色未知。\n\n请双方私聊出价。")
				skipNormalPrintFunc = true
			}
			if game.CurrentBiddingEntity.Name == "黑暗长老" {
				game.黑暗长老_天黑_countdown = 3
			}

			// 格瑞斯华尔德：只有拉卡尼休和树头木拳的拥有者可以竞拍
			if game.CurrentBiddingEntity.Name == "格瑞斯华尔德" {
				for i := 0; i < game.PlayerNum; i++ {
					if !slices.Contains(game.PlayerEntities[i], &拉卡尼休) || !slices.Contains(game.PlayerEntities[i], &树头木拳) {
						game.CurrentPlayerReady[i] = true
						game.CurrentPlayerBid[i] = 0
					}
				}
				game.PrintFunc(fmt.Sprintf("本回合拍卖的角色：%v\n\n只有【拉卡尼休】和【树头木拳】的拥有者可以竞拍，请私聊出价。", game.CurrentBiddingEntity))
				skipNormalPrintFunc = true
			}

			// 沙漠三小队：【钻地的冰虫】、【牙皮】、【疯狂血腥女巫】在一家，则只有他可以竞拍【火之眼】、【督瑞尔】
			沙漠三小队_playerid := -1
			if game.CurrentBiddingEntity.Name == "火之眼" || game.CurrentBiddingEntity.Name == "督瑞尔" {
				for i := 0; i < game.PlayerNum; i++ {
					if slices.Contains(game.PlayerEntities[i], &疯狂血腥女巫) && slices.Contains(game.PlayerEntities[i], &钻地的冰虫) && slices.Contains(game.PlayerEntities[i], &牙皮) {
						沙漠三小队_playerid = i
						break
					}
				}
			}
			if 沙漠三小队_playerid != -1 {
				for i := 0; i < game.PlayerNum; i++ {
					if i != 沙漠三小队_playerid {
						game.CurrentPlayerReady[i] = true
						game.CurrentPlayerBid[i] = 0
					}
				}
				game.PrintFunc(fmt.Sprintf("本回合拍卖的角色：%v\n\n由于【钻地的冰虫】、【牙皮】、【疯狂血腥女巫】效果，只有 %s 可以竞拍【%s】，请私聊出价。", game.CurrentBiddingEntity, game.CurrentPlayerNickname[沙漠三小队_playerid], game.CurrentBiddingEntity.Name))
				skipNormalPrintFunc = true
			} else {
				// 火之眼：如果两回合内抽到【召唤者】，只有拍得火之眼者可以竞拍【召唤者】
				if game.火之眼_countdown > 0 {
					game.火之眼_countdown--
					if game.CurrentBiddingEntity.Name == "召唤者" {
						for i := 0; i < game.PlayerNum; i++ {
							if !slices.Contains(game.PlayerEntities[i], &火之眼) {
								game.CurrentPlayerReady[i] = true
								game.CurrentPlayerBid[i] = 0
							}
						}
						game.PrintFunc(fmt.Sprintf("本回合拍卖的角色：%v\n\n由于【火之眼】效果，只有【火之眼】拥有者可以竞拍【召唤者】，请私聊出价。", game.CurrentBiddingEntity))
						skipNormalPrintFunc = true
					}
				}

				// 召唤者：如果两回合内抽到【督瑞尔】，只有拍得者可以竞拍【督瑞尔】
				if game.召唤者_countdown > 0 {
					game.召唤者_countdown--
					if game.CurrentBiddingEntity.Name == "督瑞尔" {
						for i := 0; i < game.PlayerNum; i++ {
							if !slices.Contains(game.PlayerEntities[i], &召唤者) {
								game.CurrentPlayerReady[i] = true
								game.CurrentPlayerBid[i] = 0
							}
						}
						game.PrintFunc(fmt.Sprintf("本回合拍卖的角色：%v\n\n由于【召唤者】效果，只有【召唤者】拥有者可以竞拍【督瑞尔】，请私聊出价。", game.CurrentBiddingEntity))
						skipNormalPrintFunc = true
					}
				}
			}
		}

		fmt.Println(game.CurrentPlayerReady)

		all_player_ready := true
		for id, ready := range game.CurrentPlayerReady {
			if game.PlayerLeaved[id] {
				game.CurrentPlayerReady[id] = true
			} else if !ready {
				all_player_ready = false
			}
		}
		if all_player_ready {
			game.endTurn()

			if win, _ := game.checkWinState(); win {
				return
			}

			game.startNewTurn()
			return
		}

		if !skipNormalPrintFunc {
			bidRanges := ""
			for i := 0; i < game.PlayerNum; i++ {
				bidRanges += fmt.Sprintf("%s 最高出价：%d (%+d%%)\n", game.CurrentPlayerNickname[i], game.getMaxBidValue(i), int(game.getCurrentPriceDiscount(i)*100-100.01))
			}

			game.PrintFunc(fmt.Sprintf("本回合拍卖的角色：%v\n\n%s\n请私聊出价。", game.CurrentBiddingEntity, bidRanges))
		}
	}
}

func (game *Game) handleEndTurnCredits() {
	endTurnMsg := ""
	for playerId, entities := range game.PlayerEntities {
		totalIncrement := 0
		extra := 0
		for _, entity := range entities {
			game.PlayerCredits[playerId] += entity.EndTurnCredits
			totalIncrement += entity.EndTurnCredits

			if entity.Name == "血鸟" && extra < 2 {
				extra = 2
			} else if entity.Name == "罗达门特" && extra < 3 {
				extra = 3
			} else if entity.Name == "吉得宾偷窃者" && extra < 4 {
				extra = 4
			} else if entity.Name == "督军山克" && extra < 5 {
				extra = 5
			}
		}
		game.PlayerCredits[playerId] += extra
		totalIncrement += extra

		var entityList string
		for _, entity := range entities {
			entityList += entity.Name + " "
		}
		endTurnMsg += fmt.Sprintf("%s 当前点数：%d(+%d)\n拥有：%s\n\n", game.CurrentPlayerNickname[playerId], game.PlayerCredits[playerId], totalIncrement, entityList)
	}
	game.PrintFunc(endTurnMsg)
}

func (game *Game) endTurn() {
	for i := range game.CurrentPlayerReady {
		game.CurrentPlayerReady[i] = false
	}
	winner := -1
	largest := -1
	for id, value := range game.CurrentPlayerBid {
		if largest < value {
			largest = value
			winner = id
		} else if largest == value {
			winner = -1
		}
	}

	bidDescription := ""
	for i := 0; i < game.PlayerNum; i++ {
		bidDescription += fmt.Sprintf("\n%s ：%d", game.CurrentPlayerNickname[i], game.CurrentPlayerBid[i])
	}

	if winner == -1 { // draw
		bidDescription += "\n\n最高出价相同，本回合流拍。"
	} else {
		bidDescription += fmt.Sprintf("\n\n%s 成功拍得了：%v", game.CurrentPlayerNickname[winner], game.CurrentBiddingEntity)
		game.PlayerCredits[winner] -= int(float64(game.CurrentPlayerBid[winner]) * game.getCurrentPriceDiscount(winner))
		game.PlayerEntities[winner] = append(game.PlayerEntities[winner], game.CurrentBiddingEntity)
		game.CurrentEntityPool = append(game.CurrentEntityPool[:game.CurrentBiddingEntityIndex], game.CurrentEntityPool[game.CurrentBiddingEntityIndex+1:]...)
		for _, entity := range game.CurrentBiddingEntity.UnlockEntities {
			if entity.UnlockChecker != nil {
				if !entity.UnlockChecker(game) {
					continue
				}
			}
			if entity.Name == "火之眼" {
				game.火之眼_unlocked = true
			} else if entity.Name == "墨菲斯托" {
				game.墨菲斯托_unlocked = true
			}
			game.CurrentEntityPool = append(game.CurrentEntityPool, entity)
		}
	}
	game.PrintFunc(fmt.Sprintf("第%d回合结束。\n%s", game.Turn, bidDescription))

	// 火之眼：如果两回合内抽到【召唤者】，只有拍得火之眼者可以竞拍【召唤者】
	if game.CurrentBiddingEntity.Name == "火之眼" && winner != -1 {
		game.火之眼_countdown = 2
	}

	// 召唤者：如果两回合内抽到【督瑞尔】，只有拍得者可以竞拍【督瑞尔】
	if game.CurrentBiddingEntity.Name == "召唤者" && winner != -1 {
		game.召唤者_countdown = 2
	}

	// 古代无魂之卡：拍到一方可选择移除督瑞尔，六回合后跳过抽牌，竞拍督瑞尔。【注：即使督瑞尔已被拍得，此特殊效果依然可以启动。】
	if game.CurrentBiddingEntity.Name == "古代无魂之卡" && winner != -1 {
		game.PrintFunc("【古代无魂之卡】效果触发，请拍得的玩家私聊选择是否发动效果（是/否）。")
	wait_response:
		game.WaitResponsePlayerId = winner
		for {
			if game.WaitResponsePlayerId == -1 {
				break
			}
			time.Sleep(1 * time.Second)
		}
		if game.Response == "是" {
			game.PrintFunc("【古代无魂之卡】效果已发动，6回合后将拍卖督瑞尔。")
			game.古代无魂之卡_countdown = 6
			for i := 0; i < game.PlayerNum; i++ {
				if slices.Contains(game.PlayerEntities[i], &督瑞尔) {
					index := slices.Index(game.PlayerEntities[i], &督瑞尔)
					game.PlayerEntities[i] = append(game.PlayerEntities[i][:index], game.PlayerEntities[i][index+1:]...)
				}
			}
			if slices.Contains(game.CurrentEntityPool, &督瑞尔) {
				index := slices.Index(game.CurrentEntityPool, &督瑞尔)
				game.CurrentEntityPool = append(game.CurrentEntityPool[:index], game.CurrentEntityPool[index+1:]...)
			}
		} else if game.Response == "否" {
			game.PrintFunc("【古代无魂之卡】效果未发动。")
		} else {
			game.PrintFunc("请输入【是】或【否】。")
			goto wait_response
		}
	}

	// 达克法恩：拍得者在拍得时**可以**选择跳过下次抽牌，改为抽到督军山克。
	if game.CurrentBiddingEntity.Name == "达克法恩" && winner != -1 {
		game.PrintFunc("【达克法恩】可选效果：跳过下次抽牌，改为抽到督军山克。\n\n请拍得的玩家私聊选择是否发动（是/否）。")
	wait_response2:
		game.WaitResponsePlayerId = winner
		for {
			if game.WaitResponsePlayerId == -1 {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
		if game.Response == "是" {
			game.PrintFunc("【达克法恩】效果已发动。")
			game.达克法恩_triggered = true
		} else if game.Response == "否" {
			game.PrintFunc("【达克法恩】效果未发动。")
		} else {
			game.PrintFunc("请输入【是】或【否】。")
			goto wait_response2
		}
	}

	// 西希之王：下回合不能出价boss以外的怪物
	if game.CurrentBiddingEntity.Name == "西希之王" && winner != -1 {
		game.西希之王_playerid = int8(winner)
	} else if game.西希之王_playerid != -1 {
		game.西希之王_playerid = -1
	}

	// 邪恶之手伊斯梅尔，火焰之指吉列布，冰拳托克：优先抽取议会成员
	if game.CurrentBiddingEntity.Name == "邪恶之手伊斯梅尔" || game.CurrentBiddingEntity.Name == "火焰之指吉列布" || game.CurrentBiddingEntity.Name == "冰拳托克" {
		game.优先抽取议会成员_triggered = true
	}

	// 衣卒尔
	if game.CurrentBiddingEntity.Name == "衣卒尔" && winner != -1 {
		game.衣卒尔_playerid = int8(winner)
		game.衣卒尔_player_history = nil
	} else if game.衣卒尔_playerid != -1 {
		if len(game.衣卒尔_player_history) != 4 {
			game.衣卒尔_player_history = append(game.衣卒尔_player_history, game.CurrentPlayerBid[game.衣卒尔_playerid] == 0)
		} else {
			game.衣卒尔_player_history = nil
		}
	}

	// 督军山克
	if game.CurrentBiddingEntity.Name == "督军山克" && winner != -1 {
		game.督军山克_countdown = 3
	}

	// 破坏者卡兰索, 诅咒的阿克姆, 血腥的巴特克, 不洁的凡塔, 古难记录者
	if game.CurrentBiddingEntity.Name == "破坏者卡兰索" {
		game.破坏者卡兰索_triggered = true
	}
	if game.CurrentBiddingEntity.Name == "诅咒的阿克姆" {
		game.诅咒的阿克姆_triggered = true
	}
	if game.CurrentBiddingEntity.Name == "血腥的巴特克" {
		game.血腥的巴特克_triggered = true
	}
	if game.CurrentBiddingEntity.Name == "不洁的凡塔" {
		game.不洁的凡塔_triggered = true
	}
	if game.CurrentBiddingEntity.Name == "古难记录者" {
		game.古难记录者_triggered = true
	}

	// 尼拉塞克：直接通过三个野蛮人试炼
	if game.CurrentBiddingEntity.Name == "尼拉塞克" && winner != -1 {
		game.三个野蛮人试炼_passed[winner] = true
		if !game.破坏者卡兰索_unlocked {
			game.破坏者卡兰索_unlocked = true
			game.CurrentEntityPool = append(game.CurrentEntityPool, &破坏者卡兰索)
		}
	}

	// 双拍模式
	if game.SingleMode && slices.Contains(game.CurrentBiddingEntity.Tags, "BOSS") && winner != -1 {
		game.SingleMode = false
		game.DoubleMode = true
		game.PrintFunc("由于有BOSS被拍得，【双拍模式】已开启。")
	}

	game.CurrentBiddingEntity = nil
	game.CurrentBiddingEntityIndex = -1

	game.handleEndTurnCredits()
}

func (game *Game) checkWinState() (game_ended bool, winner int) {
	// if any player has got 3 entities tagged "BOSS", he wins.
	for i := 0; i < game.PlayerNum; i++ {
		count_boss := 0
		for _, entity := range game.PlayerEntities[i] {
			if entity.IsBoss() {
				count_boss++
			}
		}
		if count_boss >= WINSTATE_BOSS_COUNT {
			game_ended = true
			winner = i
			return
		}
	}
	game_ended = false
	if len(game.CurrentEntityPool) == 0 {
		game_ended = true
	}
	winner = -1
	return
}

func (game *Game) getCurrentPriceDiscount(playerId int) float64 {
	bid_effect := 1.0

	for _, entity := range game.PlayerEntities[playerId] {
		if entity.BidEffected(game.CurrentBiddingEntity.Tags) {
			bid_effect += entity.BidAffectValue
		}

		// 女伯爵：对安达利尔出价-15%（在boss的基础上额外-10%）
		if entity.Name == "女伯爵" && game.CurrentBiddingEntity.Name == "安达利尔" {
			bid_effect -= 0.1
		}
	}

	if bid_effect < 0 {
		return 0
	}

	return bid_effect
}

func (game *Game) getMaxBidValue(playerId int) int {
	if game.西希之王_playerid == int8(playerId) && !slices.Contains(game.CurrentBiddingEntity.Tags, "BOSS") {
		return 0
	}

	// 如果是疯狂血腥女巫、钻地的冰虫、牙皮这个组合，则只有拥有者可以竞拍火之眼和督瑞尔
	if game.CurrentBiddingEntity.Name == "火之眼" || game.CurrentBiddingEntity.Name == "督瑞尔" {
		for i := 0; i < game.PlayerNum; i++ {
			if slices.Contains(game.PlayerEntities[i], &疯狂血腥女巫) && slices.Contains(game.PlayerEntities[i], &钻地的冰虫) && slices.Contains(game.PlayerEntities[i], &牙皮) && i != playerId {
				return 0
			}
		}
	}

	// 衣卒尔限制：每5回合至少出一次0能力点
	if playerId == int(game.衣卒尔_playerid) {
		if len(game.衣卒尔_player_history) == 4 && !game.衣卒尔_player_history[0] && !game.衣卒尔_player_history[1] && !game.衣卒尔_player_history[2] && !game.衣卒尔_player_history[3] {
			return 0
		}
	}

	// 破坏者卡兰索, 诅咒的阿克姆, 血腥的巴特克, 不洁的凡塔, 古难记录者：三个野蛮人试炼
	if game.CurrentBiddingEntity.Name == "破坏者卡兰索" ||
		game.CurrentBiddingEntity.Name == "诅咒的阿克姆" ||
		game.CurrentBiddingEntity.Name == "血腥的巴特克" ||
		game.CurrentBiddingEntity.Name == "不洁的凡塔" ||
		game.CurrentBiddingEntity.Name == "古难记录者" {
		if !game.三个野蛮人试炼_passed[playerId] {
			return 0
		}
	}

	bid_effect := game.getCurrentPriceDiscount(playerId)

	var max_bid int
	if bid_effect < 0.001 {
		max_bid = 10000
	} else {
		max_bid = int((float64(game.PlayerCredits[playerId]) + 0.999) / bid_effect)
	}
	return max_bid
}

func (game *Game) checkPriceValid(playerId int, price int) error {
	if maxBidValue := game.getMaxBidValue(playerId); price > maxBidValue {
		return fmt.Errorf("出价超过最大值 %d", maxBidValue)
	}
	if game.CurrentBiddingEntity.BidChecker != nil {
		if err := game.CurrentBiddingEntity.BidChecker(game.PlayerEntities[playerId], price); err != nil {
			return err
		}
	}
	return nil
}
