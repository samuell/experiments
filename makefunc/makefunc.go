package main

import (
	"fmt"
)

func main() {
	fn := makeFunc()
	fmt.Println(fn())
}

func makeFunc() func() string {
	return func() string { return "hej" }
}
