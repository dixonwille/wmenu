package wmenu

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

var clear map[string]func()

func init() {
	clear = make(map[string]func())
	clear["linux"] = func() {
		fmt.Print("\033[2J\033[H")
	}
	clear["darwin"] = func() {
		fmt.Print("\033[2J\033[H")
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			panic(err)
		}
	}
}

// Clear simply clears the command line interface (os.Stdout only).
func Clear() {
	value, ok := clear[runtime.GOOS]
	if ok {
		value()
	}
}
