package main

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// Tab represents a detail view tab.
type Tab int

const (
	TabStatus    Tab = 0
	TabInventory Tab = 1
	TabLog       Tab = 2
)

func petSprite(p *Pet) string {
	info := speciesData[p.Species]
	phase := animPhase()
	return info.Frames[phase]
}

const tmuxReset = "#[fg=#858585]"

// statusLineSparkle returns animated rarity decoration for tmux with color.
func statusLineSparkle(rarity Rarity, phase int) string {
	color := rarity.TmuxColor(phase)
	if color == "" {
		return ""
	}
	switch {
	case rarity >= Mythic:
		s := []string{"✦✧★", "✧★✦", "★✦✧", "✦★✧"}
		return color + s[phase] + tmuxReset
	case rarity >= Legendary:
		s := []string{"✦ ", "✧ ", "✦ ", " ✧"}
		return color + s[phase] + tmuxReset
	case rarity >= Epic:
		s := []string{"* ", "· ", "* ", "· "}
		return color + s[phase] + tmuxReset
	case rarity >= Rare:
		s := []string{"+ ", "  ", "+ ", "  "}
		return color + s[phase] + tmuxReset
	default:
		return ""
	}
}

// StatusLine returns a single-line string for tmux status bar.
func StatusLine(p *Pet) string {
	phase := animPhase()
	sprite := speciesData[p.Species].Frames[phase]

	sparkle := ""
	if p.Equipped != nil {
		sparkle = statusLineSparkle(p.Equipped.Rarity, phase)
	}

	if p.LastEvent != "" && time.Since(p.EventAt) < 5*time.Second {
		return fmt.Sprintf("%s%s %s", sparkle, sprite, p.LastEvent)
	}

	equip := ""
	if p.Equipped != nil {
		color := p.Equipped.Rarity.TmuxColor(phase)
		if color != "" {
			equip = " " + color + p.Equipped.Name + tmuxReset
		} else {
			equip = " " + p.Equipped.Name
		}
	}
	gen := ""
	if p.Generation > 1 {
		gen = fmt.Sprintf(" G%d", p.Generation)
	}
	return fmt.Sprintf("%s%s Lv.%d%s%s", sparkle, sprite, p.Level, equip, gen)
}

// StatusDetail returns a rich multi-line detailed status.
func StatusDetail(p *Pet) string {
	info := speciesData[p.Species]
	art := renderPetArt(p)
	age := time.Since(p.BornAt)

	var b strings.Builder

	const W = 50 // inner width in display columns
	hr := "  " + "─" + strings.Repeat("─", W) + "─"

	b.WriteString("\n")

	// Header
	row(&b, W, fmt.Sprintf("  %s  %s  %s", info.Emoji, p.Name, genLabel(p)))
	b.WriteString(hr + "\n")

	// Art + Stats (7 rows)
	stats := []string{
		"",
		fmt.Sprintf("Lv.%-2d / %d", p.Level, maxLevel),
		fmt.Sprintf("XP %s %d/%d", xpBar(p), p.XP, p.XPToNext()),
		powerLine(p),
		fmt.Sprintf("年齢 %s", formatDuration(age)),
		fmt.Sprintf("世代 第%d世代%s", p.Generation, legacyStr(p)),
		equippedShort(p),
	}
	for i := 0; i < artRows; i++ {
		artLine := ""
		if i < len(art) {
			artLine = art[i]
		}
		stat := ""
		if i < len(stats) {
			stat = stats[i]
		}
		left := padRight(artLine, 14)
		right := stat
		combined := fmt.Sprintf("  %s  %s", left, right)
		row(&b, W, combined)
	}

	// Battle stats
	b.WriteString(hr + "\n")
	row(&b, W, "  ⚔️  戦績")
	row(&b, W, "")
	winRate := 0
	if p.Wins+p.Losses > 0 {
		winRate = p.Wins * 100 / (p.Wins + p.Losses)
	}
	row(&b, W, fmt.Sprintf("  勝利 %-5d   敗北 %-5d   勝率 %3d%%", p.Wins, p.Losses, winRate))
	row(&b, W, fmt.Sprintf("  総XP %-5d   発見 %-5d アイテム", p.TotalXP, p.ItemsFound))

	if p.BestItem != nil {
		row(&b, W, "")
		row(&b, W, fmt.Sprintf("  🏆 最高: %s %s", p.BestItem.Rarity, p.BestItem.Name))
	}

	// Equipped
	b.WriteString(hr + "\n")
	if p.Equipped != nil {
		row(&b, W, fmt.Sprintf("  🗡️  装備: %s %s (力:%d)", p.Equipped.Rarity, p.Equipped.Name, p.Equipped.Power))
	} else {
		row(&b, W, "  🗡️  装備: (なし)")
	}

	// Inventory (show last 5 items; full list in inventory tab)
	b.WriteString(hr + "\n")
	row(&b, W, fmt.Sprintf("  🎒 もちもの (%d件)  [2:一覧]", len(p.Inventory)))
	row(&b, W, "")
	if len(p.Inventory) == 0 {
		row(&b, W, "    (なし)")
	} else {
		show := p.Inventory
		if len(show) > 5 {
			show = show[len(show)-5:]
		}
		for _, item := range show {
			marker := "  "
			if p.Equipped != nil && item.Name == p.Equipped.Name && item.Power == p.Equipped.Power {
				marker = "→ "
			}
			row(&b, W, fmt.Sprintf("  %s%s %s (力:%d)", marker, item.Rarity, item.Name, item.Power))
		}
		if len(p.Inventory) > 5 {
			row(&b, W, fmt.Sprintf("    ...他%d件", len(p.Inventory)-5))
		}
	}

	// Event log
	b.WriteString(hr + "\n")
	row(&b, W, "  📜 冒険の記録")
	row(&b, W, "")
	if len(p.EventLog) == 0 {
		row(&b, W, "    (まだなにもない)")
	} else {
		start := 0
		if len(p.EventLog) > 8 {
			start = len(p.EventLog) - 8
		}
		for _, entry := range p.EventLog[start:] {
			row(&b, W, "  "+entry)
		}
	}

	b.WriteString(hr + "\n")
	b.WriteString("\n")

	return b.String()
}

