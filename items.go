package main

import "math/rand"

// ---------- Visual Slot ----------

type VisualSlot int

const (
	SlotNone        VisualSlot = iota
	SlotWeaponRight            // 剣・斧など右手武器
	SlotWeaponLeft             // 盾・書物など左手
	SlotRanged                 // 弓・杖など両手
	SlotHead                   // 兜・帽子・冠
	SlotBody                   // 鎧・ローブ・マント
	SlotAccessory              // 指輪・首飾りなど
)

var itemSlotMap = map[string]VisualSlot{
	// weapon_right
	"剣": SlotWeaponRight, "大剣": SlotWeaponRight, "短剣": SlotWeaponRight,
	"槍": SlotWeaponRight, "斧": SlotWeaponRight, "太刀": SlotWeaponRight,
	"薙刀": SlotWeaponRight, "鎌": SlotWeaponRight, "鞭": SlotWeaponRight, "爪": SlotWeaponRight,
	// weapon_left
	"盾": SlotWeaponLeft, "鏡": SlotWeaponLeft, "書物": SlotWeaponLeft, "水晶玉": SlotWeaponLeft,
	// ranged
	"弓": SlotRanged, "杖": SlotRanged, "錫杖": SlotRanged, "扇": SlotRanged, "笛": SlotRanged,
	// head
	"兜": SlotHead, "帽子": SlotHead, "冠": SlotHead,
	// body
	"鎧": SlotBody, "ローブ": SlotBody, "マント": SlotBody, "羽衣": SlotBody,
	// accessory
	"指輪": SlotAccessory, "首飾り": SlotAccessory, "腕輪": SlotAccessory,
	"耳飾り": SlotAccessory, "お守り": SlotAccessory, "靴": SlotAccessory,
	"手袋": SlotAccessory, "ベルト": SlotAccessory, "灯": SlotAccessory,
	"鍵": SlotAccessory, "杯": SlotAccessory, "香炉": SlotAccessory,
	"数珠": SlotAccessory, "勾玉": SlotAccessory,
}

// ItemVisualSlot extracts the item type from a generated item name and returns its visual slot.
// Name pattern: [prefix][material]の[type][suffix]〈element〉
func ItemVisualSlot(itemName string) VisualSlot {
	// Find the first の and work from there
	runes := []rune(itemName)
	start := -1
	for i, r := range runes {
		if r == 'の' {
			start = i + 1
			break
		}
	}
	if start < 0 {
		return SlotNone
	}
	rest := string(runes[start:])
	// Match longest type name first
	bestLen := 0
	bestSlot := SlotNone
	for _, t := range itemTypes {
		tr := []rune(t.Name)
		if len(tr) > bestLen && len(rest) >= len(string(tr)) {
			nameRunes := []rune(rest)
			match := true
			for j, r := range tr {
				if nameRunes[j] != r {
					match = false
					break
				}
			}
			if match {
				bestLen = len(tr)
				if slot, ok := itemSlotMap[t.Name]; ok {
					bestSlot = slot
				}
			}
		}
	}
	return bestSlot
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
	switch r {
	case Normal:
		return "ふつう"
	case Fine:
		return "上質"
	case Rare:
		return "希少"
	case Epic:
		return "秘宝"
	case Legendary:
		return "伝説"
	case Mythic:
		return "神話"
	}
	return ""
}

// TmuxColor returns the tmux fg color code for this rarity tier.
// phase is used for Mythic rainbow cycling.
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
	Name   string `json:"name"`
	Rarity Rarity `json:"rarity"`
	Power  int    `json:"power"`
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

// ========== 65 materials ==========

var materials = []material{
	// Tier 0 — ふつう (13)
	{"木", 0, 1}, {"石", 0, 1}, {"草", 0, 1}, {"布", 0, 1}, {"骨", 0, 1},
	{"砂", 0, 1}, {"土", 0, 1}, {"銅", 0, 2}, {"竹", 0, 1}, {"粘土", 0, 1},
	{"藁", 0, 1}, {"麻", 0, 1}, {"貝", 0, 1},
	// Tier 1 — 上質 (13)
	{"鉄", 1, 3}, {"革", 1, 2}, {"水晶", 1, 3}, {"琥珀", 1, 3}, {"翡翠", 1, 4},
	{"真鍮", 1, 3}, {"珊瑚", 1, 3}, {"黒曜石", 1, 4}, {"象牙", 1, 3}, {"瑪瑙", 1, 3},
	{"青銅", 1, 3}, {"錫", 1, 2}, {"石英", 1, 3},
	// Tier 2 — 希少 (13)
	{"銀", 2, 6}, {"金", 2, 7}, {"ミスリル", 2, 8}, {"隕鉄", 2, 7}, {"蒼玉", 2, 6},
	{"紅玉", 2, 6}, {"霊木", 2, 5}, {"白金", 2, 7}, {"瑠璃", 2, 6}, {"翠玉", 2, 6},
	{"碧玉", 2, 6}, {"黄玉", 2, 6}, {"紫水晶", 2, 7},
	// Tier 3 — 秘宝 (10)
	{"オリハルコン", 3, 10}, {"竜骨", 3, 11}, {"精霊石", 3, 9}, {"月光石", 3, 10},
	{"深淵石", 3, 11}, {"天鉄", 3, 10}, {"魔晶石", 3, 9}, {"星銀", 3, 10},
	{"賢者石", 3, 10}, {"妖精鉄", 3, 9},
	// Tier 4 — 伝説 (10)
	{"星屑", 4, 14}, {"虚空石", 4, 15}, {"神樹", 4, 13}, {"冥界鉄", 4, 14},
	{"鳳凰石", 4, 15}, {"天界銀", 4, 14}, {"龍玉", 4, 16}, {"始原石", 4, 14},
	{"時の砂", 4, 15}, {"魂の結晶", 4, 14},
	// Tier 5 — 神話 (6)
	{"混沌の核", 5, 20}, {"永遠の雫", 5, 19}, {"創世の欠片", 5, 21},
	{"終焉の灰", 5, 20}, {"世界樹の実", 5, 22}, {"原初の光", 5, 21},
}

