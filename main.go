package main

import (
	"fmt"
	"os"
)

func checkArgs() bool {
	if len(os.Args) > 2 {
		fmt.Println("Use: ./calc [expr]")
		return false
	}

	return true
}

func main() {
	if !checkArgs() {
		return
	}

	result, err := Calc(os.Args[1])

	if err != nil {
		fmt.Errorf("failed to calc: %v", err)
		return
	}

	fmt.Println(result)
}
