package main

import (
	"fmt"
	"gofrec"
	"os"
)

type T1 struct {
	ID   string `Identifier:"001" Length:"3"`
	Name string `Length:"8"`
}

type T2 struct {
	ID   string `Identifier:"002" Length:"3"`
	Name string `Length:"8"`
}

type T3 struct {
	ID     string `Identifier:"003" Length:"3"`
	Number int    `Length:"8"`
}

func main() {
	// First create a list of empty types to be used by the parser: 
	recordTypes := []interface{}{
		T1{},
		T2{},
		T3{},
	}

	// Create a new instance of the parser with identifier start and end
	// as well as the list of empty types
	parser := gofrec.Parser{
		RecordTypes: recordTypes,
		IdStart: 0,
		IdEnd: 3,
	}

	// Will return bytes to be used in BytesToLines
	bytes, err := os.ReadFile("example.txt")
	if err != nil {
		panic(err)
	}

	// This can be bypassed by just passing lines directly to the parser 
	// at initialization
	numLines, err := parser.BytesToLines(bytes)
	if err != nil {
		panic(err)
	}

	fmt.Println("Number of lines:", numLines)

	// Parse the lines into records
	numRecords, err := parser.Parse()
	if err != nil {
		panic(err)
	}

	fmt.Println("Number of parsed records:", numRecords)

	// Loop through records and do stuff with the results
	for _, v := range parser.Records {
		fmt.Println(v)
	}
}
