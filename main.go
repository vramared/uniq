package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]

	fd, err := os.Open(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, "reading stderr: ", err)
	}

	scanner := bufio.NewScanner(fd)
	fmt.Println("Reading from: ", args[0])
	prev := ""
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Compare(line, prev) != 0 {
			fmt.Println(line)
		}
		prev = line
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading stderr: ", err)
	}
}
