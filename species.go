package main

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"os"
	"os/user"
)

// ArtAnchor is a position in the art grid (row, col in rune indices).
type ArtAnchor struct {
	Row int
	Col int
}

// PetAnchors defines where equipment overlays attach for a species.
type PetAnchors struct {
	Head      ArtAnchor // hat/crown position (row above head)
	LeftHand  ArtAnchor // shield/book position
	RightHand ArtAnchor // weapon position
	Body      ArtAnchor // armor overlay start position
}

// PetArt holds the rich art data for a species.
type PetArt struct {
	Base     []string // 7-line idle pose A
	IdleB    []string // 7-line idle pose B
	EyeOpen  string   // open-eye substring (e.g. "・ω・")
	EyeClose string   // closed-eye substring (e.g. "─ω─")
	Anchors  PetAnchors
}

type SpeciesInfo struct {
	Emoji  string
	Frames [4]string // compact tmux sprites (4 animation phases)
	Art    PetArt
}

var speciesList = []string{"cat", "dog", "bird", "frog", "bear", "fish"}

var speciesData = map[string]SpeciesInfo{
	"cat": {
		Emoji:  "🐱",
		Frames: [4]string{"ᓚᘏᗢ", "ᓚᘏᗢ~", "ᓚᘏ-", "ᓚᘏᗢ"},
		Art: PetArt{
			Base: []string{
				`              `,
				`    /\  /\    `,
				`   ( *ω*  )   `,
				`   (      )~  `,
				`    \    /    `,
				`     U  U    `,
				`              `,
			},
			IdleB: []string{
				`              `,
				`    /\  /\    `,
				`   ( *ω*  )   `,
				`   (      )~  `,
				`    \    /    `,
				`     U  U    `,
				`              `,
			},
			EyeOpen:  "*ω*",
			EyeClose: "-ω-",
			Anchors: PetAnchors{
				Head:      ArtAnchor{Row: 0, Col: 5},
				RightHand: ArtAnchor{Row: 2, Col: 1},
				LeftHand:  ArtAnchor{Row: 2, Col: 12},
				Body:      ArtAnchor{Row: 3, Col: 4},
			},
		},
	},
	"dog": {
		Emoji:  "🐶",
		Frames: [4]string{"ᐕᘏᐳ", "ᐕᘏᐳ~", "ᐕ-ᐳ", "ᐕᘏᐳ"},
		Art: PetArt{
			Base: []string{
				`              `,
				`     __  __   `,
				`   /( ·ᴥ· )\  `,
				`    (      )~ `,
				`    ( |||| )  `,
				`      U  U   `,
				`              `,
			},
			IdleB: []string{
				`              `,
				`     __  __   `,
				`   /( ·ᴥ· )\  `,
				`    (      )~ `,
				`    ( |||| )  `,
				`      U  U   `,
				`              `,
			},
			EyeOpen:  "·ᴥ·",
			EyeClose: "-ᴥ-",
			Anchors: PetAnchors{
				Head:      ArtAnchor{Row: 0, Col: 5},
				RightHand: ArtAnchor{Row: 3, Col: 1},
				LeftHand:  ArtAnchor{Row: 3, Col: 12},
				Body:      ArtAnchor{Row: 4, Col: 4},
			},
		},
	},
	"bird": {
		Emoji:  "🐦",
		Frames: [4]string{"ᗜΘᗜ", "ᗜΘᗜ>", "ᗜ-ᗜ", "ᗜΘᗜ"},
		Art: PetArt{
			Base: []string{
				`              `,
				`      .--.    `,
				`    ( ᗜΘᗜ )  `,
				`    /)    (\  `,
				`   (/      \) `,
				`      ^^     `,
				`              `,
			},
			IdleB: []string{
				`              `,
				`      .--.    `,
				`    ( ᗜΘᗜ )  `,
				`    /)   (\   `,
				`   (/     \)  `,
				`      ^ ^    `,
				`              `,
			},
			EyeOpen:  "ᗜΘᗜ",
			EyeClose: "ᗜ-ᗜ",
			Anchors: PetAnchors{
				Head:      ArtAnchor{Row: 0, Col: 5},
				RightHand: ArtAnchor{Row: 3, Col: 1},
				LeftHand:  ArtAnchor{Row: 3, Col: 12},
				Body:      ArtAnchor{Row: 4, Col: 4},
			},
		},
	},
	"frog": {
		Emoji:  "🐸",
		Frames: [4]string{"ᓫᘏᓫ", "ᓫᘏᓫ/", "ᓫ-ᓫ", "ᓫᘏᓫ^"},
		Art: PetArt{
			Base: []string{
				`              `,
				`    ᓫ    ᓫ   `,
				`   (  ᘏ   )  `,
				`    |    |    `,
				`  (>| () |<) `,
				`    ^    ^   `,
				`              `,
			},
			IdleB: []string{
				`              `,
				`    ᓫ    ᓫ   `,
				`   (  ᘏ   )  `,
				`    |    |    `,
				`  (>| () |<) `,
				`    ^    ^   `,
				`              `,
			},
			EyeOpen:  "ᘏ",
			EyeClose: "-",
			Anchors: PetAnchors{
				Head:      ArtAnchor{Row: 0, Col: 5},
				RightHand: ArtAnchor{Row: 3, Col: 1},
				LeftHand:  ArtAnchor{Row: 3, Col: 12},
				Body:      ArtAnchor{Row: 4, Col: 4},
			},
		},
	},
	"bear": {
		Emoji:  "🐻",
		Frames: [4]string{"ᗰᘏᗰ", "ᗰᘏᗰノ", "ᗰ-ᗰ", "ᗰᘏᗰ"},
		Art: PetArt{
			Base: []string{
				`              `,
				`   ᗰ .--. ᗰ  `,
				`  ( ·ᴥ·    ) `,
				`   |      |  `,
				`   | \  / |  `,
				`   |_/  \_|  `,
				`              `,
			},
			IdleB: []string{
				`              `,
				`   ᗰ .--. ᗰ  `,
				`  ( ·ᴥ·    ) `,
				`    |    |    `,
				`    |\  /|    `,
				`    |_/\_|    `,
				`              `,
			},
			EyeOpen:  "·ᴥ·",
			EyeClose: "-ᴥ-",
			Anchors: PetAnchors{
				Head:      ArtAnchor{Row: 0, Col: 5},
				RightHand: ArtAnchor{Row: 3, Col: 1},
				LeftHand:  ArtAnchor{Row: 3, Col: 12},
				Body:      ArtAnchor{Row: 4, Col: 4},
			},
		},
	},
	"fish": {
		Emoji:  "🐟",
		Frames: [4]string{"ᗱᘏ>", "ᗱᘏ>。", "ᗱ->", "ᗱᘏ>"},
		Art: PetArt{
			Base: []string{
				`     /|       `,
				`  .------.   `,
				` >( >| ◉  )  `,
				`  '------'   `,
				`     \|       `,
				`          。。`,
				`              `,
			},
			IdleB: []string{
				`      /|      `,
				`  .------.   `,
				` >( >| ◉  )  `,
				`  '------'   `,
				`      \|     `,
				`         。。 `,
				`              `,
			},
			EyeOpen:  "◉",
			EyeClose: "-",
			Anchors: PetAnchors{
				Head:      ArtAnchor{Row: 0, Col: 5},
				RightHand: ArtAnchor{Row: 2, Col: 1},
				LeftHand:  ArtAnchor{Row: 2, Col: 12},
				Body:      ArtAnchor{Row: 3, Col: 4},
			},
		},
	},
}

