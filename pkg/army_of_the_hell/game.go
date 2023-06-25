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

	CurrentEntityPool          []*Entity
	CurrentBiddingEntity       *Entity // current bidding entity is not removed from pool
	CurrentBiddingEntityIndex  int
	CurrentBiddingEntity2      *Entity
	CurrentBiddingEntityIndex2 int
	CurrentPlayerReady         []bool
	CurrentPlayerBid           []int
	CurrentPlayerBid2          []int
	CurrentPlayerNickname      []string

	ProcessingLock sync.Mutex

	洞穴重生的邪恶之犬_triggered bool
	海法斯特盔甲制作者_triggered bool
	黑暗长老_天黑_countdown   int8
	牙皮_triggered        bool
	火之眼_countdown       int8
	召唤者_countdown       int8
	暴躁外皮_countdown      int8
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
	三个野蛮人试炼_unlocked    bool
	三个野蛮人试炼_passed      []bool
	破坏者卡兰索_unlocked     bool
	格瑞斯华尔德_playerid     int8
	尼拉塞克_playerid       int8
	沙漠三小队_playerid      int
	火之眼_playerid        int
	督瑞尔_playerid        int
	议会成员解锁              bool
	本回合未通过野蛮人试炼         []bool

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
	game.CurrentPlayerBid2 = make([]int, number_players)
	game.CurrentPlayerNickname = make([]string, number_players)
	game.WaitResponsePlayerId = -1

	game.西希之王_playerid = -1
	game.衣卒尔_playerid = -1
	game.三个野蛮人试炼_passed = make([]bool, number_players)
	game.本回合未通过野蛮人试炼 = make([]bool, number_players)
	game.火之眼_playerid = -1
	game.督瑞尔_playerid = -1
	game.格瑞斯华尔德_playerid = -1
	game.沙漠三小队_playerid = -1
	game.尼拉塞克_playerid = -1

	game.SingleMode = true
	// game.DoubleMode = true

	game.PrintFunc = func(msg string) {
		fmt.Println(msg)
	}
	return game
}

