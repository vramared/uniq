package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"uniq/utils"
)

const (
	usage  = "usage: goniq [-c | -d | -u] [-i] [input [output]]"
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

func validateOpts(opts *flags) bool {
	if (opts.count && opts.repeated) || (opts.count && opts.unique) || (opts.repeated && opts.unique) {
		return false
	}
	return true
}

func output(out *os.File, opts *flags, counter int, line string) {
	if opts.count {
		fmt.Fprintf(out, "\t%d %s\n", counter, line)
	} else if (opts.repeated && counter > 1) || (opts.unique && counter == 1) {
		fmt.Fprintf(out, "%s\n", line)
	}
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

func uniq(scanner *bufio.Scanner, opts *flags, out *os.File) {
	prev := ""
	counter := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Compare(line, prev) != 0 && len(prev) != 0 {
			output(out, opts, counter, prev)
			counter = 1
		} else {
			counter++
		}
		prev = line
	}
	output(out, opts, counter, prev)
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading stderr: ", err)
	}
}

func main() {
	opts := processFlags()
	if !validateOpts(opts) {
		fmt.Fprintln(os.Stderr, usage)
		return
	}
	args := os.Args[1:]
	in, out := handleFiles(args[utils.Max(0, len(args)-2):])
	scanner := bufio.NewScanner(in)
	uniq(scanner, opts, out)
}
