package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var (
	sendTerminal *bool
	startTui     *bool
	sendMarkdown *string
	sendTxt      *string
)

type config struct {
	DailyItems    []string `json:"dailyItems"`
	ThreeDayItems []string `json:"threeDayItems"`
	SingleItems   []string `json:"singleItems"`
}

func NewDefaultConfig() config {
	cfg := []string{"med tablets", "pants", "socks", "shirts"}
	threeDayItems := []string{"jumpers", "jeans"}
	singleItems := []string{"deodorant", "toothbrush", "laptop", "water bottle", "headphones", "charger"}
	return config{
		DailyItems:    cfg,
		ThreeDayItems: threeDayItems,
		SingleItems:   singleItems,
	}
}

var configFileDoesntExist = errors.New("configuration file doesn't exist")
var malformedConfigFile = errors.New("malformed config file")

func readConfig(filename string) (config, error) {
	contents, err := os.Open(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return config{}, configFileDoesntExist
		}
		return config{}, err
	}
	var cfg config
	if err := json.NewDecoder(contents).Decode(&cfg); err != nil {
		if errors.Is(err, io.EOF) {
			return config{}, malformedConfigFile
		}
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

	homeDir, err := os.UserHomeDir()
	mustNotBeErr("Home dir must be valid", err)

	cfgPath := fmt.Sprintf("%s/%s", homeDir, ".go_away.json")
	cfg, err := readConfig(cfgPath)
	if errors.Is(err, configFileDoesntExist) {
		fileCreated, newConfig := runMissingConfigFileFlow(cfgPath)
		if !fileCreated {
			fmt.Println("No file created. Make a config file at", cfgPath)
			return
		} else {
			cfg = newConfig
		}
	} else {
		mustNotBeErr("Config must exist and be valid", err)
	}

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

func runMissingConfigFileFlow(path string) (fileCreated bool, cfg config) {
	fmt.Println("Looks like there's no config file at %s", path)
	fmt.Print("Would you like to create one? [y]es/[n]o (default: yes): ")
	for {
		var answer string
		_, err := bufio.NewReader(os.Stdin).ReadString('\n')
		mustNotBeErr("Reading stdin must pass", err)
		answer = strings.ToLower(strings.TrimSpace(answer))
		switch answer {
		case "yes", "ye", "y", "":
			cfg := NewDefaultConfig()
			file, err := os.Create(path)
			mustNotBeErr("Creating a file must succeed", err)
			enc := json.NewEncoder(file)
			enc.SetIndent("", "\t")
			err = enc.Encode(cfg)
			mustNotBeErr("Creating json file must succeed", err)
			return true, cfg
		case "n", "no":
			return false, config{}
		default:
			fmt.Println("Invalid input, please try again or type Ctrl-C to exit")
		}
	}
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
