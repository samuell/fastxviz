package main

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
)

func TestReadLengthsFastQ(t *testing.T) {
	testStr := `@XYZ
ACGCGCTCCC
+XYZ
----------
@ABC
ACGCGCTCCCTT
+ABC
----------
`
	reader := strings.NewReader(testStr)
	scanner := bufio.NewScanner(reader)

	haveLengths := readLengthsFastQ(scanner)

	wantLengths := []int{10, 12}
	if !reflect.DeepEqual(haveLengths, wantLengths) {
		t.Fatalf("Wanted %v but got %v\n", wantLengths, haveLengths)
	}
}

func TestReadLengthsFasta(t *testing.T) {
	testStr := `>ABC 1.2.3
ACGCGCTCCC
CGCTTAAACT
>XYZ 2.3.4
CCGTTAATTCGC
TTTTACCGGTCC
`
	reader := strings.NewReader(testStr)
	scanner := bufio.NewScanner(reader)

	haveLengths := readLengthsFasta(scanner)

	wantLengths := []int{20, 24}
	if !reflect.DeepEqual(haveLengths, wantLengths) {
		t.Fatalf("Wanted %v but got %v\n", wantLengths, haveLengths)
	}
}
