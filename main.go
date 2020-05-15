package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"uniq/pkg/utils"
)

const (
	usage  = "usage: goniq [-c | -d | -u] [-i] [-f fields] [-s chars] [input [output]]"
	hyphen = "-"
)

type flags struct {
	count    bool
	repeated bool
	unique   bool
}

func processFlags() *flags {
	count := flag.Bool("c", false, "count mode")
	repeated := flag.Bool("d", false, "output only repeated lines")
	unique := flag.Bool("u", false, "output only unique lines")
	flag.Parse()
	return &flags{*count, *repeated, *unique}
}

func validateOpts(f *flags) bool {
	if (f.count && f.repeated) || (f.count && f.unique) || (f.repeated && f.unique) {
		return false
	}
	return true
}

func output(f *flags, counter int, line string) string {
	res := ""
	if f.count {
		res += "   " + strconv.Itoa(counter) + " "
	}
	res += line
	return res
}

func handleFiles(args []string) (*os.File, *os.File) {
	input, output := os.Stdin, os.Stdout
	inputFile, outputFile := "", ""
	if strings.Contains(args[0], hyphen) {
		inputFile = args[1]
	} else {
		inputFile = args[0]
		if len(args) > 1 {
			outputFile = args[1]
		}
	}
	if len(inputFile) != 0 {
		var err error
		input, err = os.Open(inputFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error opening file: ", err)
		}
	}
	if len(outputFile) != 0 {
		var err error
		output, err = os.Open(outputFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error opening file: ", err)
		}
	}
	return input, output
}

func main() {
	opts := processFlags()
	if !validateOpts(opts) {
		fmt.Println(usage)
		return
	}
	args := os.Args[1:]

	in, out := handleFiles(args[utils.Max(0, len(args)-2):])

	scanner := bufio.NewScanner(in)
	prev := ""
	counter := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Compare(line, prev) != 0 && strings.Compare(prev, "") != 0 {
			fmt.Fprintln(out, output(opts, counter, prev))
			counter = 1
		} else {
			counter++
		}
		prev = line
	}
	fmt.Fprintln(out, output(opts, counter, prev))
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading stderr: ", err)
	}
}