func (game *Game) Start() {
	game.CurrentEntityPool = []*Entity{
		&尸体发火,
		// &三个野蛮人,
	}

	if game.PlayerNum >= 5 {
		WINSTATE_BOSS_COUNT = 1
	} else if game.PlayerNum >= 3 {
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
	for id, ready := range game.CurrentPlayerReady {
		if game.PlayerLeaved[id] {
			continue
		}
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

func (game *Game) GivePrices(playerId int, price1 int, price2 int) error {
	game.ProcessingLock.Lock()
	defer game.ProcessingLock.Unlock()

	if game.CurrentPlayerReady[playerId] {
		return errors.New("本回合已经出价了")
	}
	if err := game.checkPricesValid(playerId, price1, price2); err != nil {
		return err
	}
	game.CurrentPlayerBid[playerId] = price1
	game.CurrentPlayerBid2[playerId] = price2
	game.CurrentPlayerReady[playerId] = true

	all_player_ready := true
	for id, ready := range game.CurrentPlayerReady {
		if game.PlayerLeaved[id] {
			continue
		}
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

// return true if game ends
func (game *Game) PlayerLeave(playerId int) bool {
	game.PrintFunc(game.CurrentPlayerNickname[playerId] + "离开了游戏")

	game.PlayerLeaved[playerId] = true
	game.CurrentPlayerReady[playerId] = true
	game.CurrentPlayerBid[playerId] = 0
	if game.WaitResponsePlayerId == playerId {
		game.Response = "否"
		game.WaitResponsePlayerId = -1
	}

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

	all_player_ready := true
	for id, ready := range game.CurrentPlayerReady {
		if game.PlayerLeaved[id] {
			continue
		}
		if !ready {
			all_player_ready = false
		}
	}
	if all_player_ready {
		game.endTurn()

		if win, _ := game.checkWinState(); win {
			return true
		}

		game.startNewTurn()
	}

	return false
}

func (game *Game) startNewTurn() {
	game.Turn++

	// 暴躁外皮：被拍得后的第2回合洗入牌库
	game.暴躁外皮_countdown--
	if game.暴躁外皮_countdown == 0 {
		game.CurrentEntityPool = append(game.CurrentEntityPool, &暴躁外皮2)
	}

	// 记录回合开始时各个玩家是否通过了野蛮人试炼
	for i := 0; i < game.PlayerNum; i++ {
		game.本回合未通过野蛮人试炼[i] = game.三个野蛮人试炼_passed[i]
	}

	var poolString string
	for _, entity := range game.CurrentEntityPool {
		poolString += entity.Name + "\n"
	}
	game.PrintFunc(fmt.Sprintf("第%d回合开始。本回合抽取池：\n%s", game.Turn, poolString))

	if len(game.CurrentEntityPool) == 0 {
		game.PrintFunc("拍卖池中已经无角色，游戏结束。")
		return
	}

	game.CurrentBiddingEntityIndex = rand.Int() % len(game.CurrentEntityPool)
	game.CurrentBiddingEntity = game.CurrentEntityPool[game.CurrentBiddingEntityIndex]

	game.古代无魂之卡_countdown--
	game.督军山克_countdown--

	// 古代无魂之卡：触发后第6回合改为拍卖【督瑞尔】
	if game.古代无魂之卡_countdown == 0 {
		game.CurrentEntityPool = append(game.CurrentEntityPool, &督瑞尔)
		game.CurrentBiddingEntity = &督瑞尔
		game.CurrentBiddingEntityIndex = len(game.CurrentEntityPool) - 1
		game.PrintFunc("【古代无魂之卡】效果触发。")

		// 达克法恩：跳过下次抽牌，改为抽到督军山克。
	} else if game.达克法恩_triggered {
		game.达克法恩_triggered = false
		if slices.Contains(game.CurrentEntityPool, &督军山克) {
			game.CurrentBiddingEntityIndex = slices.Index(game.CurrentEntityPool, &督军山克)
			game.CurrentBiddingEntity = &督军山克
		}

		// 督军山克：接下来第三次有效抽牌抽到矫正者怪异，已经抽到则失效。
	} else if game.督军山克_countdown == 0 && slices.Contains(game.CurrentEntityPool, &矫正者怪异) {
		game.PrintFunc("督军山克效果已触发。")
		game.CurrentBiddingEntityIndex = slices.Index(game.CurrentEntityPool, &矫正者怪异)
		game.CurrentBiddingEntity = &矫正者怪异

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
		game.三个野蛮人试炼_unlocked = true
		game.CurrentEntityPool = append(game.CurrentEntityPool[:game.CurrentBiddingEntityIndex], game.CurrentEntityPool[game.CurrentBiddingEntityIndex+1:]...)
		if game.尼拉塞克_playerid != -1 {
			game.破坏者卡兰索_unlocked = true
			game.CurrentEntityPool = append(game.CurrentEntityPool, &破坏者卡兰索)
		}
		game.handleEndTurnCredits()
		game.startNewTurn()
		return

	}
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
	if game.CurrentBiddingEntity.Name == "黑暗长老" {
		game.黑暗长老_天黑_countdown = 3
	}

	// 格瑞斯华尔德：只有拉卡尼休和树头木拳的拥有者可以竞拍
	if game.CurrentBiddingEntity.Name == "格瑞斯华尔德" {
		for i := 0; i < game.PlayerNum; i++ {
			if slices.Contains(game.PlayerEntities[i], &拉卡尼休) && slices.Contains(game.PlayerEntities[i], &树头木拳) {
				game.格瑞斯华尔德_playerid = int8(i)
			}
		}
	}

	// 沙漠三小队：【钻地的冰虫】、【牙皮】、【疯狂血腥女巫】在一家，则只有他可以竞拍【火之眼】、【督瑞尔】
	game.沙漠三小队_playerid = -1
	if game.CurrentBiddingEntity.Name == "火之眼" || game.CurrentBiddingEntity.Name == "督瑞尔" {
		for i := 0; i < game.PlayerNum; i++ {
			if slices.Contains(game.PlayerEntities[i], &疯狂血腥女巫) && slices.Contains(game.PlayerEntities[i], &钻地的冰虫) && slices.Contains(game.PlayerEntities[i], &牙皮) {
				game.沙漠三小队_playerid = i
				break
			}
		}
	}

	// 火之眼：如果两回合内抽到【召唤者】，只有拍得火之眼者可以竞拍【召唤者】
	if game.火之眼_countdown > 0 {
		if game.SingleMode {
			game.火之眼_countdown--
		}
		if game.CurrentBiddingEntity.Name == "召唤者" {
			for i := 0; i < game.PlayerNum; i++ {
				if slices.Contains(game.PlayerEntities[i], &火之眼) {
					game.火之眼_playerid = i
				}
			}
		}
	}

	// 召唤者：如果两回合内抽到【督瑞尔】，只有拍得者可以竞拍【督瑞尔】
	if game.召唤者_countdown > 0 {
		if game.SingleMode {
			game.召唤者_countdown--
		}
		if game.CurrentBiddingEntity.Name == "督瑞尔" {
			for i := 0; i < game.PlayerNum; i++ {
				if slices.Contains(game.PlayerEntities[i], &召唤者) {
					game.督瑞尔_playerid = i
				}
			}
		}
	}

	if game.黑暗长老_天黑_countdown > 0 && !game.牙皮_triggered && game.CurrentBiddingEntity.Name != "黑暗长老" {
		if game.SingleMode {
			game.黑暗长老_天黑_countdown--
		}
		game.PrintFunc("由于黑暗长老效果，天黑了，本回合拍卖的角色未知。")
	} else {
		bidRanges := ""
		for i := 0; i < game.PlayerNum; i++ {
			bidRanges += fmt.Sprintf("%s 最高出价：%d (%+d%%)\n", game.CurrentPlayerNickname[i], game.getMaxBidValue(i), int(game.getCurrentPriceDiscount(i)*100-100.01))
		}
		game.PrintFunc(fmt.Sprintf("本回合拍卖的角色：%v\n\n%s", game.CurrentBiddingEntity, bidRanges))
	}

	if game.DoubleMode {
		if len(game.CurrentEntityPool) == 1 {
			game.CurrentBiddingEntityIndex2 = -1
			game.CurrentBiddingEntity2 = &凑数的
		} else {
			for {
				game.CurrentBiddingEntityIndex2 = rand.Intn(len(game.CurrentEntityPool))
				if game.CurrentBiddingEntityIndex2 == game.CurrentBiddingEntityIndex {
					continue
				}
				break
			}
			game.CurrentBiddingEntity2 = game.CurrentEntityPool[game.CurrentBiddingEntityIndex2]
		}

		// 达克法恩：跳过下次抽牌，改为抽到督军山克。
		if game.达克法恩_triggered {
			game.达克法恩_triggered = false
			if slices.Contains(game.CurrentEntityPool, &督军山克) {
				game.CurrentBiddingEntityIndex2 = slices.Index(game.CurrentEntityPool, &督军山克)
				game.CurrentBiddingEntity2 = &督军山克
			}

			// 督军山克：接下来第三次有效抽牌抽到矫正者怪异，已经抽到则失效。
		} else if game.督军山克_countdown == 0 && slices.Contains(game.CurrentEntityPool, &矫正者怪异) && game.CurrentBiddingEntity.Name != "矫正者怪异" {
			game.PrintFunc("督军山克效果已触发。")
			game.CurrentBiddingEntityIndex2 = slices.Index(game.CurrentEntityPool, &矫正者怪异)
			game.CurrentBiddingEntity2 = &矫正者怪异

			// 议会成员：优先抽取
		} else if game.优先抽取议会成员_triggered {
			game.优先抽取议会成员_triggered = false
			for i := 0; i < len(game.CurrentEntityPool); i++ {
				if game.CurrentEntityPool[i].Name == "邪恶之手伊斯梅尔" ||
					game.CurrentEntityPool[i].Name == "火焰之指吉列布" ||
					game.CurrentEntityPool[i].Name == "冰拳托克" {
					game.PrintFunc("优先抽取议会成员效果已触发。")
					game.CurrentBiddingEntityIndex2 = i
					game.CurrentBiddingEntity2 = game.CurrentEntityPool[i]
					break
				}
			}

			// 破坏者卡兰索：优先抽取
		} else if game.破坏者卡兰索_triggered {
			game.破坏者卡兰索_triggered = false
			for i := 0; i < len(game.CurrentEntityPool); i++ {
				if game.CurrentEntityPool[i].Name == "诅咒的阿克姆" {
					game.PrintFunc("优先抽取诅咒的阿克姆效果已触发。")
					game.CurrentBiddingEntityIndex2 = i
					game.CurrentBiddingEntity2 = game.CurrentEntityPool[i]
					break
				}
			}

			// 诅咒的阿克姆：优先抽取
		} else if game.诅咒的阿克姆_triggered {
			game.诅咒的阿克姆_triggered = false
			for i := 0; i < len(game.CurrentEntityPool); i++ {
				if game.CurrentEntityPool[i].Name == "血腥的巴特克" {
					game.PrintFunc("优先抽取血腥的巴特克效果已触发。")
					game.CurrentBiddingEntityIndex2 = i
					game.CurrentBiddingEntity2 = game.CurrentEntityPool[i]
					break
				}
			}

			// 血腥的巴特克：优先抽取
		} else if game.血腥的巴特克_triggered {
			game.血腥的巴特克_triggered = false
			for i := 0; i < len(game.CurrentEntityPool); i++ {
				if game.CurrentEntityPool[i].Name == "不洁的凡塔" {
					game.PrintFunc("优先抽取不洁的凡塔效果已触发。")
					game.CurrentBiddingEntityIndex2 = i
					game.CurrentBiddingEntity2 = game.CurrentEntityPool[i]
					break
				}
			}

			// 不洁的凡塔：优先抽取
		} else if game.不洁的凡塔_triggered {
			game.不洁的凡塔_triggered = false
			for i := 0; i < len(game.CurrentEntityPool); i++ {
				if game.CurrentEntityPool[i].Name == "古难记录者" {
					game.PrintFunc("优先抽取古难记录者效果已触发。")
					game.CurrentBiddingEntityIndex2 = i
					game.CurrentBiddingEntity2 = game.CurrentEntityPool[i]
					break
				}
			}

			// 古难记录者：优先抽取
		} else if game.古难记录者_triggered {
			game.古难记录者_triggered = false
			for i := 0; i < len(game.CurrentEntityPool); i++ {
				if game.CurrentEntityPool[i].Name == "巴尔" {
					game.PrintFunc("优先抽取巴尔效果已触发。")
					game.CurrentBiddingEntityIndex2 = i
					game.CurrentBiddingEntity2 = game.CurrentEntityPool[i]
					break
				}
			}

			// 洞穴重生的邪恶之犬：当第一次被抽到时，跳过本轮所有竞拍
		} else if !game.洞穴重生的邪恶之犬_triggered && game.CurrentBiddingEntity2.Name == "洞穴重生的邪恶之犬" {
			game.洞穴重生的邪恶之犬_triggered = true
			game.PrintFunc("洞穴重生的邪恶之犬被抽到，跳过本轮所有竞拍。")
			game.handleEndTurnCredits()
			game.startNewTurn()
			return

			// 三个野蛮人试炼：当被抽到时，改为开启试炼
		} else if game.CurrentBiddingEntity2.Name == "三个野蛮人" {
			game.PrintFunc("【三个野蛮人】试炼已开启。玩家可以随时输入【接受试炼】，支付50能力点通过三个野蛮人的试炼。一方通过试炼后解锁【破坏者卡兰索】。")
			game.三个野蛮人试炼_unlocked = true
			game.CurrentEntityPool = append(game.CurrentEntityPool[:game.CurrentBiddingEntityIndex2], game.CurrentEntityPool[game.CurrentBiddingEntityIndex2+1:]...)
			if game.尼拉塞克_playerid != -1 {
				game.破坏者卡兰索_unlocked = true
				game.CurrentEntityPool = append(game.CurrentEntityPool, &破坏者卡兰索)
			}
			game.handleEndTurnCredits()
			game.startNewTurn()
			return

		}

		// 海法斯特盔甲制作者：第一次被抽到时，改为抽到ACT5-达克法恩
		if !game.海法斯特盔甲制作者_triggered && game.CurrentBiddingEntity2.Name == "海法斯特盔甲制作者" {
			game.海法斯特盔甲制作者_triggered = true
			game.CurrentEntityPool = append(game.CurrentEntityPool, &达克法恩)
			game.CurrentBiddingEntity2 = &达克法恩
			game.CurrentBiddingEntityIndex2 = len(game.CurrentEntityPool) - 1
			game.PrintFunc("【海法斯特盔甲制作者】效果触发，本回合拍卖的角色改为【ACT5-达克法恩】。")
		}

		// 黑暗长老：抽到后三回合【天黑】；牙皮：永久取消【天黑】
		if game.CurrentBiddingEntity2.Name == "牙皮" {
			game.牙皮_triggered = true
		}
		if game.CurrentBiddingEntity2.Name == "黑暗长老" {
			game.黑暗长老_天黑_countdown = 3
		}

		// 格瑞斯华尔德：只有拉卡尼休和树头木拳的拥有者可以竞拍
		if game.CurrentBiddingEntity2.Name == "格瑞斯华尔德" {
			for i := 0; i < game.PlayerNum; i++ {
				if slices.Contains(game.PlayerEntities[i], &拉卡尼休) && slices.Contains(game.PlayerEntities[i], &树头木拳) {
					game.格瑞斯华尔德_playerid = int8(i)
				}
			}
		}

		// 沙漠三小队：【钻地的冰虫】、【牙皮】、【疯狂血腥女巫】在一家，则只有他可以竞拍【火之眼】、【督瑞尔】
		game.沙漠三小队_playerid = -1
		if game.CurrentBiddingEntity2.Name == "火之眼" || game.CurrentBiddingEntity2.Name == "督瑞尔" {
			for i := 0; i < game.PlayerNum; i++ {
				if slices.Contains(game.PlayerEntities[i], &疯狂血腥女巫) && slices.Contains(game.PlayerEntities[i], &钻地的冰虫) && slices.Contains(game.PlayerEntities[i], &牙皮) {
					game.沙漠三小队_playerid = i
					break
				}
			}
		}

		// 火之眼：如果两回合内抽到【召唤者】，只有拍得火之眼者可以竞拍【召唤者】
		if game.火之眼_countdown > 0 {
			game.火之眼_countdown--
			if game.CurrentBiddingEntity2.Name == "召唤者" {
				for i := 0; i < game.PlayerNum; i++ {
					if slices.Contains(game.PlayerEntities[i], &火之眼) {
						game.火之眼_playerid = i
					}
				}
			}
		}

		// 召唤者：如果两回合内抽到【督瑞尔】，只有拍得者可以竞拍【督瑞尔】
		if game.召唤者_countdown > 0 {
			game.召唤者_countdown--
			if game.CurrentBiddingEntity2.Name == "督瑞尔" {
				for i := 0; i < game.PlayerNum; i++ {
					if slices.Contains(game.PlayerEntities[i], &召唤者) {
						game.督瑞尔_playerid = i
					}
				}
			}
		}

		if game.黑暗长老_天黑_countdown > 0 && !game.牙皮_triggered && game.CurrentBiddingEntity.Name != "黑暗长老" && game.CurrentBiddingEntity2.Name != "黑暗长老" {
			game.黑暗长老_天黑_countdown--
			game.PrintFunc("由于黑暗长老效果，天黑了，本回合拍卖的角色未知。")
		} else {
			if game.CurrentBiddingEntity2.Name != "" {
				bidRanges := ""
				for i := 0; i < game.PlayerNum; i++ {
					bidRanges += fmt.Sprintf("%s 最高出价：%d (%+d%%)\n", game.CurrentPlayerNickname[i], game.getMaxBidValue2(i), int(game.getCurrentPriceDiscount2(i)*100-100.01))
				}
				game.PrintFunc(fmt.Sprintf("本回合拍卖的角色2：%v\n\n%s", game.CurrentBiddingEntity2, bidRanges))
			} else {
				game.PrintFunc("本回合仅有一个角色拍卖，你只需输入一个数字，表示对该角色的出价。")
				for i := 0; i < game.PlayerNum; i++ {
					game.CurrentPlayerBid2[i] = 0
				}
			}
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
	winner2 := -1
	largest := -1
	largest2 := -1
	for id, value := range game.CurrentPlayerBid {
		if largest < value {
			largest = value
			winner = id
		} else if largest == value {
			winner = -1
		}
	}
	if game.DoubleMode {
		for id, value := range game.CurrentPlayerBid2 {
			if largest2 < value {
				largest2 = value
				winner2 = id
			} else if largest2 == value {
				winner2 = -1
			}
		}
	}

	bidDescription := ""
	if game.SingleMode || game.CurrentBiddingEntity2.Name == "" {
		for i := 0; i < game.PlayerNum; i++ {
			bidDescription += fmt.Sprintf("\n%s ：%d", game.CurrentPlayerNickname[i], game.CurrentPlayerBid[i])
		}
	} else {
		for i := 0; i < game.PlayerNum; i++ {
			bidDescription += fmt.Sprintf("\n%s ：%d %d", game.CurrentPlayerNickname[i], game.CurrentPlayerBid[i], game.CurrentPlayerBid2[i])
		}
	}

	if winner == -1 { // draw
		bidDescription += "\n\n" + game.CurrentBiddingEntity.Name + "最高出价相同，本回合流拍。"
	} else {
		bidDescription += fmt.Sprintf("\n\n%s 成功拍得了：%v", game.CurrentPlayerNickname[winner], game.CurrentBiddingEntity)
	}
	if game.DoubleMode {
		if winner2 == -1 {
			if game.CurrentBiddingEntity2.Name != "" {
				bidDescription += "\n\n" + game.CurrentBiddingEntity2.Name + "最高出价相同，本回合流拍。"
			}
		} else {
			bidDescription += fmt.Sprintf("\n\n%s 成功拍得了：%v", game.CurrentPlayerNickname[winner2], game.CurrentBiddingEntity2)
		}
	}
	if winner != -1 {
		game.PlayerCredits[winner] -= int(float64(game.CurrentPlayerBid[winner]) * game.getCurrentPriceDiscount(winner))
	}
	if winner2 != -1 {
		game.PlayerCredits[winner2] -= int(float64(game.CurrentPlayerBid2[winner2]) * game.getCurrentPriceDiscount2(winner2))
	}
	if winner != -1 { // 尸体发火流拍也正常解锁后续
		if winner != -1 {
			game.PlayerEntities[winner] = append(game.PlayerEntities[winner], game.CurrentBiddingEntity)
			game.CurrentEntityPool = append(game.CurrentEntityPool[:game.CurrentBiddingEntityIndex], game.CurrentEntityPool[game.CurrentBiddingEntityIndex+1:]...)
		}
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
	if winner2 != -1 && game.CurrentBiddingEntityIndex2 != -1 {
		game.PlayerEntities[winner2] = append(game.PlayerEntities[winner2], game.CurrentBiddingEntity2)
		index2 := slices.Index(game.CurrentEntityPool, game.CurrentBiddingEntity2)
		game.CurrentEntityPool = append(game.CurrentEntityPool[:index2], game.CurrentEntityPool[index2+1:]...)
		for _, entity := range game.CurrentBiddingEntity2.UnlockEntities {
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

	// 尸体发火：第一回合无论如何解锁 毕须博须, 罗达门特
	if game.Turn == 1 {
		game.CurrentEntityPool = append(game.CurrentEntityPool, &毕须博须, &罗达门特)
	}

	// 火之眼：如果两回合内抽到【召唤者】，只有拍得火之眼者可以竞拍【召唤者】
	if game.CurrentBiddingEntity.Name == "火之眼" && winner != -1 {
		game.火之眼_countdown = 2
	}
	if winner2 != -1 && game.CurrentBiddingEntity2.Name == "火之眼" {
		game.火之眼_countdown = 2
	}

	game.火之眼_playerid = -1
	game.督瑞尔_playerid = -1

	// 召唤者：如果两回合内抽到【督瑞尔】，只有拍得者可以竞拍【督瑞尔】
	if game.CurrentBiddingEntity.Name == "召唤者" && winner != -1 {
		game.召唤者_countdown = 2
	}
	if winner2 != -1 && game.CurrentBiddingEntity2.Name == "召唤者" {
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

	if winner2 != -1 && game.CurrentBiddingEntity2.Name == "古代无魂之卡" {
		game.PrintFunc("【古代无魂之卡】效果触发，请拍得的玩家私聊选择是否发动效果（是/否）。")
	wait_response_:
		game.WaitResponsePlayerId = winner2
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
			goto wait_response_
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

	if winner2 != -1 && game.CurrentBiddingEntity2.Name == "达克法恩" {
		game.PrintFunc("【达克法恩】可选效果：跳过下次抽牌，改为抽到督军山克。\n\n请拍得的玩家私聊选择是否发动（是/否）。")
	wait_response2_:
		game.WaitResponsePlayerId = winner2
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
			goto wait_response2_
		}
	}

	// 西希之王：下回合不能出价boss以外的怪物
	if game.CurrentBiddingEntity.Name == "西希之王" && winner != -1 {
		game.西希之王_playerid = int8(winner)
	} else if winner2 != -1 && game.CurrentBiddingEntity2.Name == "西希之王" {
		game.西希之王_playerid = int8(winner2)
	} else if game.西希之王_playerid != -1 {
		game.西希之王_playerid = -1
	}

	// 邪恶之手伊斯梅尔，火焰之指吉列布，冰拳托克：优先抽取议会成员
	if winner != -1 && (game.CurrentBiddingEntity.Name == "邪恶之手伊斯梅尔" || game.CurrentBiddingEntity.Name == "火焰之指吉列布" || game.CurrentBiddingEntity.Name == "冰拳托克") {
		game.优先抽取议会成员_triggered = true
	}
	if winner2 != -1 && (game.CurrentBiddingEntity2.Name == "邪恶之手伊斯梅尔" || game.CurrentBiddingEntity2.Name == "火焰之指吉列布" || game.CurrentBiddingEntity2.Name == "冰拳托克") {
		game.优先抽取议会成员_triggered = true
	}

	// 衣卒尔
	if game.CurrentBiddingEntity.Name == "衣卒尔" && winner != -1 {
		game.衣卒尔_playerid = int8(winner)
		game.衣卒尔_player_history = nil
	} else if winner2 != -1 && game.CurrentBiddingEntity2.Name == "衣卒尔" {
		game.衣卒尔_playerid = int8(winner2)
		game.衣卒尔_player_history = nil
	} else if game.衣卒尔_playerid != -1 {
		if len(game.衣卒尔_player_history) != 4 {
			if game.SingleMode {
				game.衣卒尔_player_history = append(game.衣卒尔_player_history, game.CurrentPlayerBid[game.衣卒尔_playerid] == 0)
			} else {
				game.衣卒尔_player_history = append(game.衣卒尔_player_history, game.CurrentPlayerBid[game.衣卒尔_playerid] == 0 || game.CurrentPlayerBid2[game.衣卒尔_playerid] == 0)
			}
		} else {
			game.衣卒尔_player_history = nil
		}
	}

	// 督军山克
	if game.CurrentBiddingEntity.Name == "督军山克" && winner != -1 {
		game.督军山克_countdown = 3
	}
	if winner2 != -1 && game.CurrentBiddingEntity2.Name == "督军山克" {
		game.督军山克_countdown = 3
	}

	// 暴躁外皮
	if winner != -1 && game.CurrentBiddingEntity.Name == "暴躁外皮" {
		game.暴躁外皮_countdown = 2
	}
	if winner2 != -1 && game.CurrentBiddingEntity2.Name == "暴躁外皮" {
		game.暴躁外皮_countdown = 2
	}

	// 破坏者卡兰索, 诅咒的阿克姆, 血腥的巴特克, 不洁的凡塔, 古难记录者
	if winner != -1 && game.CurrentBiddingEntity.Name == "破坏者卡兰索" {
		game.破坏者卡兰索_triggered = true
	}
	if winner != -1 && game.CurrentBiddingEntity.Name == "诅咒的阿克姆" {
		game.诅咒的阿克姆_triggered = true
	}
	if winner != -1 && game.CurrentBiddingEntity.Name == "血腥的巴特克" {
		game.血腥的巴特克_triggered = true
	}
	if winner != -1 && game.CurrentBiddingEntity.Name == "不洁的凡塔" {
		game.不洁的凡塔_triggered = true
	}
	if winner != -1 && game.CurrentBiddingEntity.Name == "古难记录者" {
		game.古难记录者_triggered = true
	}

	if winner2 != -1 && game.CurrentBiddingEntity2.Name == "破坏者卡兰索" {
		game.破坏者卡兰索_triggered = true
	}
	if winner2 != -1 && game.CurrentBiddingEntity2.Name == "诅咒的阿克姆" {
		game.诅咒的阿克姆_triggered = true
	}
	if winner2 != -1 && game.CurrentBiddingEntity2.Name == "血腥的巴特克" {
		game.血腥的巴特克_triggered = true
	}
	if winner2 != -1 && game.CurrentBiddingEntity2.Name == "不洁的凡塔" {
		game.不洁的凡塔_triggered = true
	}
	if winner2 != -1 && game.CurrentBiddingEntity2.Name == "古难记录者" {
		game.古难记录者_triggered = true
	}

	// 尼拉塞克：直接通过三个野蛮人试炼
	if game.CurrentBiddingEntity.Name == "尼拉塞克" && winner != -1 {
		game.三个野蛮人试炼_passed[winner] = true
		if game.三个野蛮人试炼_unlocked && !game.破坏者卡兰索_unlocked {
			game.破坏者卡兰索_unlocked = true
			game.CurrentEntityPool = append(game.CurrentEntityPool, &破坏者卡兰索)
		}
		game.尼拉塞克_playerid = int8(winner)
	}
	if winner2 != -1 && game.CurrentBiddingEntity2.Name == "尼拉塞克" {
		game.三个野蛮人试炼_passed[winner2] = true
		if game.三个野蛮人试炼_unlocked && !game.破坏者卡兰索_unlocked {
			game.破坏者卡兰索_unlocked = true
			game.CurrentEntityPool = append(game.CurrentEntityPool, &破坏者卡兰索)
		}
		game.尼拉塞克_playerid = int8(winner2)
	}

	// 双拍模式
	if game.SingleMode && slices.Contains(game.CurrentBiddingEntity.Tags, "BOSS") && winner != -1 {
		game.SingleMode = false
		game.DoubleMode = true
		game.PrintFunc("由于有BOSS被拍得，【双拍模式】已开启，请输入两个数字，分别代表对两个怪物的出价，以空格分隔。")
	}

	game.CurrentBiddingEntity = nil
	game.CurrentBiddingEntityIndex = -1
	game.CurrentBiddingEntity2 = nil
	game.CurrentBiddingEntityIndex2 = -1

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

func (game *Game) getCurrentPriceDiscount2(playerId int) float64 {
	bid_effect := 1.0

	for _, entity := range game.PlayerEntities[playerId] {
		if entity.BidEffected(game.CurrentBiddingEntity2.Tags) {
			bid_effect += entity.BidAffectValue
		}

		// 女伯爵：对安达利尔出价-15%（在boss的基础上额外-10%）
		if entity.Name == "女伯爵" && game.CurrentBiddingEntity2.Name == "安达利尔" {
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

	if game.格瑞斯华尔德_playerid != -1 && game.格瑞斯华尔德_playerid != int8(playerId) && game.CurrentBiddingEntity.Name == "格瑞斯华尔德" {
		return 0
	}
	if game.沙漠三小队_playerid != -1 && game.沙漠三小队_playerid != playerId && (game.CurrentBiddingEntity.Name == "火之眼" || game.CurrentBiddingEntity.Name == "督瑞尔") {
		return 0
	}
	if game.火之眼_playerid != -1 && game.火之眼_playerid != playerId && game.CurrentBiddingEntity.Name == "召唤者" {
		return 0
	}
	if game.督瑞尔_playerid != -1 && game.督瑞尔_playerid != playerId && game.CurrentBiddingEntity.Name == "督瑞尔" {
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
	if game.SingleMode && playerId == int(game.衣卒尔_playerid) {
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
		if !game.本回合未通过野蛮人试炼[playerId] {
			return 0
		}
	}

	bid_effect := game.getCurrentPriceDiscount(playerId)

	var max_bid int
	if bid_effect < 0.0001 {
		max_bid = 20000
	} else {
		max_bid = int((float64(game.PlayerCredits[playerId]) + 0.9999) / bid_effect)
	}
	return max_bid
}

func (game *Game) getMaxBidValue2(playerId int) int {
	if game.西希之王_playerid == int8(playerId) && !slices.Contains(game.CurrentBiddingEntity2.Tags, "BOSS") {
		return 0
	}

	if game.格瑞斯华尔德_playerid != -1 && game.格瑞斯华尔德_playerid != int8(playerId) && game.CurrentBiddingEntity2.Name == "格瑞斯华尔德" {
		return 0
	}
	if game.沙漠三小队_playerid != -1 && game.沙漠三小队_playerid != playerId && (game.CurrentBiddingEntity2.Name == "火之眼" || game.CurrentBiddingEntity2.Name == "督瑞尔") {
		return 0
	}
	if game.火之眼_playerid != -1 && game.火之眼_playerid != playerId && game.CurrentBiddingEntity2.Name == "召唤者" {
		return 0
	}
	if game.督瑞尔_playerid != -1 && game.督瑞尔_playerid != playerId && game.CurrentBiddingEntity2.Name == "督瑞尔" {
		return 0
	}

	// 如果是疯狂血腥女巫、钻地的冰虫、牙皮这个组合，则只有拥有者可以竞拍火之眼和督瑞尔
	if game.CurrentBiddingEntity2.Name == "火之眼" || game.CurrentBiddingEntity2.Name == "督瑞尔" {
		for i := 0; i < game.PlayerNum; i++ {
			if slices.Contains(game.PlayerEntities[i], &疯狂血腥女巫) && slices.Contains(game.PlayerEntities[i], &钻地的冰虫) && slices.Contains(game.PlayerEntities[i], &牙皮) && i != playerId {
				return 0
			}
		}
	}

	// 破坏者卡兰索, 诅咒的阿克姆, 血腥的巴特克, 不洁的凡塔, 古难记录者：三个野蛮人试炼
	if game.CurrentBiddingEntity2.Name == "破坏者卡兰索" ||
		game.CurrentBiddingEntity2.Name == "诅咒的阿克姆" ||
		game.CurrentBiddingEntity2.Name == "血腥的巴特克" ||
		game.CurrentBiddingEntity2.Name == "不洁的凡塔" ||
		game.CurrentBiddingEntity2.Name == "古难记录者" {
		if !game.三个野蛮人试炼_passed[playerId] {
			return 0
		}
	}

	bid_effect := game.getCurrentPriceDiscount2(playerId)

	var max_bid int
	if bid_effect < 0.0001 {
		max_bid = 20000
	} else {
		max_bid = int((float64(game.PlayerCredits[playerId]) + 0.9999) / bid_effect)
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

func (game *Game) checkPricesValid(playerId int, price1 int, price2 int) error {
	if maxBidValue := game.getMaxBidValue(playerId); price1 > maxBidValue {
		return fmt.Errorf("对 %s 的出价超过最大值 %d", game.CurrentBiddingEntity.Name, maxBidValue)
	}
	if maxBidValue := game.getMaxBidValue2(playerId); price2 > maxBidValue {
		return fmt.Errorf("对 %s 的出价超过最大值 %d", game.CurrentBiddingEntity2.Name, maxBidValue)
	}
	if game.衣卒尔_playerid == int8(playerId) &&
		len(game.衣卒尔_player_history) == 4 &&
		!game.衣卒尔_player_history[0] &&
		!game.衣卒尔_player_history[1] &&
		!game.衣卒尔_player_history[2] &&
		!game.衣卒尔_player_history[3] &&
		price1 > 0 &&
		(price2 > 0 || game.CurrentBiddingEntity2.Name == "") {
		return fmt.Errorf("由于衣卒尔效果限制，至少得为一个单位出价0点")
	}
	if game.CurrentBiddingEntity.BidChecker != nil {
		if err := game.CurrentBiddingEntity.BidChecker(game.PlayerEntities[playerId], price1); err != nil {
			return err
		}
	}
	if game.CurrentBiddingEntity2.BidChecker != nil {
		if err := game.CurrentBiddingEntity2.BidChecker(game.PlayerEntities[playerId], price2); err != nil {
			return err
		}
	}
	discount := game.getCurrentPriceDiscount(playerId)
	discount2 := game.getCurrentPriceDiscount2(playerId)
	total := int(float64(price1)*discount) + int(float64(price2)*discount2)
	if total > game.PlayerCredits[playerId] {
		return fmt.Errorf("你当前的总花费为 %d，超过能力点总数 %d", total, game.PlayerCredits[playerId])
	}
	return nil
}
