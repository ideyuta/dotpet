package main

import "math/rand"

// ---------- Visual Slot ----------

type VisualSlot int

const (
	SlotNone        VisualSlot = iota
	SlotWeaponRight            // sword, axe, etc.
	SlotWeaponLeft             // shield, tome, etc.
	SlotRanged                 // bow, staff, etc.
	SlotHead                   // helm, hat, crown
	SlotBody                   // armor, robe, cloak
	SlotAccessory              // ring, necklace, etc.
)

// ItemVisualSlot returns the visual slot for an item.
// Uses the Slot field if set; falls back to name-parsing for legacy items.
func ItemVisualSlot(item *Item) VisualSlot {
	if item.Slot != SlotNone {
		return item.Slot
	}
	return itemVisualSlotFromName(item.Name)
}

// itemVisualSlotFromName parses a generated item name to determine its slot.
// This handles legacy items saved before the Slot field was added.
func itemVisualSlotFromName(itemName string) VisualSlot {
	// Try all languages' slot maps
	for _, vocab := range itemVocabs {
		for typeName, slot := range vocab.SlotMap {
			if containsType(itemName, typeName) {
				return slot
			}
		}
	}
	return SlotNone
}

func containsType(name, typeName string) bool {
	runes := []rune(name)
	typeRunes := []rune(typeName)
	for i := 0; i <= len(runes)-len(typeRunes); i++ {
		match := true
		for j, r := range typeRunes {
			if runes[i+j] != r {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

// ---------- Rarity (6 tiers) ----------

type Rarity int

const (
	Normal    Rarity = iota // ★
	Fine                    // ★★
	Rare                    // ★★★
	Epic                    // ★★★★
	Legendary               // ★★★★★
	Mythic                  // ★★★★★★
)

func (r Rarity) String() string {
	s := ""
	for i := 0; i <= int(r); i++ {
		s += "★"
	}
	return s
}

func (r Rarity) Label() string {
	keys := []string{"rarity_normal", "rarity_fine", "rarity_rare", "rarity_epic", "rarity_legendary", "rarity_mythic"}
	if int(r) < len(keys) {
		return T(keys[r])
	}
	return ""
}

// TmuxColor returns the tmux fg color code for this rarity tier.
func (r Rarity) TmuxColor(phase int) string {
	switch r {
	case Fine:
		return "#[fg=#D4D4D4]"
	case Rare:
		return "#[fg=#569CD6]"
	case Epic:
		return "#[fg=#C586C0]"
	case Legendary:
		return "#[fg=#DCDCAA]"
	case Mythic:
		colors := []string{"#[fg=#CE9178]", "#[fg=#DCDCAA]", "#[fg=#569CD6]", "#[fg=#C586C0]"}
		return colors[phase%len(colors)]
	default:
		return ""
	}
}

// AnsiColor returns the ANSI escape code for this rarity tier.
func (r Rarity) AnsiColor(phase int) string {
	switch r {
	case Fine:
		return "\033[97m"
	case Rare:
		return "\033[94m"
	case Epic:
		return "\033[95m"
	case Legendary:
		return "\033[93m"
	case Mythic:
		colors := []string{"\033[91m", "\033[93m", "\033[94m", "\033[95m"}
		return colors[phase%len(colors)]
	default:
		return ""
	}
}

func clampRarity(r int) Rarity {
	if r < 0 {
		return Normal
	}
	if r > int(Mythic) {
		return Mythic
	}
	return Rarity(r)
}

// ---------- Item ----------

type Item struct {
	Name   string     `json:"name"`
	Rarity Rarity     `json:"rarity"`
	Power  int        `json:"power"`
	Slot   VisualSlot `json:"slot,omitempty"`
}

// ---------- Parts ----------

type material struct {
	Name  string
	Tier  int
	Power int
}

type itemType struct {
	Name  string
	Power int
}

type prefix struct {
	Name      string
	TierMod   int
	PowerMult float64
}

type suffix struct {
	Name     string
	TierMod  int
	PowerAdd int
}

type element struct {
	Name     string
	TierMod  int
	PowerAdd int
}

// ---------- Generation ----------

// RollItem generates a procedural item scaled to pet level.
func RollItem(level int) Item {
	vocab := Vocab()

	mat := rollMaterial(level, vocab.Materials)
	typ := vocab.Types[rand.Intn(len(vocab.Types))]

	pfxChance := 10 + (level-1)*25/19
	var pfx *prefix
	if rand.Intn(100) < pfxChance {
		p := rollPrefix(level, vocab.Prefixes)
		pfx = &p
	}

	sfxChance := 5 + (level-1)*15/19
	var sfx *suffix
	if rand.Intn(100) < sfxChance {
		s := rollSuffix(level, vocab.Suffixes)
		sfx = &s
	}

	elmChance := 3 + (level-1)*12/19
	var elm *element
	if rand.Intn(100) < elmChance {
		e := vocab.Elements[rand.Intn(len(vocab.Elements))]
		elm = &e
	}

	// Build name
	pfxName := ""
	if pfx != nil {
		pfxName = pfx.Name
	}
	sfxName := ""
	if sfx != nil {
		sfxName = sfx.Name
	}
	elmName := ""
	if elm != nil {
		elmName = elm.Name
	}
	name := vocab.BuildName(pfxName, mat.Name, typ.Name, sfxName, elmName)

	// Calculate rarity
	tier := mat.Tier
	if pfx != nil {
		tier += pfx.TierMod
	}
	if sfx != nil {
		tier += sfx.TierMod
	}
	if elm != nil {
		tier += elm.TierMod
	}
	rarity := clampRarity(tier)

	// Calculate power
	power := float64(mat.Power + typ.Power)
	if pfx != nil {
		power *= pfx.PowerMult
	}
	if sfx != nil {
		power += float64(sfx.PowerAdd)
	}
	if elm != nil {
		power += float64(elm.PowerAdd)
	}
	power += float64(rand.Intn(7) - 3)
	if power < 1 {
		power = 1
	}

	// Determine slot from type name
	slot := SlotNone
	if s, ok := vocab.SlotMap[typ.Name]; ok {
		slot = s
	}

	return Item{
		Name:   name,
		Rarity: rarity,
		Power:  int(power),
		Slot:   slot,
	}
}

func maxTierForLevel(level int) int {
	switch {
	case level <= 5:
		return 0
	case level <= 9:
		return 1
	case level <= 13:
		return 2
	case level <= 16:
		return 3
	case level <= 19:
		return 4
	default:
		return 5
	}
}

func rollMaterial(level int, mats []material) material {
	maxTier := maxTierForLevel(level)

	r := rand.Intn(100000)
	var tier int
	switch {
	case r < 90000:
		tier = 0
	case r < 98000:
		tier = 1
	case r < 99500:
		tier = 2
	case r < 99900:
		tier = 3
	case r < 99999:
		tier = 4
	default:
		tier = 5
	}
	if tier > maxTier {
		tier = maxTier
	}

	var pool []material
	for _, m := range mats {
		if m.Tier == tier {
			pool = append(pool, m)
		}
	}
	return pool[rand.Intn(len(pool))]
}

func rollPrefix(level int, pfxs []prefix) prefix {
	for {
		p := pfxs[rand.Intn(len(pfxs))]
		if p.TierMod >= 1 && level < 8 {
			continue
		}
		if p.TierMod >= 2 && level < 18 {
			continue
		}
		return p
	}
}

func rollSuffix(level int, sfxs []suffix) suffix {
	for {
		s := sfxs[rand.Intn(len(sfxs))]
		if s.TierMod >= 1 && level < 12 {
			continue
		}
		return s
	}
}
