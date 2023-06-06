package arena60

import (
	"math/rand"
)

type Entity struct {
	Health         float64
	Attack         float64
	Defence        float64
	Minus          float64 // 减伤
	Plus           float64 // 增伤
	CounterDefence float64 // 破防
}

func New() *Entity {
	return &Entity{
		Health:         100,
		Attack:         10,
		Defence:        0,
		Minus:          0,
		Plus:           0,
		CounterDefence: 0.2,
	}
}

func Create(health_points int, attack_points int, defence_points int, minus_points int, plus_points int, counter_defence_points int) *Entity {
	return &Entity{
		Health:         100 + 10*float64(health_points),
		Attack:         10 + float64(attack_points),
		Defence:        float64(defence_points),
		Minus:          float64(minus_points),
		Plus:           float64(plus_points),
		CounterDefence: 0.2 + 0.2*float64(counter_defence_points),
	}
}

func Rand() *Entity {
	return &Entity{
		Health:         100 + 10*float64(rand.Intn(220)),
		Attack:         10 + float64(rand.Intn(220)),
		Defence:        0 + float64(rand.Intn(220)),
		Minus:          0 + float64(rand.Intn(220)),
		Plus:           0 + float64(rand.Intn(220)),
		CounterDefence: 0.2 + 0.2*float64(rand.Intn(220)),
	}
}

func Max(a float64, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func (host *Entity) Fight(target *Entity) int {
	host_health := host.Health
	target_health := target.Health
	for host_health > 0 && target_health > 0 {
		target_health -= Max(host.CounterDefence, host.Attack-target.Defence) * (1 + host.Plus/10) / (1 + target.Minus/10)
		host_health -= Max(target.CounterDefence, target.Attack-host.Defence) * (1 + target.Plus/10) / (1 + host.Minus/10)
	}
	if host_health > 0 {
		return 1
	}
	if target_health > 0 {
		return -1
	}
	return 0
}

func (host *Entity) FightStageOne() int {
	count := 0
	for i := 0; i < 1000; i++ {
		target := Rand()
		if host.Fight(target) == 1 {
			count++
		}
	}
	return count
}
