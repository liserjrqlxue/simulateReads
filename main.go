package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/liserjrqlxue/goUtil/osUtil"
	"github.com/liserjrqlxue/goUtil/simpleUtil"
)

var (
	length = flag.Int(
		"len",
		35,
		"reads length",
	)
	slide = flag.Int(
		"slide",
		1,
		"slide of reads",
	)
	ref = flag.String(
		"ref",
		"-",
		"ref fasta file, default from stdin",
	)
	output = flag.String(
		"output",
		"-",
		"output file,default to stdout",
	)
)

var (
	isHeader = regexp.MustCompile(`^>`)
)

func main() {
	flag.Parse()

	var in, out *os.File
	if *ref == "-" {
		in = os.Stdin
	} else {
		in = osUtil.Open(*ref)
	}
	if *output == "-" {
		out = os.Stdout
	} else {
		out = osUtil.Create(*output)
	}
	defer simpleUtil.DeferClose(in)
	defer simpleUtil.DeferClose(out)

	var scanner = bufio.NewScanner(in)
	var prefix string
	var seqBuffer string
	var offset int
	for scanner.Scan() {
		var line = scanner.Text()
		if isHeader.MatchString(line) {
			prefix = line[1:]
			seqBuffer = ""
			offset = 0
		} else {
			seqBuffer += line
			var l = len(seqBuffer)
			if l >= *length {
				for i := 0; i+*length <= l; i += *slide {
					printFQ(out, prefix, seqBuffer[i:i+*length], offset+i+1)
				}
				seqBuffer = seqBuffer[l-*length+*slide:]
				offset += l - *length + *slide
			}
		}
	}
}

func printFQ(out *os.File, prefix, seq string, i int) {
	var _, e = fmt.Fprintf(out, "@%s_%d\n%s\n+\n%s\n", prefix, i, seq, seq)
	simpleUtil.CheckErr(e)
}
