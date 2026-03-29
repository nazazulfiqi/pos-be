package main

import (
	"fmt"
	"os"
)

func main() {
	path := "cmd/server/main.go"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("open error:", err)
		return
	}
	defer f.Close()
	b := make([]byte, 16)
	n, _ := f.Read(b)
	fmt.Println("bytes:", b[:n])
	fmt.Println("hex:")
	for i := 0; i < n; i++ {
		fmt.Printf("%02x ", b[i])
	}
	fmt.Println()
}
