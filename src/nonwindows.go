//go:build !windows

package main

import "fmt"

func main() {
	fmt.Println("MiniBin is a Windows application. Build with GOOS=windows GOARCH=386.")
}

func executableDir() string { return "." }
