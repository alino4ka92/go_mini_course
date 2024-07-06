package main

import "fmt"

func main() {
	var n int
	fmt.Scanln(&n)
	prime := make([]bool, n+1)
	for i := 0; i <= n; i++ {
		prime[i] = true
	}
	prime[0] = false
	prime[1] = false
	for i := 2; i <= n; i++ {
		if prime[i] {
			for j := i * i; j <= n; j += i {
				prime[j] = false
			}
		}
	}
	for i, v := range prime {
		if v {
			fmt.Println(i)
		}
	}

}
