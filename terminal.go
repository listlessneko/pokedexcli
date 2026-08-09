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
	keyh         = '\x68'
	keyj         = '\x6A'
	keyk         = '\x6B'
	keyl         = '\x6C'
	enterSeq     = "\x0d\x0a"
	eraseSeq     = "\x1b[K"
	cursorHide   = "\x1b[?25l"
	cursorShow   = "\x1b[?25h"
	cursorUp     = "\x1b[A"
	cursorDown   = "\x1b[B"
	cursorFwd    = "\x1b[C"
	cursorBckwd  = "\x1b[D"
)

func redrawCurrentLine(currentLine []byte, cursor int) {
	if cursor > 0 {
		seq := fmt.Sprintf("\x1b[%dD", cursor)
		os.Stdout.Write([]byte(seq))
	}
	os.Stdout.Write([]byte(eraseSeq))
	os.Stdout.Write(currentLine)
}

func addCharBeforeCursor(currentLine []byte, cursor *int, b byte) []byte {
	currentLine = append(currentLine, 0)
	copy(currentLine[*cursor+1:], currentLine[*cursor:])
	currentLine[*cursor] = b
	*cursor += 1
	os.Stdout.Write(currentLine[*cursor-1:])
	nCol := len(currentLine) - *cursor
	if nCol > 0 {
		seq := fmt.Sprintf("\x1b[%dD", nCol)
		os.Stdout.Write([]byte(seq))
	}
	return currentLine
}

func deleteCharBeforeCursor(currentLine []byte, cursor *int) []byte {
	if *cursor > 0 {
		copy(currentLine[*cursor-1:], currentLine[*cursor:])
		currentLine = currentLine[:len(currentLine)-1]
		*cursor -= 1
		os.Stdout.Write([]byte(cursorBckwd))
		os.Stdout.Write([]byte(currentLine[*cursor:]))
		os.Stdout.Write([]byte(" "))
		nCol := len(currentLine) - *cursor + 1
		seq := fmt.Sprintf("\x1b[%dD", nCol)
		os.Stdout.Write([]byte(seq))
	}
	return currentLine
}

func moveUpHistory(currentLine []byte, cursor *int, history []string, historyIndex *int) []byte {
	if *historyIndex > 0 {
		*historyIndex -= 1
		currentLine = []byte(history[*historyIndex])
		redrawCurrentLine(currentLine, *cursor)
		*cursor = len(currentLine)
	}
	return currentLine
}

func moveDownHistory(currentLine []byte, cursor *int, history []string, historyIndex *int) []byte {
	if *historyIndex < len(history) {
		*historyIndex += 1
		if *historyIndex == len(history) {
			currentLine = []byte("")
		} else {
			currentLine = []byte(history[*historyIndex])
		}
		redrawCurrentLine(currentLine, *cursor)
		*cursor = len(currentLine)
	}
	return currentLine
}

func moveCursorFwd(cursor *int, currentLine []byte) {
	if *cursor < len(currentLine) {
		*cursor += 1
		os.Stdout.Write([]byte(cursorFwd))
	}
}

func moveCursorBckwd(cursor *int, currentLine []byte) {
	if *cursor > 0 {
		*cursor -= 1
		os.Stdout.Write([]byte(cursorBckwd))
	}
}

func drawCommandsWithPrefix(currentLine []byte, prompt string, commands []string) {
	os.Stdout.Write([]byte(enterSeq))
	for i, command := range commands {
		os.Stdout.Write([]byte(command))
		if i < len(commands)-1 {
			os.Stdout.Write([]byte("  "))
		}
	}
	os.Stdout.Write([]byte(enterSeq))
	os.Stdout.Write([]byte(prompt))
	os.Stdout.Write((currentLine))
}

func autocompleteCommand(currentLine []byte, cursor *int, commands []string) []byte {
	currentLine = []byte(commands[0])
	redrawCurrentLine(currentLine, *cursor)
	*cursor = len(currentLine)
	return currentLine
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
				currentLine = deleteCharBeforeCursor(currentLine, &cursor)
				continue
			case keyEscape:
				if i+2 < n && buf[i+1] == keyLSqBrckt {
					switch buf[i+2] {
					case keyA:
						currentLine = moveUpHistory(currentLine, &cursor, history, &historyIndex)
					case keyB:
						currentLine = moveDownHistory(currentLine, &cursor, history, &historyIndex)
					case keyC:
						moveCursorFwd(&cursor, currentLine)
					case keyD:
						moveCursorBckwd(&cursor, currentLine)
					}
					i += 2
				}
				continue
			case keyTab:
				commands := autocomplete.searchByPrefix(string(currentLine))
				if len(commands) > 1 {
					drawCommandsWithPrefix(currentLine, prompt, commands)
				} else if len(commands) == 1 {
					currentLine = autocompleteCommand(currentLine, &cursor, commands)
				}
				continue
			}
			currentLine = addCharBeforeCursor(currentLine, &cursor, b)
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

func moveDownPromptSelector(prompts []string, selected *int) {
	if *selected < len(prompts)-1 {
		*selected++
	} else if *selected == len(prompts)-1 {
		*selected = 0
	}
	moveCursorUp := fmt.Sprintf("\r\x1b[%dA", len(prompts)-1)
	os.Stdout.Write([]byte(moveCursorUp))
	drawList(prompts, *selected)
}

func moveUpPromptSelector(prompts []string, selected *int) {
	if *selected > 0 {
		*selected--
	} else if *selected == 0 {
		*selected = len(prompts) - 1
	}
	moveCursorUp := fmt.Sprintf("\r\x1b[%dA", len(prompts)-1)
	os.Stdout.Write([]byte(moveCursorUp))
	drawList(prompts, *selected)

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
			case keyj:
				moveDownPromptSelector(prompts, &selected)
			case keyk:
				moveUpPromptSelector(prompts, &selected)
			case keyEscape:
				if i+2 < n && buf[i+1] == keyLSqBrckt {
					switch buf[i+2] {
					case keyB:
						moveDownPromptSelector(prompts, &selected)
					case keyA:
						moveUpPromptSelector(prompts, &selected)
					}
					i += 2
				}
				continue
			}
		}
	}
}
