package main

import (
	"fmt"
	"strings"
)

func main() {
	var c string
	fmt.Scanln(&c)
	if strings.Contains("aoiuey", strings.ToLower(c)) {
		fmt.Printf("vowel")
	} else {
		fmt.Printf("not vowel")
	}
}
