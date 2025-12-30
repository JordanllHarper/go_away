package main

import (
	"fmt"
	"strings"
)

type stdoutFormatter struct{ builder *strings.Builder }

func newStdoutFormatter() stdoutFormatter {
	return stdoutFormatter{&strings.Builder{}}
}

func (s stdoutFormatter) FormatSection(header string, count int, items []string) {
	headerLen := len(header)
	fmt.Fprintln(s.builder, header)
	for range headerLen {
		s.builder.WriteRune('-')
	}
	fmt.Fprintln(s.builder)
	fmt.Fprintln(s.builder)

	printStrategy := func(item string) { fmt.Fprintln(s.builder, "-", count, item) }

	if count == 0 {
		printStrategy = func(item string) {
			fmt.Fprintln(s.builder, "-", item)
		}
	}

	for _, item := range items {
		printStrategy(item)
	}

	fmt.Fprintln(s.builder)
}

func (s stdoutFormatter) Output() error {
	_, err := fmt.Println(s.builder)
	return err
}
