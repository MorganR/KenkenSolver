package main

import (
	"fmt"
	"kenken"
)

func main() {
	size := getSize()
	p := kenken.RequestPuzzle(size)
	p.Solve()
	p.Print()
}

func getSize() uint8 {
	var size uint8
	for {
		fmt.Println("Enter the size of your grid (eg. for 5x5, enter 5):")
		fmt.Scan(&size)
		if size > 0 {
			break
		}
		fmt.Println("ERROR: Size must be between 1 and 255")
	}
	return size
}
