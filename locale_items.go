package main

import "strings"

// ItemVocab holds all item-related vocabulary for a language.
type ItemVocab struct {
	Materials []material
	Types     []itemType
	Prefixes  []prefix
	Suffixes  []suffix
	Elements  []element
	BuildName func(pfx, mat, typ, sfx, elm string) string
	SlotMap   map[string]VisualSlot
}

// Vocab returns the item vocabulary for the current language.
func Vocab() *ItemVocab {
	if v, ok := itemVocabs[currentLang]; ok {
		return v
	}
	return itemVocabs[LangEN]
}

var itemVocabs = map[Lang]*ItemVocab{
	LangJA: {
		Materials: materialsJA,
		Types:     itemTypesJA,
		Prefixes:  prefixesJA,
		Suffixes:  suffixesJA,
		Elements:  elementsJA,
		BuildName: func(pfx, mat, typ, sfx, elm string) string {
			name := pfx + mat + "の" + typ + sfx
			if elm != "" {
				name += "〈" + elm + "〉"
			}
			return name
		},
		SlotMap: slotMapJA,
	},
	LangEN: {
		Materials: materialsEN,
		Types:     itemTypesEN,
		Prefixes:  prefixesEN,
		Suffixes:  suffixesEN,
		Elements:  elementsEN,
		BuildName: func(pfx, mat, typ, sfx, elm string) string {
			parts := []string{}
			if pfx != "" {
				parts = append(parts, pfx)
			}
			parts = append(parts, mat, typ)
			name := strings.Join(parts, " ")
			if sfx != "" {
				name += " " + sfx
			}
			if elm != "" {
				name += " <" + elm + ">"
			}
			return name
		},
		SlotMap: slotMapEN,
	},
	LangZH: {
		Materials: materialsZH,
		Types:     itemTypesZH,
		Prefixes:  prefixesZH,
		Suffixes:  suffixesZH,
		Elements:  elementsZH,
		BuildName: func(pfx, mat, typ, sfx, elm string) string {
			name := pfx + mat + typ + sfx
			if elm != "" {
				name += "〈" + elm + "〉"
			}
			return name
		},
		SlotMap: slotMapZH,
	},
}

// ========== Japanese ==========

var materialsJA = []material{
	{"木", 0, 1}, {"石", 0, 1}, {"草", 0, 1}, {"布", 0, 1}, {"骨", 0, 1},
	{"砂", 0, 1}, {"土", 0, 1}, {"銅", 0, 2}, {"竹", 0, 1}, {"粘土", 0, 1},
	{"藁", 0, 1}, {"麻", 0, 1}, {"貝", 0, 1},
	{"鉄", 1, 3}, {"革", 1, 2}, {"水晶", 1, 3}, {"琥珀", 1, 3}, {"翡翠", 1, 4},
	{"真鍮", 1, 3}, {"珊瑚", 1, 3}, {"黒曜石", 1, 4}, {"象牙", 1, 3}, {"瑪瑙", 1, 3},
	{"青銅", 1, 3}, {"錫", 1, 2}, {"石英", 1, 3},
	{"銀", 2, 6}, {"金", 2, 7}, {"ミスリル", 2, 8}, {"隕鉄", 2, 7}, {"蒼玉", 2, 6},
	{"紅玉", 2, 6}, {"霊木", 2, 5}, {"白金", 2, 7}, {"瑠璃", 2, 6}, {"翠玉", 2, 6},
	{"碧玉", 2, 6}, {"黄玉", 2, 6}, {"紫水晶", 2, 7},
	{"オリハルコン", 3, 10}, {"竜骨", 3, 11}, {"精霊石", 3, 9}, {"月光石", 3, 10},
	{"深淵石", 3, 11}, {"天鉄", 3, 10}, {"魔晶石", 3, 9}, {"星銀", 3, 10},
	{"賢者石", 3, 10}, {"妖精鉄", 3, 9},
	{"星屑", 4, 14}, {"虚空石", 4, 15}, {"神樹", 4, 13}, {"冥界鉄", 4, 14},
	{"鳳凰石", 4, 15}, {"天界銀", 4, 14}, {"龍玉", 4, 16}, {"始原石", 4, 14},
	{"時の砂", 4, 15}, {"魂の結晶", 4, 14},
	{"混沌の核", 5, 20}, {"永遠の雫", 5, 19}, {"創世の欠片", 5, 21},
	{"終焉の灰", 5, 20}, {"世界樹の実", 5, 22}, {"原初の光", 5, 21},
}

