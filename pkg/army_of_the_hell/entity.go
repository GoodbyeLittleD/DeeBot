package army_of_the_hell

import (
	"fmt"

	"golang.org/x/exp/slices"
)

type Entity struct {
	Name string
	Tags []string

	EndTurnCredits int

	BidAffectValue float64
	BidAffectTags  []string

	BidChecker    func(currentEntityList []*Entity, currentBid int) error
	UnlockChecker func(game *Game) bool

	UnlockEntities []*Entity

	Desc string
}

var (
	尸体发火 = Entity{
		Name:           "尸体发火",
		Tags:           []string{"all"},
		EndTurnCredits: 2,
		UnlockEntities: []*Entity{&毕须博须, &罗达门特},
		Desc: `
效果：每回合额外给予所属玩家1能力点。
特殊：如果流拍也正常解锁后续
解锁：毕须博须、ACT2-罗达门特`,
	}
	毕须博须 = Entity{
		Name:           "毕须博须",
		Tags:           []string{"all"},
		BidAffectTags:  []string{"andariel"},
		BidAffectValue: -0.05,
		EndTurnCredits: 1,
		UnlockEntities: []*Entity{&血鸟, &拉卡尼休, &冰冷乌鸦},
		Desc: `
效果：安达利尔出价-5%。
特殊：无
解锁：血鸟、拉卡尼休、冰冷乌鸦`,
	}
	血鸟 = Entity{
		Name:           "血鸟",
		Tags:           []string{"all"},
		BidAffectTags:  []string{"andariel"},
		BidAffectValue: -0.15,
		EndTurnCredits: 1,
		UnlockEntities: []*Entity{&碎骨者},
		Desc: `
效果：安达利尔出价-15%。
特殊：每回合额外给予所属玩家2能力点，不可与ACT2-罗达门特，ACT3-吉得宾偷窃者，ACT5-督军山克效果重叠，若同时持有，则只取效果最高的一个。
解锁：碎骨者`,
	}
	碎骨者 = Entity{
		Name:           "碎骨者",
		Tags:           []string{"all"},
		EndTurnCredits: 1,
		Desc: `
无特殊效果。`,
	}
	冰冷乌鸦 = Entity{
		Name:           "冰冷乌鸦",
		Tags:           []string{"all"},
		EndTurnCredits: 1,
		Desc: `
无特殊效果。`,
	}
	拉卡尼休 = Entity{
		Name:           "拉卡尼休",
		Tags:           []string{"all"},
		EndTurnCredits: 1,
		UnlockEntities: []*Entity{&树头木拳, &格瑞斯华尔德},
		Desc: `
效果：无
特殊：如果同时拥有拉卡尼休和树头木拳，则解锁格瑞斯华尔德
解锁：树头木拳`,
	}
	树头木拳 = Entity{
		Name:           "树头木拳",
		Tags:           []string{"all"},
		EndTurnCredits: 1,
		UnlockEntities: []*Entity{&女伯爵, &洞穴重生的邪恶之犬, &铁匠, &骨灰, &格瑞斯华尔德},
		Desc: `
效果：无
特殊：如果同时拥有拉卡尼休和树头木拳，则解锁格瑞斯华尔德
解锁：女伯爵、洞穴重生的邪恶之犬、铁匠、骨灰`,
	}
	格瑞斯华尔德 = Entity{
		Name:           "格瑞斯华尔德",
		Tags:           []string{"all"},
		EndTurnCredits: 5,
		Desc: `
效果：每回合额外给予所属玩家4能力点。
特殊：只有拉卡尼休和树头木拳的拥有者可以竞拍，但是底价为10。
解锁：无`,
	}
	洞穴重生的邪恶之犬 = Entity{
		Name:           "洞穴重生的邪恶之犬",
		Tags:           []string{"all"},
		EndTurnCredits: 1,
		Desc: `
效果：无
特殊：当第一次被抽到时，跳过本轮所有竞拍。
解锁：无`,
	}
	女伯爵 = Entity{
		Name:           "女伯爵",
		Tags:           []string{"all"},
		BidAffectValue: -0.05,
		BidAffectTags:  []string{"BOSS"},
		EndTurnCredits: 1,
		UnlockEntities: []*Entity{&安达利尔},
		Desc: `
效果：安达利尔出价-15%，其他四个boss出价-5%
特殊：如果铁匠或者骨灰之一已被拍得，解锁安达利尔
解锁：无`,
	}
	铁匠 = Entity{
		Name:           "铁匠",
		Tags:           []string{"all"},
		BidAffectValue: -0.15,
		BidAffectTags:  []string{"andariel"},
		EndTurnCredits: 1,
		UnlockEntities: []*Entity{&安达利尔},
		Desc: `
效果：安达利尔出价-15%。
特殊：如果女伯爵或者骨灰之一已被拍得，解锁安达利尔
解锁：无`,
	}
	骨灰 = Entity{
		Name:           "骨灰",
		Tags:           []string{"all"},
		BidAffectValue: -0.1,
		BidAffectTags:  []string{"andariel"},
		EndTurnCredits: 1,
		UnlockEntities: []*Entity{&安达利尔},
		Desc: `
效果：安达利尔出价-10%。
特殊：如果女伯爵或者铁匠之一已被拍得，解锁安达利尔
解锁：无`,
	}
	安达利尔 = Entity{
		Name:           "安达利尔",
		Tags:           []string{"all", "andariel", "BOSS"},
		BidAffectValue: 0.2,
		BidAffectTags:  []string{"BOSS"},
		EndTurnCredits: 5,
		Desc: `
【BOSS】
效果：其他boss出价+20%，每回合给予所属玩家5能力点。`,
	}
	罗达门特 = Entity{
		Name:           "罗达门特",
		Tags:           []string{"all", "act2"},
		BidAffectValue: -0.1,
		BidAffectTags:  []string{"durial"},
		EndTurnCredits: 1,
		UnlockEntities: []*Entity{&爬行的容貌, &疯狂血腥女巫, &爆开的甲虫, &钻地的冰虫, &黑暗长老, &牙皮, &燃烧者韦布},
		Desc: `
效果：督瑞尔出价-10%。
特殊：每回合额外给予所属玩家3能力点，不可与ACT1-血鸟，ACT3-吉得宾偷窃者，ACT5-督军山克效果重叠，若同时持有，则只取效果最高的一个。
解锁：【沙漠六小队】、ACT3-燃烧者韦布`,
	}
	爬行的容貌 = Entity{
		Name:           "爬行的容貌",
		Tags:           []string{"all", "act2"},
		EndTurnCredits: 1,
		UnlockEntities: []*Entity{&火之眼},
		Desc: `
【沙漠六小队】
无特殊效果。`,
	}
	疯狂血腥女巫 = Entity{
		Name:           "疯狂血腥女巫",
		Tags:           []string{"all", "act2"},
		BidAffectValue: -0.1,
		BidAffectTags:  []string{"all"},
		EndTurnCredits: 1,
		UnlockEntities: []*Entity{&火之眼},
		Desc: `
【沙漠六小队】
效果：所有怪物出价-10%。
特殊：无
解锁：无`,
	}
	爆开的甲虫 = Entity{
		Name:           "爆开的甲虫",
		Tags:           []string{"all", "act2"},
		EndTurnCredits: 1,
		UnlockEntities: []*Entity{&火之眼},
		Desc: `
【沙漠六小队】
无特殊效果。`,
	}
	钻地的冰虫 = Entity{
		Name:           "钻地的冰虫",
		Tags:           []string{"all", "act2"},
		EndTurnCredits: 1,
		BidAffectValue: -0.2,
		BidAffectTags:  []string{"durial"},
		UnlockEntities: []*Entity{&火之眼},
		Desc: `
【沙漠六小队】
效果：督瑞尔出价-20%。
特殊：无
解锁：无`,
	}
	黑暗长老 = Entity{
		Name:           "黑暗长老",
		Tags:           []string{"all", "act2"},
		EndTurnCredits: 1,
		BidAffectValue: -0.1,
		BidAffectTags:  []string{"durial"},
		UnlockEntities: []*Entity{&火之眼},
		Desc: `
【沙漠六小队】
效果：督瑞尔出价-10%。
特殊：抽到时触发【天黑】，【天黑】效果为：接下来三回合双方要在不知道抽到什么的情况下竞拍。
解锁：无`,
	}
	牙皮 = Entity{
		Name:           "牙皮",
		Tags:           []string{"all", "act2"},
		EndTurnCredits: 1,
		BidAffectValue: -0.1,
		BidAffectTags:  []string{"durial"},
		UnlockEntities: []*Entity{&火之眼},
		Desc: `
【沙漠六小队】
效果：督瑞尔出价-10%。
特殊：当牙皮被抽到时，本局对决之后的时间将不会触发【天黑】，如当前处于【天黑】，则立即解除。
解锁：无`,
	}
	火之眼 = Entity{
		Name:           "火之眼",
		Tags:           []string{"all", "act2"},
		EndTurnCredits: 1,
		UnlockEntities: []*Entity{&召唤者},
		Desc: `
效果：无
特殊：如果两回合内抽到赫拉森（即召唤者），只有拍得【火之眼】者可以竞拍【召唤者】
解锁：召唤者`,
	}
	召唤者 = Entity{
		Name:           "召唤者",
		Tags:           []string{"all", "act2"},
		BidAffectValue: -0.1,
		BidAffectTags:  []string{"durial"},
		EndTurnCredits: 1,
		UnlockEntities: []*Entity{&督瑞尔, &古代无魂之卡},
		Desc: `
效果：督瑞尔出价-10%。
特殊：如果两回合内抽到督瑞尔，只有拍得者可以竞拍
解锁：督瑞尔、古代无魂之卡`,
	}
	古代无魂之卡 = Entity{
		Name:           "古代无魂之卡",
		Tags:           []string{"all"},
		BidAffectValue: -0.4,
		BidAffectTags:  []string{"act2"},
		EndTurnCredits: 1,
		Desc: `
效果：ACT2其他非boss出价-40%
特殊：拍到一方可选择移除督瑞尔，六回合后跳过抽牌，竞拍督瑞尔。【注：即使督瑞尔已被拍得，此特殊效果依然可以启动。】
解锁：无`,
	}
	督瑞尔 = Entity{
		Name:           "督瑞尔",
		Tags:           []string{"all", "durial", "BOSS"},
		EndTurnCredits: 8,
		BidAffectValue: 0.3,
		BidAffectTags:  []string{"BOSS"},
		Desc: `
【BOSS】
效果：其他boss出价+30%，每回合给予所属玩家8能力点。`,
	}
	燃烧者韦布 = Entity{
		Name:           "燃烧者韦布",
		Tags:           []string{"all"},
		BidAffectValue: -0.2,
		BidAffectTags:  []string{"parliament"},
		EndTurnCredits: 1,
		UnlockEntities: []*Entity{&吉得宾偷窃者, &古巫医印都, &暴风之树, &衣卒尔},
		Desc: `
效果：所有【议会成员】出价-20%
特殊：无
解锁：吉得宾偷窃者、古巫医印都、暴风之树、ACT4-衣卒尔。`,
	}
	吉得宾偷窃者 = Entity{
		Name:           "吉得宾偷窃者",
		Tags:           []string{"all"},
		EndTurnCredits: 1,
		Desc: `
效果：无
特殊：每回合额外给予所属玩家4能力点，不可与ACT1-血鸟，ACT2-罗达门特，ACT5-督军山克效果重叠，若同时持有，则只取效果最高的一个。
解锁：无`,
	}
	古巫医印都 = Entity{
		Name:           "古巫医印都",
		Tags:           []string{"all"},
		BidAffectValue: -0.2,
		BidAffectTags:  []string{"parliament"},
		EndTurnCredits: 1,
		UnlockEntities: []*Entity{&邪恶之手伊斯梅尔, &火焰之指吉列布, &冰拳托克},
		Desc: `
效果：所有议会成员出价-20%
特殊：古巫医印都和裂缝之翼冰鹰都被抽取后，解锁【议会成员第一组】
解锁：无`,
	}
	暴风之树 = Entity{
		Name:           "暴风之树",
		Tags:           []string{"all"},
		EndTurnCredits: 1,
		UnlockEntities: []*Entity{&战场处子沙丽娜, &裂缝之翼冰鹰},
		Desc: `
效果：无
特殊：无
解锁：战场处子沙丽娜、裂缝之翼冰鹰`,
	}
	战场处子沙丽娜 = Entity{
		Name:           "战场处子沙丽娜",
		Tags:           []string{"all"},
		EndTurnCredits: 3,
		Desc: `
效果：每回合额外给予所属玩家2能力点。
特殊：无
解锁：无`,
	}
	裂缝之翼冰鹰 = Entity{
		Name:           "裂缝之翼冰鹰",
		Tags:           []string{"all"},
		BidAffectValue: -0.2,
		BidAffectTags:  []string{"parliament"},
		EndTurnCredits: 1,
		UnlockEntities: []*Entity{&邪恶之手伊斯梅尔, &火焰之指吉列布, &冰拳托克},
		Desc: `
效果：所有议会成员出价-20%
特殊：古巫医印都和裂缝之翼冰鹰都被抽取后，解锁【议会成员第一组】
解锁：无`,
	}
	邪恶之手伊斯梅尔 = Entity{
		Name:           "邪恶之手伊斯梅尔",
		Tags:           []string{"all", "parliament"},
		EndTurnCredits: 1,
		BidAffectValue: -0.2,
		BidAffectTags:  []string{"mofisite"},
		UnlockEntities: []*Entity{&火花之拳布瑞姆},
		Desc: `
效果：墨菲斯托出价-20%
特殊：下回合【优先抽取】另一个第一组议会成员
解锁：无`,
	}
	火焰之指吉列布 = Entity{
		Name:           "火焰之指吉列布",
		Tags:           []string{"all", "parliament"},
		EndTurnCredits: 1,
		BidAffectValue: -0.2,
		BidAffectTags:  []string{"mofisite"},
		UnlockEntities: []*Entity{&火花之拳布瑞姆},
		Desc: `
效果：墨菲斯托出价-20%
特殊：下回合【优先抽取】另一个第一组议会成员
解锁：无`,
	}
	冰拳托克 = Entity{
		Name:           "冰拳托克",
		Tags:           []string{"all", "parliament"},
		EndTurnCredits: 1,
		BidAffectValue: -0.2,
		BidAffectTags:  []string{"mofisite"},
		UnlockEntities: []*Entity{&火花之拳布瑞姆},
		Desc: `
效果：墨菲斯托出价-20%
特殊：下回合【优先抽取】另一个第一组议会成员
解锁：无`,
	}
	火花之拳布瑞姆 = Entity{
		Name:           "火花之拳布瑞姆",
		Tags:           []string{"all", "parliament"},
		EndTurnCredits: 1,
		BidAffectValue: -0.2,
		BidAffectTags:  []string{"mofisite"},
		UnlockEntities: []*Entity{&空虚之指维恩, &龙手马弗},
		Desc: `
效果：墨菲斯托出价-20%
特殊：无
解锁：空虚之指维恩、龙手马弗`,
	}
	空虚之指维恩 = Entity{
		Name:           "空虚之指维恩",
		Tags:           []string{"all", "parliament"},
		EndTurnCredits: 1,
		BidAffectValue: -0.2,
		BidAffectTags:  []string{"mofisite"},
		UnlockEntities: []*Entity{&墨菲斯托},
		Desc: `
效果：墨菲斯托出价-20%
特殊：无
解锁：墨菲斯托`,
	}
	龙手马弗 = Entity{
		Name:           "龙手马弗",
		Tags:           []string{"all", "parliament"},
		EndTurnCredits: 1,
		BidAffectValue: -0.2,
		BidAffectTags:  []string{"mofisite"},
		UnlockEntities: []*Entity{&墨菲斯托},
		Desc: `
效果：墨菲斯托出价-20%
特殊：无
解锁：墨菲斯托`,
	}
	墨菲斯托 = Entity{
		Name:           "墨菲斯托",
		Tags:           []string{"all", "mofisite", "BOSS"},
		EndTurnCredits: 10,
		BidAffectValue: 0.5,
		BidAffectTags:  []string{"BOSS"},
		Desc: `
【BOSS】
效果：其他boss出价+50%，每回合给予所属玩家10能力点。`,
	}
	衣卒尔 = Entity{
		Name:           "衣卒尔",
		Tags:           []string{"all"},
		BidAffectValue: -0.2,
		BidAffectTags:  []string{"all"},
		EndTurnCredits: 1,
		UnlockEntities: []*Entity{&海法斯特盔甲制作者},
		Desc: `
效果：所有出价-20%
特殊：拍得【衣卒尔】的玩家，从下一回合起，每5个回合必须至少出价一次0能力点。
解锁：海法斯特盔甲制作者`,
	}
	海法斯特盔甲制作者 = Entity{
		Name:           "海法斯特盔甲制作者",
		Tags:           []string{"all"},
		EndTurnCredits: 5,
		UnlockEntities: []*Entity{&宏伟的混沌大臣, &西希之王, &灵魂传播者},
		Desc: `
效果：每回合给予所属玩家5能力点
特殊：第一次抽到该牌时，改为抽到ACT5-达克法恩
解锁：宏伟的混沌大臣、西希之王、灵魂传播者`,
	}
	宏伟的混沌大臣 = Entity{
		Name:           "宏伟的混沌大臣",
		Tags:           []string{"all"},
		BidAffectValue: -0.35,
		BidAffectTags:  []string{"diabolo"},
		EndTurnCredits: 1,
		UnlockEntities: []*Entity{&迪亚波罗},
		Desc: `
效果：迪亚波罗出价-35%
特殊：宏伟的混沌大臣、西希之王、灵魂传播者全被抽到后解锁迪亚波罗
解锁：无`,
	}
	西希之王 = Entity{
		Name:           "西希之王",
		Tags:           []string{"all"},
		BidAffectValue: -0.5,
		BidAffectTags:  []string{"diabolo"},
		EndTurnCredits: 1,
		UnlockEntities: []*Entity{&迪亚波罗},
		Desc: `
效果：迪亚波罗出价-50%
特殊：宏伟的混沌大臣、西希之王、灵魂传播者全被抽到后解锁迪亚波罗，拍得【西希之王】者下回合不能出价除boss以外的其他怪物。
解锁：无`,
	}
	灵魂传播者 = Entity{
		Name:           "灵魂传播者",
		Tags:           []string{"all"},
		BidAffectValue: -0.35,
		BidAffectTags:  []string{"diabolo"},
		EndTurnCredits: 1,
		UnlockEntities: []*Entity{&迪亚波罗},
		Desc: `
效果：迪亚波罗出价-35%
特殊：宏伟的混沌大臣、西希之王、灵魂传播者全被抽到后解锁迪亚波罗
解锁：无`,
	}
	迪亚波罗 = Entity{
		Name:           "迪亚波罗",
		Tags:           []string{"all", "BOSS", "diabolo"},
		BidAffectValue: 0.5,
		BidAffectTags:  []string{"BOSS"},
		EndTurnCredits: 10,
		Desc: `
【BOSS】
效果：其他boss出价+50%，每回合给予所属玩家10能力点。`,
	}

	达克法恩 = Entity{
		Name:           "达克法恩",
		Tags:           []string{"all", "act5"},
		EndTurnCredits: 1,
		UnlockEntities: []*Entity{&督军山克},
		Desc: `
效果：无
特殊：拍得者在拍得时可以选择跳过下次抽牌，改为抽到督军山克。
解锁：督军山克`,
	}
	督军山克 = Entity{
		Name:           "督军山克",
		Tags:           []string{"all", "act5"},
		EndTurnCredits: 1,
		BidAffectValue: -0.1,
		BidAffectTags:  []string{"act5"},
		UnlockEntities: []*Entity{&矫正者怪异},
		Desc: `
效果：ACT5所有非boss出价-10%
特殊1：每回合额外给予所属玩家5能力点，不可与ACT1-血鸟，ACT2-罗达门特，ACT3-吉得宾偷窃者效果重叠，若同时持有，则只取效果最高的一个。
特殊2：接下来第三次有效抽牌抽到矫正者怪异，已经抽到则失效。
解锁：矫正者怪异`,
	}
	矫正者怪异 = Entity{
		Name:           "矫正者怪异",
		Tags:           []string{"all", "act5"},
		EndTurnCredits: 1,
		UnlockEntities: []*Entity{&剥壳凹槽, &利牙杀手, &狂暴者眼魔},
		Desc: `
效果：无
特殊：无
解锁：利牙杀手、狂暴者眼魔、剥壳凹槽`,
	}
	利牙杀手 = Entity{
		Name:           "利牙杀手",
		Tags:           []string{"all", "act5"},
		EndTurnCredits: 1,
		BidAffectValue: -0.1,
		BidAffectTags:  []string{"act5"},
		Desc: `
效果：ACT5所有非boss出价-10%
特殊：无
解锁：无`,
	}
	狂暴者眼魔 = Entity{
		Name:           "狂暴者眼魔",
		Tags:           []string{"all", "act5"},
		EndTurnCredits: 1,
		Desc: `
无效果特殊解锁。`,
	}
	剥壳凹槽 = Entity{
		Name:           "剥壳凹槽",
		Tags:           []string{"all", "act5"},
		EndTurnCredits: 1,
		UnlockEntities: []*Entity{&骨据破坏者, &粉碎者, &冰冻魔怪, &三个野蛮人},
		Desc: `
效果：无
特殊：无
解锁：冰冻魔怪、骨据破坏者、粉碎者、三个野蛮人`,
	}
	三个野蛮人 = Entity{
		Name: "三个野蛮人",
		Desc: `
特殊：此牌被抽出时不进行拍卖，改为使双方玩家进入试炼，当每回合拍卖结束后到下一回合开始前，玩家可选择给出50能力点通过试炼，当一方通过试炼时，解锁【破坏者卡兰索】。
解锁：无`,
	}
	骨据破坏者 = Entity{
		Name:           "骨据破坏者",
		Tags:           []string{"all", "act5"},
		EndTurnCredits: 1,
		Desc: `
无效果特殊解锁。`,
	}
	粉碎者 = Entity{
		Name:           "粉碎者",
		Tags:           []string{"all", "act5"},
		EndTurnCredits: 1,
		Desc: `
无效果特殊解锁。`,
	}
	冰冻魔怪 = Entity{
		Name:           "冰冻魔怪",
		Tags:           []string{"all", "act5"},
		EndTurnCredits: 3,
		UnlockEntities: []*Entity{&暴躁外皮},
		Desc: `
效果：每回合额外给予所属玩家2能力点
特殊：无
解锁：暴躁外皮`,
	}
	暴躁外皮 = Entity{
		Name:           "暴躁外皮",
		Tags:           []string{"all", "act5"},
		EndTurnCredits: 2,
		UnlockEntities: []*Entity{&尼拉塞克},
		Desc: `
效果：每回合给予所属玩家2能力点
特殊：下回合结束时，将一个暴躁外皮洗入牌库
解锁：尼拉塞克`,
	}
	暴躁外皮2 = Entity{
		Name:           "暴躁外皮",
		Tags:           []string{"all", "act5"},
		EndTurnCredits: 2,
		Desc: `
效果：每回合给予所属玩家2能力点
特殊：下回合结束时，将一个暴躁外皮洗入牌库`,
	}
	尼拉塞克 = Entity{
		Name:           "尼拉塞克",
		Tags:           []string{"all", "act5"},
		EndTurnCredits: 4,
		UnlockEntities: []*Entity{},
		Desc: `
效果：每回合额外给予所属玩家3能力点
特殊：可以直接通过三个野蛮人试炼
解锁：无`,
	}
	破坏者卡兰索 = Entity{
		Name:           "破坏者卡兰索",
		Tags:           []string{"all", "act5"},
		BidAffectValue: -0.15,
		BidAffectTags:  []string{"bar"},
		EndTurnCredits: 1,
		UnlockEntities: []*Entity{&诅咒的阿克姆},
		Desc: `
效果：巴尔出价-15%
特殊1：未通过三个野蛮人不得竞拍
特殊2：下回合【优先抽取】诅咒的阿克姆
解锁：诅咒的阿克姆`,
	}
	诅咒的阿克姆 = Entity{
		Name:           "诅咒的阿克姆",
		Tags:           []string{"all", "act5"},
		BidAffectValue: -0.2,
		BidAffectTags:  []string{"bar"},
		EndTurnCredits: 1,
		UnlockEntities: []*Entity{&血腥的巴特克},
		Desc: `
效果：巴尔出价-20%
特殊1：未通过三个野蛮人不得竞拍
特殊2：下回合【优先抽取】血腥的巴特克
解锁：血腥的巴特克`,
	}
	血腥的巴特克 = Entity{
		Name:           "血腥的巴特克",
		Tags:           []string{"all", "act5"},
		BidAffectValue: -0.25,
		BidAffectTags:  []string{"bar"},
		EndTurnCredits: 1,
		UnlockEntities: []*Entity{&不洁的凡塔},
		Desc: `
效果：巴尔出价-25%
特殊1：未通过三个野蛮人不得竞拍
特殊2：下回合【优先抽取】不洁的凡塔
解锁：不洁的凡塔`,
	}
	不洁的凡塔 = Entity{
		Name:           "不洁的凡塔",
		Tags:           []string{"all", "act5"},
		BidAffectValue: -0.3,
		BidAffectTags:  []string{"bar"},
		EndTurnCredits: 1,
		UnlockEntities: []*Entity{&古难记录者},
		Desc: `
效果：巴尔出价-30%
特殊1：未通过三个野蛮人不得竞拍
特殊2：下回合【优先抽取】古难记录者
解锁：古难记录者`,
	}
	古难记录者 = Entity{
		Name:           "古难记录者",
		Tags:           []string{"all", "act5"},
		BidAffectValue: -0.35,
		BidAffectTags:  []string{"bar"},
		EndTurnCredits: 1,
		UnlockEntities: []*Entity{&巴尔},
		Desc: `
效果：巴尔出价-35%
特殊1：未通过三个野蛮人不得竞拍
特殊2：下回合【优先抽取】巴尔
解锁：巴尔`,
	}
	巴尔 = Entity{
		Name:           "巴尔",
		Tags:           []string{"all", "bar", "BOSS"},
		BidAffectValue: 0.5,
		BidAffectTags:  []string{"BOSS"},
		EndTurnCredits: 10,
		Desc: `
【BOSS】
效果：其他boss出价+50%，每回合给予所属玩家10能力点。`,
	}
	凑数的 = Entity{
		Name: "",
	}
)

