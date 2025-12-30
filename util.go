package main

import (
	"log"
	"os"
)

func mustPass(desc string, pred bool) {
	if !pred {
		log.Fatalln(desc, " -> did not pass")
	}
}
func mustNotBeErr(desc string, err error) {
	if err != nil {
		log.Fatalln(desc, "-> Failed with err:", err.Error())
	}
}

func writeToFile(filename string, s string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(s)
	return err
}