var itemTypesJA = []itemType{
	{"剣", 2}, {"大剣", 3}, {"短剣", 1}, {"槍", 2}, {"斧", 3},
	{"弓", 2}, {"杖", 2}, {"錫杖", 3}, {"盾", 2}, {"兜", 2},
	{"鎧", 3}, {"ローブ", 1}, {"マント", 1}, {"指輪", 1}, {"首飾り", 1},
	{"腕輪", 1}, {"耳飾り", 1}, {"お守り", 1}, {"靴", 1}, {"手袋", 1},
	{"帽子", 1}, {"冠", 3}, {"ベルト", 1}, {"書物", 2}, {"水晶玉", 2},
	{"扇", 1}, {"笛", 1}, {"鏡", 2}, {"灯", 1}, {"鍵", 1},
	{"太刀", 3}, {"薙刀", 2}, {"鎌", 2}, {"爪", 1}, {"鞭", 2},
	{"杯", 1}, {"香炉", 1}, {"数珠", 1}, {"羽衣", 2}, {"勾玉", 2},
}

var prefixesJA = []prefix{
	{"古びた", -1, 0.7}, {"錆びた", -1, 0.6}, {"粗末な", -1, 0.7},
	{"壊れかけの", -1, 0.5}, {"汚れた", -1, 0.7}, {"朽ちた", -1, 0.6},
	{"磨かれた", 0, 1.1}, {"頑丈な", 0, 1.1}, {"鋭い", 0, 1.2},
	{"美しい", 0, 1.0}, {"軽い", 0, 1.0}, {"重厚な", 0, 1.1}, {"精巧な", 0, 1.1},
	{"凛とした", 0, 1.1},
	{"炎の", 1, 1.3}, {"氷の", 1, 1.3}, {"雷の", 1, 1.3}, {"風の", 1, 1.2},
	{"水の", 1, 1.2}, {"大地の", 1, 1.2}, {"光の", 1, 1.3}, {"闇の", 1, 1.3},
	{"聖なる", 1, 1.4}, {"呪われた", 1, 1.4}, {"輝く", 1, 1.3},
	{"凍てつく", 1, 1.3}, {"燃え盛る", 1, 1.4}, {"嵐の", 1, 1.3},
	{"影の", 1, 1.3}, {"魔性の", 1, 1.3},
	{"蒼き", 1, 1.2}, {"紅き", 1, 1.2}, {"黄金の", 1, 1.3},
	{"白銀の", 1, 1.3}, {"漆黒の", 1, 1.3},
	{"伝説の", 2, 1.5}, {"神々の", 2, 1.6}, {"始まりの", 2, 1.5},
	{"終わりの", 2, 1.5}, {"永遠の", 2, 1.6}, {"禁断の", 2, 1.5}, {"至高の", 2, 1.6},
	{"眠れる", 2, 1.4}, {"忘れられた", 2, 1.5}, {"覚めた", 2, 1.4},
}

