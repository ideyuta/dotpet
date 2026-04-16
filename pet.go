package main

import "time"

type Pet struct {
	Name       string    `json:"name"`
	Species    string    `json:"species"`
	Level      int       `json:"level"`
	XP         int       `json:"xp"`
	Power      int       `json:"power"`
	Generation int       `json:"generation"`
	Legacy     int       `json:"legacy"`
	Inventory  []Item    `json:"inventory"`
	Equipped   *Item     `json:"equipped"`
	LastEvent  string    `json:"last_event"`
	EventAt    time.Time `json:"event_at"`
	EventLog   []string  `json:"event_log"`
	BornAt     time.Time `json:"born_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	// Stats
	Wins       int       `json:"wins"`
	Losses     int       `json:"losses"`
	ItemsFound int       `json:"items_found"`
	TotalXP    int       `json:"total_xp"`
	BestItem   *Item     `json:"best_item"`
}

const maxLevel = 20

func NewPet(generation, legacy int) *Pet {
	species := DetermineSpecies()
	now := time.Now()
	return &Pet{
		Name:       GenerateName(species),
		Species:    species,
		Level:      1,
		XP:         0,
		Power:      1 + legacy,
		Generation: generation,
		Legacy:     legacy,
		Inventory:  []Item{},
		BornAt:     now,
		UpdatedAt:  now,
	}
}

func (p *Pet) XPToNext() int {
	// Quadratic curve: early levels are quick, later levels take much longer.
	// Lv.1→2: 15 XP, Lv.10→11: 150 XP, Lv.19→20: 435 XP
	// Total 1→20: ~4,500 XP (~3 days of idle play)
	return p.Level*(p.Level+2) + 12
}

func (p *Pet) TotalPower() int {
	pw := p.Power
	if p.Equipped != nil {
		pw += p.Equipped.Power
	}
	return pw
}

func (p *Pet) AddXP(xp int) bool {
	p.XP += xp
	p.TotalXP += xp
	leveled := false
	for p.Level < maxLevel && p.XP >= p.XPToNext() {
		p.XP -= p.XPToNext()
		p.Level++
		p.Power++
		leveled = true
	}
	return leveled
}

func (p *Pet) AddItem(item Item) {
	p.ItemsFound++
	if p.BestItem == nil || item.Rarity > p.BestItem.Rarity || (item.Rarity == p.BestItem.Rarity && item.Power > p.BestItem.Power) {
		best := item
		p.BestItem = &best
	}
	if len(p.Inventory) >= 10 {
		// replace weakest item
		weakest := 0
		for i, it := range p.Inventory {
			if it.Power < p.Inventory[weakest].Power {
				weakest = i
			}
		}
		if item.Power > p.Inventory[weakest].Power {
			p.Inventory[weakest] = item
		}
	} else {
		p.Inventory = append(p.Inventory, item)
	}
	p.AutoEquip()
}

func (p *Pet) AutoEquip() {
	if len(p.Inventory) == 0 {
		return
	}
	best := 0
	for i, it := range p.Inventory {
		if it.Power > p.Inventory[best].Power {
			best = i
		}
	}
	eq := p.Inventory[best]
	p.Equipped = &eq
}

func (p *Pet) CanReincarnate() bool {
	return p.Level >= maxLevel
}

func (p *Pet) Reincarnate() *Pet {
	return NewPet(p.Generation+1, p.Legacy+1)
}

func (p *Pet) LogEvent(msg string) {
	p.LastEvent = msg
	p.EventAt = time.Now()
	p.EventLog = append(p.EventLog, time.Now().Format("15:04")+" "+msg)
	if len(p.EventLog) > 20 {
		p.EventLog = p.EventLog[len(p.EventLog)-20:]
	}
}
