package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/spf13/afero"
	"os"
	"runtime"
	"sync/atomic"
)

const (
	BUFSIZE = 128
)

func main() {
	// Parse flags
	inFileNamePtr := flag.String("in", "", "The input file name")
	flag.Parse()

	inFileName := *inFileNamePtr
	if inFileName == "" {
		fmt.Println("No filename specified to --in")
		os.Exit(1)
	}

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
		frd.In_FileName <- string(inFileName)
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
	In_FastaLine  *DisruptorChan
	Out_FastaLine *DisruptorChan
}

func NewBaseComplementer() *BaseComplementer {
	return &BaseComplementer{
		In_FastaLine:  NewDisruptorChan(),
		Out_FastaLine: NewDisruptorChan(),
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
	defer p.Out_FastaLine.Close()

	for {
		line := p.In_FastaLine.Recv()
		if line == nil {
			break
		}
		l := line.Content
		if l[0] != '>' {
			for pos := range l {
				l[pos] = convTable[l[pos]]
			}
		}
		p.Out_FastaLine.Send(NewFastaLine(l))
	}
}

// --------------------------------------------------------------------------------
// FileReader
// --------------------------------------------------------------------------------

type FileReader struct {
	In_FileName chan string
	Out_Line    *DisruptorChan
	fs          afero.Fs
}

func NewOsFileReader() *FileReader {
	return NewFileReader(afero.NewOsFs())
}

func NewFileReader(fileSystem afero.Fs) *FileReader {
	return &FileReader{
		In_FileName: make(chan string, BUFSIZE),
		Out_Line:    NewDisruptorChan(),
		fs:          fileSystem,
	}
}

func (p *FileReader) Run() {
	defer p.Out_Line.Close()

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
			p.Out_Line.Send(NewFastaLine(sc.Bytes()))
		}
	}
}

// ------------------------------------------------
// Printer
// ------------------------------------------------

type Printer struct {
	In_Line *DisruptorChan
}

func NewPrinter() *Printer {
	return &Printer{
		In_Line: NewDisruptorChan(),
	}
}

func (p *Printer) Run() {
	for {
		line := p.In_Line.Recv()
		if line == nil {
			break
		}
		l := line.Content
		fmt.Println(string(l))
	}
}

// ------------------------------------------------
// DisruptorChan
// ------------------------------------------------
// The code below is adapted from Gringo (https://github.com/textnode/gringo)
// by Darren Elwood (@textnode)
// Which is licensed under Apache 2.0 (http://www.apache.org/licenses/LICENSE-2.0)
//
// Changes made:
// - Use *FastaFinder as data type instead of *PayLoad
// - Added some closing logic
// - Some naming changes
// ------------------------------------------------

// Masking is faster than division
const idxMask uint64 = BUFSIZE - 1

type FastaLine struct {
	Content []byte
}

func NewFastaLine(newContent []byte) *FastaLine {
	return &FastaLine{
		Content: newContent,
	}
}

type DisruptorChan struct {
	pad1             [8]uint64
	lastCommittedIdx uint64
	pad2             [8]uint64
	nextFreeIdx      uint64
	pad3             [8]uint64
	recvIdx          uint64
	pad4             [8]uint64
	contents         [BUFSIZE]*FastaLine
	pad5             [8]uint64
	open             bool
}

func NewDisruptorChan() *DisruptorChan {
	return &DisruptorChan{lastCommittedIdx: 0, nextFreeIdx: 1, recvIdx: 1, open: true}
}

func (c *DisruptorChan) Send(value *FastaLine) {
	var anIdx = atomic.AddUint64(&c.nextFreeIdx, 1) - 1
	for anIdx > (c.recvIdx + BUFSIZE - 2) {
		runtime.Gosched()
	}
	c.contents[anIdx&idxMask] = value
	for !atomic.CompareAndSwapUint64(&c.lastCommittedIdx, anIdx-1, anIdx) {
		runtime.Gosched()
	}
}

func (c *DisruptorChan) Recv() *FastaLine {
	var anIdx = atomic.AddUint64(&c.recvIdx, 1) - 1
	if anIdx > c.lastCommittedIdx && !c.IsOpen() {
		return nil
	}
	for anIdx > c.lastCommittedIdx {
		runtime.Gosched()
	}
	return c.contents[anIdx&idxMask]
}

func (c *DisruptorChan) IsOpen() bool {
	return c.open
}

func (c *DisruptorChan) Close() {
	c.open = false
}