var suffixesJA = []suffix{
	{"の欠片", -1, -2}, {"の残骸", -1, -3}, {"の写し", -1, -1},
	{"の原型", 0, 0}, {"の試作", 0, 0}, {"の証", 0, 1},
	{"・改", 0, 2}, {"・真", 0, 3}, {"・極", 1, 4}, {"・天", 1, 5}, {"・零", 0, 2},
	{"・壱", 0, 1}, {"・陸", 0, 2},
	{"─覚醒─", 1, 5}, {"─幻影─", 1, 4}, {"─残光─", 0, 3},
	{"─咆哮─", 1, 5}, {"─黎明─", 1, 4}, {"─黄昏─", 0, 3},
	{"─深淵─", 1, 5}, {"─追憶─", 0, 3}, {"─胎動─", 0, 3},
	{"─輪廻─", 1, 5}, {"─天命─", 1, 4}, {"─因果─", 0, 3},
	{"+1", 0, 1}, {"+2", 0, 2}, {"+3", 0, 3}, {"+4", 0, 4}, {"+5", 1, 5},
	{"[聖]", 1, 4}, {"[魔]", 1, 4}, {"[王]", 1, 5}, {"[天]", 1, 5}, {"[地]", 1, 4}, {"[人]", 0, 3},
	{"の傑作", 1, 4}, {"の祝福", 1, 3}, {"の加護", 0, 3}, {"の呪い", 1, 4},
	{"の余韻", 0, 2},
}

var elementsJA = []element{
	{"火", 0, 2}, {"水", 0, 2}, {"雷", 0, 3}, {"氷", 0, 2}, {"風", 0, 1},
	{"土", 0, 1}, {"光", 1, 3}, {"闇", 1, 3}, {"毒", 0, 2}, {"聖", 1, 3},
	{"時", 1, 4}, {"空", 1, 3}, {"夢", 0, 2}, {"命", 1, 4}, {"無", 0, 1},
}

var slotMapJA = map[string]VisualSlot{
	"剣": SlotWeaponRight, "大剣": SlotWeaponRight, "短剣": SlotWeaponRight,
	"槍": SlotWeaponRight, "斧": SlotWeaponRight, "太刀": SlotWeaponRight,
	"薙刀": SlotWeaponRight, "鎌": SlotWeaponRight, "鞭": SlotWeaponRight, "爪": SlotWeaponRight,
	"盾": SlotWeaponLeft, "鏡": SlotWeaponLeft, "書物": SlotWeaponLeft, "水晶玉": SlotWeaponLeft,
	"弓": SlotRanged, "杖": SlotRanged, "錫杖": SlotRanged, "扇": SlotRanged, "笛": SlotRanged,
	"兜": SlotHead, "帽子": SlotHead, "冠": SlotHead,
	"鎧": SlotBody, "ローブ": SlotBody, "マント": SlotBody, "羽衣": SlotBody,
	"指輪": SlotAccessory, "首飾り": SlotAccessory, "腕輪": SlotAccessory,
	"耳飾り": SlotAccessory, "お守り": SlotAccessory, "靴": SlotAccessory,
	"手袋": SlotAccessory, "ベルト": SlotAccessory, "灯": SlotAccessory,
	"鍵": SlotAccessory, "杯": SlotAccessory, "香炉": SlotAccessory,
	"数珠": SlotAccessory, "勾玉": SlotAccessory,
}

// ========== English ==========

var materialsEN = []material{
	{"Wood", 0, 1}, {"Stone", 0, 1}, {"Herb", 0, 1}, {"Cloth", 0, 1}, {"Bone", 0, 1},
	{"Sand", 0, 1}, {"Clay", 0, 1}, {"Copper", 0, 2}, {"Bamboo", 0, 1}, {"Mud", 0, 1},
	{"Straw", 0, 1}, {"Hemp", 0, 1}, {"Shell", 0, 1},
	{"Iron", 1, 3}, {"Leather", 1, 2}, {"Crystal", 1, 3}, {"Amber", 1, 3}, {"Jade", 1, 4},
	{"Brass", 1, 3}, {"Coral", 1, 3}, {"Obsidian", 1, 4}, {"Ivory", 1, 3}, {"Agate", 1, 3},
	{"Bronze", 1, 3}, {"Tin", 1, 2}, {"Quartz", 1, 3},
	{"Silver", 2, 6}, {"Gold", 2, 7}, {"Mithril", 2, 8}, {"Meteor", 2, 7}, {"Sapphire", 2, 6},
	{"Ruby", 2, 6}, {"Spirit Wood", 2, 5}, {"Platinum", 2, 7}, {"Lapis", 2, 6}, {"Emerald", 2, 6},
	{"Jasper", 2, 6}, {"Topaz", 2, 6}, {"Amethyst", 2, 7},
	{"Orichalcum", 3, 10}, {"Dragon Bone", 3, 11}, {"Spirit Stone", 3, 9}, {"Moonstone", 3, 10},
	{"Abyssal Stone", 3, 11}, {"Sky Iron", 3, 10}, {"Mana Crystal", 3, 9}, {"Star Silver", 3, 10},
	{"Sage Stone", 3, 10}, {"Fairy Iron", 3, 9},
	{"Stardust", 4, 14}, {"Void Stone", 4, 15}, {"World Tree", 4, 13}, {"Nether Iron", 4, 14},
	{"Phoenix Stone", 4, 15}, {"Celestial Silver", 4, 14}, {"Dragon Orb", 4, 16}, {"Primal Stone", 4, 14},
	{"Sands of Time", 4, 15}, {"Soul Crystal", 4, 14},
	{"Chaos Core", 5, 20}, {"Eternal Drop", 5, 19}, {"Genesis Shard", 5, 21},
	{"Doom Ash", 5, 20}, {"World Fruit", 5, 22}, {"Primal Light", 5, 21},
}

