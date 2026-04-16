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

var speciesList = []string{"cat", "dog", "bird", "frog"}

var speciesData = map[string]SpeciesInfo{
	"cat": {
		Emoji:  "🐱",
		Frames: [4]string{"⣾⣷", "⣷⣾", "⣾⣷", "⣶⣦"},
		Art: PetArt{
			Base: []string{
				`              `,
				`    ∧___∧     `,
				`   (・ω・)    `,
				`   /|   |\    `,
				`  / |   | \   `,
				`     UU U     `,
				`              `,
			},
			IdleB: []string{
				`              `,
				`    ∧___∧     `,
				`   (・ω・)    `,
				`    /|  |\    `,
				`   / |  | \   `,
				`     U  U     `,
				`              `,
			},
			EyeOpen:  "・ω・",
			EyeClose: "─ω─",
			Anchors: PetAnchors{
				Head:      ArtAnchor{Row: 0, Col: 5},
				RightHand: ArtAnchor{Row: 3, Col: 1},
				LeftHand:  ArtAnchor{Row: 3, Col: 12},
				Body:      ArtAnchor{Row: 4, Col: 4},
			},
		},
	},
	"dog": {
		Emoji:  "🐶",
		Frames: [4]string{"⢸⡇", "⣿⣿", "⢸⡇", "⣤⣤"},
		Art: PetArt{
			Base: []string{
				`              `,
				`   ∪・ω・∪    `,
				`    (    )    `,
				`   /|    |\   `,
				`  / |    | \  `,
				`     UU UU    `,
				`              `,
			},
			IdleB: []string{
				`              `,
				`   ∪・ω・∪    `,
				`    (    )    `,
				`    /|   |\   `,
				`   / |   | \  `,
				`     U   U    `,
				`              `,
			},
			EyeOpen:  "・ω・",
			EyeClose: "─ω─",
			Anchors: PetAnchors{
				Head:      ArtAnchor{Row: 0, Col: 5},
				RightHand: ArtAnchor{Row: 3, Col: 1},
				LeftHand:  ArtAnchor{Row: 3, Col: 12},
				Body:      ArtAnchor{Row: 2, Col: 5},
			},
		},
	},
	"bird": {
		Emoji:  "🐦",
		Frames: [4]string{"⣴⣦", "⣦⣴", "⣴⣦", "⣤⣴"},
		Art: PetArt{
			Base: []string{
				`              `,
				`     .-.      `,
				`   (・v・)    `,
				`  /)    (\    `,
				` (/      \)   `,
				`     ^^       `,
				`              `,
			},
			IdleB: []string{
				`              `,
				`     .-.      `,
				`   (・v・)    `,
				`   /)   (\    `,
				`  (/     \)   `,
				`     ^ ^      `,
				`              `,
			},
			EyeOpen:  "・v・",
			EyeClose: "─v─",
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
		Frames: [4]string{"⣤⣶", "⣶⣤", "⣤⣶", "⣴⣤"},
		Art: PetArt{
			Base: []string{
				`              `,
				`   .------.   `,
				`  ( ・__・ )  `,
				`  /|      |\  `,
				`  ||  ()  ||  `,
				`   '------'   `,
				`              `,
			},
			IdleB: []string{
				`              `,
				`   .------.   `,
				`  ( ・__・ )  `,
				`   |      |   `,
				`   |  ()  |   `,
				`   '------'   `,
				`              `,
			},
			EyeOpen:  "・__・",
			EyeClose: "─__─",
			Anchors: PetAnchors{
				Head:      ArtAnchor{Row: 0, Col: 5},
				RightHand: ArtAnchor{Row: 3, Col: 1},
				LeftHand:  ArtAnchor{Row: 3, Col: 12},
				Body:      ArtAnchor{Row: 4, Col: 4},
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

var catNames = []string{"Mochi", "Kuro", "Tama", "Hana", "Sora", "Yuki", "Maru", "Niko"}
var dogNames = []string{"Pochi", "Shiro", "Riku", "Hachi", "Koro", "Taro", "Gon", "Momo"}
var birdNames = []string{"Piyo", "Tori", "Kaze", "Hane", "Uzu", "Sui", "Rin", "Fuu"}
var frogNames = []string{"Kero", "Gama", "Ame", "Numa", "Suzu", "Tsuyu", "Mizu", "Kawa"}

func GenerateName(species string) string {
	var names []string
	switch species {
	case "cat":
		names = catNames
	case "dog":
		names = dogNames
	case "bird":
		names = birdNames
	case "frog":
		names = frogNames
	default:
		names = catNames
	}
	return fmt.Sprintf("%s-%d", names[rand.Intn(len(names))], rand.Intn(100))
}
