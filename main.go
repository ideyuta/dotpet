package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"golang.org/x/sys/unix"
	"golang.org/x/term"
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
		runWatchLoop(pet)
		return
	case "log":
		fmt.Print(EventLog(pet))
	case "reset":
		p := NewPet(1, 0)
		if err := SavePet(p); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Reset! New pet: %s (%s)\n", p.Name, p.Species)
		return
	default:
		fmt.Fprintf(os.Stderr, "usage: dotpet [status|log|watch|reset]\n")
		os.Exit(1)
	}

	if err := SavePet(pet); err != nil {
		fmt.Fprintf(os.Stderr, "error saving: %v\n", err)
		os.Exit(1)
	}
}

func runWatchLoop(pet *Pet) {
	input, needClose := getInput()
	if input == nil {
		runSimpleLoop(pet)
		return
	}
	if needClose {
		defer input.Close()
	}

	fd := int(input.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		runSimpleLoop(pet)
		return
	}
	defer term.Restore(fd, oldState)

	// Re-enable output post-processing so \n produces \r\n.
	// MakeRaw disables OPOST which breaks line rendering.
	if t, err := unix.IoctlGetTermios(fd, unix.TIOCGETA); err == nil {
		t.Oflag |= unix.OPOST
		unix.IoctlSetTermios(fd, unix.TIOCSETA, t)
	}

	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h\033[2J\033[H\r\n")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	showInventory := false
	scroll := 0
	cursor := 0

	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()

	keys := make(chan byte, 16)
	go func() {
		buf := make([]byte, 1)
		for {
			n, err := input.Read(buf)
			if err != nil || n == 0 {
				return
			}
			keys <- buf[0]
		}
	}()

	lastOutput := ""
	draw := func() {
		var out string
		if showInventory {
			out = InventoryView(pet, scroll, cursor)
		} else {
			out = StatusDetail(pet)
		}
		if out == lastOutput {
			return
		}
		lastOutput = out
		fmt.Print("\033[H")
		fmt.Print(out)
		fmt.Print("\033[J")
	}
	draw()

	for {
		select {
		case <-sig:
			return
		case key := <-keys:
			switch key {
			case '2', 'i':
				showInventory = true
				scroll = 0
				cursor = 0
			case '1', 27: // 1 or Esc
				showInventory = false
			case 'j':
				if showInventory {
					total := len(pet.Inventory)
					if cursor < total-1 {
						cursor++
					}
					// auto-scroll: keep cursor visible
					const pageSize = 15
					if cursor >= scroll+pageSize {
						scroll = cursor - pageSize + 1
					}
				}
			case 'k':
				if showInventory {
					if cursor > 0 {
						cursor--
					}
					if cursor < scroll {
						scroll = cursor
					}
				}
			case 'e', 13: // e or Enter
				if showInventory && len(pet.Inventory) > 0 {
					sorted := SortedInventory(pet)
					if cursor >= 0 && cursor < len(sorted) {
						eq := sorted[cursor]
						pet.Equipped = &eq
						if err := SavePet(pet); err != nil {
							fmt.Fprintf(os.Stderr, "error saving: %v\n", err)
						}
					}
				}
			case 'q', 3:
				return
			default:
				continue
			}
			draw()
		case <-ticker.C:
			draw()
		}
	}
}

func getInput() (*os.File, bool) {
	if term.IsTerminal(int(os.Stdin.Fd())) {
		return os.Stdin, false
	}
	f, err := os.Open("/dev/tty")
	if err != nil {
		return nil, false
	}
	return f, true
}

func runSimpleLoop(pet *Pet) {
	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()

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