var itemTypesEN = []itemType{
	{"Sword", 2}, {"Greatsword", 3}, {"Dagger", 1}, {"Spear", 2}, {"Axe", 3},
	{"Bow", 2}, {"Staff", 2}, {"Crosier", 3}, {"Shield", 2}, {"Helm", 2},
	{"Armor", 3}, {"Robe", 1}, {"Cloak", 1}, {"Ring", 1}, {"Necklace", 1},
	{"Bracelet", 1}, {"Earring", 1}, {"Charm", 1}, {"Boots", 1}, {"Gloves", 1},
	{"Hat", 1}, {"Crown", 3}, {"Belt", 1}, {"Tome", 2}, {"Orb", 2},
	{"Fan", 1}, {"Flute", 1}, {"Mirror", 2}, {"Lantern", 1}, {"Key", 1},
	{"Katana", 3}, {"Glaive", 2}, {"Scythe", 2}, {"Claw", 1}, {"Whip", 2},
	{"Chalice", 1}, {"Censer", 1}, {"Rosary", 1}, {"Veil", 2}, {"Magatama", 2},
}

var prefixesEN = []prefix{
	{"Old", -1, 0.7}, {"Rusted", -1, 0.6}, {"Crude", -1, 0.7},
	{"Cracked", -1, 0.5}, {"Soiled", -1, 0.7}, {"Decayed", -1, 0.6},
	{"Polished", 0, 1.1}, {"Sturdy", 0, 1.1}, {"Keen", 0, 1.2},
	{"Beautiful", 0, 1.0}, {"Light", 0, 1.0}, {"Heavy", 0, 1.1}, {"Fine", 0, 1.1},
	{"Noble", 0, 1.1},
	{"Flame", 1, 1.3}, {"Frost", 1, 1.3}, {"Thunder", 1, 1.3}, {"Wind", 1, 1.2},
	{"Aqua", 1, 1.2}, {"Earth", 1, 1.2}, {"Radiant", 1, 1.3}, {"Shadow", 1, 1.3},
	{"Holy", 1, 1.4}, {"Cursed", 1, 1.4}, {"Gleaming", 1, 1.3},
	{"Frozen", 1, 1.3}, {"Blazing", 1, 1.4}, {"Storm", 1, 1.3},
	{"Phantom", 1, 1.3}, {"Demonic", 1, 1.3},
	{"Azure", 1, 1.2}, {"Crimson", 1, 1.2}, {"Golden", 1, 1.3},
	{"Argent", 1, 1.3}, {"Obsidian", 1, 1.3},
	{"Legendary", 2, 1.5}, {"Divine", 2, 1.6}, {"Genesis", 2, 1.5},
	{"Final", 2, 1.5}, {"Eternal", 2, 1.6}, {"Forbidden", 2, 1.5}, {"Supreme", 2, 1.6},
	{"Dormant", 2, 1.4}, {"Forgotten", 2, 1.5}, {"Awakened", 2, 1.4},
}

