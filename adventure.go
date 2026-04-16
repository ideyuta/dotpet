package main

import (
	"fmt"
	"math/rand"
	"time"
)

var monsters = []struct {
	Name      string
	BasePower int
}{
	{"スライム", 1},
	{"ゴブリン", 3},
	{"オオカミ", 5},
	{"スケルトン", 7},
	{"オーク", 10},
	{"ドラゴン", 14},
	{"デーモン", 18},
}

var specialEvents = []string{
	"回復の泉を見つけた",
	"旅の商人に出会った",
	"古代遺跡を発見した",
	"流れ星を見た",
	"隠し通路を見つけた",
}

func RollEvent() int {
	r := rand.Intn(100)
	switch {
	case r < 55:
		return 0 // nothing
	case r < 80:
		return 1 // battle
	case r < 90:
		return 2 // item
	default:
		return 3 // special
	}
}

// eventCooldown is the minimum duration between random events.
const eventCooldown = 10 * time.Minute

// Tick processes one game tick. Returns the (possibly new) pet and whether reincarnation happened.
func Tick(p *Pet) (*Pet, bool) {
	now := time.Now()

	// XP from time passage (1 XP per minute elapsed)
	elapsed := int(now.Sub(p.UpdatedAt).Minutes())
	if elapsed > 0 {
		leveled := p.AddXP(elapsed)
		if leveled {
			p.LogEvent(fmt.Sprintf("✨ Lv.%d!", p.Level))
		}
	}

	p.UpdatedAt = now

	// Check reincarnation
	if p.CanReincarnate() {
		gen := p.Generation
		p.LogEvent(fmt.Sprintf("🔄 転生! 第%d世代へ", gen+1))
		newPet := p.Reincarnate()
		newPet.EventLog = p.EventLog
		newPet.LastEvent = p.LastEvent
		newPet.EventAt = p.EventAt
		return newPet, true
	}

	// Random event (only if cooldown has passed since last event)
	if now.Sub(p.EventAt) >= eventCooldown {
		switch RollEvent() {
		case 1:
			doBattle(p)
		case 2:
			doItemFind(p)
		case 3:
			doSpecial(p)
		}
	}

	return p, false
}

func doBattle(p *Pet) {
	maxIdx := p.Level / 3
	if maxIdx >= len(monsters) {
		maxIdx = len(monsters) - 1
	}
	m := monsters[rand.Intn(maxIdx+1)]

	power := p.TotalPower()
	monsterPower := m.BasePower + rand.Intn(3)

	if power+rand.Intn(5) >= monsterPower {
		p.Wins++
		xp := m.BasePower + 2
		leveled := p.AddXP(xp)
		msg := fmt.Sprintf("⚔️ %sを倒した! +%dXP", m.Name, xp)
		if leveled {
			msg = fmt.Sprintf("✨ Lv.%d! (%sを倒した)", p.Level, m.Name)
		}
		p.LogEvent(msg)

		if rand.Intn(100) < 10 {
			item := RollItem(p.Level)
			p.AddItem(item)
			p.LogEvent(fmt.Sprintf("🌟 %s[%s]を手に入れた!", item.Name, item.Rarity.Label()))
		}
	} else {
		p.Losses++
		p.LogEvent(fmt.Sprintf("⚔️ %sから逃げた", m.Name))
	}
}

func doItemFind(p *Pet) {
	item := RollItem(p.Level)
	p.AddItem(item)
	p.LogEvent(fmt.Sprintf("🌟 %s[%s]を見つけた!", item.Name, item.Rarity.Label()))
}

func doSpecial(p *Pet) {
	evt := specialEvents[rand.Intn(len(specialEvents))]
	bonus := rand.Intn(3) + 1
	p.AddXP(bonus)
	p.LogEvent(fmt.Sprintf("🔮 %s +%dXP", evt, bonus))

}
