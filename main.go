package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/nsf/termbox-go"
)

// current holds the current key press
var current string

// curev holds the current event
var curev termbox.Event

// coldef is the default color
const coldef = termbox.ColorDefault

// tbprint prints to the termbox
func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		if c == '\n' {
			x = 0
			y++
		} else {
			termbox.SetCell(x, y, c, fg|termbox.AttrBold, bg)
			x++
		}
	}
}

// tberror shows an error in termbox
func tberror(title, msg string) {
	tbbox(termbox.ColorWhite, termbox.ColorRed, title, msg)
}

// tberror shows info in a termbox
func tbinfo(title, msg string) {
	tbbox(termbox.ColorWhite, termbox.ColorBlue, title, msg)
}

func tbbox(fg, bg termbox.Attribute, title, msg string) {
	x, y := centerCoordinates()
	msg = fmt.Sprintf(" %s ", msg)
	x = x - (len(msg) / 2)
	format := fmt.Sprintf("%%-%ds", len(msg))
	errorHeader := fmt.Sprintf(format, title)
	tbprint(x, y-2, fg|termbox.AttrBold, bg, errorHeader)
	tbprint(x, y-1, fg, bg, strings.Repeat(" ", len(msg)))
	tbprint(x, y, fg, bg, msg)
	tbprint(x, y+1, fg, bg, strings.Repeat(" ", len(msg)))
	tbprint(x, y+2, fg, bg, strings.Repeat(" ", len(msg)))
}

// centerCoordinates finds the middle point of the terminal
func centerCoordinates() (int, int) {
	width, height := termbox.Size()
	x := width / 2
	y := height / 2
	return x, y
}

// showCommandOutput attempts to show the output from a given command
func showCommandOutput(entryPoint string, args ...string) {
	cmd := exec.Command(entryPoint, args...)
	if output, err := cmd.Output(); err != nil {
		tberror("Error", err.Error())
	} else {
		termbox.Clear(coldef, coldef)
		tbprint(0, 0, coldef, coldef, string(output))
	}
	termbox.Flush()
	time.Sleep(3 * time.Second)
}

// redraw redraws and refreshes the main menu
func redraw() {
	termbox.Clear(coldef, coldef)

	x, y := centerCoordinates()
	tbprint(x, y, coldef, coldef, "0. Exit")
	tbprint(x, y+1, coldef, coldef, "1. Shutdown")
	tbprint(x, y+2, coldef, coldef, "2. Reboot")
	tbprint(x, y+3, coldef, coldef, "3. Openshift Status")
	tbprint(x, y+4, coldef, coldef, "4. OSTree Status")
	tbprint(x, y+5, coldef, coldef, "5. Shell")

	termbox.Flush()
}

// replaceProcess replaces this execs process with another one
func replaceProcess(entryPoint string, argv, env []string) {
	termbox.Clear(coldef, coldef)
	termbox.Close()
	syscall.Exec(entryPoint, argv, env)
	// Should never get this far ... but just in case...
	os.Exit(0)
}

// mainloop provides the main menu loop
func mainloop() {
	termbox.SetInputMode(termbox.InputCurrent)
	// initial draw
	redraw()

mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			current := ev.Ch
			if current == '0' {
				break mainloop
			} else if current == '1' {
				exec.Command("/usr/bin/bash", "-c", "systemctl halt -i").Start()
			} else if current == '2' {
				exec.Command("/usr/bin/bash", "-c", "systemctl reboot -i").Start()
			} else if current == '3' {
				showCommandOutput("oc status")
			} else if current == '4' {
				showCommandOutput("/usr/bin/rpm-ostree", "status", "-v")
			} else if current == '5' {
				tbinfo("Note", "Modifying the base system is not recommended")
				termbox.Flush()
				time.Sleep(5 * time.Second)
				replaceProcess("/bin/bash", []string{"bash"}, os.Environ())
			}
		}
		redraw()
	}
}

// main is the main entry point for the binary
func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	mainloop()
}
