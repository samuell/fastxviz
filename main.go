package main

import (
	"bufio"
	"compress/gzip"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {
	infile := flag.String("input", "", "Input file in fastq format")
	flag.Parse()
	if *infile == "" {
		fmt.Println("You have to specify an input filename!\n")
		flag.Usage()
		os.Exit(1)
	}

	file, err := os.Open(*infile)
	checkMsg(err, "Could not open file")

	var scanner *bufio.Scanner
	if strings.HasSuffix(*infile, ".gz") {
		ungzipper, err := gzip.NewReader(file)
		checkMsg(err, "Could not create GZip reader")
		scanner = bufio.NewScanner(ungzipper)
	} else {
		scanner = bufio.NewScanner(file)
	}

	// Compute length of lines
	lengths := []int{}
	lineNo := 1
	for scanner.Scan() {
		line := scanner.Text()
		if lineNo%4 == 2 {
			//fmt.Printf("Line %7d: %s\n", lineNo, line)
			lengths = append(lengths, len(line))
		}
		lineNo++
	}
	checkMsg(scanner.Err(), "Error scanning text")
	sort.Ints(lengths)

	for i, length := range lengths {
		lenStr := ""
		for l := 0; l < length; l++ {
			lenStr = lenStr + "*"
		}
		fmt.Printf("Length %7d: %s\n", i, lenStr)
	}
}

func checkMsg(err error, msg string) {
	if err != nil {
		fmt.Println(msg, err)
	}
}
