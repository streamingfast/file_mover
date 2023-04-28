package main

import "file_mover"

func main() {
	folders := map[file_mover.SourceFolder]*file_mover.DestinationFolder{
		file_mover.SourceFolder("/tmp/pic"): file_mover.NewDestinationFolder("/mnt/data/pic", 1000000000),
		file_mover.SourceFolder("/tmp/4k"):  file_mover.NewDestinationFolder("/mnt/data/4k", 1000000000),
	}

	mover := file_mover.NewMover(folders)
	err := mover.Move() //blocking call
	panic(err)
}
