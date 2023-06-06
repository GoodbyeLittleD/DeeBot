package arena60

import (
	groupmsghelper "dvabot/pkg/groupmsg_helper"
	"fmt"

	zero "github.com/wdvxdr1123/ZeroBot"
	"golang.org/x/exp/slices"
)

var groupWhiteList = []int64{
	541065568,  // 测试群
	730561889,  // 测试群
	757553235,  // 一费奶五战队
	1006990128, // 单挑联赛群
}

func Handle() {
	zero.OnMessage().
		Handle(func(ctx *zero.Ctx) {
			if !slices.Contains(groupWhiteList, ctx.Event.GroupID) {
				fmt.Println("group not in white list")
				return
			}
			var a, b, c, d, e, f int
			if n, err := fmt.Sscanf(ctx.Event.RawMessage, "%d %d %d %d %d %d\n", &a, &b, &c, &d, &e, &f); err != nil || n != 6 {
				return
			}
			host := Create(a, b, c, d, e, f)
			result := fmt.Sprintf("血量：%d 攻击：%d 防御：%d 减伤：%d 增伤：%d 破防：%d\n正在对战1000个bot.....\n胜场：\n", a, b, c, d, e, f)
			for i := 0; i < 10; i++ {
				result += fmt.Sprintf("%d\n", host.FightStageOne())
			}
			groupmsghelper.SendText(result)
		})
}