func DetermineSpecies() string {
	hostname, _ := os.Hostname()
	u, _ := user.Current()
	username := ""
	if u != nil {
		username = u.Username
	}
	h := sha256.Sum256([]byte(hostname + ":" + username))
	idx := int(h[0]) % len(speciesList)
	return speciesList[idx]
}

var speciesNames = map[string][]string{
	"cat":  {"Mochi", "Kuro", "Tama", "Hana", "Sora", "Yuki", "Maru", "Niko"},
	"dog":  {"Pochi", "Shiro", "Riku", "Hachi", "Koro", "Taro", "Gon", "Momo"},
	"bird": {"Piyo", "Tori", "Kaze", "Hane", "Uzu", "Sui", "Rin", "Fuu"},
	"frog": {"Kero", "Gama", "Ame", "Numa", "Suzu", "Tsuyu", "Mizu", "Kawa"},
	"bear": {"Kuma", "Goro", "Mori", "Don", "Yama", "Tetsu", "Genta", "Dai"},
	"fish": {"Sakana", "Tai", "Nami", "Umi", "Shio", "Kai", "Fugu", "Koi"},
}

func GenerateName(species string) string {
	names := speciesNames[species]
	if len(names) == 0 {
		names = speciesNames["cat"]
	}
	return fmt.Sprintf("%s-%d", names[rand.Intn(len(names))], rand.Intn(100))
}
