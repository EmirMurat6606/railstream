// Package main is the entry point of the backend program
package main

import (
	gtfs "github.com/EmirMurat6606/railstream/internal/gtfs"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	err := gtfs.StartLoader()
	check(err)

	// Let the program run infinitely long, untill it is killed explicitely
	gtfs.Waitgroup.Wait()
}
