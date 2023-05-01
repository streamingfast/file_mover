package main

import (
	"file_mover"
	"fmt"
	"os"
	"strconv"
)

func main() {

	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) == 0 {
		panic("Expected at least one source and destination folders and destination max size")
	}

	if len(argsWithoutProg)%3 != 0 {
		panic("Wrong number of arguments")
	}

	folders := map[file_mover.SourceFolder]*file_mover.DestinationFolder{}

	for i := 0; i < len(argsWithoutProg); i += 3 {
		destinationMaxSize, err := strconv.Atoi(argsWithoutProg[i+2])
		if err != nil {
			panic(fmt.Sprintf("Failed to parse destination max size: %s", argsWithoutProg[i+2]))
		}
		folders[file_mover.SourceFolder(argsWithoutProg[i])] = file_mover.NewDestinationFolder(argsWithoutProg[i+1], int64(destinationMaxSize))
	}

	mover := file_mover.NewMover(folders)
	err := mover.Move() //blocking call
	panic(err)
}
