//package main is to show how multiple menus can be used with eachother
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/dixonwille/wmenu"
)

type menuItem int

const (
	pizza menuItem = iota
	iceCream
	tacos
)

var menuItemStrings = map[menuItem]string{
	pizza:    "Pizza",
	iceCream: "Ice Cream",
	tacos:    "Tacos",
}

func main() {
	mm := mainMenu()
	err := mm.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func mainMenu() *wmenu.Menu {
	menu := wmenu.NewMenu("What is your favorite food?")
	menu.Option(menuItemStrings[pizza], pizza, true, nil)
	menu.Option(menuItemStrings[iceCream], iceCream, false, nil)
	menu.Option(menuItemStrings[tacos], tacos, false, func(opt wmenu.Opt) error {
		fmt.Printf("Tacos are great!\n")
		return nil
	})
	menu.Action(func(opts []wmenu.Opt) error {
		if len(opts) != 1 {
			return errors.New("wrong number of options chosen")
		}

		tm := toppingsMenu(opts[0].Value.(menuItem))
		return tm.Run()
	})
	return menu
}

func toppingsMenu(favorite menuItem) *wmenu.Menu {
	menu := wmenu.NewMenu(fmt.Sprintf("What is your favorite topping for %s?", menuItemStrings[favorite]))
	if favorite == pizza {
		menu.Option("Meat", nil, true, nil)
		menu.Option("Cheese", nil, false, nil)
		menu.Option("Vegitables", nil, false, nil)
	}
	if favorite == iceCream {
		menu.Option("Fruit", nil, true, nil)
		menu.Option("Chocolate Syrup", nil, false, nil)
		menu.Option("Caramel Syrup", nil, false, nil)
	}
	menu.Action(func(opts []wmenu.Opt) error {
		if len(opts) != 1 {
			return errors.New("wrong number of options chosen")
		}

		fmt.Printf("Your favorite food is %s with %s on top.\n", menuItemStrings[favorite], opts[0].Text)

		return nil
	})
	return menu
}
