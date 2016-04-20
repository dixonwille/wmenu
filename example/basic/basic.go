package main

import (
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/dixonwille/wmenu"
)

func main() {
	menu := wmenu.NewMenu("Choose an option")
	menu.SetSeperator(",")
	menu.LoopOnInvalid()
	menu.Action(func(opt []wmenu.Option) { spew.Dump(opt) })
	menu.MultipleAction(func(opt []wmenu.Option) { fmt.Println("multiple"); spew.Dump(opt) })
	menu.Option("This is option 1.", true, func() { fmt.Println("Choose option 1") })
	menu.Option("This is option 2.", false, func() { fmt.Println("Choose option 2") })
	menu.Option("This is option 3.", false, nil)
	err := menu.Run()
	if err != nil {
		log.Fatal(err)
	}
}
