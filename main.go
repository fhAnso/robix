package main

import (
	"fmt"
	"os"
	"robx/lib"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Syntax: %s <URL>\n", os.Args[0])
		os.Exit(-1)
	}
	targetUrl := os.Args[1]
	if err := lib.ReadRobots(targetUrl); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
