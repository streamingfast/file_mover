package file_mover

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type FileInfo struct {
	name             string
	size             int64
	modificationTime time.Time
}

type SourceFolder string

type DestinationFolder struct {
	lock        sync.Mutex
	path        string
	maxSize     int64
	currentSize int64
	files       []*FileInfo
	knowFiles   map[string]bool
}

func NewDestinationFolder(path string, maxSize int64) *DestinationFolder {
	return &DestinationFolder{
		path:      path,
		maxSize:   maxSize,
		knowFiles: map[string]bool{},
	}
}

func (d *DestinationFolder) AddFile(f *FileInfo) {
	d.lock.Lock()
	defer d.lock.Unlock()
	if _, found := d.knowFiles[f.name]; !found {
		d.knowFiles[f.name] = true
		d.files = append(d.files, f)

		d.currentSize += f.size
	}
}

func (d *DestinationFolder) loadInitialState() error {
	files, err := os.ReadDir(d.path)
	if err != nil {
		return fmt.Errorf("reading gps data path: %w", err)
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		i, err := f.Info()
		if err != nil {
			return fmt.Errorf("getting file info: %s , %w", f.Name(), err)
		}

		fi := &FileInfo{
			name:             f.Name(),
			size:             i.Size(),
			modificationTime: i.ModTime(),
		}
		d.AddFile(fi)
	}
	return nil
}

func (d *DestinationFolder) freeUpSpace(nextFileSize int64) error {
	d.lock.Lock()
	defer d.lock.Unlock()

	if d.currentSize+nextFileSize > d.maxSize {
		spaceToReclaim := d.maxSize * 10 / 100
		for spaceToReclaim > 0 {
			fi := d.files[0]
			d.files = d.files[1:]

			fp := path.Join(d.path, fi.name)
			if fileExists(fp) {
				err := os.Remove(fp)
				if err != nil {
					return fmt.Errorf("removing file: %s, %w", fp, err)
				}
			} else {
				log.Println("free space: skipping file that does not exist anymore: ", fp)
			}

			delete(d.knowFiles, fi.name)
			d.currentSize -= fi.size
			spaceToReclaim -= fi.size
		}
	}

	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}

type Mover struct {
	folders map[SourceFolder]*DestinationFolder
	watcher *fsnotify.Watcher
}

func NewMover(folders map[SourceFolder]*DestinationFolder) *Mover {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("NewWatcher failed: ", err)
	}
	return &Mover{
		folders: folders,
		watcher: watcher,
	}
}

func (m *Mover) Move() error {

	for sourceFolder, destination := range m.folders {

		if !fileExists(destination.path) {
			fmt.Printf("Creating sourceFolder: %s\n", destination.path)
			err := os.MkdirAll(destination.path, os.ModePerm)
			if err != nil {
				return fmt.Errorf("creating sourceFolder %s : %w", destination.path, err)
			}
		}

		if !fileExists(string(sourceFolder)) {
			fmt.Printf("Creating sourceFolder: %s\n", sourceFolder)
			err := os.MkdirAll(string(sourceFolder), os.ModePerm)
			if err != nil {
				return fmt.Errorf("creating sourceFolder %s : %w", sourceFolder, err)
			}
		} else {
			go func() {
				files, err := os.ReadDir(string(sourceFolder))
				if err != nil {
					panic(fmt.Errorf("reading sourceFolder %s : %w", sourceFolder, err))
				}
				for _, file := range files {
					if file.IsDir() {
						continue
					}
					fmt.Println("Moving existing file: ", file.Name())
					err := m.moveFile(path.Join(string(sourceFolder), file.Name()))
					if err != nil {
						panic(fmt.Errorf("moving existing file %s : %w", file.Name(), err))
					}
				}
			}()
		}

		fmt.Printf("About to move sourceFolder: %s\n", sourceFolder)
		err := m.watcher.Add(string(sourceFolder))
		if err != nil {
			log.Fatal(fmt.Sprintf("adding sourceFolder %s: %s", sourceFolder, err))
		}

		err = destination.loadInitialState()
		if err != nil {
			return fmt.Errorf("loading initial state of destination %s : %w", destination.path, err)
		}
	}

	err := m.move() //blocking call
	if err != nil {
		return fmt.Errorf("moving files: %w", err)
	}

	return nil
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

	if destination, ok := m.folders[SourceFolder(dir)]; ok {
		if stat, err := os.Stat(f); err == nil {

			fp := path.Join(destination.path, fileName)
			if fileExists(f) {
				destination.AddFile(&FileInfo{
					name:             fileName,
					size:             stat.Size(),
					modificationTime: stat.ModTime(),
				})
				err := destination.freeUpSpace(stat.Size())
				if err != nil {
					return fmt.Errorf("freeing space: %w", err)
				}

				err = moveFile(f, fp)
				if err != nil {
					return fmt.Errorf("moving file %w", err)
				}
			} else {
				log.Println("move: skipping file that does not exist anymore: ", fp)
			}

		}
	}
	return nil
}

func (m *Mover) handleErr(err error) {
	panic(err)
}
