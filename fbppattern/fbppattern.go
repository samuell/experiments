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
	hs := NewHiSayer()
	ss := NewStringSplitter()
	lc := NewLowerCaser()
	uc := NewUpperCaser()

	// Network definition
	ss.In = hs.Out
	lc.In = ss.OutLeft
	uc.In = ss.OutRight

	// Set up processes for running (spawn go-routines)
	hs.Init()
	ss.Init()
	lc.Init()
	uc.Init()

	// Drive the processing
	for {
		l, okl := <-lc.Out
		r, okr := <-uc.Out
		if !okl && !okr {
			break
		}
		println(l, r)
	}
	println("Finished program!")
}

// ======= HiGenerator =======

type hiSayer struct {
	Out chan string
}

func NewHiSayer() *hiSayer {
	t := new(hiSayer)
	t.Out = make(chan string, BUFSIZE)
	return t
}

func (t *hiSayer) Init() {
	go func() {
		defer close(t.Out)
		for _, i := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10} {
			t.Out <- fmt.Sprintf("Hi for the %d:th time!", i)
		}
	}()
}

// ======= StringSplitter =======

type stringSplitter struct {
	In       chan string
	OutLeft  chan string
	OutRight chan string
}

func NewStringSplitter() *stringSplitter {
	ss := new(stringSplitter)
	ss.OutLeft = make(chan string, BUFSIZE)
	ss.OutRight = make(chan string, BUFSIZE)
	return ss
}

func (ss *stringSplitter) Init() {
	go func() {
		defer close(ss.OutLeft)
		defer close(ss.OutRight)
		for s := range ss.In {
			halfLen := int(math.Floor(float64(len(s)) / float64(2)))
			ss.OutLeft <- s[0:halfLen]
			ss.OutRight <- s[halfLen:len(s)]
		}
	}()
}

// ======= LowerCaser =======

type lowerCaser struct {
	In  chan string
	Out chan string
}

func NewLowerCaser() *lowerCaser {
	lc := new(lowerCaser)
	lc.Out = make(chan string, BUFSIZE)
	return lc
}

func (lc *lowerCaser) Init() {
	go func() {
		defer close(lc.Out)
		for s := range lc.In {
			lc.Out <- strings.ToLower(s)
		}
	}()
}

// ======= UpperCaser =======

type upperCaser struct {
	In  chan string
	Out chan string
}

func NewUpperCaser() *upperCaser {
	uc := new(upperCaser)
	uc.Out = make(chan string, BUFSIZE)
	return uc
}

func (uc *upperCaser) Init() {
	go func() {
		defer close(uc.Out)
		for s := range uc.In {
			uc.Out <- strings.ToUpper(s)
		}
	}()
}
