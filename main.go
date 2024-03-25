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
	inpath := flag.String("input", "", "Input file in fastq format")
	plotpath := flag.String("plot", "", "Output plot file in .pdf or .png format (depending on file extension)")
	flag.Parse()
	if *inpath == "" || *plotpath == "" {
		fmt.Println("You have to specify an input filename!\n")
		flag.Usage()
		os.Exit(1)
	}

	file, err := os.Open(*inpath)
	checkMsg(err, "Could not open file")

	var scanner *bufio.Scanner
	if strings.HasSuffix(*inpath, ".gz") {
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
			lengths = append(lengths, len(line))
		}
		lineNo++
	}
	checkMsg(scanner.Err(), "Error scanning text")
	sort.Ints(lengths)

	//for i, length := range lengths {
	//	lenStr := ""
	//	for l := 0; l < length; l++ {
	//		lenStr = lenStr + "*"
	//	}
	//	fmt.Printf("Length %7d: %s\n", i, lenStr)
	//}

	plotLengths(lengths, *plotpath)
}

func plotLengths(lengths []int, plotPath string) {
	p := plot.New()
	p.Title.Text = "Read length distribution"
	p.X.Label.Text = "Reads"
	p.Y.Label.Text = "Length in bases"

	lengthBars := plotter.Values{}
	for _, l := range lengths {
		lengthBars = append(lengthBars, float64(l))
	}

	w := vg.Points(20)
	bars, err := plotter.NewBarChart(lengthBars, w)
	bars.LineStyle.Width = 0
	bars.Color = plotutil.Color(1)

	checkMsg(err, "Could not create bar chart")

	p.Add(bars)

	err = p.Save(15*vg.Centimeter, 10*vg.Centimeter, plotPath)
	checkMsg(err, "Could not save plot")
}

func checkMsg(err error, msg string) {
	if err != nil {
		fmt.Println(msg, err)
	}
}
