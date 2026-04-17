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

// sanitizeTmux prevents tmux format injection by breaking #[ sequences.
func sanitizeTmux(s string) string {
	return strings.ReplaceAll(s, "#[", "# [")
}

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
		return fmt.Sprintf("%s%s %s", sparkle, sprite, sanitizeTmux(p.LastEvent))
	}

	equip := ""
	if p.Equipped != nil {
		name := sanitizeTmux(p.Equipped.Name)
		color := p.Equipped.Rarity.TmuxColor(phase)
		if color != "" {
			equip = " " + color + name + tmuxReset
		} else {
			equip = " " + name
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

	const W = 50
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
		fmt.Sprintf("%s %s", T("age"), formatDuration(age)),
		genDetail(p),
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
	row(&b, W, fmt.Sprintf("  %s", T("record_header")))
	row(&b, W, "")
	winRate := 0
	if p.Wins+p.Losses > 0 {
		winRate = p.Wins * 100 / (p.Wins + p.Losses)
	}
	row(&b, W, fmt.Sprintf("  %-6s %-5d   %-6s %-5d   %-4s %3d%%",
		T("wins"), p.Wins, T("losses"), p.Losses, T("win_rate"), winRate))
	row(&b, W, fmt.Sprintf("  %-6s %-5d   %-6s %-5d %s",
		T("total_xp"), p.TotalXP, T("found"), p.ItemsFound, T("items_suffix")))

	if p.BestItem != nil {
		row(&b, W, "")
		row(&b, W, fmt.Sprintf("  🏆 %s: %s %s", T("best"), p.BestItem.Rarity, p.BestItem.Name))
	}

	// Equipped
	b.WriteString(hr + "\n")
	if p.Equipped != nil {
		row(&b, W, fmt.Sprintf("  🗡️  %s: %s %s (%s:%d)",
			T("equip_label"), p.Equipped.Rarity, p.Equipped.Name, T("pow_label"), p.Equipped.Power))
	} else {
		row(&b, W, fmt.Sprintf("  🗡️  %s: (%s)", T("equip_label"), T("none")))
	}

	// Inventory (show last 5 items)
	b.WriteString(hr + "\n")
	row(&b, W, fmt.Sprintf("  🎒 %s (%s)  %s",
		T("inventory"), fmt.Sprintf(T("items_count"), len(p.Inventory)), T("inv_keys")))
	row(&b, W, "")
	if len(p.Inventory) == 0 {
		row(&b, W, fmt.Sprintf("    (%s)", T("none")))
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
			row(&b, W, fmt.Sprintf("  %s%s %s (%s:%d)", marker, item.Rarity, item.Name, T("pow_label"), item.Power))
		}
		if len(p.Inventory) > 5 {
			row(&b, W, fmt.Sprintf("    "+T("others_fmt"), len(p.Inventory)-5))
		}
	}

	// Event log
	b.WriteString(hr + "\n")
	row(&b, W, fmt.Sprintf("  %s", T("log_header")))
	row(&b, W, "")
	if len(p.EventLog) == 0 {
		row(&b, W, fmt.Sprintf("    (%s)", T("nothing_yet")))
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

	row(&b, W, fmt.Sprintf("  🎒 %s (%s)  %s",
		T("inventory"), fmt.Sprintf(T("all_items"), len(p.Inventory)), T("inv_full_keys")))
	b.WriteString(hr + "\n")

	if p.Equipped != nil {
		row(&b, W, fmt.Sprintf("  🗡️  %s: %s %s (%s:%d)",
			T("equip_label"), p.Equipped.Rarity, p.Equipped.Name, T("pow_label"), p.Equipped.Power))
	} else {
		row(&b, W, fmt.Sprintf("  🗡️  %s: (%s)", T("equip_label"), T("none")))
	}
	b.WriteString(hr + "\n")

	if len(p.Inventory) == 0 {
		row(&b, W, fmt.Sprintf("    (%s)", T("none")))
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
		row(&b, W, fmt.Sprintf("  %s%s %s (%s:%d)", marker, item.Rarity, item.Name, T("pow_label"), item.Power))
	}

	if total > pageSize {
		row(&b, W, "")
		row(&b, W, fmt.Sprintf("  %d-%d / "+T("items_count"), scroll+1, end, total))
	}

	b.WriteString(hr + "\n")
	b.WriteString("\n")
	return b.String()
}

// EventLog returns recent event log (standalone).
func EventLog(p *Pet) string {
	var b strings.Builder
	fmt.Fprintf(&b, "\n  %s\n", T("log_header"))
	b.WriteString("  ───────────────\n")
	if len(p.EventLog) == 0 {
		fmt.Fprintf(&b, "  (%s)\n", T("nothing_yet"))
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

// row writes a padded line.
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
		return fmt.Sprintf("%s %d (%d+%d)", T("power"), p.TotalPower(), p.Power, p.Equipped.Power)
	}
	return fmt.Sprintf("%s %d", T("power"), p.TotalPower())
}

func genLabel(p *Pet) string {
	if p.Generation > 1 {
		return fmt.Sprintf(T("gen_nth"), p.Generation)
	}
	return "      "
}

func genDetail(p *Pet) string {
	return fmt.Sprintf(T("gen_nth"), p.Generation) + legacyStr(p)
}

func equippedShort(p *Pet) string {
	if p.Equipped != nil {
		return fmt.Sprintf("%s %s", T("equipment"), p.Equipped.Rarity)
	}
	return ""
}

func legacyStr(p *Pet) string {
	if p.Legacy > 0 {
		return fmt.Sprintf(T("legacy_fmt"), p.Legacy)
	}
	return ""
}

// padRight pads a string with spaces to the given display width.
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
	if r >= 0x4E00 && r <= 0x9FFF {
		return true
	}
	if r >= 0xF900 && r <= 0xFAFF {
		return true
	}
	if r >= 0x3040 && r <= 0x309F {
		return true
	}
	if r >= 0x30A0 && r <= 0x30FF {
		return true
	}
	if r >= 0xFF00 && r <= 0xFFEF {
		return true
	}
	if r >= 0x3000 && r <= 0x303F {
		return true
	}
	if r >= 0x3200 && r <= 0x32FF {
		return true
	}
	if r >= 0x3300 && r <= 0x33FF {
		return true
	}
	if r >= 0x2500 && r <= 0x257F {
		return false
	}
	if r >= 0x2800 && r <= 0x28FF {
		return false
	}
	if r >= 0x1F300 && r <= 0x1FAFF {
		return true
	}
	if r >= 0x2600 && r <= 0x27BF {
		return true
	}
	if r >= 0xAC00 && r <= 0xD7AF {
		return true
	}
	if r >= 0xFF65 && r <= 0xFF9F {
		return false
	}
	if r >= 0xFF01 && r <= 0xFF60 {
		return true
	}
	return false
}
