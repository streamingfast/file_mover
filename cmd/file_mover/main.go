package main

import "file_mover"

func main() {
	folders := map[file_mover.SourceFolder]*file_mover.DestinationFolder{
		file_mover.SourceFolder("/tmp/2k"): file_mover.NewDestinationFolder("/mnt/data/pic", 21474836480),
		file_mover.SourceFolder("/tmp/4k"): file_mover.NewDestinationFolder("/mnt/data/4k", 21474836480),
	}

	//todo: create source and destination folders if they don't exist

	mover := file_mover.NewMover(folders)
	err := mover.Move() //blocking call
	panic(err)
}
