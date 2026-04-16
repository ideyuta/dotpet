package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {
	pet, err := LoadPet()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	cmd := ""
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	switch cmd {
	case "":
		pet, _ = Tick(pet)
		fmt.Print(StatusLine(pet))
	case "status":
		pet, _ = Tick(pet)
		fmt.Print(StatusDetail(pet))
	case "watch":
		pet, _ = Tick(pet)
		if err := SavePet(pet); err != nil {
			fmt.Fprintf(os.Stderr, "error saving: %v\n", err)
		}
		runStatusLoop(pet)
		return
	case "log":
		fmt.Print(EventLog(pet))
	default:
		fmt.Fprintf(os.Stderr, "usage: dotpet [status|log|watch]\n")
		os.Exit(1)
	}

	if err := SavePet(pet); err != nil {
		fmt.Fprintf(os.Stderr, "error saving: %v\n", err)
		os.Exit(1)
	}
}

// runStatusLoop continuously redraws the status screen with animation.
// Exits on Ctrl+C.
func runStatusLoop(pet *Pet) {
	// Hide cursor
	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h") // Show cursor on exit

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()

	// Initial draw
	fmt.Print("\033[2J\033[H")
	fmt.Print(StatusDetail(pet))

	for {
		select {
		case <-sig:
			fmt.Print("\033[2J\033[H")
			return
		case <-ticker.C:
			fmt.Print("\033[H")
			fmt.Print(StatusDetail(pet))
		}
	}
}