var suffixesEN = []suffix{
	{"Shard", -1, -2}, {"Wreck", -1, -3}, {"Replica", -1, -1},
	{"Proto", 0, 0}, {"Draft", 0, 0}, {"Mark", 0, 1},
	{"Mk.II", 0, 2}, {"Mk.III", 0, 3}, {"EX", 1, 4}, {"Omega", 1, 5}, {"Zero", 0, 2},
	{"Alpha", 0, 1}, {"Sigma", 0, 2},
	{"-Awoken-", 1, 5}, {"-Mirage-", 1, 4}, {"-Twilight-", 0, 3},
	{"-Roar-", 1, 5}, {"-Dawn-", 1, 4}, {"-Dusk-", 0, 3},
	{"-Abyss-", 1, 5}, {"-Memory-", 0, 3}, {"-Pulse-", 0, 3},
	{"-Samsara-", 1, 5}, {"-Fate-", 1, 4}, {"-Karma-", 0, 3},
	{"+1", 0, 1}, {"+2", 0, 2}, {"+3", 0, 3}, {"+4", 0, 4}, {"+5", 1, 5},
	{"[Holy]", 1, 4}, {"[Arcane]", 1, 4}, {"[Royal]", 1, 5}, {"[Celestial]", 1, 5}, {"[Terra]", 1, 4}, {"[Mortal]", 0, 3},
	{"Opus", 1, 4}, {"Blessed", 1, 3}, {"Aegis", 0, 3}, {"Hex", 1, 4},
	{"Echo", 0, 2},
}

var elementsEN = []element{
	{"Fire", 0, 2}, {"Water", 0, 2}, {"Thunder", 0, 3}, {"Ice", 0, 2}, {"Wind", 0, 1},
	{"Earth", 0, 1}, {"Light", 1, 3}, {"Dark", 1, 3}, {"Poison", 0, 2}, {"Holy", 1, 3},
	{"Time", 1, 4}, {"Void", 1, 3}, {"Dream", 0, 2}, {"Life", 1, 4}, {"Null", 0, 1},
}

var slotMapEN = map[string]VisualSlot{
	"Sword": SlotWeaponRight, "Greatsword": SlotWeaponRight, "Dagger": SlotWeaponRight,
	"Spear": SlotWeaponRight, "Axe": SlotWeaponRight, "Katana": SlotWeaponRight,
	"Glaive": SlotWeaponRight, "Scythe": SlotWeaponRight, "Whip": SlotWeaponRight, "Claw": SlotWeaponRight,
	"Shield": SlotWeaponLeft, "Mirror": SlotWeaponLeft, "Tome": SlotWeaponLeft, "Orb": SlotWeaponLeft,
	"Bow": SlotRanged, "Staff": SlotRanged, "Crosier": SlotRanged, "Fan": SlotRanged, "Flute": SlotRanged,
	"Helm": SlotHead, "Hat": SlotHead, "Crown": SlotHead,
	"Armor": SlotBody, "Robe": SlotBody, "Cloak": SlotBody, "Veil": SlotBody,
	"Ring": SlotAccessory, "Necklace": SlotAccessory, "Bracelet": SlotAccessory,
	"Earring": SlotAccessory, "Charm": SlotAccessory, "Boots": SlotAccessory,
	"Gloves": SlotAccessory, "Belt": SlotAccessory, "Lantern": SlotAccessory,
	"Key": SlotAccessory, "Chalice": SlotAccessory, "Censer": SlotAccessory,
	"Rosary": SlotAccessory, "Magatama": SlotAccessory,
}

// ========== Chinese ==========

