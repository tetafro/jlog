package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/hokaccha/go-prettyjson"
)

// Flags for modifying input structs.
var (
	blacklistFlag = flag.String("b", "", "Black list of fields")
	whitelistFlag = flag.String("w", "", "White list of fields")
)

// Formatters for different log levels.
var (
	debugFormat = prettyjson.Formatter{
		KeyColor:    color.New(color.FgMagenta, color.Bold),
		StringColor: color.New(color.FgWhite, color.Concealed),
		BoolColor:   color.New(color.FgWhite, color.Concealed),
		NumberColor: color.New(color.FgWhite, color.Concealed),
		NullColor:   color.New(color.FgWhite, color.Concealed),
		Indent:      4,
		Newline:     "\n",
	}
	infoFormat = prettyjson.Formatter{
		KeyColor:    color.New(color.FgBlue, color.Bold),
		StringColor: color.New(color.FgWhite, color.Concealed),
		BoolColor:   color.New(color.FgWhite, color.Concealed),
		NumberColor: color.New(color.FgWhite, color.Concealed),
		NullColor:   color.New(color.FgWhite, color.Concealed),
		Indent:      4,
		Newline:     "\n",
	}
	warnFormat = prettyjson.Formatter{
		KeyColor:    color.New(color.FgYellow, color.Bold),
		StringColor: color.New(color.FgWhite, color.Concealed),
		BoolColor:   color.New(color.FgWhite, color.Concealed),
		NumberColor: color.New(color.FgWhite, color.Concealed),
		NullColor:   color.New(color.FgWhite, color.Concealed),
		Indent:      4,
		Newline:     "\n",
	}
	errorFormat = prettyjson.Formatter{
		KeyColor:    color.New(color.FgRed, color.Bold),
		StringColor: color.New(color.FgWhite, color.Concealed),
		BoolColor:   color.New(color.FgWhite, color.Concealed),
		NumberColor: color.New(color.FgWhite, color.Concealed),
		NullColor:   color.New(color.FgWhite, color.Concealed),
		Indent:      4,
		Newline:     "\n",
	}
	defaultFormat = prettyjson.Formatter{
		DisabledColor: true,
		Indent:        4,
		Newline:       "\n",
	}
)

func main() {
	flag.Parse()

	blacklist := strings.Split(*blacklistFlag, ",")
	whitelist := strings.Split(*whitelistFlag, ",")

	if len(blacklist) == 1 && blacklist[0] == "" {
		blacklist = []string{}
	}
	if len(whitelist) == 1 && whitelist[0] == "" {
		whitelist = []string{}
	}

	if len(blacklist) > 0 && len(whitelist) > 0 {
		fmt.Printf("Use only -b or -w flag, not both [%v, %v]\n", blacklist, whitelist)
		os.Exit(1)
	}

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

		fmt.Println(format(text, blacklist, whitelist))
	}
}

func format(s string, blacklist, whitelist []string) string {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		return s
	}

	for _, bk := range blacklist {
		delete(m, bk)
	}
	if len(whitelist) > 0 {
		tmp := map[string]interface{}{}
		for _, wk := range whitelist {
			tmp[wk] = m[wk]
		}
		m = tmp
	}

	level := getLevel(m)
	fmtr := getFormatter(level)

	f, err := fmtr.Marshal(m)
	if err != nil {
		return s
	}
	return string(f)
}

func getLevel(msg map[string]interface{}) string {
	for _, k := range []string{"level", "lvl", "lev", "l"} {
		if val, ok := msg[k]; ok {
			if level, ok := val.(string); ok {
				return level
			}
		}
	}
	return ""
}

func getFormatter(level string) prettyjson.Formatter {
	switch strings.ToLower(level) {
	case "debug", "dbg", "d":
		return debugFormat
	case "info", "inf", "i":
		return infoFormat
	case "warning", "warn", "wrn", "w":
		return warnFormat
	case "error", "err", "e":
		return errorFormat
	}
	return defaultFormat
}
