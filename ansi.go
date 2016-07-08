package main

import "github.com/nsf/termbox-go"

type ansiMode int

const (
	modeNormal ansiMode = iota
	modeEscaped
	modeEscapeSequence
)

type setAttribute func(fg *termbox.Attribute, bg *termbox.Attribute)

func setFg(a termbox.Attribute) setAttribute {
	return func(fg *termbox.Attribute, bg *termbox.Attribute) {
		*fg = a
	}
}

func setBg(a termbox.Attribute) setAttribute {
	return func(fg *termbox.Attribute, bg *termbox.Attribute) {
		*bg = a
	}
}

func setFgBg(a termbox.Attribute) setAttribute {
	return func(fg *termbox.Attribute, bg *termbox.Attribute) {
		*fg, *bg = a, a
	}
}

func modFg(a termbox.Attribute) setAttribute {
	return func(fg *termbox.Attribute, bg *termbox.Attribute) {
		*fg |= a
	}
}

var ansiColors = map[string]setAttribute{
	"":   setFgBg(termbox.ColorDefault),
	"0":  setFgBg(termbox.ColorDefault),
	"1":  modFg(termbox.AttrBold),
	"4":  modFg(termbox.AttrUnderline),
	"7":  modFg(termbox.AttrReverse),
	"30": setFg(termbox.ColorBlack),
	"31": setFg(termbox.ColorRed),
	"32": setFg(termbox.ColorGreen),
	"33": setFg(termbox.ColorYellow),
	"34": setFg(termbox.ColorBlue),
	"35": setFg(termbox.ColorMagenta),
	"36": setFg(termbox.ColorCyan),
	"37": setFg(termbox.ColorWhite),
	"40": setBg(termbox.ColorBlack),
	"41": setBg(termbox.ColorRed),
	"42": setBg(termbox.ColorGreen),
	"43": setBg(termbox.ColorYellow),
	"44": setBg(termbox.ColorBlue),
	"45": setBg(termbox.ColorMagenta),
	"46": setBg(termbox.ColorCyan),
	"47": setBg(termbox.ColorWhite),
}
