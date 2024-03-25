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

func main() {
	inPath := flag.String("in", "", "Input file in FastQ format")
	plotType := flag.String("type", "png", "Decide plot type, between: png, pdf and cli")
	plotPath := flag.String("out", "", "Output plot file in .pdf or .png format (depending on file extension). If not specified, it will be the input path, with .png appended.")
	flag.Parse()

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

	// Compute length of lines
	fmt.Println("Computing lengths of reads ...")
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
	err = p.Save(15*vg.Centimeter, 10*vg.Centimeter, plotPath)
	checkMsg(err, "Could not save plot")
}

func checkMsg(err error, msg string) {
	if err != nil {
		fmt.Println(msg, err)
	}
}
