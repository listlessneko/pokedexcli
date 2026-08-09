package main

import (
	"fmt"
	"golang.org/x/term"
	"io"
	"os"
)

const (
	bufSize      = 3
	newLine      = '\x0a'
	keyTab       = '\x09'
	keyEnter     = '\x0d'
	keyBackspace = '\x7f'
	keyCtrlC     = '\x03'
	keyCtrlD     = '\x04'
	keyEscape    = '\x1b'
	keyLSqBrckt  = '\x5b'
	keyA         = '\x41'
	keyB         = '\x42'
	keyC         = '\x43'
	keyD         = '\x44'
	enterSeq     = "\x0d\x0a"
	eraseSeq     = "\x1b[K"
	cursorHide   = "\x1b[?25l"
	cursorShow   = "\x1b[?25h"
	cursorUp     = "\x1b[A"
	cursorDown   = "\x1b[B"
	cursorFwd    = "\x1b[C"
	cursorBckwd  = "\x1b[D"
)

func redraw(new []byte, cursor int) {
	if cursor > 0 {
		seq := fmt.Sprintf("\x1b[%dD", cursor)
		os.Stdout.Write([]byte(seq))
	}
	os.Stdout.Write([]byte(eraseSeq))
	os.Stdout.Write(new)
}

func readLine(prompt string, history []string, autocomplete *trie) (string, error) {
	os.Stdout.Write([]byte(prompt))

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return "", fmt.Errorf("error enabling raw mode: %w", err)
	}

	defer term.Restore(fd, oldState)

	historyIndex := len(history)
	var cursor int
	var currentLine []byte
	buf := make([]byte, bufSize)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			return "", fmt.Errorf("error reading bytes: %w", err)
		}
		for i := 0; i < n; i++ {
			b := buf[i]
			switch b {
			case keyCtrlC, keyCtrlD:
				return "", io.EOF
			case keyEnter:
				os.Stdout.Write([]byte(enterSeq))
				return string(currentLine), nil
			case keyBackspace:
				if cursor > 0 {
					copy(currentLine[cursor-1:], currentLine[cursor:])
					currentLine = currentLine[:len(currentLine)-1]
					cursor -= 1
					os.Stdout.Write([]byte(cursorBckwd))
					os.Stdout.Write([]byte(currentLine[cursor:]))
					os.Stdout.Write([]byte(" "))
					nCol := len(currentLine) - cursor + 1
					seq := fmt.Sprintf("\x1b[%dD", nCol)
					os.Stdout.Write([]byte(seq))
				}
				continue
			case keyEscape:
				if i+2 < n && buf[i+1] == keyLSqBrckt {
					switch buf[i+2] {
					case keyA:
						if historyIndex > 0 {
							historyIndex -= 1
							currentLine = []byte(history[historyIndex])
							redraw(currentLine, cursor)
							cursor = len(currentLine)
						}
					case keyB:
						if historyIndex < len(history) {
							historyIndex += 1
							if historyIndex == len(history) {
								currentLine = []byte("")
							} else {
								currentLine = []byte(history[historyIndex])
							}
							redraw(currentLine, cursor)
							cursor = len(currentLine)
						}
					case keyC:
						if cursor < len(currentLine) {
							cursor += 1
							os.Stdout.Write([]byte(cursorFwd))
						}
					case keyD:
						if cursor > 0 {
							cursor -= 1
							os.Stdout.Write([]byte(cursorBckwd))
						}
					}
					i += 2
				}
				continue
			}
			currentLine = append(currentLine, 0)
			copy(currentLine[cursor+1:], currentLine[cursor:])
			currentLine[cursor] = b
			cursor += 1
			os.Stdout.Write(currentLine[cursor-1:])
			nCol := len(currentLine) - cursor
			if nCol > 0 {
				seq := fmt.Sprintf("\x1b[%dD", nCol)
				os.Stdout.Write([]byte(seq))
			}
		}
	}
}

func drawList(keys []string, selected int) {
	for i, k := range keys {
		prefix := " "
		if i == selected {
			prefix = "> "
		}
		line := eraseSeq + prefix + k
		if i < len(keys)-1 {
			line += "\r\n"
		}
		os.Stdout.Write([]byte(line))
	}
}

func selectPrompt(prompts []string) (string, error) {
	selected := 0
	drawList(prompts, selected)

	os.Stdout.Write([]byte(cursorHide))
	defer os.Stdout.Write([]byte(cursorShow))

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return "", fmt.Errorf("error enabling raw mode: %w", err)
	}

	defer term.Restore(fd, oldState)

	buf := make([]byte, bufSize)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			return "", fmt.Errorf("error reading bytes: %w", err)
		}
		for i := 0; i < n; i++ {
			b := buf[i]
			switch b {
			case keyCtrlC, keyCtrlD:
				return "", io.EOF
			case keyEnter:
				os.Stdout.Write([]byte(enterSeq))
				return prompts[selected], nil
			case keyEscape:
				if i+2 < n && buf[i+1] == keyLSqBrckt {
					switch buf[i+2] {
					case keyA:
						if selected > 0 {
							selected--
						} else if selected == 0 {
							selected = len(prompts) - 1
						}
					case keyB:
						if selected < len(prompts)-1 {
							selected++
						} else if selected == len(prompts)-1 {
							selected = 0
						}
					}
					if buf[i+2] == keyA || buf[i+2] == keyB {
						moveCursorUp := fmt.Sprintf("\r\x1b[%dA", len(prompts)-1)
						os.Stdout.Write([]byte(moveCursorUp))
						drawList(prompts, selected)
					}
					i += 2
				}
				continue
			}
		}
	}
}
