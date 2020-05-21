package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/fatih/color"
)

// started is a flag that indicates, that the first line was processed,
// and all further lines need a separator "--".
var started = false

func main() {
	go handleSignals()

	reader := bufio.NewReader(os.Stdin)
	for {
		text, err := reader.ReadString('\n')
		if err == io.EOF {
			os.Exit(0)
		}
		if err != nil {
			fmt.Printf("Failed to read line: %v", err)
			os.Exit(1)
		}
		text = strings.TrimSpace(text)
		display(text)
	}
}

func display(s string) {
	if started {
		fmt.Println("--")
	} else {
		started = true
	}

	var m map[string]interface{}
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		fmt.Println(s)
		return
	}

	level := getLevel(m)
	keyColor, txtColor := getColors(level)
	fields := getFields(m)

	for _, field := range fields {
		keyColor.Printf("%s: ", field)
		txtColor.Println(m[field])
	}
}

func getLevel(msg map[string]interface{}) string {
	for _, k := range []string{"level", "lvl", "lev", "l", "type"} {
		if val, ok := msg[k]; ok {
			if level, ok := val.(string); ok {
				return level
			}
		}
	}
	return ""
}

func getColors(level string) (*color.Color, *color.Color) {
	switch strings.ToLower(level) {
	case "debug", "dbg", "d":
		return color.New(color.FgMagenta), color.New(color.FgWhite)
	case "info", "inf", "i":
		return color.New(color.FgBlue), color.New(color.FgWhite)
	case "warning", "warn", "wrn", "w":
		return color.New(color.FgYellow), color.New(color.FgWhite)
	case "error", "err", "e":
		return color.New(color.FgRed), color.New(color.FgWhite)
	case "fatal", "f":
		return color.New(color.FgRed), color.New(color.FgWhite)
	default:
		return color.New(color.FgWhite), color.New(color.FgWhite)
	}
}

func getFields(m map[string]interface{}) []string {
	// Determine order of fields - some fields should always be on top,
	// other - always at bottom
	firstFields := []string{"time", "level", "type"}
	lastFields := []string{"message"}
	blacklistedFields := []string{"lineno", "function", "env", "tag"}

	var fields []string

	for _, f := range firstFields {
		if _, ok := m[f]; ok {
			fields = append(fields, f)
		}
	}

	for f := range m {
		if in(f, blacklistedFields) || in(f, firstFields) || in(f, lastFields) {
			continue
		}
		fields = append(fields, f)
	}

	for _, f := range lastFields {
		if _, ok := m[f]; ok {
			fields = append(fields, f)
		}
	}

	return fields
}

func in(str string, ss []string) bool {
	for _, s := range ss {
		if str == s {
			return true
		}
	}
	return false
}

func handleSignals() {
	trap := make(chan os.Signal, 1)
	signal.Notify(trap, syscall.SIGINT, os.Interrupt, syscall.SIGTERM)
	for range trap {
		// Do nothing
	}
}