// SortedInventory returns a copy of the inventory sorted by rarity (desc), then power (desc).
func SortedInventory(p *Pet) []Item {
	sorted := make([]Item, len(p.Inventory))
	copy(sorted, p.Inventory)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Rarity != sorted[j].Rarity {
			return sorted[i].Rarity > sorted[j].Rarity
		}
		return sorted[i].Power > sorted[j].Power
	})
	return sorted
}

// InventoryView returns the inventory screen as a string.
func InventoryView(p *Pet, scroll, cursor int) string {
	const W = 50
	hr := "  " + "─" + strings.Repeat("─", W) + "─"

	var b strings.Builder
	b.WriteString("\n")

	row(&b, W, fmt.Sprintf("  🎒 もちもの (全%d件)  [j/k:選択] [e:装備] [1:戻る]", len(p.Inventory)))
	b.WriteString(hr + "\n")

	if p.Equipped != nil {
		row(&b, W, fmt.Sprintf("  🗡️  装備: %s %s (力:%d)", p.Equipped.Rarity, p.Equipped.Name, p.Equipped.Power))
	} else {
		row(&b, W, "  🗡️  装備: (なし)")
	}
	b.WriteString(hr + "\n")

	if len(p.Inventory) == 0 {
		row(&b, W, "    (なし)")
		b.WriteString(hr + "\n")
		b.WriteString("\n")
		return b.String()
	}

	sorted := SortedInventory(p)

	const pageSize = 15
	total := len(sorted)
	if scroll < 0 {
		scroll = 0
	}
	if scroll > total-1 {
		scroll = total - 1
	}
	end := scroll + pageSize
	if end > total {
		end = total
	}

	for idx, item := range sorted[scroll:end] {
		absIdx := scroll + idx
		marker := "  "
		if absIdx == cursor {
			marker = "▶ "
		} else if p.Equipped != nil && item.Name == p.Equipped.Name && item.Power == p.Equipped.Power {
			marker = "→ "
		}
		row(&b, W, fmt.Sprintf("  %s%s %s (力:%d)", marker, item.Rarity, item.Name, item.Power))
	}

	if total > pageSize {
		row(&b, W, "")
		row(&b, W, fmt.Sprintf("  %d-%d / %d件", scroll+1, end, total))
	}

	b.WriteString(hr + "\n")
	b.WriteString("\n")
	return b.String()
}

