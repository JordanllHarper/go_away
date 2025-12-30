package main

import (
	"fmt"
	"strings"
)

type (
	markdownFormatter struct {
		filename string
		*strings.Builder
	}
	textFormatter struct {
		filename string
		*strings.Builder
	}
)

func newTextFormatter(filename string) textFormatter {
	return textFormatter{
		filename: filename,
		Builder:  &strings.Builder{},
	}
}

func newMarkdownFormatter(filename string) markdownFormatter {
	return markdownFormatter{
		filename: filename,
		Builder:  &strings.Builder{},
	}
}

func (mf markdownFormatter) FormatSection(header string, count int, items []string) {
	fmt.Fprintln(mf, "#", header)
	fmt.Fprintln(mf)

	printStrategy := func(v string) { fmt.Fprintln(mf, "- [ ]", count, v) }

	if count <= 0 {
		printStrategy = func(v string) { fmt.Fprintln(mf, "- [ ]", v) }
	}

	for _, item := range items {
		printStrategy(item)
	}

	fmt.Fprintln(mf)
}

func (tf textFormatter) FormatSection(header string, count int, items []string) {
	fmt.Fprintln(tf, header)
	printStrategy := func(v string) { fmt.Fprintln(tf, count, v) }

	if count <= 0 {
		printStrategy = func(v string) { fmt.Fprintln(tf, v) }
	}

	for _, item := range items {
		printStrategy(item)
	}
	fmt.Fprintln(tf)
}

func (mf markdownFormatter) Output() error { return writeToFile(mf.filename, mf.String()) }
func (tf textFormatter) Output() error     { return writeToFile(tf.filename, tf.String()) }
