package main

import (
	"fmt"
	"math"
	"strings"
)

const (
	BUFSIZE = 16
)

// ======= Main =======

func main() {
	// Init processes
	hisay := NewHiSayer()
	split := NewStringSplitter()
	lower := NewLowerCaser()
	upper := NewUpperCaser()

	// Network definition *** This is where to look! ***
	split.In = hisay.Out
	lower.In = split.OutLeft
	upper.In = split.OutRight

	// Set up processes for running (spawn go-routines)
	go hisay.Run()
	go split.Run()
	go lower.Run()
	go upper.Run()

	// Drive the processing
	for {
		left, okLeft := <-lower.Out
		right, okRight := <-upper.Out
		if !okLeft && !okRight {
			break
		}
		println(left, right)
	}
	println("Finished program!")
}

// ======= HiSayer =======

type hiSayer struct {
	Out chan string
}

func NewHiSayer() *hiSayer {
	t := new(hiSayer)
	t.Out = make(chan string, BUFSIZE)
	return t
}

func (proc *hiSayer) Run() {
	defer close(proc.Out)
	for _, i := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10} {
		proc.Out <- fmt.Sprintf("Hi for the %d:th time!", i)
	}
}

// ======= StringSplitter =======

type stringSplitter struct {
	In       chan string
	OutLeft  chan string
	OutRight chan string
}

func NewStringSplitter() *stringSplitter {
	proc := new(stringSplitter)
	proc.OutLeft = make(chan string, BUFSIZE)
	proc.OutRight = make(chan string, BUFSIZE)
	return proc
}

func (proc *stringSplitter) Run() {
	defer close(proc.OutLeft)
	defer close(proc.OutRight)
	for s := range proc.In {
		halfLen := int(math.Floor(float64(len(s)) / float64(2)))
		proc.OutLeft <- s[0:halfLen]
		proc.OutRight <- s[halfLen:len(s)]
	}
}

// ======= LowerCaser =======

type lowerCaser struct {
	In  chan string
	Out chan string
}

func NewLowerCaser() *lowerCaser {
	proc := new(lowerCaser)
	proc.Out = make(chan string, BUFSIZE)
	return proc
}

func (proc *lowerCaser) Run() {
	defer close(proc.Out)
	for s := range proc.In {
		proc.Out <- strings.ToLower(s)
	}
}

// ======= UpperCaser =======

type upperCaser struct {
	In  chan string
	Out chan string
}

func NewUpperCaser() *upperCaser {
	proc := new(upperCaser)
	proc.Out = make(chan string, BUFSIZE)
	return proc
}

func (proc *upperCaser) Run() {
	defer close(proc.Out)
	for s := range proc.In {
		proc.Out <- strings.ToUpper(s)
	}
}