// ========== 40 types ==========

var itemTypes = []itemType{
	{"剣", 2}, {"大剣", 3}, {"短剣", 1}, {"槍", 2}, {"斧", 3},
	{"弓", 2}, {"杖", 2}, {"錫杖", 3}, {"盾", 2}, {"兜", 2},
	{"鎧", 3}, {"ローブ", 1}, {"マント", 1}, {"指輪", 1}, {"首飾り", 1},
	{"腕輪", 1}, {"耳飾り", 1}, {"お守り", 1}, {"靴", 1}, {"手袋", 1},
	{"帽子", 1}, {"冠", 3}, {"ベルト", 1}, {"書物", 2}, {"水晶玉", 2},
	{"扇", 1}, {"笛", 1}, {"鏡", 2}, {"灯", 1}, {"鍵", 1},
	{"太刀", 3}, {"薙刀", 2}, {"鎌", 2}, {"爪", 1}, {"鞭", 2},
	{"杯", 1}, {"香炉", 1}, {"数珠", 1}, {"羽衣", 2}, {"勾玉", 2},
}

// ========== 45 prefixes (+ none = 46) ==========

var prefixes = []prefix{
	// Negative (-1 tier)
	{"古びた", -1, 0.7}, {"錆びた", -1, 0.6}, {"粗末な", -1, 0.7},
	{"壊れかけの", -1, 0.5}, {"汚れた", -1, 0.7}, {"朽ちた", -1, 0.6},
	// Neutral (0 tier)
	{"磨かれた", 0, 1.1}, {"頑丈な", 0, 1.1}, {"鋭い", 0, 1.2},
	{"美しい", 0, 1.0}, {"軽い", 0, 1.0}, {"重厚な", 0, 1.1}, {"精巧な", 0, 1.1},
	{"凛とした", 0, 1.1},
	// Elemental (+1 tier)
	{"炎の", 1, 1.3}, {"氷の", 1, 1.3}, {"雷の", 1, 1.3}, {"風の", 1, 1.2},
	{"水の", 1, 1.2}, {"大地の", 1, 1.2}, {"光の", 1, 1.3}, {"闇の", 1, 1.3},
	// Powerful (+1 tier)
	{"聖なる", 1, 1.4}, {"呪われた", 1, 1.4}, {"輝く", 1, 1.3},
	{"凍てつく", 1, 1.3}, {"燃え盛る", 1, 1.4}, {"嵐の", 1, 1.3},
	{"影の", 1, 1.3}, {"魔性の", 1, 1.3},
	// Color (+1 tier)
	{"蒼き", 1, 1.2}, {"紅き", 1, 1.2}, {"黄金の", 1, 1.3},
	{"白銀の", 1, 1.3}, {"漆黒の", 1, 1.3},
	// Mythic (+2 tier)
	{"伝説の", 2, 1.5}, {"神々の", 2, 1.6}, {"始まりの", 2, 1.5},
	{"終わりの", 2, 1.5}, {"永遠の", 2, 1.6}, {"禁断の", 2, 1.5}, {"至高の", 2, 1.6},
	{"眠れる", 2, 1.4}, {"忘れられた", 2, 1.5}, {"覚めた", 2, 1.4},
}

// ========== 40 suffixes (+ none = 41) ==========

