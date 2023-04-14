package main

import "file_mover"

func main() {
	folders := map[string]string{"/tmp/pic": "/mnt/data/pic", "/tmp/4k": "/mnt/data/4k"}

	mover := file_mover.NewMover(folders)
	mover.Move()
}