// EventLog returns recent event log (standalone).
func EventLog(p *Pet) string {
	var b strings.Builder
	b.WriteString("\n  📜 冒険の記録\n")
	b.WriteString("  ───────────────\n")
	if len(p.EventLog) == 0 {
		b.WriteString("  (まだなにもない)\n")
	} else {
		start := 0
		if len(p.EventLog) > 15 {
			start = len(p.EventLog) - 15
		}
		for _, entry := range p.EventLog[start:] {
			fmt.Fprintf(&b, "  %s\n", entry)
		}
	}
	b.WriteString("\n")
	return b.String()
}

// row writes a padded line with no box-drawing borders.
func row(b *strings.Builder, width int, content string) {
	b.WriteString("  ")
	b.WriteString(padRight(content, width))
	b.WriteString("\n")
}

func xpBar(p *Pet) string {
	needed := p.XPToNext()
	if needed == 0 {
		return "[██████████]"
	}
	filled := p.XP * 10 / needed
	if filled > 10 {
		filled = 10
	}
	empty := 10 - filled
	return "[" + strings.Repeat("█", filled) + strings.Repeat("░", empty) + "]"
}

func formatDuration(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	if days > 0 {
		return fmt.Sprintf("%dd %dh", days, hours)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, int(d.Minutes())%60)
	}
	return fmt.Sprintf("%dm", int(d.Minutes()))
}

func powerLine(p *Pet) string {
	if p.Equipped != nil {
		return fmt.Sprintf("つよさ %d (%d+%d)", p.TotalPower(), p.Power, p.Equipped.Power)
	}
	return fmt.Sprintf("つよさ %d", p.TotalPower())
}

func genLabel(p *Pet) string {
	if p.Generation > 1 {
		return fmt.Sprintf("第%d世代", p.Generation)
	}
	return "      "
}

func equippedShort(p *Pet) string {
	if p.Equipped != nil {
		return fmt.Sprintf("装備 %s", p.Equipped.Rarity)
	}
	return ""
}

func legacyStr(p *Pet) string {
	if p.Legacy > 0 {
		return fmt.Sprintf(" (遺産+%d)", p.Legacy)
	}
	return ""
}

// padRight pads a string with spaces to the given display width.
// Accounts for wide (CJK) characters and emoji taking 2 columns.
func padRight(s string, width int) string {
	w := displayWidth(s)
	if w >= width {
		return s
	}
	return s + strings.Repeat(" ", width-w)
}

func displayWidth(s string) int {
	w := 0
	inEsc := false
	for _, r := range s {
		if r == '\033' {
			inEsc = true
			continue
		}
		if inEsc {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inEsc = false
			}
			continue
		}
		switch {
		case r >= 0x1100 && isWide(r):
			w += 2
		default:
			w++
		}
	}
	return w
}

func isWide(r rune) bool {
	// CJK Unified Ideographs
	if r >= 0x4E00 && r <= 0x9FFF {
		return true
	}
	// CJK Compatibility Ideographs
	if r >= 0xF900 && r <= 0xFAFF {
		return true
	}
	// Hiragana
	if r >= 0x3040 && r <= 0x309F {
		return true
	}
	// Katakana
	if r >= 0x30A0 && r <= 0x30FF {
		return true
	}
	// Halfwidth and Fullwidth Forms
	if r >= 0xFF00 && r <= 0xFFEF {
		return true
	}
	// CJK Symbols and Punctuation
	if r >= 0x3000 && r <= 0x303F {
		return true
	}
	// Enclosed CJK Letters
	if r >= 0x3200 && r <= 0x32FF {
		return true
	}
	// CJK Compatibility
	if r >= 0x3300 && r <= 0x33FF {
		return true
	}
	// Box Drawing and Miscellaneous symbols (some are wide)
	if r >= 0x2500 && r <= 0x257F {
		return false // box drawing is narrow
	}
	// Braille
	if r >= 0x2800 && r <= 0x28FF {
		return false
	}
	// Emoji / Symbols
	if r >= 0x1F300 && r <= 0x1FAFF {
		return true
	}
	// Misc symbols
	if r >= 0x2600 && r <= 0x27BF {
		return true
	}
	// Hangul
	if r >= 0xAC00 && r <= 0xD7AF {
		return true
	}
	// Katakana Halfwidth (NOT wide)
	if r >= 0xFF65 && r <= 0xFF9F {
		return false
	}
	// Fullwidth Latin
	if r >= 0xFF01 && r <= 0xFF60 {
		return true
	}
	return false
}