var suffixes = []suffix{
	// Downgrade
	{"の欠片", -1, -2}, {"の残骸", -1, -3}, {"の写し", -1, -1},
	// Neutral
	{"の原型", 0, 0}, {"の試作", 0, 0}, {"の証", 0, 1},
	// Grade
	{"・改", 0, 2}, {"・真", 0, 3}, {"・極", 1, 4}, {"・天", 1, 5}, {"・零", 0, 2},
	{"・壱", 0, 1}, {"・陸", 0, 2},
	// Awakening
	{"─覚醒─", 1, 5}, {"─幻影─", 1, 4}, {"─残光─", 0, 3},
	{"─咆哮─", 1, 5}, {"─黎明─", 1, 4}, {"─黄昏─", 0, 3},
	{"─深淵─", 1, 5}, {"─追憶─", 0, 3}, {"─胎動─", 0, 3},
	{"─輪廻─", 1, 5}, {"─天命─", 1, 4}, {"─因果─", 0, 3},
	// Enhancement
	{"+1", 0, 1}, {"+2", 0, 2}, {"+3", 0, 3}, {"+4", 0, 4}, {"+5", 1, 5},
	// Title
	{"[聖]", 1, 4}, {"[魔]", 1, 4}, {"[王]", 1, 5}, {"[天]", 1, 5}, {"[地]", 1, 4}, {"[人]", 0, 3},
	// Quality
	{"の傑作", 1, 4}, {"の祝福", 1, 3}, {"の加護", 0, 3}, {"の呪い", 1, 4},
	{"の余韻", 0, 2},
}

// ========== 15 elements (+ none = 16) ==========

var elements = []element{
	{"火", 0, 2}, {"水", 0, 2}, {"雷", 0, 3}, {"氷", 0, 2}, {"風", 0, 1},
	{"土", 0, 1}, {"光", 1, 3}, {"闇", 1, 3}, {"毒", 0, 2}, {"聖", 1, 3},
	{"時", 1, 4}, {"空", 1, 3}, {"夢", 0, 2}, {"命", 1, 4}, {"無", 0, 1},
}

// ---------- Generation ----------

// RollItem generates a procedural item.
//
// Name: [prefix][material]の[type][suffix]〈element〉
// Combinations: 46 × 65 × 40 × 41 × 16 × 7(power variance) ≈ 5.5億
// RollItem generates a procedural item scaled to pet level.
// Low levels only see low-tier materials; higher tiers unlock gradually.
//
//   Lv.1-5:   Tier 0 only
//   Lv.6-9:   Tier 0-1
//   Lv.10-13: Tier 0-2
//   Lv.14-16: Tier 0-3
//   Lv.17-19: Tier 0-4
//   Lv.20:    Tier 0-5
//
// Prefix/suffix/element chance also scales with level.
func RollItem(level int) Item {
	mat := rollMaterial(level)
	typ := itemTypes[rand.Intn(len(itemTypes))]

	// Prefix chance: 10% at Lv.1 → 35% at Lv.20
	pfxChance := 10 + (level-1)*25/19
	var pfx *prefix
	if rand.Intn(100) < pfxChance {
		p := rollPrefix(level)
		pfx = &p
	}

	// Suffix chance: 5% at Lv.1 → 20% at Lv.20
	sfxChance := 5 + (level-1)*15/19
	var sfx *suffix
	if rand.Intn(100) < sfxChance {
		s := rollSuffix(level)
		sfx = &s
	}

	// Element chance: 3% at Lv.1 → 15% at Lv.20
	elmChance := 3 + (level-1)*12/19
	var elm *element
	if rand.Intn(100) < elmChance {
		e := elements[rand.Intn(len(elements))]
		elm = &e
	}

	// Build name
	name := ""
	if pfx != nil {
		name += pfx.Name
	}
	name += mat.Name + "の" + typ.Name
	if sfx != nil {
		name += sfx.Name
	}
	if elm != nil {
		name += "〈" + elm.Name + "〉"
	}

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

	return Item{
		Name:   name,
		Rarity: rarity,
		Power:  int(power),
	}
}

// maxTierForLevel returns the highest material tier available at this level.
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
		return 5 // Lv.20 only
	}
}

func rollMaterial(level int) material {
	maxTier := maxTierForLevel(level)

	// Weighted roll, but capped at maxTier
	// Tier 0: 90%, 1: 8%, 2: 1.5%, 3: 0.4%, 4: 0.099%, 5: 0.001%
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
	for _, m := range materials {
		if m.Tier == tier {
			pool = append(pool, m)
		}
	}
	return pool[rand.Intn(len(pool))]
}

// rollPrefix picks a prefix, limiting mythic (+2) prefixes to higher levels.
func rollPrefix(level int) prefix {
	for {
		p := prefixes[rand.Intn(len(prefixes))]
		// +1 tier prefixes only at Lv.8+
		if p.TierMod >= 1 && level < 8 {
			continue
		}
		// Mythic prefixes (+2 tier) only at Lv.18+
		if p.TierMod >= 2 && level < 18 {
			continue
		}
		return p
	}
}

// rollSuffix picks a suffix, limiting +1 tier suffixes to higher levels.
func rollSuffix(level int) suffix {
	for {
		s := suffixes[rand.Intn(len(suffixes))]
		// +1 tier suffixes only at Lv.12+
		if s.TierMod >= 1 && level < 12 {
			continue
		}
		return s
	}
}
