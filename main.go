package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		log.Fatal(fmt.Sprintf("Incorrect number of arguments provided. Usage:\n %v input.csv", args[0]))
	}
	file := args[1]
	fmt.Printf("Reading %v\n", file)

	err := ProcessFile(file)
	if err != nil {
		log.Fatal(fmt.Sprintf("Cannot find or open file %v", file))
	}

	WriteResults()
}
