package file_mover

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"path"
	"path/filepath"
	"strings"
)

type Mover struct {
	folders map[string]string
	watcher *fsnotify.Watcher
}

func NewMover(folders map[string]string) *Mover {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("NewWatcher failed: ", err)
	}

	for folder, _ := range folders {
		fmt.Printf("About to move folder: %s\n", folder)
		err = watcher.Add(folder)
		if err != nil {
			log.Fatal(fmt.Sprintf("adding folder %s: %s", folder, err))
		}
	}

	return &Mover{
		folders: folders,
		watcher: watcher,
	}
}

func (m *Mover) Move() {
	//todo: cleanup source folder on startup by manually moving file
	err := m.move()
	panic(err) //should never stop
}

func (m *Mover) move() error {
	for {
		select {
		case event, ok := <-m.watcher.Events:
			if !ok {
				return fmt.Errorf("watcher channel closed")
			}
			if event.Op == fsnotify.Create {
				if strings.HasSuffix(event.Name, "jpg") {
					err := m.moveFile(event.Name)
					if err != nil {
						m.handleErr(err)
					}
				}
			}
		case err, ok := <-m.watcher.Errors:
			if !ok {
				return fmt.Errorf("watcher channel closed")
			}
			log.Println("error:", err)
		}
	}
}

func (m *Mover) moveFile(f string) error {
	dir := filepath.Dir(f)
	fileName := filepath.Base(f)

	if dest, ok := m.folders[dir]; ok {
		err := moveFile(f, path.Join(dest, fileName))
		if err != nil {
			return fmt.Errorf("moving file %w", err)
		}
	}
	return nil
}

func (m *Mover) handleErr(err error) {
	panic(err)
}
