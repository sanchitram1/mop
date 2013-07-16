// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
package mop

import (
	"fmt"
	"sort"
	"regexp"
	"strings"
	"github.com/nsf/termbox-go"
)

type LineEditor struct {
	command   rune
	prompt    string
	cursor    int
	input     string
	profile	 *Profile
}

//-----------------------------------------------------------------------------
func (self *LineEditor) Prompt(command rune, profile *Profile) {
	prompts := map[rune]string{'+': `Add tickers: `, '-': `Remove tickers: `}
	if prompt, ok := prompts[command]; ok {
		self.prompt = prompt
		self.command = command
		self.profile = profile

		DrawLine(0, 3, "<white>"+self.prompt+"</white>")
		termbox.SetCursor(len(self.prompt), 3)
		termbox.Flush()
	}
}

//-----------------------------------------------------------------------------
func (self *LineEditor) Handle(ev termbox.Event) bool {
	defer termbox.Flush()

	switch ev.Key {
	case termbox.KeyEsc:
		self.done()
		return true

	case termbox.KeyEnter:
		self.execute()
		self.done()
		return true

        case termbox.KeyBackspace, termbox.KeyBackspace2:
		self.delete_previous_character()

	case termbox.KeyCtrlB, termbox.KeyArrowLeft:
		self.move_left()

	case termbox.KeyCtrlF, termbox.KeyArrowRight:
		self.move_right()

	case termbox.KeyCtrlA:
		self.jump_to_beginning()

	case termbox.KeyCtrlE:
		self.jump_to_end()

	case termbox.KeySpace:
		self.insert_character(' ')

	default:
		if ev.Ch != 0 {
			self.insert_character(ev.Ch)
		}
	}
	DrawLine(20,20, fmt.Sprintf("cursor: %02d [%s] %08d", self.cursor, self.input, ev.Ch))
	return false
}

//-----------------------------------------------------------------------------
func (self *LineEditor) delete_previous_character() {
	if self.cursor > 0 {
		if self.cursor < len(self.input) {
			// Remove character in the middle of the input string.
			self.input = self.input[0 : self.cursor-1] + self.input[self.cursor : len(self.input)]
		} else {
			// Remove last input character.
			self.input = self.input[ : len(self.input)-1]
		}
		DrawLine(len(self.prompt), 3, self.input + " ") // Erase last character.
		self.move_left()
	}
}

//-----------------------------------------------------------------------------
func (self *LineEditor) insert_character(ch rune) {
	if self.cursor < len(self.input) {
		// Insert the character in the middle of the input string.
		self.input = self.input[0 : self.cursor] + string(ch) + self.input[self.cursor : len(self.input)]
	} else {
		// Append the character to the end of the input string.
		self.input += string(ch)
	}
	DrawLine(len(self.prompt), 3, self.input)
	self.move_right()
}

//-----------------------------------------------------------------------------
func (self *LineEditor) move_left() {
	if self.cursor > 0 {
		self.cursor -= 1
		termbox.SetCursor(len(self.prompt) + self.cursor, 3)
	}
}

//-----------------------------------------------------------------------------
func (self *LineEditor) move_right() {
	if self.cursor < len(self.input) {
		self.cursor += 1
		termbox.SetCursor(len(self.prompt) + self.cursor, 3)
	}
}

//-----------------------------------------------------------------------------
func (self *LineEditor) jump_to_beginning() {
	self.cursor = 0
	termbox.SetCursor(len(self.prompt) + self.cursor, 3)
}

//-----------------------------------------------------------------------------
func (self *LineEditor) jump_to_end() {
	self.cursor = len(self.input)
	termbox.SetCursor(len(self.prompt) + self.cursor, 3)
}

//-----------------------------------------------------------------------------
func (self *LineEditor) done() {
	ClearLine(0, 3)
	termbox.HideCursor()
}

//-----------------------------------------------------------------------------
func (self *LineEditor) execute() {
	switch self.command {
	case '+', '-':
		input := strings.ToUpper(strings.TrimSpace(self.input))
		tickers := regexp.MustCompile(`[,\s]+`).Split(input, -1)
		if len(tickers) > 0 {
			if self.command == '+' {
				self.add_tickers(tickers)
			} else {
				self.remove_tickers(tickers)
			}
		}
	}
}

//-----------------------------------------------------------------------------
func  (self *LineEditor) add_tickers(tickers []string) {
	existing := make(map[string]bool)
	for _, ticker := range self.profile.Tickers {
		existing[ticker] = true
	}

	for _, ticker := range tickers {
		if _, found := existing[ticker]; !found {
			self.profile.Tickers = append(self.profile.Tickers, ticker)
		}
	}

	sort.Strings(self.profile.Tickers)
	self.profile.Save()
	DrawQuotes(self.profile.Quotes())
}

//-----------------------------------------------------------------------------
func  (self *LineEditor) remove_tickers(tickers []string) {
}