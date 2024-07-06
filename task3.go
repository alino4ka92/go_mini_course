package main

import "fmt"

func main() {
	var x int
	fmt.Scanln(&x)
	if x%2 == 0 {
		fmt.Printf("even")
	} else {
		fmt.Printf("odd")
	}
}