var materialsZH = []material{
	{"木", 0, 1}, {"石", 0, 1}, {"草", 0, 1}, {"布", 0, 1}, {"骨", 0, 1},
	{"砂", 0, 1}, {"土", 0, 1}, {"铜", 0, 2}, {"竹", 0, 1}, {"黏土", 0, 1},
	{"稻草", 0, 1}, {"麻", 0, 1}, {"贝", 0, 1},
	{"铁", 1, 3}, {"皮革", 1, 2}, {"水晶", 1, 3}, {"琥珀", 1, 3}, {"翡翠", 1, 4},
	{"黄铜", 1, 3}, {"珊瑚", 1, 3}, {"黑曜石", 1, 4}, {"象牙", 1, 3}, {"玛瑙", 1, 3},
	{"青铜", 1, 3}, {"锡", 1, 2}, {"石英", 1, 3},
	{"银", 2, 6}, {"金", 2, 7}, {"秘银", 2, 8}, {"陨铁", 2, 7}, {"蓝宝石", 2, 6},
	{"红宝石", 2, 6}, {"灵木", 2, 5}, {"白金", 2, 7}, {"琉璃", 2, 6}, {"祖母绿", 2, 6},
	{"碧玉", 2, 6}, {"黄玉", 2, 6}, {"紫水晶", 2, 7},
	{"奥利哈钢", 3, 10}, {"龙骨", 3, 11}, {"精灵石", 3, 9}, {"月光石", 3, 10},
	{"深渊石", 3, 11}, {"天铁", 3, 10}, {"魔晶石", 3, 9}, {"星银", 3, 10},
	{"贤者石", 3, 10}, {"妖精铁", 3, 9},
	{"星尘", 4, 14}, {"虚空石", 4, 15}, {"神树", 4, 13}, {"冥界铁", 4, 14},
	{"凤凰石", 4, 15}, {"天界银", 4, 14}, {"龙珠", 4, 16}, {"始原石", 4, 14},
	{"时之砂", 4, 15}, {"魂之结晶", 4, 14},
	{"混沌核心", 5, 20}, {"永恒之滴", 5, 19}, {"创世碎片", 5, 21},
	{"终焉之灰", 5, 20}, {"世界树果实", 5, 22}, {"原初之光", 5, 21},
}

var itemTypesZH = []itemType{
	{"剑", 2}, {"大剑", 3}, {"短剑", 1}, {"枪", 2}, {"斧", 3},
	{"弓", 2}, {"杖", 2}, {"锡杖", 3}, {"盾", 2}, {"头盔", 2},
	{"铠甲", 3}, {"法袍", 1}, {"斗篷", 1}, {"戒指", 1}, {"项链", 1},
	{"手镯", 1}, {"耳环", 1}, {"护符", 1}, {"靴子", 1}, {"手套", 1},
	{"帽子", 1}, {"王冠", 3}, {"腰带", 1}, {"典籍", 2}, {"水晶球", 2},
	{"扇", 1}, {"笛", 1}, {"镜", 2}, {"灯", 1}, {"钥匙", 1},
	{"太刀", 3}, {"薙刀", 2}, {"镰刀", 2}, {"爪", 1}, {"鞭", 2},
	{"圣杯", 1}, {"香炉", 1}, {"念珠", 1}, {"羽衣", 2}, {"勾玉", 2},
}

var prefixesZH = []prefix{
	{"陈旧的", -1, 0.7}, {"锈蚀的", -1, 0.6}, {"粗糙的", -1, 0.7},
	{"破损的", -1, 0.5}, {"污浊的", -1, 0.7}, {"腐朽的", -1, 0.6},
	{"精磨的", 0, 1.1}, {"坚固的", 0, 1.1}, {"锋利的", 0, 1.2},
	{"华美的", 0, 1.0}, {"轻盈的", 0, 1.0}, {"厚重的", 0, 1.1}, {"精工的", 0, 1.1},
	{"凛然的", 0, 1.1},
	{"炎之", 1, 1.3}, {"冰之", 1, 1.3}, {"雷之", 1, 1.3}, {"风之", 1, 1.2},
	{"水之", 1, 1.2}, {"大地之", 1, 1.2}, {"光之", 1, 1.3}, {"暗之", 1, 1.3},
	{"神圣的", 1, 1.4}, {"被诅咒的", 1, 1.4}, {"闪耀的", 1, 1.3},
	{"冻彻的", 1, 1.3}, {"燃烧的", 1, 1.4}, {"暴风之", 1, 1.3},
	{"幻影的", 1, 1.3}, {"魔性的", 1, 1.3},
	{"苍蓝的", 1, 1.2}, {"绯红的", 1, 1.2}, {"黄金的", 1, 1.3},
	{"白银的", 1, 1.3}, {"漆黑的", 1, 1.3},
	{"传说的", 2, 1.5}, {"神明的", 2, 1.6}, {"起源的", 2, 1.5},
	{"终焉的", 2, 1.5}, {"永恒的", 2, 1.6}, {"禁忌的", 2, 1.5}, {"至高的", 2, 1.6},
	{"沉睡的", 2, 1.4}, {"被遗忘的", 2, 1.5}, {"觉醒的", 2, 1.4},
}

