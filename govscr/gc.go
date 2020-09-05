package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	lineChan := make(chan string, 16)

	// ------------------------------------------------------------------------------
	// Loop over the input file in a separate fiber
	// ------------------------------------------------------------------------------
	go func() {
		defer close(lineChan)

		gcFile, err := os.Open("Homo_sapiens.GRCh37.67.dna_rm.chromosome.Y.fa")
		defer gcFile.Close()
		if err != nil {
			panic(err)
		}

		scan := bufio.NewScanner(gcFile)
		for scan.Scan() {
			line := scan.Text()
			lineChan <- line
		}
	}()

	at := 0
	gc := 0

	for line := range lineChan {
		if line[0] == '>' {
			continue
		}

		for _, chr := range line {
			switch chr {
			case 'A', 'T':
				at += 1
				continue
			case 'G', 'C':
				gc += 1
				continue
			}
		}
	}

	var gcFrac float64
	gcFrac = float64(gc) / float64(at+gc)
	fmt.Printf("GC fraction: %f\n", gcFrac)
}
