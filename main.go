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
		"",
		"ref fasta",
	)
	output = flag.String(
		"output",
		"",
		"output file",
	)
)

var (
	isHeader = regexp.MustCompile(`^>`)
)

func main() {
	flag.Parse()
	var out = osUtil.Create(*output)
	defer simpleUtil.DeferClose(out)
	var refF = osUtil.Open(*ref)
	defer simpleUtil.DeferClose(refF)
	var scanner = bufio.NewScanner(refF)
	var prefix string
	var seqbuffer string
	var offset int
	for scanner.Scan() {
		var line = scanner.Text()
		if isHeader.MatchString(line) {
			prefix = line[1:]
			seqbuffer = ""
			offset = 0
		} else {
			seqbuffer += line
			var l = len(seqbuffer)
			if l >= *length {
				for i := 0; i+*length <= l; i += *slide {
					printFQ(out, prefix, seqbuffer[i:i+*length], offset+i)
				}
				seqbuffer = seqbuffer[l-*length+*slide:]
				offset += l - *length + *slide
			}
		}
	}
}

func printFQ(out *os.File, prefix, seq string, i int) {
	var _, e = fmt.Fprintf(out, "@%s_%d\n%s\n+\n%s\n", prefix, i, seq, seq)
	simpleUtil.CheckErr(e)
}
