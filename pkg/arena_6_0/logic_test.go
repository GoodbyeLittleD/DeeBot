package arena60

import "testing"

func TestFight(t *testing.T) {
	host := New()
	target := New()
	result := host.Fight(target)
	if result != 0 {
		t.Error("result should be 0")
	}
}

func TestArena(t *testing.T) {
	host := Create(180, 0, 130, 190, 150, 300)
	total := 0
	for i := 0; i < 16; i++ {
		result := host.FightStageOne()
		total += result
		t.Logf("result: %d", result)
	}
	t.Logf("average: %d", total/16)
}

func Test1v1(t *testing.T) {
	host := Create(
		180+0,
		0,
		130+0,
		190+16,
		150+0,
		300+0,
	)
	targets := []*Entity{
		Create(203, 322, 186, 144, 144, 0),  // 滑稽怪
		Create(120, 425, 178, 120, 156, 0),  // 匿名
		Create(170, 270, 230, 170, 170, 0),  // 防御塔
		Create(191, 0, 220, 192, 192, 200),  // 睦月
		Create(180, 0, 250, 180, 180, 180),  // 日召
		Create(200, 250, 170, 200, 130, 10), // 烧饼
		Create(122, 452, 0, 191, 191, 0),    // 问号
		Create(250, 210, 180, 70, 100, 200), // 新鱼 +107
	}
	for _, target := range targets {
		result := host.Fight(target)
		t.Logf("result: %d", result)
	}
}
