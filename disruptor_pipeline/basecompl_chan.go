package main

import (
	"bufio"
	"github.com/spf13/afero"
	//"os"
	"fmt"
)

const (
	BUFSIZE   = 16
	FASTAFILE = "Homo_sapiens.GRCh37.67.dna_rm.chromosome.Y.fa"
)

func main() {
	// Init
	frd := NewOsFileReader()
	bc1 := NewBaseComplementer()
	bc2 := NewBaseComplementer()
	bc3 := NewBaseComplementer()
	bc4 := NewBaseComplementer()
	bc5 := NewBaseComplementer()
	bc6 := NewBaseComplementer()
	bc7 := NewBaseComplementer()
	bc8 := NewBaseComplementer()
	prt := NewPrinter()

	// Connect
	frd.Out_Line = bc1.In_FastaLine
	bc1.Out_FastaLine = bc2.In_FastaLine
	bc2.Out_FastaLine = bc3.In_FastaLine
	bc3.Out_FastaLine = bc4.In_FastaLine
	bc4.Out_FastaLine = bc5.In_FastaLine
	bc5.Out_FastaLine = bc6.In_FastaLine
	bc6.Out_FastaLine = bc7.In_FastaLine
	bc7.Out_FastaLine = bc8.In_FastaLine
	bc8.Out_FastaLine = prt.In_Line

	// Run
	go func() {
		defer close(frd.In_FileName)
		frd.In_FileName <- FASTAFILE
	}()
	go frd.Run()
	go bc1.Run()
	go bc2.Run()
	go bc3.Run()
	go bc4.Run()
	go bc5.Run()
	go bc6.Run()
	go bc7.Run()
	go bc8.Run()
	prt.Run()
}

// ------------------------------------------------
// BaseComplementer
// ------------------------------------------------

type BaseComplementer struct {
	In_FastaLine  chan []byte
	Out_FastaLine chan []byte
}

func NewBaseComplementer() *BaseComplementer {
	return &BaseComplementer{
		In_FastaLine:  make(chan []byte, BUFSIZE),
		Out_FastaLine: make(chan []byte, BUFSIZE),
	}
}

var convTable = [256]byte{
	'A':  'T',
	'T':  'A',
	'C':  'G',
	'G':  'C',
	'N':  'N',
	'\n': '\n',
}

func (p *BaseComplementer) Run() {
	defer close(p.Out_FastaLine)

	for line := range p.In_FastaLine {
		if line[0] != '>' {
			for pos := range line {
				line[pos] = convTable[line[pos]]
			}
		}
		p.Out_FastaLine <- line
	}
}

// --------------------------------------------------------------------------------
// FileReader
// --------------------------------------------------------------------------------

type FileReader struct {
	In_FileName chan string
	Out_Line    chan []byte
	fs          afero.Fs
}

func NewOsFileReader() *FileReader {
	return NewFileReader(afero.NewOsFs())
}

func NewFileReader(fileSystem afero.Fs) *FileReader {
	return &FileReader{
		In_FileName: make(chan string, BUFSIZE),
		Out_Line:    make(chan []byte, BUFSIZE),
		fs:          fileSystem,
	}
}

func (p *FileReader) Run() {
	defer close(p.Out_Line)

	for fileName := range p.In_FileName {
		fh, err := p.fs.Open(fileName)
		if err != nil {
			panic(err)
		}
		defer fh.Close()

		sc := bufio.NewScanner(fh)
		for sc.Scan() {
			if err := sc.Err(); err != nil {
				panic(err)
			}
			p.Out_Line <- append([]byte(nil), sc.Bytes()...)
		}
	}
}

// ------------------------------------------------
// Printer
// ------------------------------------------------

type Printer struct {
	In_Line chan []byte
}

func NewPrinter() *Printer {
	return &Printer{
		In_Line: make(chan []byte, BUFSIZE),
	}
}

func (p *Printer) Run() {
	//f := bufio.NewWriter(os.Stdout)
	for line := range p.In_Line {
		//f.Write(line)
		fmt.Println(string(line))
	}
	//f.Flush()
}
