# dotpet

```
                                                                            
     /\  /\        __  __         .--.        ᓫ    ᓫ     ᗰ .--. ᗰ       /|       
    ( *ω*  )     /( ·ᴥ· )\      ( ᗜΘᗜ )      (  ᘏ   )    ( ·ᴥ·    )   .------.   
    (      )~     (      )~     /)    (\      |    |      |      |    >( >| ◉  )  
     \    /       ( |||| )     (/      \)   (>| () |<)    | \  / |    '------'   
      U  U          U  U          ^^          ^    ^      |_/  \_|       \|       
       cat           dog          bird         frog         bear          fish     
                                                                            
```

> An idle terminal pet that lives in your tmux status bar — battling, collecting items, leveling up, and reincarnating on its own.

```
ᓚᘏᗢ~ Lv.3 Iron Sword    ᐕᘏᐳ~ Lv.5 Silver Shield    ᗜΘᗜ> Lv.2 Wood Staff
ᓫᘏᓫ/ Lv.8 Gold Ring     ᗰᘏᗰノ Lv.6 Shadow Axe      ᗱᘏ>。Lv.2 Aqua Orb
```

## Features

- **Status bar resident** — animated silhouette kaomoji in the right side of tmux
- **Fully idle** — XP gain, battles, and item discovery happen automatically over time
- **550M item combos** — 65 materials × 40 types × prefixes × suffixes × elements, procedurally generated
- **6-tier rarity** — ★ Common → ★★★★★★ Mythic (color-coded display)
- **Equipment overlay** — weapons and armor reflected on pet ASCII art, rarity effects around it
- **Manual equipment** — browse inventory with cursor, equip from the list
- **Reincarnation** — auto-reincarnate at Lv.20, next generation starts with a legacy bonus
- **6 species** — determined by hostname × username hash
- **3 languages** — Japanese, English, Chinese — auto-detected from `LANG`

## Demo

### Status bar

```
ᓚᘏᗢ~ Lv.8 Dragon Bone Greatsword -Awoken-
```

### Detail view (`Prefix+P`)

```
  🐱  Maru-2
  ────────────────────────────────────────────────────
              ▽
        /\  /\      Lv.8  / 20
       ( *ω*  )     XP [████████░░] 82/92
       (      )~    Power 14 (8+6)
        \    /      Age 2d 5h
         U  U       Gen 1
                    Equip ★★★
  ────────────────────────────────────────────────────
    ⚔️  Record

    Wins   58      Losses 3       Rate  95%
    TotalXP 271     Found  22    items

    🏆 Best: ★★★ Mithril Greatsword
  ────────────────────────────────────────────────────
```

## Setup

```bash
go build -o ~/.local/bin/dotpet .
```

### tmux config

```tmux
# Show pet in status bar
set -g status-right "#(~/.local/bin/dotpet)"

# Prefix+P opens detail popup
bind P display-popup -w 58 -h 80% -E "~/.local/bin/dotpet watch"

# Click status bar to open popup
bind -n MouseDown1StatusRight display-popup -w 58 -h 80% -E "~/.local/bin/dotpet watch"
```

## Usage

```bash
dotpet          # status line (1 line, for tmux)
dotpet status   # detail view
dotpet watch    # interactive mode
dotpet log      # adventure log
dotpet reset    # reset pet
```

### Language

Language is auto-detected from `LANG` / `LC_ALL` environment variable.

| Variable | Language |
|----------|----------|
| `ja_*` | 日本語 |
| `en_*` | English |
| `zh_*` | 中文 |

All text including item names, monster names, and UI labels adapts to the detected language. Items keep the name from the language they were generated in.

### Watch mode

| Key | Action |
|-----|--------|
| `1` / `Esc` | Back to status |
| `2` / `i` | Inventory |
| `j` / `k` | Cursor move |
| `e` / `Enter` | Equip item |
| `q` | Quit |

<details>
<summary>Game System</summary>

### Level and XP

1 XP per minute elapsed. Required XP scales quadratically (Lv.1→2: 15XP, Lv.19→20: 435XP). Roughly 3 days to reach Lv.20.

### Adventure events

Checked every 10 minutes:

| Event | Chance |
|-------|--------|
| Battle | 25% |
| Item find | 10% |
| Special event | 10% |
| Nothing | 55% |

### Item generation

Pattern varies by language:
- **JA:** `[接頭辞][素材]の[種類][接尾辞]〈属性〉`
- **EN:** `[Prefix] [Material] [Type] [Suffix] <Element>`
- **ZH:** `[前缀][素材][种类][后缀]〈属性〉`

- 65 materials — Wood, Iron, Mithril, Dragon Bone, Chaos Core...
- 40 types — Sword, Staff, Armor, Ring...
- 45 prefixes — Old, Flame, Divine...
- 40 suffixes — Mk.II, -Awoken-, [Holy]...
- 15 elements — Fire, Ice, Light, Dark, Time...

### Rarity

| Rarity | Color |
|--------|-------|
| ★ Common | — |
| ★★ Fine | White |
| ★★★ Rare | Blue |
| ★★★★ Epic | Purple |
| ★★★★★ Legendary | Gold |
| ★★★★★★ Mythic | Rainbow |

High-rarity equipment adds sparkle effects around the pet. Mythic tier animates in rainbow colors.

### Reincarnation

Auto-reincarnate at Lv.20. Next generation gets a legacy bonus (+1 power) and starts from Lv.1.

</details>

## Data

```
~/.config/dotpet/pet.json
```
