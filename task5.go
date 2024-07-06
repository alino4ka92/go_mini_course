package main

import "fmt"

func main() {
	var x int
	fmt.Scanln(&x)
	for i := x - 1; i > 0; i-- {
		x *= i
	}
	fmt.Println(x)
}
