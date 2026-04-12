// Package main is the entry point of the backend program
package main

import (
	"encoding/json"
	sloader "github.com/EmirMurat6606/railstream/internal/gtfs"
	"os"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	feed, err := sloader.ParseStaticGTFS()
	check(err)

	data, err := json.MarshalIndent(feed, "", "  ")
	check(err)

	err = os.WriteFile("feed.json", data, 0644)
	check(err)

}
