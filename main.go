package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"math"
	"os"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

var (
	tabWidth = flag.Int("tab-width", 8, "Tab width")
)

func main() {
	flag.Parse()

	var err error
	var in *os.File
	if !terminal.IsTerminal(0) {
		in = os.Stdin
	} else {
		if flag.NArg() != 1 {
			log.Fatal("Invalid arguments")
		}
		if filename := flag.Args()[0]; filename == "-" {
			in = os.Stdin
		} else {
			if in, err = os.Open(filename); err != nil {
				log.Fatal("Failed to open file: ", err)
			}
			defer in.Close()
		}
	}

	var pager Pager
	if err = pager.Init(); err != nil {
		log.Fatal("Failed to initialize terminal: ", err)
	}
	defer pager.Close()

	b, err := ioutil.ReadAll(in)
	if err != nil {
		log.Fatal("Failed to read content: ", err)
	}

	pager.AddContent(b)
	pager.Redraw()
	pager.PollEvent()
}

type Pager struct {
	lines          []string
	incompleteLine bool
	width          int
	height         int
	viewX          int
}

func (p *Pager) Init() error {
	if err := termbox.Init(); err != nil {
		return err
	}
	return nil
}

func (p *Pager) Close() {
	termbox.Close()
}

func (p *Pager) PollEvent() {
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if !p.handleKeyEvent(ev) {
				return
			}
		case termbox.EventResize:
			p.Redraw()
		}
	}
}

func (p *Pager) handleKeyEvent(ev termbox.Event) bool {
	p.width, p.height = termbox.Size()
	switch ev.Key {
	case termbox.KeyCtrlN, termbox.KeyArrowDown, termbox.KeyEnter:
		p.scrollDown(1)
	case termbox.KeyCtrlP, termbox.KeyArrowUp:
		p.scrollUp(1)
	case termbox.KeyCtrlU:
		p.scrollUp(int(math.Ceil(float64(p.height) / 2)))
	case termbox.KeyCtrlD:
		p.scrollDown(int(math.Ceil(float64(p.height) / 2)))
	case termbox.KeyCtrlF, termbox.KeySpace, termbox.KeyPgdn:
		p.scrollDown(max(0, p.height-3))
	case termbox.KeyCtrlB, termbox.KeyPgup:
		p.scrollUp(max(0, p.height-3))
	default:
		switch ev.Ch {
		case 'j':
			p.scrollDown(1)
		case 'k':
			p.scrollUp(1)
		case 'g':
			p.viewX = 0
		case 'G':
			p.viewX = max(0, len(p.lines)-p.height)
		case 'q':
			return false
		}
	}
	p.Redraw()
	return true
}

func (p *Pager) scrollUp(n int) {
	p.viewX = max(0, p.viewX-n)
}

func (p *Pager) scrollDown(n int) {
	p.viewX = min(p.viewX+n, max(0, len(p.lines)-p.height))
}

func (p *Pager) Redraw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	var x, y int
	var buf []rune
	fg, bg := termbox.ColorDefault, termbox.ColorDefault
	mode := modeNormal
	for _, line := range p.lines[p.viewX:] {
		for _, r := range line {
			switch mode {
			case modeNormal:
				switch r {
				case '\t':
					x += *tabWidth
					continue
				case '\033':
					mode = modeEscaped
					continue
				}
			case modeEscaped:
				switch r {
				case '[':
					mode = modeEscapeSequence
					continue
				default:
					mode = modeNormal
				}
			case modeEscapeSequence:
				switch r {
				case ';':
					if fn, ok := ansiColors[string(buf)]; ok {
						fn(&fg, &bg)
					}
					buf = nil
				case 'm':
					if fn, ok := ansiColors[string(buf)]; ok {
						fn(&fg, &bg)
					}
					buf = nil
					mode = modeNormal
				case 'K':
					mode = modeNormal
				default:
					buf = append(buf, r)
				}
				continue
			}

			termbox.SetCell(x, y, r, fg, bg)
			w := runewidth.RuneWidth(r)
			if w == 0 || (w == 2 && runewidth.IsAmbiguousWidth(r)) {
				w = 1
			}
			x += w
		}
		x = 0
		y++
	}
	termbox.Flush()
}

func (p *Pager) AddContent(b []byte) {
	for len(b) > 0 {
		n := bytes.IndexByte(b, '\n')
		if n == -1 {
			p.appendLine(b)
			p.incompleteLine = true
		} else {
			p.appendLine(b[:n])
			b = b[n+1:]
			p.incompleteLine = false
		}
	}
}

func (p *Pager) appendLine(b []byte) {
	if p.incompleteLine {
		p.lines[len(p.lines)-1] += string(b)
	} else {
		p.lines = append(p.lines, string(b))
	}
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
