package main

import (
	"crypto/sha256"
	"fmt"
)

func main() {
	for i := 1; i < 3000; i++ {
		n := sha256.Sum256([]byte(string(i)))
		fmt.Println(n)
		fmt.Println(i)
	}
}
