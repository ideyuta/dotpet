package main

import (
	"strings"
	"time"
)

const (
	artRows = 7
	artCols = 14
)

// animPhase returns 0-3 on a ~3.3s cycle:
//
//	0: idle A  (1000ms)
//	1: idle B  (1000ms)
//	2: blink   (300ms)
//	3: idle A' (1000ms)
func animPhase() int {
	ms := time.Now().UnixMilli()
	cycle := ms % 3300
	switch {
	case cycle < 1000:
		return 0
	case cycle < 2000:
		return 1
	case cycle < 2300:
		return 2
	default:
		return 3
	}
}

// renderPetArt returns the fully composed 7-line art for the pet,
// including animation phase, equipment overlay, and rarity effects.
func renderPetArt(p *Pet) []string {
	info := speciesData[p.Species]
	phase := animPhase()

	// 1. Pick base art for this phase
	art := baseArt(&info.Art, phase)

	// 2. Apply equipment overlay
	if p.Equipped != nil {
		applyEquipment(art, &info.Art.Anchors, p.Equipped)
	}

	// 3. Apply rarity effects
	if p.Equipped != nil {
		applyRarityEffects(art, p.Equipped.Rarity, phase)
	}

	// 4. Convert back to strings
	return runeGridToStrings(art)
}

// baseArt builds a [][]rune grid from the species art for the given phase.
func baseArt(pa *PetArt, phase int) [][]rune {
	var lines []string
	switch phase {
	case 0, 3:
		lines = pa.Base
	case 1:
		lines = pa.IdleB
	case 2:
		// Blink: use Base with eyes closed
		lines = pa.Base
	}

	grid := make([][]rune, artRows)
	for i := 0; i < artRows; i++ {
		if i < len(lines) {
			grid[i] = []rune(lines[i])
		} else {
			grid[i] = []rune(strings.Repeat(" ", artCols))
		}
		// Pad to artCols
		for len(grid[i]) < artCols {
			grid[i] = append(grid[i], ' ')
		}
	}

	// Apply blink (replace open eyes with closed eyes)
	if phase == 2 && pa.EyeOpen != "" && pa.EyeClose != "" {
		applyBlink(grid, pa.EyeOpen, pa.EyeClose)
	}

	return grid
}

// applyBlink replaces the eye-open substring with eye-closed in the grid.
func applyBlink(grid [][]rune, eyeOpen, eyeClose string) {
	openRunes := []rune(eyeOpen)
	closeRunes := []rune(eyeClose)
	if len(openRunes) != len(closeRunes) {
		return
	}
	for r := range grid {
		for c := 0; c <= len(grid[r])-len(openRunes); c++ {
			match := true
			for k, or := range openRunes {
				if grid[r][c+k] != or {
					match = false
					break
				}
			}
			if match {
				for k, cr := range closeRunes {
					grid[r][c+k] = cr
				}
				return
			}
		}
	}
}

// ---------- Equipment Overlay ----------

// Equipment overlay characters by slot and rarity tier.
type slotOverlay struct {
	Chars    [2]string // [normal, high-rarity] variants
	AnchorFn func(anchors *PetAnchors) ArtAnchor
}

var slotOverlays = map[VisualSlot]slotOverlay{
	SlotWeaponRight: {
		Chars:    [2]string{"/ ", "† "},
		AnchorFn: func(a *PetAnchors) ArtAnchor { return a.RightHand },
	},
	SlotWeaponLeft: {
		Chars:    [2]string{"[]", "◇ "},
		AnchorFn: func(a *PetAnchors) ArtAnchor { return a.LeftHand },
	},
	SlotRanged: {
		Chars:    [2]string{") ", "⌒ "},
		AnchorFn: func(a *PetAnchors) ArtAnchor { return a.RightHand },
	},
	SlotHead: {
		Chars:    [2]string{"▽", "♦ "},
		AnchorFn: func(a *PetAnchors) ArtAnchor { return a.Head },
	},
	SlotBody: {
		Chars:    [2]string{"##", "%%"},
		AnchorFn: func(a *PetAnchors) ArtAnchor { return a.Body },
	},
}

