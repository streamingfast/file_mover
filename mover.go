package file_mover

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"path"
	"path/filepath"
	"strings"
)

func Move(folders map[string]string) {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("NewWatcher failed: ", err)
	}
	defer watcher.Close()

	done := make(chan bool)
	newFiles := make(chan string)
	go func() {
		defer close(done)

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op == fsnotify.Create {
					if strings.HasSuffix(event.Name, "jpg") {
						newFiles <- event.Name
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}

	}()

	go func() {
		for {
			select {
			case f := <-newFiles:
				dir := filepath.Dir(f)
				fileName := filepath.Base(f)
				fmt.Println("about to move file:", dir, fileName)

				if dest, ok := folders[dir]; ok {
					err := moveFile(f, path.Join(dest, fileName))
					if err != nil {
						log.Println(fmt.Sprintf("moving file: %s", err))
						continue
					}
					fmt.Println("file move to :", dest, fileName)
				}

			case <-done:
			}
		}
	}()

	for folder, _ := range folders {
		fmt.Printf("About to watch folder: %s\n", folder)
		err = watcher.Add(folder)
		if err != nil {
			log.Fatal(fmt.Sprintf("adding folder %s: %s", folder, err))
		}
	}
	<-done
}
