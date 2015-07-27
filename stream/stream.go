package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

func main() {
	exb := make([]byte, 0, 128*128)
	buf := bytes.NewBuffer(exb)

	f, err := os.Open("tmp.txt")
	if err != nil {
		panic(err)
	}
	s := bufio.NewScanner(f)
	start := make(chan int)
	go func() {
		<-start
		for s.Scan() {
			if s.Err() != nil {
				panic(err)
			}
			_, err := buf.WriteString(s.Text() + "\n")
			if err != nil {
				panic(err)
			}
		}
	}()

	start <- 1
	s2 := bufio.NewScanner(buf)
	for s2.Scan() {
		if err := s2.Err(); err != nil {
			panic(err)
		}
		fmt.Println("Received line:", s2.Text())
	}
}
