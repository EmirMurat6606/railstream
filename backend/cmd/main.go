// Package main is the entry point of the backend program
package main

import (
	"fmt"
	sloader "github.com/EmirMurat6606/railstream/internal/gtfs"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	feed, err := sloader.ParseStaticGTFS()
	check(err)

	fmt.Println(feed)
}
