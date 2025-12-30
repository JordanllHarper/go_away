package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	sendTerminal *bool
	startTui     *bool
	sendMarkdown *string
	sendTxt      *string

	// TODO: parameterize these in a file
	// cfg    = []string{"med tablets", "pants", "socks", "shirts"}
	// threeDayItems = []string{"jumpers", "jeans"}
	// singleItems   = []string{"deodorant", "toothbrush", "laptop", "water bottle", "headphones", "charger"}
)

type config struct {
	DailyItems    []string `json:"dailyItems"`
	ThreeDayItems []string `json:"threeDayItems"`
	SingleItems   []string `json:"singleItems"`
}

func readConfig(filename string) (config, error) {
	contents, err := os.Open(filename)
	if err != nil {
		return config{}, err
	}
	var cfg config
	if err := json.NewDecoder(contents).Decode(&cfg); err != nil {
		return config{}, err
	}
	return cfg, nil
}

type (
	FormatOutputter interface {
		FormatSection(header string, count int, items []string)
		Output() error
	}
)

func main() {
	sendTerminal = flag.Bool("term", true, "Send to the terminal")
	startTui = flag.Bool("tui", false, "Start an interactive tui")
	sendMarkdown = flag.String("markdown", "", "Send to a markdown file")
	sendTxt = flag.String("text", "", "Send to a txt file")
	flag.Parse()

	cfg, err := readConfig("config.json")

	fmt.Print("Days away for: ")
	var input string
	nValues, err := fmt.Scanln(&input)
	mustNotBeErr("Reading input", err)
	mustPass("Input must not be empty", nValues > 0)

	days, err := strconv.Atoi(input)
	mustNotBeErr("Days is a number", err)
	mustPass("Days > 0", days > 0)

	writeToAllOutputs(days, cfg)
}

// writeToAllOutputs function    writes to all the outputs specified in command line
func writeToAllOutputs(days int, cfg config) {
	write := func(fo FormatOutputter) {
		fo.FormatSection("Daily items", days, cfg.DailyItems)
		fo.FormatSection("Three day items", max(days/3, 1), cfg.ThreeDayItems)
		fo.FormatSection("Regular items", 0, cfg.SingleItems)
		mustNotBeErr("Writing items to output", fo.Output())
	}

	if *startTui {
		write(newTuiFormatter())
	}
	if *sendTerminal {
		write(newStdoutFormatter())
	}
	if strings.TrimSpace(*sendMarkdown) != "" {
		write(newMarkdownFormatter(*sendMarkdown))
	}
	if strings.TrimSpace(*sendTxt) != "" {
		write(newTextFormatter(*sendTxt))
	}

}
