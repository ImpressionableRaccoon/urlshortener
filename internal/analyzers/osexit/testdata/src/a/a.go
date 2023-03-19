package main

import "os"

func main() {
	defer os.Exit(0)
	os.Exit(0) // want "using os.Exit()"
}

func nonMain() {
	os.Exit(0)
}
