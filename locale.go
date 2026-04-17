package main

import (
	"os"
	"strings"
)

// Lang represents a supported language.
type Lang string

const (
	LangJA Lang = "ja"
	LangEN Lang = "en"
	LangZH Lang = "zh"
)

var currentLang Lang = LangEN

// InitLocale detects the language from LC_ALL or LANG environment variables.
func InitLocale() {
	env := os.Getenv("LC_ALL")
	if env == "" || env == "C" || env == "POSIX" {
		env = os.Getenv("LANG")
	}
	code := strings.SplitN(env, "_", 2)[0]
	switch code {
	case "ja":
		currentLang = LangJA
	case "zh":
		currentLang = LangZH
	default:
		currentLang = LangEN
	}
}

// T returns the translated string for the given key in the current language.
func T(key string) string {
	if m, ok := messages[currentLang]; ok {
		if s, ok := m[key]; ok {
			return s
		}
	}
	if m, ok := messages[LangEN]; ok {
		if s, ok := m[key]; ok {
			return s
		}
	}
	return key
}

var messages = map[Lang]map[string]string{
	LangJA: {
		// Status labels
		"age":        "年齢",
		"generation": "世代",
		"gen_nth":    "第%d世代",
		"legacy_fmt": " (遺産+%d)",
		"power":      "つよさ",
		"equipment":  "装備",
		"none":       "なし",

		// Battle stats
		"record_header": "⚔️  戦績",
		"wins":          "勝利",
		"losses":        "敗北",
		"win_rate":      "勝率",
		"total_xp":      "総XP",
		"found":         "発見",
		"items_suffix":  "アイテム",
		"best":          "最高",

		// Inventory
		"inventory":       "もちもの",
		"items_count":     "%d件",
		"all_items":       "全%d件",
		"others_fmt":      "...他%d件",
		"inv_keys":        "[2:一覧]",
		"inv_full_keys":   "[j/k:選択] [e:装備] [1:戻る]",
		"equip_label":     "装備",
		"pow_label":       "力",

		// Event log
		"log_header":  "📜 冒険の記録",
		"nothing_yet": "まだなにもない",

		// Events
		"reincarnate": "🔄 転生! 第%d世代へ",
		"level_up":    "✨ Lv.%d!",
		"defeated":    "⚔️ %sを倒した! +%dXP",
		"defeated_lv": "✨ Lv.%d! (%sを倒した)",
		"fled":        "⚔️ %sから逃げた",
		"item_found":  "🌟 %s[%s]を見つけた!",
		"item_got":    "🌟 %s[%s]を手に入れた!",
		"special_fmt": "🔮 %s +%dXP",

		// Rarity
		"rarity_normal":    "ふつう",
		"rarity_fine":      "上質",
		"rarity_rare":      "希少",
		"rarity_epic":      "秘宝",
		"rarity_legendary": "伝説",
		"rarity_mythic":    "神話",

		// Commands
		"reset_msg": "リセット! 新しいペット: %s (%s)\n",
	},
	LangEN: {
		"age":        "Age",
		"generation": "Gen",
		"gen_nth":    "Gen %d",
		"legacy_fmt": " (Legacy+%d)",
		"power":      "Power",
		"equipment":  "Equip",
		"none":       "none",

		"record_header": "⚔️  Record",
		"wins":          "Wins",
		"losses":        "Losses",
		"win_rate":      "Rate",
		"total_xp":      "TotalXP",
		"found":         "Found",
		"items_suffix":  "items",
		"best":          "Best",

		"inventory":       "Inventory",
		"items_count":     "%d",
		"all_items":       "%d items",
		"others_fmt":      "...%d more",
		"inv_keys":        "[2:list]",
		"inv_full_keys":   "[j/k:select] [e:equip] [1:back]",
		"equip_label":     "Equip",
		"pow_label":       "Pow",

		"log_header":  "📜 Adventure Log",
		"nothing_yet": "Nothing yet",

		"reincarnate": "🔄 Reborn! Gen %d",
		"level_up":    "✨ Lv.%d!",
		"defeated":    "⚔️ Defeated %s! +%dXP",
		"defeated_lv": "✨ Lv.%d! (Defeated %s)",
		"fled":        "⚔️ Fled from %s",
		"item_found":  "🌟 Found %s [%s]!",
		"item_got":    "🌟 Got %s [%s]!",
		"special_fmt": "🔮 %s +%dXP",

		"rarity_normal":    "Common",
		"rarity_fine":      "Fine",
		"rarity_rare":      "Rare",
		"rarity_epic":      "Epic",
		"rarity_legendary": "Legendary",
		"rarity_mythic":    "Mythic",

		"reset_msg": "Reset! New pet: %s (%s)\n",
	},
	LangZH: {
		"age":        "年龄",
		"generation": "世代",
		"gen_nth":    "第%d世代",
		"legacy_fmt": " (遗产+%d)",
		"power":      "战力",
		"equipment":  "装备",
		"none":       "无",

		"record_header": "⚔️  战绩",
		"wins":          "胜利",
		"losses":        "败北",
		"win_rate":      "胜率",
		"total_xp":      "总XP",
		"found":         "发现",
		"items_suffix":  "物品",
		"best":          "最佳",

		"inventory":       "物品栏",
		"items_count":     "%d件",
		"all_items":       "共%d件",
		"others_fmt":      "...另有%d件",
		"inv_keys":        "[2:列表]",
		"inv_full_keys":   "[j/k:选择] [e:装备] [1:返回]",
		"equip_label":     "装备",
		"pow_label":       "力",

		"log_header":  "📜 冒险记录",
		"nothing_yet": "还没有记录",

		"reincarnate": "🔄 转生! 第%d世代",
		"level_up":    "✨ Lv.%d!",
		"defeated":    "⚔️ 击败了%s! +%dXP",
		"defeated_lv": "✨ Lv.%d! (击败了%s)",
		"fled":        "⚔️ 从%s逃跑了",
		"item_found":  "🌟 发现了%s [%s]!",
		"item_got":    "🌟 获得了%s [%s]!",
		"special_fmt": "🔮 %s +%dXP",

		"rarity_normal":    "普通",
		"rarity_fine":      "优良",
		"rarity_rare":      "稀有",
		"rarity_epic":      "秘宝",
		"rarity_legendary": "传说",
		"rarity_mythic":    "神话",

		"reset_msg": "重置! 新宠物: %s (%s)\n",
	},
}
