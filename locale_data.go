package main

// Monster names per language (indices match monsterPowers).
var monstersLocale = map[Lang][]string{
	LangJA: {"スライム", "ゴブリン", "オオカミ", "スケルトン", "オーク", "ドラゴン", "デーモン"},
	LangEN: {"Slime", "Goblin", "Wolf", "Skeleton", "Orc", "Dragon", "Demon"},
	LangZH: {"史莱姆", "哥布林", "狼", "骷髅", "兽人", "龙", "恶魔"},
}

// monsterPowers holds BasePower for each monster (language-independent).
var monsterPowers = []int{1, 3, 5, 7, 10, 14, 18}

// Special event descriptions per language.
var specialEventsLocale = map[Lang][]string{
	LangJA: {
		"回復の泉を見つけた",
		"旅の商人に出会った",
		"古代遺跡を発見した",
		"流れ星を見た",
		"隠し通路を見つけた",
	},
	LangEN: {
		"Found a healing spring",
		"Met a traveling merchant",
		"Discovered ancient ruins",
		"Saw a shooting star",
		"Found a hidden passage",
	},
	LangZH: {
		"发现了恢复之泉",
		"遇到了旅行商人",
		"发现了古代遗迹",
		"看到了流星",
		"发现了隐藏通道",
	},
}

// Pet names per language and species.
var speciesNamesLocale = map[Lang]map[string][]string{
	LangJA: {
		"cat":  {"モチ", "クロ", "タマ", "ハナ", "ソラ", "ユキ", "マル", "ニコ"},
		"dog":  {"ポチ", "シロ", "リク", "ハチ", "コロ", "タロウ", "ゴン", "モモ"},
		"bird": {"ピヨ", "トリ", "カゼ", "ハネ", "ウズ", "スイ", "リン", "フウ"},
		"frog": {"ケロ", "ガマ", "アメ", "ヌマ", "スズ", "ツユ", "ミズ", "カワ"},
		"bear": {"クマ", "ゴロウ", "モリ", "ドン", "ヤマ", "テツ", "ゲンタ", "ダイ"},
		"fish": {"サカナ", "タイ", "ナミ", "ウミ", "シオ", "カイ", "フグ", "コイ"},
	},
	LangEN: {
		"cat":  {"Mochi", "Kuro", "Tama", "Hana", "Sora", "Yuki", "Maru", "Niko"},
		"dog":  {"Pochi", "Shiro", "Riku", "Hachi", "Koro", "Taro", "Gon", "Momo"},
		"bird": {"Piyo", "Tori", "Kaze", "Hane", "Uzu", "Sui", "Rin", "Fuu"},
		"frog": {"Kero", "Gama", "Ame", "Numa", "Suzu", "Tsuyu", "Mizu", "Kawa"},
		"bear": {"Kuma", "Goro", "Mori", "Don", "Yama", "Tetsu", "Genta", "Dai"},
		"fish": {"Sakana", "Tai", "Nami", "Umi", "Shio", "Kai", "Fugu", "Koi"},
	},
	LangZH: {
		"cat":  {"年糕", "小黑", "小玉", "小花", "天空", "雪儿", "丸子", "妮可"},
		"dog":  {"阿旺", "小白", "小陆", "小八", "滚滚", "太郎", "阿权", "桃子"},
		"bird": {"小啾", "鸟鸟", "风儿", "羽毛", "漩涡", "翠翠", "铃铃", "清风"},
		"frog": {"呱呱", "蛤蟆", "雨滴", "沼沼", "铃蛙", "露珠", "水水", "河河"},
		"bear": {"熊熊", "咕噜", "森森", "阿咚", "山山", "铁铁", "源太", "大大"},
		"fish": {"鱼鱼", "鲷鲷", "波浪", "海海", "盐盐", "贝贝", "河豚", "锦鲤"},
	},
}