var suffixesZH = []suffix{
	{"碎片", -1, -2}, {"残骸", -1, -3}, {"仿品", -1, -1},
	{"原型", 0, 0}, {"试作", 0, 0}, {"之证", 0, 1},
	{"·改", 0, 2}, {"·真", 0, 3}, {"·极", 1, 4}, {"·天", 1, 5}, {"·零", 0, 2},
	{"·壹", 0, 1}, {"·陆", 0, 2},
	{"─觉醒─", 1, 5}, {"─幻影─", 1, 4}, {"─残光─", 0, 3},
	{"─咆哮─", 1, 5}, {"─黎明─", 1, 4}, {"─黄昏─", 0, 3},
	{"─深渊─", 1, 5}, {"─追忆─", 0, 3}, {"─胎动─", 0, 3},
	{"─轮回─", 1, 5}, {"─天命─", 1, 4}, {"─因果─", 0, 3},
	{"+1", 0, 1}, {"+2", 0, 2}, {"+3", 0, 3}, {"+4", 0, 4}, {"+5", 1, 5},
	{"[圣]", 1, 4}, {"[魔]", 1, 4}, {"[王]", 1, 5}, {"[天]", 1, 5}, {"[地]", 1, 4}, {"[人]", 0, 3},
	{"杰作", 1, 4}, {"祝福", 1, 3}, {"守护", 0, 3}, {"诅咒", 1, 4},
	{"余韵", 0, 2},
}

var elementsZH = []element{
	{"火", 0, 2}, {"水", 0, 2}, {"雷", 0, 3}, {"冰", 0, 2}, {"风", 0, 1},
	{"土", 0, 1}, {"光", 1, 3}, {"暗", 1, 3}, {"毒", 0, 2}, {"圣", 1, 3},
	{"时", 1, 4}, {"空", 1, 3}, {"梦", 0, 2}, {"命", 1, 4}, {"无", 0, 1},
}

var slotMapZH = map[string]VisualSlot{
	"剑": SlotWeaponRight, "大剑": SlotWeaponRight, "短剑": SlotWeaponRight,
	"枪": SlotWeaponRight, "斧": SlotWeaponRight, "太刀": SlotWeaponRight,
	"薙刀": SlotWeaponRight, "镰刀": SlotWeaponRight, "鞭": SlotWeaponRight, "爪": SlotWeaponRight,
	"盾": SlotWeaponLeft, "镜": SlotWeaponLeft, "典籍": SlotWeaponLeft, "水晶球": SlotWeaponLeft,
	"弓": SlotRanged, "杖": SlotRanged, "锡杖": SlotRanged, "扇": SlotRanged, "笛": SlotRanged,
	"头盔": SlotHead, "帽子": SlotHead, "王冠": SlotHead,
	"铠甲": SlotBody, "法袍": SlotBody, "斗篷": SlotBody, "羽衣": SlotBody,
	"戒指": SlotAccessory, "项链": SlotAccessory, "手镯": SlotAccessory,
	"耳环": SlotAccessory, "护符": SlotAccessory, "靴子": SlotAccessory,
	"手套": SlotAccessory, "腰带": SlotAccessory, "灯": SlotAccessory,
	"钥匙": SlotAccessory, "圣杯": SlotAccessory, "香炉": SlotAccessory,
	"念珠": SlotAccessory, "勾玉": SlotAccessory,
}
