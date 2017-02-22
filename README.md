# WMenu [![Build Status](https://travis-ci.org/dixonwille/wmenu.svg?branch=v1)](https://travis-ci.org/dixonwille/wmenu) [![Go Report Card](https://goreportcard.com/badge/github.com/dixonwille/wmenu)](https://goreportcard.com/report/github.com/dixonwille/wmenu) [![GoDoc](https://godoc.org/github.com/dixonwille/wmenu?status.svg)](https://godoc.org/github.com/dixonwille/wmenu)

Package wmenu creates menus for cli programs. It uses wlog for it's interface
with the command line. It uses os.Stdin, os.Stdout, and os.Stderr with
concurrency by default. wmenu allows you to change the color of the different
parts of the menu. This package also creates it's own error structure so you can
type assert if you need to. wmenu will validate all responses before calling any function. It will also figure out which function should be called so you don't have to.

## Import
    import "github.com/dixonwille/wmenu"

## Features
* Force single selection
* Allow multiple selection
* Change the delimiter
* Change the color of different parts of the menu
* Easily see which option(s) are default
* Change the symbol used for default option(s)
* Ask simple yes and no questions
* Validate all responses before calling any functions
* With yes and no can accept:
  * yes, Yes, YES, y, Y
  * no, No, NO, n, N
* Figure out which Action should be called (Options, Default, or Multiple Action)
* Re-ask question if invalid response up to a certain number of times
* Can change max number of times to ask before failing output
* Change reader and writer
* Clear the screen whenever the menu is brought up
* Has its own error structure so you can type assert menu errors

## Usage
This is a simple use of the package.
``` go
menu := wmenu.NewMenu("What is your favorite food?")
menu.Action(func (opt Opt) error {fmt.Printf(opt.Text + " is your favorite food."); return nil})
menu.Option("Pizza", true, nil)
menu.Option("Ice Cream", false, nil)
menu.Option("Tacos", false, func() error {
  fmt.Printf("Tacos are great")
})
err := menu.Run()
if err != nil{
  log.Fatal(err)
}
```
The output would look like this:
```
0) *Pizza
1) Ice Cream
2) Tacos
What is your favorite food?
```
If the user just presses `[Enter]` then the option(s) with the `*` will be selected. This indicates that it is a default function. If they choose `1` then they would see `Ice Cream is your favorite food.`. This used the Action's function because the option selected didn't have a function along with it. But if they choose `2` they would see `Tacos are great`. That option did have a function with it which take precedence over Action.

You can you also use:
``` go
menu.MultipleAction(func (opt []Opt) error {return nil})
```
This will allow the user to select multiple options. The default delimiter is a `[space]`, but can be changed by using:
``` go
menu.SetSeperator("some string")
```

Another feature is the ability to ask yes or no questions.
``` go
menu.IsYesNo(0)
```
This will remove any options previously added options and hide the ones used for the menu. It will simply just ask yes or no. Menu will parse and validate the response for you. This option will always call the Action's function and pass in the option that was selected.

## Further Reading
This whole package has been documented and has a few examples in the [godocs](https://godoc.org/github.com/dixonwille/wmenu). You should read the docs to find all functions and structures at your finger tips.
