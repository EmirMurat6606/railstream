// Package main is the entry point of the backend program
package main

import (
	"time"
	sloader "github.com/EmirMurat6606/railstream/internal/gtfs"
)


func check(err error){
	if (err != nil){
		panic(err)
	}
}

func main() {
	err := sloader.StartLoader()
	check(err)

	time.Sleep(100000000000000000)

}
