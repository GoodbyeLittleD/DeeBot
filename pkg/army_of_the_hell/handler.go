package army_of_the_hell

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"golang.org/x/exp/slices"
)

var groupWhiteList = []int64{
	541065568,  // 测试群
	730561889,  // 测试群
	757553235,  // 一费奶五战队
	1006990128, // 单挑联赛Week8
}

var sendGroupWelcome = true
var sendGroupMsg = false
var sendPrivateMsg = true

func Handle() {
	var currentGroupId int64
	var currentPlayerIds []int64
	var currentWatchingPlayerIds []int64
	var currentPlayerNames []string
	var gameStarted bool = false
	var gameLock sync.Mutex
	var game *Game

	var groupTicket = make(chan struct{})
	var groupMsgQueue = make(chan string, 102400)
	go func() {
		for {
			groupTicket <- struct{}{}
			time.Sleep(3000 * time.Millisecond)
		}
	}()
	var privateTicket = make(chan struct{})
	var privateMsgQueue = make(chan string, 102400)
	go func() {
		for {
			privateTicket <- struct{}{}
			time.Sleep(2800 * time.Millisecond)
		}
	}()

	zero.OnCommand("地狱大军").
		Handle(func(ctx *zero.Ctx) {
			if ctx.Event.GroupID == 0 {
				ctx.Send("地狱大军只能在群聊中行动！")
				return
			}
			if !slices.Contains(groupWhiteList, ctx.Event.GroupID) {
				fmt.Println("group not in white list")
				return
			}
			gameLock.Lock()
			defer gameLock.Unlock()
			if currentGroupId != 0 {
				if gameStarted {
					ctx.Send("地狱大军正在别处行动，无法同时出动！")
				}
				return
			}
			if sendPrivateMsg {
				go ctx.SendPrivateMessage(ctx.Event.UserID, message.Message{
					message.Text("地狱大军即将出动！"),
				})
			}
			if sendGroupWelcome {
				go ctx.SendChain(message.At(ctx.Event.UserID), message.Text(
					"地狱大军即将出动！\n其他玩家可输入 #加入 #退出 参与游戏。"),
				)
				// groupmsghelper.SendText("地狱大军即将出动！\n其他玩家可输入：\n #加入  加入游戏\n #退出  退出游戏")
			}
			currentGroupId = ctx.Event.GroupID
			currentPlayerIds = []int64{ctx.Event.UserID}
			currentPlayerNames = []string{ctx.Event.Sender.Name()}
		})
	zero.OnCommand("加入").
		Handle(func(ctx *zero.Ctx) {
			gameLock.Lock()
			defer gameLock.Unlock()
			if currentGroupId == 0 {
				return
			}
			if currentGroupId != ctx.Event.GroupID {
				return
			}
			if gameStarted {
				return
			}
			if slices.Contains(currentPlayerIds, ctx.Event.UserID) {
				return
			}
			if slices.Contains(currentWatchingPlayerIds, ctx.Event.UserID) {
				currentWatchingPlayerIds = append(currentWatchingPlayerIds[:slices.Index(currentWatchingPlayerIds, ctx.Event.UserID)], currentWatchingPlayerIds[slices.Index(currentWatchingPlayerIds, ctx.Event.UserID)+1:]...)
			}
			if sendGroupWelcome {
				// groupmsghelper.SendText("地狱大军即将出动！")
				go ctx.SendChain(message.At(ctx.Event.UserID), message.Text("地狱大军即将出动！"))
			}
			if sendPrivateMsg {
				go ctx.SendPrivateMessage(ctx.Event.UserID, message.Message{
					message.Text("地狱大军即将出动！"),
				})
			}
			currentPlayerIds = append(currentPlayerIds, ctx.Event.UserID)
			currentPlayerNames = append(currentPlayerNames, ctx.Event.Sender.Name())
		})
	zero.OnCommand("围观").Handle(func(ctx *zero.Ctx) {
		gameLock.Lock()
		defer gameLock.Unlock()
		if currentGroupId == 0 {
			return
		}
		if currentGroupId != ctx.Event.GroupID {
			return
		}
		if slices.Contains(currentPlayerIds, ctx.Event.UserID) {
			return
		}
		if slices.Contains(currentWatchingPlayerIds, ctx.Event.UserID) {
			return
		}
		currentWatchingPlayerIds = append(currentWatchingPlayerIds, ctx.Event.UserID)
	})
	quit := func(ctx *zero.Ctx, playerId int64) {
		if slices.Contains(currentWatchingPlayerIds, playerId) {
			index := slices.Index(currentWatchingPlayerIds, playerId)
			currentWatchingPlayerIds = append(currentWatchingPlayerIds[:index], currentWatchingPlayerIds[index+1:]...)
			return
		}
		if gameStarted {
			if slices.Contains(currentPlayerIds, playerId) {
				index := slices.Index(currentPlayerIds, playerId)
				if game.PlayerLeave(index) {
					gameStarted = false
					currentGroupId = 0
					return
				}
			}
			return
		}
		if slices.Contains(currentPlayerIds, playerId) {
			index := slices.Index(currentPlayerIds, playerId)
			currentPlayerIds = append(currentPlayerIds[:index], currentPlayerIds[index+1:]...)
			currentPlayerNames = append(currentPlayerNames[:index], currentPlayerNames[index+1:]...)
			if sendGroupWelcome {
				// groupmsghelper.SendText("已退出地狱大军！")
				go ctx.SendChain(message.At(playerId), message.Text("已退出地狱大军！"))
			}
			if sendPrivateMsg {
				go ctx.SendPrivateMessage(playerId, message.Message{
					message.Text("已退出地狱大军！"),
				})
			}
			if len(currentPlayerIds) == 0 {
				gameStarted = false
				currentGroupId = 0
			}
		}
	}
	zero.OnCommand("退出").
		Handle(func(ctx *zero.Ctx) {
			gameLock.Lock()
			defer gameLock.Unlock()
			quit(ctx, ctx.Event.UserID)
		})
	zero.OnCommand("开始").
		Handle(func(ctx *zero.Ctx) {
			gameLock.Lock()
			defer gameLock.Unlock()
			if currentGroupId == 0 {
				return
			}
			if currentGroupId != ctx.Event.GroupID {
				return
			}
			if !slices.Contains(currentPlayerIds, ctx.Event.UserID) {
				return
			}
			if gameStarted {
				return
			}
			gameStarted = true
			game = New(len(currentPlayerIds))
			game.PrintFunc = func(msg string) {
				fmt.Println(msg)
				msg = strings.TrimSpace(msg)
				if msg == "" {
					return
				}
				if sendGroupMsg {
					groupMsgQueue <- msg
					go func() {
						<-groupTicket
						ctx.Send(<-groupMsgQueue)
						// groupmsghelper.SendText(<-groupMsgQueue)
					}()
				}
				if sendPrivateMsg {
					privateMsgQueue <- msg
					go func() {
						<-privateTicket
						msg := <-privateMsgQueue
						for _, id := range currentPlayerIds {
							ctx.SendPrivateMessage(id, msg)
							time.Sleep(300 * time.Millisecond)
						}
						for _, id := range currentWatchingPlayerIds {
							ctx.SendPrivateMessage(id, msg)
							time.Sleep(300 * time.Millisecond)
						}
					}()
				}
			}
			for index, name := range currentPlayerNames {
				game.SetName(index, name)
			}
			game.Start()
		})
	zero.OnCommand("帮助").
		Handle(func(ctx *zero.Ctx) {
			ctx.Send(message.Text("地狱大军游戏帮助：\n游戏目标是通过拍卖地狱随从获取能力，并逐渐解锁更强大的地狱生物。首先招募三个BOSS的玩家取胜。\n详细帮助请查看单挑联赛群文件。\n\n相关指令：\n#加入  #退出  #围观  #开始"))
			// ctx.SendGroupForwardMessage(ctx.Event.GroupID, message.Message{
			// 	message.CustomNode("天下缟素", 2700582117, []message.MessageSegment{
			// 		message.Text("地狱大军游戏帮助："),
			// 	}),
			// 	message.CustomNode("dva", 2446629225, []message.MessageSegment{
			// 		message.Text("游戏目标是通过拍卖地狱随从获取能力，并逐渐解锁更强大的地狱生物。首先招募三个BOSS的玩家取胜。"),
			// 	}),
			// 	message.CustomNode("睦月mutsuki", 3182618911, []message.MessageSegment{
			// 		message.Text("详细游戏规则请在《单挑联赛》群文件查看。"),
			// 	}),
			// 	message.CustomNode("含墨", 2154799006, []message.MessageSegment{
			// 		message.Text("可用操作：\n#加入\n#退出\n#开始\n#帮助"),
			// 	}),
			// })
		})
	zero.OnMessage().
		Handle(func(ctx *zero.Ctx) {
			if ctx.Event.GroupID != currentGroupId {
				return
			}
			if !gameStarted {
				return
			}
			if strings.HasPrefix(ctx.Event.RawMessage, "#帮他退出") {
				for _, message := range ctx.Event.Message {
					if message.Type == "at" {
						qq, ok := message.Data["qq"]
						if !ok {
							continue
						}
						userId, err := strconv.ParseInt(qq, 10, 64)
						if err != nil {
							fmt.Printf("帮他退出失败: 无法解析 %s\n", qq)
							continue
						}
						quit(ctx, userId)
					}
				}
				return
			}
			if strings.HasPrefix(ctx.Event.RawMessage, "#赛况") {
				status := fmt.Sprintf("游戏进行中，当前第%d回合。\n\n", game.Turn)

				ctx.Send(status)
				return
			}

			id := slices.Index(currentPlayerIds, ctx.Event.UserID)
			if id == -1 {
				return
			}
			fmt.Printf("id: %d group_msg: %v\n", id, ctx.Event.Message)

			// handle public message here.
			gameLock.Lock()
			defer gameLock.Unlock()
		})
	zero.OnMessage().
		Handle(func(ctx *zero.Ctx) {
			if ctx.Event.GroupID != 0 {
				return
			}
			if !gameStarted {
				return
			}
			id := slices.Index(currentPlayerIds, ctx.Event.UserID)
			if id == -1 {
				return
			}
			fmt.Printf("id: %d private_msg: %v\n", id, ctx.Event.Message)

			// handle private message here.
			if ctx.Event.RawMessage == "接受试炼" || ctx.Event.RawMessage == "通过试炼" {
				if err := game.AcceptTrial(id); err != nil {
					ctx.Send(err.Error())
				}
				return
			}

			if game.WaitResponsePlayerId == id {
				game.GiveResponse(id, ctx.Event.RawMessage)
				return
			}

			gameLock.Lock()
			defer gameLock.Unlock()
			if game.SingleMode || game.CurrentBiddingEntity2.Name == "" {
				if !game.CurrentPlayerReady[id] {
					price, err := strconv.Atoi(ctx.Event.RawMessage)
					if err != nil {
						fmt.Println("waiting for number, got ", ctx.Event.RawMessage)
						return
					}
					if err := game.GivePrice(id, price); err != nil {
						ctx.Send(err.Error())
					}
				}
			} else {
				if !game.CurrentPlayerReady[id] {
					var price1, price2 int
					_, err := fmt.Sscanf(ctx.Event.RawMessage, "%d %d", &price1, &price2)
					if err != nil {
						fmt.Println("waiting for 2 number, got ", ctx.Event.RawMessage)
						return
					}
					if err := game.GivePrices(id, price1, price2); err != nil {
						ctx.Send(err.Error())
					} else {
						ctx.Send("出价成功。")
					}
				}
			}

			if scores := game.GetScores(); scores != nil {
				gameStarted = false
				currentGroupId = 0
				currentPlayerIds = nil
				currentWatchingPlayerIds = nil
			}
		})
}
