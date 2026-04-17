package main

import (
	"fmt"
	"math/rand"
	"time"
)

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

// Tick processes one game tick.
func Tick(p *Pet) (*Pet, bool) {
	now := time.Now()

	elapsed := int(now.Sub(p.UpdatedAt).Minutes())
	if elapsed > 0 {
		leveled := p.AddXP(elapsed)
		if leveled {
			p.LogEvent(fmt.Sprintf(T("level_up"), p.Level))
		}
	}

	p.UpdatedAt = now

	if p.CanReincarnate() {
		gen := p.Generation
		p.LogEvent(fmt.Sprintf(T("reincarnate"), gen+1))
		newPet := p.Reincarnate()
		newPet.EventLog = p.EventLog
		newPet.LastEvent = p.LastEvent
		newPet.EventAt = p.EventAt
		return newPet, true
	}

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

func monsterName(idx int) string {
	names := monstersLocale[currentLang]
	if names == nil {
		names = monstersLocale[LangEN]
	}
	if idx < len(names) {
		return names[idx]
	}
	return names[0]
}

func doBattle(p *Pet) {
	maxIdx := p.Level / 3
	if maxIdx >= len(monsterPowers) {
		maxIdx = len(monsterPowers) - 1
	}
	idx := rand.Intn(maxIdx + 1)
	name := monsterName(idx)
	basePower := monsterPowers[idx]

	power := p.TotalPower()
	monsterPow := basePower + rand.Intn(3)

	if power+rand.Intn(5) >= monsterPow {
		p.Wins++
		xp := basePower + 2
		leveled := p.AddXP(xp)
		msg := fmt.Sprintf(T("defeated"), name, xp)
		if leveled {
			msg = fmt.Sprintf(T("defeated_lv"), p.Level, name)
		}
		p.LogEvent(msg)

		if rand.Intn(100) < 10 {
			item := RollItem(p.Level)
			p.AddItem(item)
			p.LogEvent(fmt.Sprintf(T("item_got"), item.Name, item.Rarity.Label()))
		}
	} else {
		p.Losses++
		p.LogEvent(fmt.Sprintf(T("fled"), name))
	}
}

func doItemFind(p *Pet) {
	item := RollItem(p.Level)
	p.AddItem(item)
	p.LogEvent(fmt.Sprintf(T("item_found"), item.Name, item.Rarity.Label()))
}

func doSpecial(p *Pet) {
	events := specialEventsLocale[currentLang]
	if events == nil {
		events = specialEventsLocale[LangEN]
	}
	evt := events[rand.Intn(len(events))]
	bonus := rand.Intn(3) + 1
	p.AddXP(bonus)
	p.LogEvent(fmt.Sprintf(T("special_fmt"), evt, bonus))
}
