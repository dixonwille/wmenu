package main

import (
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/dixonwille/wlog"
	"github.com/dixonwille/wmenu"
)

func main() {
	menu := wmenu.NewMenu("Choose an option")
	menu.SetSeperator(",")
	menu.LoopOnInvalid()
	menu.Action(def)
	menu.MultipleAction(multiDef)
	menu.Option("This is option 1.", true, option1)
	menu.Option("This is option 2.", false, option2)
	menu.Option("This is option 3.", false, nil)
	err := menu.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func option1() {
	menu := wmenu.NewMenu("This is option 1's menu")
	menu.LoopOnInvalid()
	menu.Action(func(opt wmenu.Option) { spew.Dump(opt) })
	menu.Option("1.1", false, nil)
	menu.Option("1.2", false, nil)
	menu.Option("1.3", false, nil)
	err := menu.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func option2() {
	menu := wmenu.NewMenu("This is option 2's menu")
	menu.LoopOnInvalid()
	menu.AddColor(wlog.BrightGreen, wlog.BrightMagenta, wlog.BrightCyan, wlog.BrightRed)
	menu.Action(func(opt wmenu.Option) { spew.Dump(opt) })
	menu.Option("2.1", false, nil)
	menu.Option("2.2", true, nil)
	menu.Option("2.3", false, nil)
	err := menu.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func def(opt wmenu.Option) {
	fmt.Println("This is the default action")
	fmt.Printf("Do something with %v\n", opt)
	fmt.Println("You can create another menu if you wanted to.")
}

func multiDef(opt []wmenu.Option) {
	fmt.Println("This is the multiple default action")
	for _, o := range opt {
		fmt.Printf("Do something with %v\n", o)
	}
	fmt.Println("You can create another menu if you wanted to.")
}