func init() {
	格瑞斯华尔德.UnlockChecker = func(game *Game) bool {
		for _, entities := range game.PlayerEntities {
			if slices.Contains(entities, &拉卡尼休) && slices.Contains(entities, &树头木拳) {
				return true
			}
		}
		return false
	}
	格瑞斯华尔德.BidChecker = func(currentEntityList []*Entity, currentBid int) error {
		if currentBid < 10 && currentBid != 0 {
			return fmt.Errorf("底价为10")
		}
		return nil
	}
	安达利尔.UnlockChecker = func(game *Game) bool {
		count := 0
		for _, entities := range game.PlayerEntities {
			if slices.Contains(entities, &女伯爵) {
				count += 1
			}
			if slices.Contains(entities, &铁匠) {
				count += 1
			}
			if slices.Contains(entities, &骨灰) {
				count += 1
			}
		}
		if count == 2 {
			game.PrintFunc("安达利尔已解锁。")
			return true
		}
		return false
	}
	火之眼.UnlockChecker = func(game *Game) bool {
		if game.火之眼_unlocked {
			return false
		}
		for _, entities := range game.PlayerEntities {
			count := 0
			if slices.Contains(entities, &爬行的容貌) {
				count += 1
			}
			if slices.Contains(entities, &疯狂血腥女巫) {
				count += 1
			}
			if slices.Contains(entities, &爆开的甲虫) {
				count += 1
			}
			if slices.Contains(entities, &钻地的冰虫) {
				count += 1
			}
			if slices.Contains(entities, &黑暗长老) {
				count += 1
			}
			if slices.Contains(entities, &牙皮) {
				count += 1
			}
			if count == 3 {
				game.PrintFunc("火之眼已解锁。")
				return true
			}
		}
		// 对于多人模式的fix
		count := 0
		for _, entities := range game.PlayerEntities {
			if slices.Contains(entities, &爬行的容貌) {
				count += 1
			}
			if slices.Contains(entities, &疯狂血腥女巫) {
				count += 1
			}
			if slices.Contains(entities, &爆开的甲虫) {
				count += 1
			}
			if slices.Contains(entities, &钻地的冰虫) {
				count += 1
			}
			if slices.Contains(entities, &黑暗长老) {
				count += 1
			}
			if slices.Contains(entities, &牙皮) {
				count += 1
			}
		}
		if count == 6 {
			game.PrintFunc("由于沙漠六小队已全部被拍完，火之眼已强制解锁。")
			return true
		}
		return false
	}
	沙漠六小队 := []*Entity{
		&爬行的容貌,
		&疯狂血腥女巫,
		&爆开的甲虫,
		&钻地的冰虫,
		&黑暗长老,
		&牙皮,
	}
	沙漠六小队限制 := func(currentEntityList []*Entity, currentBid int) error {
		count := 0
		for _, entity := range currentEntityList {
			if slices.Contains(沙漠六小队, entity) {
				count += 1
			}
		}
		if count >= 3 && currentBid > 0 {
			return fmt.Errorf("拍得【沙漠六小队】的其中三个后，不得参与其他沙漠六小队的竞拍")
		}
		return nil
	}
	for _, entity := range 沙漠六小队 {
		entity.BidChecker = 沙漠六小队限制
	}

	议会第一组_UnlockChecker := func(game *Game) bool {
		count := 0
		for _, entities := range game.PlayerEntities {
			if slices.Contains(entities, &古巫医印都) {
				count += 1
			}
			if slices.Contains(entities, &裂缝之翼冰鹰) {
				count += 1
			}
		}
		if count == 2 {
			game.PrintFunc("议会成员第一组已解锁。")
			return true
		}
		return false
	}
	邪恶之手伊斯梅尔.UnlockChecker = 议会第一组_UnlockChecker
	火焰之指吉列布.UnlockChecker = 议会第一组_UnlockChecker
	冰拳托克.UnlockChecker = 议会第一组_UnlockChecker

	火花之拳布瑞姆.UnlockChecker = func(game *Game) bool {
		count := 0
		for _, entities := range game.PlayerEntities {
			if slices.Contains(entities, &邪恶之手伊斯梅尔) {
				count += 1
			}
			if slices.Contains(entities, &火焰之指吉列布) {
				count += 1
			}
			if slices.Contains(entities, &冰拳托克) {
				count += 1
			}
		}
		if count == 3 {
			game.PrintFunc("火花之拳布瑞姆已解锁。")
			return true
		}
		return false
	}

	迪亚波罗.UnlockChecker = func(game *Game) bool {
		count := 0
		for _, entities := range game.PlayerEntities {
			if slices.Contains(entities, &宏伟的混沌大臣) {
				count += 1
			}
			if slices.Contains(entities, &西希之王) {
				count += 1
			}
			if slices.Contains(entities, &灵魂传播者) {
				count += 1
			}
		}
		if count == 3 {
			game.PrintFunc("迪亚波罗已解锁。")
			return true
		}
		return false
	}

	墨菲斯托.UnlockChecker = func(game *Game) bool {
		return !game.墨菲斯托_unlocked
	}
}

func (entity *Entity) String() string {
	return entity.Name + entity.Desc
}

func (entity *Entity) IsBoss() bool {
	for _, tag := range entity.Tags {
		if tag == "BOSS" {
			return true
		}
	}
	return false
}

func (entity *Entity) BidEffected(currentEntityTags []string) bool {
	for _, tag := range entity.BidAffectTags {
		if slices.Contains(currentEntityTags, tag) {
			return true
		}
	}
	return false
}