func applyEquipment(grid [][]rune, anchors *PetAnchors, item *Item) {
	slot := ItemVisualSlot(item.Name)
	overlay, ok := slotOverlays[slot]
	if !ok {
		return
	}

	anchor := overlay.AnchorFn(anchors)
	if anchor.Row < 0 || anchor.Row >= artRows {
		return
	}

	// Pick variant based on rarity (Epic+ gets fancier chars)
	variant := 0
	if item.Rarity >= Epic {
		variant = 1
	}
	chars := []rune(overlay.Chars[variant])

	row := grid[anchor.Row]
	for i, ch := range chars {
		col := anchor.Col + i
		if col >= 0 && col < len(row) {
			row[col] = ch
		}
	}
}

// ---------- Rarity Effects ----------

// Rarity effect definitions: characters placed around the pet art.
type rarityEffect struct {
	Chars   []rune // effect characters to cycle through
	Count   int    // number of effect chars to place
	Surround bool  // true = place around edges, false = just a few spots
}

var rarityEffects = map[Rarity]rarityEffect{
	Rare:      {Chars: []rune{'+', ' '}, Count: 2, Surround: false},
	Epic:      {Chars: []rune{'*', '·'}, Count: 3, Surround: false},
	Legendary: {Chars: []rune{'✦', '✧', '·'}, Count: 5, Surround: false},
	Mythic:    {Chars: []rune{'✦', '✧', '★', '·'}, Count: 0, Surround: true},
}

// Predefined sparkle positions (row, col) for non-surround effects.
var sparklePositions = []ArtAnchor{
	{0, 1}, {0, 12}, {0, 7},
	{6, 2}, {6, 11}, {6, 6},
	{3, 0}, {3, 13},
	{1, 0}, {5, 13},
}

func applyRarityEffects(grid [][]rune, rarity Rarity, phase int) {
	effect, ok := rarityEffects[rarity]
	if !ok {
		return
	}

	if effect.Surround {
		// Mythic: fill empty edge positions with cycling chars
		applySurroundEffect(grid, effect.Chars, phase)
		return
	}

	// Place sparkle chars at predefined positions
	for i := 0; i < effect.Count && i < len(sparklePositions); i++ {
		pos := sparklePositions[i]
		if pos.Row >= 0 && pos.Row < len(grid) && pos.Col >= 0 && pos.Col < len(grid[pos.Row]) {
			if grid[pos.Row][pos.Col] == ' ' {
				// Cycle through effect chars based on phase and position
				ci := (phase + i) % len(effect.Chars)
				grid[pos.Row][pos.Col] = effect.Chars[ci]
			}
		}
	}
}

func applySurroundEffect(grid [][]rune, chars []rune, phase int) {
	if len(chars) == 0 {
		return
	}
	idx := 0
	// Top and bottom rows
	for _, r := range []int{0, 6} {
		if r >= len(grid) {
			continue
		}
		for c := 0; c < len(grid[r]); c++ {
			if grid[r][c] == ' ' {
				ci := (phase + idx) % len(chars)
				grid[r][c] = chars[ci]
				idx++
			}
		}
	}
	// Left and right edges of middle rows
	for r := 1; r < 6 && r < len(grid); r++ {
		if len(grid[r]) > 0 && grid[r][0] == ' ' {
			ci := (phase + idx) % len(chars)
			grid[r][0] = chars[ci]
			idx++
		}
		last := len(grid[r]) - 1
		if last > 0 && grid[r][last] == ' ' {
			ci := (phase + idx) % len(chars)
			grid[r][last] = chars[ci]
			idx++
		}
	}
}

// ---------- Helpers ----------

func runeGridToStrings(grid [][]rune) []string {
	result := make([]string, len(grid))
	for i, row := range grid {
		result[i] = string(row)
	}
	return result
}
