package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

const usage = "usage: go run main.go [-c | -d | -u] [-i] [-f fields] [-s chars] [input [output]]"

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

func main() {
	opts := processFlags()
	if !validateOpts(opts) {
		fmt.Println(usage)
		return
	}
	fmt.Println(opts)
	args := os.Args[2:]

	fd, err := os.Open(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, "reading stderr: ", err)
	}

	scanner := bufio.NewScanner(fd)
	prev := ""
	counter := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Compare(line, prev) != 0 && strings.Compare(prev, "") != 0 {
			fmt.Println("  ", counter, prev)
			counter = 1
		} else {
			counter++
		}
		prev = line
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading stderr: ", err)
	}
}
