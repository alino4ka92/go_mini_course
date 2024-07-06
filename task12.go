package main

import "fmt"

func main() {
	var n int
	fmt.Scanln(&n)
	for ; n >= 1; n-- {
		fmt.Println(n)
	}
}
