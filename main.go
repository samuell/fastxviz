package main

import (
	"bufio"
	"compress/gzip"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

const (
	FASTXVIZ_VERSION = "v0.3.0"
)

func main() {
	inPath := flag.String("in", "", "Input file in FASTA or FastQ format")
	plotType := flag.String("type", "png", "Decide plot type, between: png, pdf and cli")
	plotPath := flag.String("out", "", "Output plot file in .pdf or .png format (depending on file extension). If not specified, it will be the input path, with .png appended.")
	printVersion := flag.Bool("version", false, "prints the current version")
	flag.Parse()

	if *printVersion {
		fmt.Printf("FastXViz %s\n", FASTXVIZ_VERSION)
		os.Exit(1)
	}

	if *inPath == "" || *plotType == "" {
		fmt.Println("You have to specify an input filename!")
		flag.Usage()
		os.Exit(1)
	}

	file, err := os.Open(*inPath)
	checkMsg(err, "Could not open file")

	var scanner *bufio.Scanner
	if strings.HasSuffix(*inPath, ".gz") {
		ungzipper, err := gzip.NewReader(file)
		checkMsg(err, "Could not create GZip reader")
		scanner = bufio.NewScanner(ungzipper)
	} else {
		scanner = bufio.NewScanner(file)
	}

	basePath := strings.Replace(*inPath, ".gz", "", 1)

	var lengths []int
	if strings.HasSuffix(basePath, ".fastq") || strings.HasSuffix(basePath, ".fq") {
		lengths = readLengthsFastQ(scanner)
	} else if strings.HasSuffix(basePath, ".fasta") || strings.HasSuffix(basePath, ".fa") {
		lengths = readLengthsFasta(scanner)
	}

	if *plotType == "cli" {
		for i, length := range lengths {
			lenStr := ""
			for l := 0; l < length; l++ {
				lenStr = lenStr + "*"
			}
			fmt.Printf("%7d [%5d] %s\n", i, length, lenStr)
		}
	} else {
		if *plotPath == "" {
			*plotPath = *inPath + "." + *plotType
		}

		fmt.Println("Plotting ...")
		plotLengths(lengths, *plotPath)
	}
}

func readLengthsFastQ(scanner *bufio.Scanner) []int {
	// Compute length of lines
	fmt.Println("Computing lengths of reads from FastQ file ...")
	lengths := []int{}
	lineNo := 1
	for scanner.Scan() {
		line := scanner.Text()
		if lineNo%4 == 2 {
			lengths = append(lengths, len(line))
		}
		lineNo++
	}
	checkMsg(scanner.Err(), "Error scanning text")
	fmt.Println("Sorting lengths ...")
	sort.Ints(lengths)
	return lengths
}

func readLengthsFasta(scanner *bufio.Scanner) []int {
	// Compute length of lines
	fmt.Println("Computing lengths of reads of Fasta...")
	lengths := []int{}
	lineNo := 1
	currLen := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, ">") {
			if currLen > 0 {
				lengths = append(lengths, currLen)
				currLen = 0
			}
		} else {
			currLen += len(strings.TrimSuffix(line, "\n"))
		}
		lineNo++
	}
	// Do count the last item as well
	if currLen > 0 {
		lengths = append(lengths, currLen)
	}
	checkMsg(scanner.Err(), "Error scanning text")
	fmt.Println("Sorting lengths ...")
	sort.Ints(lengths)
	return lengths
}

func plotLengths(lengths []int, plotPath string) {
	p := plot.New()
	p.Title.Text = "Read length distribution"
	p.X.Label.Text = "Reads"
	p.Y.Label.Text = "Length in bases"

	lengthBars := plotter.Values{}
	fmt.Println("Adding lengths to plot ...")
	for _, l := range lengths {
		lengthBars = append(lengthBars, float64(l))
	}

	fmt.Println("Creating plot ...")
	w := vg.Points(20)
	bars, err := plotter.NewBarChart(lengthBars, w)
	bars.LineStyle.Width = 0
	bars.Color = plotutil.Color(1)

	checkMsg(err, "Could not create bar chart")

	p.Add(bars)

	fmt.Println("Saving plot ...")
	err = p.Save(29*vg.Centimeter, 21*vg.Centimeter, plotPath)
	checkMsg(err, "Could not save plot")
}

func checkMsg(err error, msg string) {
	if err != nil {
		fmt.Println(msg, err)
	}
}
